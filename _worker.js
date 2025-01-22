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
      reqHeader.set("User-Agent", "okhttp");
      reqHeader.set("Accept-Encoding", request.headers.get("Accept-Encoding") || "");
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
      newResponse.headers.set("Cache-Control", "s-maxage=3600");
      newResponse.headers.set("Last-Modified", new Date().toUTCString());
      newResponse.headers.set("Content-Type", "text/html; charset=utf-8");
      if (response.headers.get("Content-Encoding")) {
        newResponse.headers.set("Content-Encoding", response.headers.get("Content-Encoding"));
      }

      ctx.waitUntil(cache.put(cacheKey, newResponse.clone()));
      console.log("Cache new request");
      return newResponse;
    } catch (error) {
      return new Response(`Invalid remote URL: ${error.message}`, { status: 400 });
    }
  },
};

