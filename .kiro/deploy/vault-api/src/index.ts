import { S3Client, PutObjectCommand, GetObjectCommand, ListObjectsV2Command, DeleteObjectCommand } from "@aws-sdk/client-s3";
import Fastify from "fastify";
import cors from "@fastify/cors";

const fastify = Fastify({ logger: true });

await fastify.register(cors, { origin: true });

const s3Client = new S3Client({
  endpoint: process.env.S3_ENDPOINT || "http://localhost:3900",
  region: process.env.S3_REGION || "us-east-1",
  credentials: {
    accessKeyId: process.env.S3_ACCESS_KEY_ID || "garage",
    secretAccessKey: process.env.S3_SECRET_ACCESS_KEY || "garage",
  },
  forcePathStyle: true,
});

const BUCKET = process.env.S3_BUCKET || "study-vault";

fastify.get("/", async () => ({ status: "ok", service: "vault-mcp-api" }));

fastify.get("/health", async () => ({ status: "healthy" }));

fastify.post<{ Params: { "*": string }; Body: { content: string } }>(
  "/vault/:path(*)",
  async (request, reply) => {
    const key = request.params.path;
    const { content } = request.body || {};
    
    if (!content) {
      return reply.status(400).send({ error: "content required" });
    }

    await s3Client.send(
      new PutObjectCommand({
        Bucket: BUCKET,
        Key: key,
        Body: content,
        ContentType: "text/markdown",
      })
    );

    return { success: true, key };
  }
);

fastify.get<{ Params: { "*": string } }>("/vault/:path(*)", async (request) => {
  const key = request.params.path;

  try {
    const response = await s3Client.send(
      new GetObjectCommand({ Bucket: BUCKET, Key: key })
    );
    const body = await response.Body?.transformToString();
    return { key, content: body };
  } catch (e: any) {
    if (e.name === "NoSuchKey") {
      return { error: "not found", key };
    }
    throw e;
  }
});

fastify.get("/vault", async () => {
  const response = await s3Client.send(
    new ListObjectsV2Command({ Bucket: BUCKET })
  );
  return { files: response.Contents?.map((c) => c.Key) || [] };
});

fastify.delete<{ Params: { "*": string } }>("/vault/:path(*)", async (request) => {
  const key = request.params.path;
  await s3Client.send(new DeleteObjectCommand({ Bucket: BUCKET, Key: key }));
  return { success: true, deleted: key };
});

const start = async () => {
  try {
    const port = parseInt(process.env.PORT || "10001");
    await fastify.listen({ port, host: "0.0.0.0" });
    console.log(`Vault API listening on port ${port}`);
  } catch (err) {
    fastify.log.error(err);
    process.exit(1);
  }
};

start();
