package m3u8

import (
	"context"
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

func TestProxyCacheExpiry(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	params := ProxyParams{Url: "http://example.com/expired"}
	key := calcCacheKey(params)
	holder := &cacheHolder{lastReadTime: time.Now().Add(-CacheNoReadExpiredDuration), cancel: cancel}
	proxyCacheMu.Lock()
	proxyCacheTasks[key] = holder
	proxyCacheMu.Unlock()
	maintainCache(ctx, key, params, holder)
	proxyCacheMu.Lock()
	defer proxyCacheMu.Unlock()
	if len(proxyCacheTasks) != 0 || ctx.Err() == nil {
		t.Fatal("expired worker did not release capacity and cancel its context")
	}
}

func TestProxyCacheEvictsLeastRecentlyRead(t *testing.T) {
	canceled := make(chan struct{})
	refreshing := make(chan struct{})
	var bRequests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/b" && bRequests.Add(1) == 2 {
			close(refreshing)
			<-r.Context().Done()
			close(canceled)
			return
		}
		fmt.Fprint(w, "#EXTM3U")
	}))
	defer server.Close()
	defer func() {
		proxyCacheMu.Lock()
		defer proxyCacheMu.Unlock()
		for key, holder := range proxyCacheTasks {
			holder.cancel()
			delete(proxyCacheTasks, key)
		}
	}()
	params := func(name string) ProxyParams { return ProxyParams{Url: server.URL + "/" + name} }
	pushRead := func(name string) {
		t.Helper()
		if err := PushCacheTask(params(name)); err != nil {
			t.Fatal(err)
		}
		if _, _, err := ReadCacheContent(params(name)); err != nil {
			t.Fatal(err)
		}
	}
	for _, name := range []string{"a", "b", "c"} {
		pushRead(name)
	}
	// Reading the oldest-created task makes b the least recently read.
	if _, _, err := ReadCacheContent(params("a")); err != nil {
		t.Fatal(err)
	}
	select {
	case <-refreshing:
	case <-time.After(5 * time.Second):
		t.Fatal("refresh did not start")
	}
	// A duplicate at capacity must neither evict nor renew the read timestamp.
	if err := PushCacheTask(params("b")); err != nil {
		t.Fatal(err)
	}
	pushRead("d")
	proxyCacheMu.Lock()
	count := len(proxyCacheTasks)
	_, hasB := proxyCacheTasks[calcCacheKey(params("b"))]
	proxyCacheMu.Unlock()
	if count != 3 || hasB {
		t.Fatalf("expected 3 tasks without b, count=%d hasB=%v", count, hasB)
	}
	select {
	case <-canceled:
	case <-time.After(5 * time.Second):
		t.Fatal("eviction did not cancel the in-flight refresh")
	}
	// An evicted subscription can start again on the next request.
	pushRead("b")
	proxyCacheMu.Lock()
	count = len(proxyCacheTasks)
	_, hasC := proxyCacheTasks[calcCacheKey(params("c"))]
	proxyCacheMu.Unlock()
	if count != 3 || hasC {
		t.Fatalf("expected replacement to evict c, count=%d hasC=%v", count, hasC)
	}
}

func TestOldProxyWorkerDoesNotDeleteReplacement(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	params := ProxyParams{Url: "http://example.com/replaced"}
	key := calcCacheKey(params)
	old := &cacheHolder{lastReadTime: time.Now().Add(-CacheNoReadExpiredDuration), cancel: cancel}
	replacement := &cacheHolder{lastReadTime: time.Now()}
	proxyCacheMu.Lock()
	proxyCacheTasks[key] = replacement
	proxyCacheMu.Unlock()
	maintainCache(ctx, key, params, old)
	proxyCacheMu.Lock()
	defer proxyCacheMu.Unlock()
	defer delete(proxyCacheTasks, key)
	if proxyCacheTasks[key] != replacement {
		t.Fatal("old worker deleted its replacement")
	}
}
