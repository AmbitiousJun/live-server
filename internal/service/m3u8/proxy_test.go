package m3u8

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

func TestCalcCacheKey(t *testing.T) {
	a := ProxyParams{Url: "https://example.com/ab", Header: http.Header{"A": {"one", "two"}, "B": {"three"}}}
	b := ProxyParams{Url: a.Url, Header: http.Header{"B": {"three"}, "A": {"one", "two"}}}
	if calcCacheKey(a) != calcCacheKey(b) {
		t.Fatal("header insertion order changed the cache key")
	}
	b.Url = "https://example.com/ba"
	if calcCacheKey(a) == calcCacheKey(b) {
		t.Fatal("different URLs share a key")
	}
	b.Url = a.Url
	b.Header["A"] = []string{"two", "one"}
	if calcCacheKey(a) == calcCacheKey(b) {
		t.Fatal("header value order was lost")
	}
}

func TestReadCacheContentTimeout(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		start := time.Now()
		if _, _, err := ReadCacheContent(ProxyParams{Url: "missing"}); err == nil {
			t.Fatal("expected timeout")
		}
		if time.Since(start) != 10*time.Second {
			t.Fatalf("unexpected wait: %v", time.Since(start))
		}
	})
}

func waitForProxy(t *testing.T, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for !condition() {
		if time.Now().After(deadline) {
			t.Fatal("timed out waiting for cache worker")
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func TestProxyCacheLifecycle(t *testing.T) {
	var requests atomic.Int32
	var fail atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := requests.Add(1)
		if r.Header.Get("X-Test") != "original" {
			t.Error("worker did not retain its own headers")
		}
		if fail.Load() {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		fmt.Fprintf(w, "#EXTM3U\n# version %d\n", n)
	}))
	defer server.Close()
	params := ProxyParams{Url: server.URL, Header: http.Header{"X-Test": {"original"}}}
	key := calcCacheKey(params)
	t.Cleanup(func() {
		proxyCacheMu.Lock()
		if holder, ok := proxyCacheTasks[key]; ok {
			holder.lastReadTime = time.Now().Add(-CacheNoReadExpiredDuration)
		}
		proxyCacheMu.Unlock()
		waitForProxy(t, func() bool {
			proxyCacheMu.Lock()
			defer proxyCacheMu.Unlock()
			_, ok := proxyCacheTasks[key]
			return !ok
		})
	})
	var wg sync.WaitGroup
	for range 100 {
		wg.Go(func() {
			if err := PushCacheTask(params); err != nil {
				t.Error(err)
			}
		})
	}
	wg.Wait()
	_, first, err := ReadCacheContent(params)
	if err != nil || first == "" {
		t.Fatalf("first read: %q, %v", first, err)
	}
	proxyCacheMu.Lock()
	count := len(proxyCacheTasks)
	proxyCacheMu.Unlock()
	if count != 1 || requests.Load() != 1 {
		t.Fatalf("duplicate workers: tasks=%d requests=%d", count, requests.Load())
	}
	params.Header.Set("X-Test", "modified")
	readParams := ProxyParams{Url: server.URL, Header: http.Header{"X-Test": {"original"}}}
	waitForProxy(t, func() bool {
		_, content, err := ReadCacheContent(readParams)
		return err == nil && content != first
	})
	fail.Store(true)
	_, previous, err := ReadCacheContent(readParams)
	if err != nil {
		t.Fatal(err)
	}
	before := requests.Load()
	waitForProxy(t, func() bool { return requests.Load() > before })
	_, content, err := ReadCacheContent(readParams)
	if err != nil || content != previous {
		t.Fatalf("failed refresh discarded cache: %q, %v", content, err)
	}
	proxyCacheMu.Lock()
	lastRead := proxyCacheTasks[key].lastReadTime
	proxyCacheMu.Unlock()
	if time.Since(lastRead) > time.Second {
		t.Fatal("read did not renew cache lifetime")
	}
}

func TestProxyCacheWorkerLimitAndExpiry(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// Use already-expired holders so workers exit without performing network IO.
		proxyCacheMu.Lock()
		for i := range CacheMaintainWorkerMaxCount {
			params := ProxyParams{Url: fmt.Sprintf("http://example.com/%d", i)}
			key := calcCacheKey(params)
			holder := &cacheHolder{lastReadTime: time.Now().Add(-CacheNoReadExpiredDuration)}
			proxyCacheTasks[key] = holder
		}
		proxyCacheMu.Unlock()
		if err := PushCacheTask(ProxyParams{Url: "http://example.com/overflow"}); err == nil {
			t.Error("accepted a task above the worker limit")
		}
		if err := PushCacheTask(ProxyParams{Url: "http://example.com/0"}); err != nil {
			t.Errorf("duplicate at capacity: %v", err)
		}
		for i := range CacheMaintainWorkerMaxCount {
			params := ProxyParams{Url: fmt.Sprintf("http://example.com/%d", i)}
			key := calcCacheKey(params)
			maintainCache(key, params, proxyCacheTasks[key])
		}
		if len(proxyCacheTasks) != 0 {
			t.Fatal("expired workers did not release capacity")
		}
	})
}
