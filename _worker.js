export default {
  async fetch(request, env, ctx) {
    const url = new URL(request.url);
    const remoteParam = url.searchParams.get("remote");

    // 检查 `remote` 参数是否存在
    if (!remoteParam) {
      return new Response("Empty remote", { status: 400 });
    }

    try {
      // Base64 解码 `remote` 参数
      const remoteUrl = new URL(atob(remoteParam));
      const clientHeaders = atob(url.searchParams.get("headers") || '').split('[[[:]]]');

      // 获取 Cache 对象
      const cacheKey = new Request(remoteUrl.toString(), request);
      const cache = caches.default;

      // 检查缓存是否命中
      const cachedResponse = await cache.match(cacheKey);
      if (cachedResponse) {
        console.log("Cache hit");
        return cachedResponse;
      }

      // 未命中缓存，发起代理请求
      const reqHeader = new Headers();
      reqHeader.set("User-Agent", "libmpv");
      reqHeader.set("Accept-Encoding", request.headers.get("Accept-Encoding") || "");
      for (let i = 0; i + 1 < clientHeaders.length; i += 2) {
        reqHeader.set(clientHeaders[i], clientHeaders[i+1]);
      }

      console.log(`Client Accept-Encoding header: ${reqHeader.get("Accept-Encoding")}`);
      const body = request.method === 'GET' || request.method === 'HEAD' ? null : request.body;
      const response = await fetch(remoteUrl, {
        method: request.method,
        headers: reqHeader,
        body,
      });

      // 确保响应体是可缓存的
      if (!response.ok || !response.body) {
        return new Response("Failed to fetch remote URL", { status: 500 });
      }

      const newResponse = new Response(response.body, response);
      // 设置 CORS 头
      newResponse.headers.set("Access-Control-Allow-Origin", "*");
      newResponse.headers.set("Access-Control-Allow-Methods", "GET,HEAD,POST,OPTIONS");
      newResponse.headers.set("Access-Control-Allow-Headers", "Content-Type");
      newResponse.headers.set("Cache-Control", "s-maxage=300");
      newResponse.headers.set("Last-Modified", new Date().toUTCString());
      newResponse.headers.set("Content-Type", "text/html");
      if (response.headers.get("Content-Encoding")) {
        newResponse.headers.set("Content-Encoding", response.headers.get("Content-Encoding"));
        console.log(`Server Content-Encoding header: ${newResponse.headers.get("Accept-Encoding")}`);
      }

      ctx.waitUntil(cache.put(cacheKey, newResponse.clone()));
      // cache.put(cacheKey, newResponse.clone());
      console.log("Cache new request");
      return newResponse;
    } catch (error) {
      console.error(`代理异常：${error}`);
      return new Response(`Invalid remote URL: ${error.message}`, { status: 400 });
    }
  },
};


