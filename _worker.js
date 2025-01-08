export default {
  async fetch(request) {
    const url = new URL(request.url);
    const remoteParam = url.searchParams.get("remote");

    // 检查 `remote` 参数是否存在
    if (!remoteParam) {
      return new Response("Empty remote", { status: 400 });
    }

    try {
      // Base64 解码 `remote` 参数
      const remoteUrl = atob(remoteParam);

      // 获取 Cache 对象
      const cache = caches.default;

      // 检查缓存是否命中
      const cacheKey = new Request(request.url, request);
      const cachedResponse = await cache.match(cacheKey);
      if (cachedResponse) {
        console.log("Cache hit");
        return cachedResponse;
      }

      // 未命中缓存，发起代理请求
      const reqHeader = new Headers();
      reqHeader.set("User-Agent", "okhttp");
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

      // 设置 CORS 头
      const headers = new Headers(response.headers);
      headers.set("Access-Control-Allow-Origin", "*");
      headers.set("Access-Control-Allow-Methods", "GET,HEAD,POST,OPTIONS");
      headers.set("Access-Control-Allow-Headers", "Content-Type");
      headers.set("Cache-Control", "max-age=3600");
      headers.set("Last-Modified", new Date().toUTCString());
      headers.set("Content-Type", "text/html; charset=utf-8");

      const proxyResponse = new Response(response.body, {
        status: response.status,
        headers,
      });

      // 将响应缓存
      console.log("Cache miss, caching response");
      await cache.put(cacheKey, proxyResponse.clone());

      return proxyResponse;
    } catch (error) {
      return new Response(`Invalid remote URL: ${error.message}`, { status: 400 });
    }
  },
};

