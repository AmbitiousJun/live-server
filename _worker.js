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

      // 验证解码后的结果是否是有效 URL
      const targetUrl = new URL(remoteUrl);

      // 代理请求到目标 URL
      request.headers.set("User-Agent", "libmpv");
      const response = await fetch(targetUrl.toString(), {
        method: request.method,
        headers: request.headers,
        body: request.body,
      });

      // 返回目标 URL 的响应
      const headers = new Headers(response.headers);
      headers.set("Access-Control-Allow-Origin", "*");
      headers.set("Access-Control-Allow-Methods", "GET,HEAD,POST,OPTIONS");
      headers.set("Access-Control-Allow-Headers", "Content-Type");
      return new Response(response.body, {
        status: response.status,
        headers,
      });
    } catch (error) {
      return new Response(`Invalid remote URL: ${error.message}`, { status: 400 });
    }
  },
};
