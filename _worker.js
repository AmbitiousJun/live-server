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

      // 代理请求到目标 URL
      const reqHeader = new Headers();
      reqHeader.set("User-Agent", "okhttp");
      const response = await fetch(remoteUrl, {
        method: request.method,
        headers: reqHeader,
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
