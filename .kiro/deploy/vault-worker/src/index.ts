export default {
  async fetch(request: Request, env: VAULT_ENV): Promise<Response> {
    const url = new URL(request.url);

    if (url.pathname === "/") {
      return Response.json({ status: "ok", service: "vault-api", deployed: "cloudflare" });
    }

    if (url.pathname === "/health") {
      return Response.json({ status: "healthy" });
    }

    if (url.pathname === "/vault" && request.method === "GET") {
      const list = await env.VAULT_STORAGE.list();
      return Response.json({ files: list.keys.map((k) => k.name) });
    }

    const path = url.pathname.replace("/vault/", "");

    if (path && url.pathname.startsWith("/vault/")) {
      if (request.method === "GET") {
        const value = await env.VAULT_STORAGE.get(path);
        if (!value) {
          return Response.json({ error: "not found", key: path }, { status: 404 });
        }
        return Response.json({ key: path, content: value });
      }

      if (request.method === "POST") {
        const { content } = await request.json();
        if (!content) {
          return Response.json({ error: "content required" }, { status: 400 });
        }
        await env.VAULT_STORAGE.put(path, content);
        return Response.json({ success: true, key: path });
      }

      if (request.method === "DELETE") {
        await env.VAULT_STORAGE.delete(path);
        return Response.json({ success: true, deleted: path });
      }
    }

    return Response.json({ error: "not found" }, { status: 404 });
  },
} satisfies ExportedHandler<VAULT_ENV>;

interface VAULT_ENV {
  VAULT_STORAGE: KVNamespace;
}
