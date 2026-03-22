// ============================================================
// BAC Study RAG System — Cloudflare Worker
// Architecture: Supervisor + Specialist Agents via RAG endpoints
// Stack: Workers AI (embeddings + synthesis), Upstash Vector (search), D1 (storage)
// Spaced Repetition: FSRS v4 (ts-fsrs)
// ============================================================

import { fsrs, createEmptyCard, Rating, State } from 'ts-fsrs';

const f = fsrs();

export interface Env {
  DB: D1Database;
  AI: Ai;
  UPSTASH_URL: string;
  UPSTASH_TOKEN: string;
  JINA_API_KEY: string;
  OPENROUTER_API_KEY?: string;
  MISTRAL_API_KEY?: string;
  YANDEX_IAM_TOKEN?: string;
  YANDEX_FOLDER_ID?: string;
}

async function runSynthesis(
  messages: { role: string; content: string }[],
  env: Env,
  maxTokens = 1024,
): Promise<{ text: string; source: string }> {
  const { AI, OPENROUTER_API_KEY } = env;

  try {
    const response = await AI.run("@cf/meta/llama-3.2-3b-instruct", { messages, max_tokens: maxTokens });
    return { text: (response as { response: string }).response, source: "workers-ai" };
  } catch (primaryErr) {
    console.warn(`Workers AI failed (${primaryErr}), trying OpenRouter...`);
  }

  if (OPENROUTER_API_KEY) {
    try {
      const resp = await fetch("https://openrouter.ai/api/v1/chat/completions", {
        method: "POST",
        headers: {
          "Authorization": `Bearer ${OPENROUTER_API_KEY}`,
          "Content-Type": "application/json",
          "HTTP-Referer": "https://bac-api.workers.dev",
          "X-Title": "BAC Study RAG",
        },
        body: JSON.stringify({
          model: "google/gemini-2.0-flash-lite-001",
          messages,
          max_tokens: maxTokens,
        }),
      });
      const data = await resp.json() as {
        choices?: {
          message?: {
            content?: string | null;
            refusal?: string | null;
            reasoning?: string | null;
            reasoning_details?: { text?: string }[];
          };
        }[];
        error?: { message?: string };
      };
      const msg = data.choices?.[0]?.message;
      let text = msg?.content;
      if (!text && msg?.reasoning) {
        text = msg.reasoning;
      }
      if (!text && msg?.reasoning_details?.length) {
        text = msg.reasoning_details.map((r) => r.text || "").join("\n");
      }
      if (text) {
        return { text, source: "openrouter" };
      }
      if (msg?.refusal) {
        console.warn(`OpenRouter refusal: ${msg.refusal}`);
      }
      if (data.error) {
        console.warn(`OpenRouter error: ${data.error.message}`);
      }
    } catch (e) {
      console.warn(`OpenRouter call failed: ${e}`);
    }
  }

  return { text: "", source: "offline" };
}

// ============================================================
// HELPERS
// ============================================================

function extractEmbedding(dataObj: Record<string, number>): number[] {
  return Object.values(dataObj);
}

type WorkersAIEmbedding = { data: Record<string, number>[]; shape: number[] };

const SUBJECT_MAP: Record<string, string> = {
  "biologie": "biology",
  "biology": "biology",
  "chimie": "chemistry",
  "chemistry": "chemistry",
  "physique": "physics",
  "physics": "physics",
  "mathématiques": "math",
  "maths": "math",
  "math": "math",
  "physiologie": "biology",
};
// ============================================================

const CORS = {
  "Access-Control-Allow-Origin": "*",
  "Access-Control-Allow-Methods": "GET, POST, OPTIONS",
  "Access-Control-Allow-Headers": "Content-Type, Authorization",
};

function json(data: unknown, status = 200) {
  return new Response(JSON.stringify(data), {
    status,
    headers: { "Content-Type": "application/json", ...CORS },
  });
}

function workerError(msg: string, status = 400) {
  return json({ error: msg }, status);
}

async function upstashQuery(
  url: string,
  token: string,
  vector: number[],
  topK = 5,
  subject?: string
) {
  const body: Record<string, unknown> = { vector, topK, includeMetadata: true, includeVectors: false };
  if (subject) body.filter = { subject };

  const res = await fetch(`${url}/query`, {
    method: "POST",
    headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  return res.json() as Promise<UpstashQueryResult>;
}

async function upstashUpsert(url: string, token: string, vectors: UpstashVector[]) {
  const payload = vectors.map((v) => ({
    id: v.id,
    vector: v.vector.map((n) => parseFloat(n.toFixed(6))),
    metadata: v.metadata,
  }));
  const res = await fetch(`${url}/upsert`, {
    method: "POST",
    headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
  return res.json() as Promise<{ status: string } | { error: string }>;
}

// ============================================================
// TYPES
// ============================================================

interface UpstashVector {
  id: string;
  vector: number[];
  metadata: {
    subject: string;
    path: string;
    chunkIndex: number;
    text: string;
  };
}

interface UpstashMatch {
  id: string;
  score: number;
  metadata: { subject: string; path: string; chunkIndex: number; text: string };
}

interface UpstashQueryResult {
  result: UpstashMatch[];
}

// ============================================================
// ROUTE — classify query to subject(s)
// ============================================================

async function routeQuery(query: string, env: Env): Promise<{ subjects: string[]; confidence: number }> {
  const messages = [
    {
      role: "system",
      content: 'Classify into: biology, chemistry, math, physics, general. ' +
        'Respond ONLY with JSON: {"subjects":["..."],"confidence":0.9}',
    },
    { role: "user", content: query },
  ];

  const { text, source } = await runSynthesis(messages, env, 64);

  try {
    const parsed = JSON.parse(text.trim());
    return { subjects: parsed.subjects || ["general"], confidence: parsed.confidence || 0.5 };
  } catch {
    const q = query.toLowerCase();
    const kw: Record<string, string[]> = {
      biology: ["dna", "cell", "organism", "photosynth\u00e8se", "mitose", "m\u00e9iose", "g\u00e8ne", "enzyme", "biologie", "svt"],
      chemistry: ["chimie", "mol\u00e9cule", "atome", "raction", "p\u00e9riodique", "acide", "base", "ph", "solution"],
      math: ["quation", "fonction", "int\u00e9gral", "d\u00e9riv\u00e9", "limite", "alg\u00e8bre", "g\u00e9om\u00e9trie", "probabilit\u00e9", "math"],
      physics: ["force", "\u00e9nergie", "physique", "mouvement", "vitesse", "acc\u00e9l\u00e9ration", "newton", "\u00e9lectricit\u00e9", "onde", "thermodynamique"],
    };
    const found = Object.entries(kw).filter(([, kws]) => kws.some((k) => q.includes(k))).map(([s]) => s);
    return { subjects: found.length > 0 ? found : ["general"], confidence: found.length > 0 ? 0.6 : 0.3 };
  }
}

// ============================================================
// EXPAND — generate search variants
// ============================================================

async function expandQuery(query: string, env: Env): Promise<{ variants: string[] }> {
  const messages = [
    {
      role: "system",
      content: "Generate 2-3 search-friendly phrasings. JSON: {\"variants\":[\"...\"]}",
    },
    { role: "user", content: query },
  ];

  const { text } = await runSynthesis(messages, env, 128);

  try {
    const parsed = JSON.parse(text.trim());
    return { variants: parsed.variants || [query] };
  } catch {
    return { variants: [query] };
  }
}

// ============================================================
// SEARCH — embed + vector search
// ============================================================

async function embedWithFallback(texts: string[], env: Env): Promise<number[][] | null> {
  const { AI, JINA_API_KEY } = env;
  try {
    const r = await AI.run("@cf/baai/bge-m3", { text: texts }) as WorkersAIEmbedding;
    return texts.map((_, i) => extractEmbedding(r.data[i]));
  } catch (e) {
    console.warn(`Workers AI embed failed (${e}), trying Jina...`);
  }
  if (JINA_API_KEY) {
    try {
      const resp = await fetch("https://api.jina.ai/v1/embeddings", {
        method: "POST",
        headers: { Authorization: `Bearer ${JINA_API_KEY}`, "Content-Type": "application/json" },
        body: JSON.stringify({ model: "jina-embeddings-v2-base-en", input: texts }),
      });
      const data = await resp.json() as { data?: { embedding: number[] }[] };
      if (data.data) {
        const vecs = data.data.map(e => e.embedding);
        return vecs.map(v => v.concat(new Array(1024 - v.length).fill(0)));
      }
    } catch (e2) {
      console.warn(`Jina embed failed: ${e2}`);
    }
  }
  return null;
}

async function searchChunks(
  variants: string[],
  subjects: string[],
  env: Env,
  topK = 10
): Promise<{ query: string; chunks: UpstashMatch[]; total: number }> {
  const { UPSTASH_URL, UPSTASH_TOKEN } = env;

  const embeddings = await embedWithFallback(variants, env);
  if (!embeddings) {
    return { query: variants[0], chunks: [], total: 0 };
  }

  const seen = new Map<string, UpstashMatch>();
  const allSubjects = new Set<string>();
  for (const s of subjects) {
    allSubjects.add(s);
    if (SUBJECT_MAP[s]) allSubjects.add(SUBJECT_MAP[s]);
  }
  for (const english of Object.values(SUBJECT_MAP)) allSubjects.add(english);

  for (const embedding of embeddings as number[][]) {
    for (const subject of allSubjects) {
      try {
        const result = await upstashQuery(UPSTASH_URL, UPSTASH_TOKEN, embedding, topK, subject);
        for (const m of result.result) seen.set(m.id, m);
      } catch (e) { console.error(`Upstash query failed: ${e}`); }
    }
    try {
      const result = await upstashQuery(UPSTASH_URL, UPSTASH_TOKEN, embedding, topK);
      for (const m of result.result) { if (!seen.has(m.id)) seen.set(m.id, m); }
    } catch (e) { console.error(`Upstash query failed: ${e}`); }
  }

  const unique = Array.from(seen.values()).sort((a, b) => b.score - a.score);

  // Rerank top 20 with Jina if available
  let final = unique.slice(0, 5);
  if (unique.length > 5 && env.JINA_API_KEY) {
    try {
      const rerankTexts = unique.slice(0, 20).map((m) => m.metadata.text);
      const rRes = await fetch("https://api.jina.ai/v1/rerank", {
        method: "POST",
        headers: { Authorization: `Bearer ${env.JINA_API_KEY}`, "Content-Type": "application/json" },
        body: JSON.stringify({ model: "jina-reranker-v1-base-en", query: variants[0], documents: rerankTexts, top_n: 5 }),
      });
      const rData = await rRes.json() as { results: { index: number }[] };
      final = rData.results.map((r) => unique[r.index]).filter(Boolean);
    } catch (e) { console.error(`Rerank failed: ${e}`); }
  }

  return { query: variants[0], chunks: final, total: unique.length };
}

// ============================================================
// QUERY — full RAG with synthesis
// ============================================================

async function queryVault(
  query: string,
  subjects: string[],
  variants: string[],
  env: Env,
): Promise<{ answer: string; sources: { source: string; subject: string; score: number }[]; model: string; subjects: string[] }> {
  const search = await searchChunks(variants, subjects, env);

  if (search.chunks.length === 0) {
    return {
      answer: "Je n'ai pas trouv\u00e9 d'informations pertinentes dans la base de connaissances.",
      sources: [],
      model: "none",
      subjects,
    };
  }

  const context = search.chunks.map((c, i) =>
    `[${i + 1}] ${c.metadata.path} (${c.metadata.subject}):\n${c.metadata.text}`
  ).join("\n\n");

  const { text, source } = await runSynthesis([
    { role: "system", content: "Tu es un assistant p\u00e9dagogique pour les \u00e9tudiants du BAC en Mauritanie." },
    { role: "user", content: `Contexte:\n${context}\n\nQuestion: ${query}\n\nR\u00e9ponds en citant tes sources [1], [2], etc.` },
  ], env, 1024);

  return {
    answer: text,
    sources: search.chunks.map((c) => ({ source: c.metadata.path, subject: c.metadata.subject, score: c.score })),
    model: source,
    subjects,
  };
}

// ============================================================
// ADD — ingest new content
// ============================================================

async function addContent(text: string, subject: string, docPath: string, env: Env): Promise<{ chunksAdded: number; message: string }> {
  const { AI, UPSTASH_URL, UPSTASH_TOKEN } = env;

  const CHUNK_SIZE = 2000;
  const CHUNK_OVERLAP = 200;
  const chunks: string[] = [];
  for (let start = 0; start < text.length; start += CHUNK_SIZE - CHUNK_OVERLAP) {
    const chunk = text.slice(start, start + CHUNK_SIZE);
    if (chunk.trim()) chunks.push(chunk);
  }
  if (chunks.length === 0) return { chunksAdded: 0, message: "No content" };

  const vectors: UpstashVector[] = [];
  const batchSize = 25;
  for (let i = 0; i < chunks.length; i += batchSize) {
    const batch = chunks.slice(i, i + batchSize);
    try {
      const r = await AI.run("@cf/baai/bge-m3", { text: batch }) as WorkersAIEmbedding;
      for (let j = 0; j < batch.length; j++) {
        vectors.push({
          id: `${subject}-${docPath.replace(/\//g, "-")}-${i + j}`,
          vector: extractEmbedding(r.data[j]),
          metadata: { subject, path: docPath, chunkIndex: i + j, text: batch[j] },
        });
      }
    } catch (e) { return { chunksAdded: 0, message: `Embedding failed: ${String(e)}` }; }
  }

  if (vectors.length === 0) return { chunksAdded: 0, message: "No vectors generated" };

  try {
    const result = await upstashUpsert(UPSTASH_URL, UPSTASH_TOKEN, vectors) as { status?: string; error?: string };
    if ("error" in result) return { chunksAdded: 0, message: `Upsert error: ${result.error}` };
  } catch (e) {
    return { chunksAdded: 0, message: `Upsert failed: ${e}` };
  }

  return { chunksAdded: vectors.length, message: `Indexed ${vectors.length} chunks` };
}

// ============================================================
// PRACTICE — spaced repetition + grading
// ============================================================

interface CardProgress {
  id: number;
  stability: number;
  difficulty: number;
  state: string;
  reps: number;
  lapses: number;
  interval: number;
  due_date: string | null;
  elapsed_days: number;
  scheduled_days: number;
  last_attempt: string | null;
}

function sm2StateToFSRS(s: string): State {
  switch (s) {
    case "new": return State.New;
    case "learning": case "again": return State.Learning;
    case "review": return State.Review;
    case "relearning": return State.Learning;
    default: return State.New;
  }
}

function fsrsStateToStr(s: State): string {
  switch (s) {
    case State.New: return "new";
    case State.Learning: return "learning";
    case State.Review: return "review";
    case State.Relearning: return "relearning";
  }
}

function smScoreToRating(score: number): number {
  if (score <= 2) return Rating.Again;
  if (score === 3) return Rating.Hard;
  if (score === 4) return Rating.Good;
  return Rating.Easy;
}

interface Question {
  id: number;
  question_text: string;
  solution_text: string | null;
  topic_tags: string | null;
  difficulty: number;
}

function buildFSRSCard(p: CardProgress, now: Date): { card: ReturnType<typeof createEmptyCard>; delta_t: number } {
  if (!p.last_attempt) {
    const card = createEmptyCard(now);
    card.stability = p.stability;
    card.difficulty = p.difficulty;
    card.state = State.New;
    card.reps = 0;
    card.lapses = 0;
    card.scheduled_days = 0;
    return { card, delta_t: 0 };
  }

  const lastReview = new Date(p.last_attempt);
  const delta_t = Math.max(0, Math.round((now.getTime() - lastReview.getTime()) / 86400000));

  const card = createEmptyCard(now);
  card.due = new Date(p.due_date || now);
  card.stability = p.stability;
  card.difficulty = p.difficulty;
  card.state = sm2StateToFSRS(p.state);
  card.reps = p.reps;
  card.lapses = p.lapses;
  card.elapsed_days = delta_t;
  card.scheduled_days = p.interval;
  card.last_review = lastReview;

  return { card, delta_t };
}

function applyFSRSResult(r: ReturnType<typeof f.next>, now: Date) {
  const due_date = r.card.due.toISOString();
  const elapsed_days = Math.max(0, Math.round((now.getTime() - (r.card.last_review?.getTime() || now.getTime())) / 86400000));
  return {
    stability: r.card.stability,
    difficulty: r.card.difficulty,
    interval: r.card.scheduled_days,
    reps: r.card.reps,
    lapses: r.card.lapses,
    state: fsrsStateToStr(r.card.state),
    due_date,
    elapsed_days,
    scheduled_days: r.card.scheduled_days,
  };
}

async function getPracticeQuestion(userId: number, subject: string | undefined, db: D1Database): Promise<Question | null> {
  const now = new Date().toISOString();
  const q = subject ? SUBJECT_MAP[subject] || subject : undefined;
  const tagFilter = q ? ` AND topic_tags LIKE '%${q}%'` : "";

  const due = await db.prepare(
    `SELECT q.* FROM questions q
     INNER JOIN user_progress up ON q.id = up.question_id AND up.user_id = ?
     WHERE up.due_date IS NOT NULL AND up.due_date <= ?${tagFilter}
     ORDER BY up.due_date ASC LIMIT 1`
  ).bind(userId, now).first<Question>();

  if (due) return due;

  const newQ = await db.prepare(
    `SELECT q.* FROM questions q
     WHERE q.id NOT IN (SELECT question_id FROM user_progress WHERE user_id = ?)${tagFilter}
     ORDER BY RANDOM() LIMIT 1`
  ).bind(userId).first<Question>();

  if (newQ) return newQ;

  const fallback = await db.prepare(
    `SELECT * FROM questions WHERE difficulty <= 2${tagFilter ? ` AND topic_tags LIKE '%${q}%'` : ''} ORDER BY RANDOM() LIMIT 1`
  ).first<Question>();

  return fallback || null;
}

function extractNumbers(str: string): number[] {
  const nums = str.replace(/,/g, "").match(/-?\d+\.?\d*/g);
  if (!nums) return [];
  return nums.map(n => parseFloat(n)).filter(n => !isNaN(n) && isFinite(n));
}

function extractFinalNumber(str: string): number | null {
  const cleaned = str.replace(/,/g, "");
  const nums = cleaned.match(/-?\d+\.?\d*/g);
  if (!nums || nums.length === 0) return null;
  const last = nums[nums.length - 1];
  const parsed = parseFloat(last);
  return isFinite(parsed) ? parsed : null;
}

function numMatch(a: number, b: number): boolean {
  if (a === b) return true;
  const avg = (Math.abs(a) + Math.abs(b)) / 2;
  if (avg === 0) return a === 0 && b === 0;
  return Math.abs(a - b) / avg < 0.01;
}

function digitsOnly(s: string): string {
  return s.replace(/\D/g, "");
}

function checkAnswer(userAnswer: string, solution: string): { score: number; feedback: string; partial: boolean } {
  const ua = userAnswer.trim();
  const sol = solution.trim();
  if (!ua || !sol) return { score: 3, feedback: "Réponse enregistrée.", partial: false };

  const userDigits = digitsOnly(ua);
  const solDigits = digitsOnly(sol);
  if (userDigits && solDigits) {
    const userLast6 = userDigits.slice(-6);
    if (solDigits.includes(userLast6)) return { score: 5, feedback: "Parfait ! Réponse exacte.", partial: false };
  }

  const finalSol = extractFinalNumber(sol);
  const userNums = extractNumbers(ua);
  const solNums = extractNumbers(sol);

  if (finalSol !== null && userNums.length > 0) {
    const found = userNums.some(u => numMatch(u, finalSol));
    if (found) return { score: 5, feedback: "Parfait ! Réponse exacte.", partial: false };

    if (solNums.length > 1) {
      let matches = 0;
      for (const u of userNums) {
        if (solNums.some(s => numMatch(u, s))) matches++;
      }
      const ratio = matches / solNums.length;
      if (ratio >= 0.8) return { score: 4, feedback: `Très bien !`, partial: false };
      if (ratio >= 0.4) return { score: 3, feedback: `Partiellement correct.`, partial: true };
    }
  }

  const solLower = sol.toLowerCase().replace(/[^\w\s]/g, " ");
  const uaLower = ua.toLowerCase().replace(/[^\w\s]/g, " ");
  const solWords = solLower.split(/\s+/).filter(w => w.length > 2);
  const uaWords = uaLower.split(/\s+/).filter(w => w.length > 2);
  if (solWords.length > 0 && uaWords.length > 0) {
    let matches = 0;
    for (const w of uaWords) {
      if (solWords.includes(w)) matches++;
    }
    const ratio = matches / solWords.length;
    if (ratio >= 0.7) return { score: 4, feedback: "Bonne réponse !", partial: false };
    if (ratio >= 0.4) return { score: 3, feedback: "Partiellement correct.", partial: true };
    if (ratio >= 0.2) return { score: 2, feedback: "Quelques éléments corrects.", partial: true };
  }

  return { score: 1, feedback: "Réponse incorrecte.", partial: true };
}

async function gradeAnswer(questionId: number, userAnswer: string, question: Question, env: Env): Promise<{ score: number; feedback: string; partial: boolean }> {
  const solution = question.solution_text || "Pas de solution disponible.";

  try {
    const messages = [
      { role: "system", content: "Tu es un évaluateur pédagogique strict mais juste. Réponds uniquement en JSON." },
      { role: "user", content: `Question: ${question.question_text}\nSolution: ${solution}\nRéponse de l'élève: ${userAnswer}\n\nÉvalue sur 0-5 et réponds au format JSON: {"score": N, "feedback": "commentaire"}` },
    ];
    const { text } = await runSynthesis(messages, env, 256);
    const jsonMatch = text.match(/\{[\s\S]*?\}/);
    if (jsonMatch) {
      const parsed = JSON.parse(jsonMatch[0]);
      return {
        score: Math.max(0, Math.min(5, Number(parsed.score) || 3)),
        feedback: parsed.feedback || "Réponse enregistrée.",
        partial: false,
      };
    }
  } catch (e) {
    // fallback to pattern matching
  }

  return checkAnswer(userAnswer, solution);
}

// ============================================================
// SOLVE — auto-solve questions
// ============================================================

async function solveQuestion(question: string, subject: string, env: Env) {
  const messages = [
    { role: "system", content: "Tu es un assistant p\u00e9dagogique expert. R\u00e9sous \u00e9tape par \u00e9tape en fran\u00e7ais. Utilise le LaTeX pour les formules." },
    { role: "user", content: `[${subject}] ${question}` },
  ];
  const { text } = await runSynthesis(messages, env, 2048);
  return { solution: text };
}

// ============================================================
// OCR — MULTI-PROVIDER FALLBACK CHAIN
// ============================================================

interface OCRResult {
  text: string;
  source: string;
  confidence?: number;
  chars: number;
  error?: string;
  latency_ms?: number;
  provider: string;
}

async function ocrMistral(imageBase64: string, apiKey: string): Promise<OCRResult> {
  const start = Date.now();
  try {
    const cleanBase64 = imageBase64.replace(/^data:image\/[^;]+;base64,/, "");
    const resp = await fetch("https://api.mistral.ai/v1/ocr", {
      method: "POST",
      headers: {
        "Authorization": `Bearer ${apiKey}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        model: "mistral-ocr-latest",
        document: {
          type: "image_url",
          image_url: `data:image/jpeg;base64,${cleanBase64}`,
        },
      }),
    });

    if (!resp.ok) {
      const errText = await resp.text();
      return { text: "", source: "mistral", chars: 0, provider: "mistral", error: `HTTP ${resp.status}: ${errText}`, latency_ms: Date.now() - start };
    }

    const data = await resp.json() as {
      pages?: { markdown?: string }[];
      error?: string;
    };

    if (data.error) {
      return { text: "", source: "mistral", chars: 0, provider: "mistral", error: data.error, latency_ms: Date.now() - start };
    }

    const text = data.pages?.map(p => p.markdown || "").join("\n\n") || "";
    return { text, source: "mistral", chars: text.length, provider: "mistral", latency_ms: Date.now() - start };
  } catch (e) {
    return { text: "", source: "mistral", chars: 0, provider: "mistral", error: String(e), latency_ms: Date.now() - start };
  }
}

async function ocrOCRspace(imageBase64: string, apiKey = "helloworld"): Promise<OCRResult> {
  const start = Date.now();
  try {
    const cleanBase64 = imageBase64.replace(/^data:image\/[^;]+;base64,/, "");
    const formData = new FormData();
    formData.append("base64Image", `data:image/jpeg;base64,${cleanBase64}`);
    formData.append("language", "fra"); // French + auto-detect
    formData.append("isOverlayRequired", "false");
    formData.append("detectOrientation", "true");
    formData.append("scale", "true");
    formData.append("OCREngine", "2"); // Engine 2 = more accurate

    const resp = await fetch("https://api.ocr.space/parse/image", {
      method: "POST",
      headers: apiKey && apiKey !== "helloworld" ? { "apikey": apiKey } : {},
      body: formData,
    });

    if (!resp.ok) {
      return { text: "", source: "ocrspace", chars: 0, provider: "ocrspace", error: `HTTP ${resp.status}`, latency_ms: Date.now() - start };
    }

    const data = await resp.json() as {
      ParsedResults?: { ParsedText: string }[];
      ErrorMessage?: string[];
      IsErroredOnProcessing?: boolean;
    };

    if (data.IsErroredOnProcessing || (data.ErrorMessage && data.ErrorMessage.length > 0)) {
      return { text: "", source: "ocrspace", chars: 0, provider: "ocrspace", error: (data.ErrorMessage || []).join(", "), latency_ms: Date.now() - start };
    }

    const text = data.ParsedResults?.map(r => r.ParsedText).join("\n") || "";
    return { text, source: "ocrspace", chars: text.length, provider: "ocrspace", latency_ms: Date.now() - start };
  } catch (e) {
    return { text: "", source: "ocrspace", chars: 0, provider: "ocrspace", error: String(e), latency_ms: Date.now() - start };
  }
}

async function ocrYandex(imageBase64: string, iamToken: string, folderId: string): Promise<OCRResult> {
  const start = Date.now();
  try {
    const cleanBase64 = imageBase64.replace(/^data:image\/[^;]+;base64,/, "");
    const resp = await fetch("https://ocr.api.cloud.yandex.net/ocr/v1/recognizeText", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${iamToken}`,
        "x-folder-id": folderId,
      },
      body: JSON.stringify({
        mimeType: "JPEG",
        languageCodes: ["fr", "ar", "en"],
        model: "page",
        content: cleanBase64,
      }),
    });

    if (!resp.ok) {
      const errText = await resp.text();
      return { text: "", source: "yandex", chars: 0, provider: "yandex", error: `HTTP ${resp.status}: ${errText}`, latency_ms: Date.now() - start };
    }

    const data = await resp.json() as {
      result?: {
        textAnnotation?: { words?: { text?: string }[]; fullText?: string };
        pages?: { blocks?: { paragraphs?: { words?: { text?: string }[] }[] }[] }[];
      };
      error?: { message?: string };
    };

    if (data.error?.message) {
      return { text: "", source: "yandex", chars: 0, provider: "yandex", error: data.error.message, latency_ms: Date.now() - start };
    }

    const text = data.result?.textAnnotation?.fullText ||
      data.result?.textAnnotation?.words?.map(w => w.text).join(" ") ||
      "";
    return { text, source: "yandex", chars: text.length, provider: "yandex", latency_ms: Date.now() - start };
  } catch (e) {
    return { text: "", source: "yandex", chars: 0, provider: "yandex", error: String(e), latency_ms: Date.now() - start };
  }
}

async function ocrLlamaVision(imageData: string, ai: Ai): Promise<OCRResult> {
  const start = Date.now();
  try {
    const base64 = imageData.startsWith("data:") ? imageData : `data:image/jpeg;base64,${imageData}`;
    const messages = [
      {
        role: "user",
        content: [
          { type: "text", text: "Extract ALL text from this image word-for-word. Preserve French (é, è, ê, à, ç, etc.) and Arabic characters exactly. Maintain line breaks. If no readable text is found, say 'No text detected'." },
          { type: "image_url", image_url: { url: base64 } }
        ]
      }
    ];

    const result = await ai.run("@cf/meta/llama-3.2-11b-vision-instruct", { messages, max_tokens: 1024 }) as { response?: string };
    const text = result.response || "";

    if (text.toLowerCase().includes("no text detected")) {
      return { text: "", source: "llama-vision", chars: 0, provider: "llama-vision", error: "No text detected", latency_ms: Date.now() - start };
    }

    return { text, source: "llama-vision", chars: text.length, provider: "llama-vision", latency_ms: Date.now() - start };
  } catch (e) {
    return { text: "", source: "llama-vision", chars: 0, provider: "llama-vision", error: String(e), latency_ms: Date.now() - start };
  }
}

async function ocrUform(imageData: string, ai: Ai): Promise<OCRResult> {
  const start = Date.now();
  try {
    const base64 = imageData.startsWith("data:") ? imageData.replace(/^data:image\/[^;]+;base64,/, "") : imageData;
    const result = await ai.run("@cf/unum/uform-gen2-qwen-500m", {
      image: imageData,
      prompt: "Extract all text from this image word for word. Preserve French (é, è, ê, à, ç) and Arabic characters exactly. Maintain line breaks.",
    }) as { response?: string };

    const text = result.response || "";
    return { text, source: "uform", chars: text.length, provider: "uform", latency_ms: Date.now() - start };
  } catch (e) {
    return { text: "", source: "uform", chars: 0, provider: "uform", error: String(e), latency_ms: Date.now() - start };
  }
}

async function runOCRChain(imageBase64: string, env: Env): Promise<{ results: OCRResult[]; best: OCRResult | null }> {
  const results: OCRResult[] = [];

  // Provider chain: Mistral → OCR.space → Yandex → LlamaVision → UForm
  const providers = [
    { fn: async () => {
      if (!env.MISTRAL_API_KEY) return null;
      return await ocrMistral(imageBase64, env.MISTRAL_API_KEY);
    }, name: "mistral" },
    { fn: async () => await ocrOCRspace(imageBase64), name: "ocrspace" },
    { fn: async () => {
      if (!env.YANDEX_IAM_TOKEN || !env.YANDEX_FOLDER_ID) return null;
      return await ocrYandex(imageBase64, env.YANDEX_IAM_TOKEN, env.YANDEX_FOLDER_ID);
    }, name: "yandex" },
    { fn: async () => await ocrLlamaVision(imageBase64, env.AI), name: "llama-vision" },
    { fn: async () => await ocrUform(imageBase64, env.AI), name: "uform" },
  ];

  for (const { fn, name } of providers) {
    const result = await fn();
    if (result) {
      results.push(result);
      // Return first successful result with meaningful text (>= 5 chars)
      if (result.text && result.text.length >= 5) {
        return { results, best: result };
      }
    }
  }

  return { results, best: results.find(r => r.text && r.text.length >= 5) || null };
}

// ============================================================
// MAIN FETCH
// ============================================================

export default {
  async fetch(request: Request, env: Env): Promise<Response> {
    const url = new URL(request.url);
    const path = url.pathname;

    if (request.method === "OPTIONS") return new Response(null, { headers: CORS });

    try {
      // HEALTH
      if (path === "/rag/health") {
        try {
          const r = await fetch(`${env.UPSTASH_URL}/info`, { headers: { Authorization: `Bearer ${env.UPSTASH_TOKEN}` } });
          const d = await r.json() as { result?: { vectorCount: number } };
          return json({ status: "ok", upstash_vectors: d.result?.vectorCount || 0, workers_ai: "bound" });
        } catch { return json({ status: "degraded", workers_ai: "bound" }); }
      }

      // STATUS
      if (path === "/rag/status") {
        try {
          const r = await fetch(`${env.UPSTASH_URL}/info`, { headers: { Authorization: `Bearer ${env.UPSTASH_TOKEN}` } });
          const d = await r.json() as { result?: { vectorCount: number; dimension: number } };
          return json({ total_vectors: d.result?.vectorCount || 0, dimension: d.result?.dimension || 0 });
        } catch (e) { return workerError(`Upstash error: ${e}`, 500); }
      }

      // SUBJECTS — per-subject vector counts from Upstash
      if (path === "/rag/subjects" && request.method === "GET") {
        const placeholder = new Array(1024).fill(1 / 32);
        const counts: Record<string, number> = { subject_biology: 0, subject_chemistry: 0, subject_math: 0, subject_physics: 0 };
        try {
          const r = await fetch(`${env.UPSTASH_URL}/query`, {
            method: "POST",
            headers: { Authorization: `Bearer ${env.UPSTASH_TOKEN}`, "Content-Type": "application/json" },
            body: JSON.stringify({ vector: placeholder, topK: 1000, includeMetadata: true }),
          });
          const d = await r.json() as { result?: { metadata?: { subject?: string } }[]; error?: string };
          if (d.error) {
            console.error("Upstash error:", d.error);
          }
          for (const item of d.result || []) {
            const subj = item.metadata?.subject;
            if (subj && counts[`subject_${subj}`] !== undefined) {
              counts[`subject_${subj}`]++;
            }
          }
        } catch (e) { console.error("Subjects query failed:", e); }
        return json(counts);
      }

      // ROUTE
      if (path === "/rag/route" && request.method === "POST") {
        const { query } = await request.json() as { query?: string };
        if (!query) return workerError("query required");
        return json(await routeQuery(query, env));
      }

      // EXPAND
      if (path === "/rag/expand" && request.method === "POST") {
        const { query } = await request.json() as { query?: string };
        if (!query) return workerError("query required");
        return json(await expandQuery(query, env));
      }

      // SEARCH
      if (path === "/rag/search" && request.method === "POST") {
        const body = await request.json() as { query?: string; subjects?: string[]; topK?: number };
        if (!body.query) return workerError("query required");
        const route = body.subjects?.length
          ? { subjects: body.subjects, confidence: 1 }
          : await routeQuery(body.query, env);
        const expand = await expandQuery(body.query, env);
        const result = await searchChunks(expand.variants, route.subjects, env, body.topK || 10);
        return json({ ...result, subjects: route.subjects });
      }

      // QUERY — full RAG (falls back to direct LLM if embeddings unavailable)
      if (path === "/rag/query" && request.method === "POST") {
        const body = await request.json() as { query?: string; subjects?: string[]; direct?: boolean };
        if (!body.query) return workerError("query required");

        const route = body.subjects?.length
          ? { subjects: body.subjects, confidence: 1 }
          : await routeQuery(body.query, env);

        try {
          const expand = await expandQuery(body.query, env);
          const result = await queryVault(body.query, route.subjects, expand.variants, env);
          return json(result);
        } catch (e) {
          // Fallback: direct LLM answer without RAG context
          if (body.direct !== false) {
            const { text, source } = await runSynthesis([
              { role: "system", content: "Tu es un assistant pédagogique pour les étudiants du BAC en Mauritanie. Réponds de manière concise et précise en français." },
              { role: "user", content: body.query },
            ], env, 1024);
            return json({ answer: text, sources: [], model: source, subjects: route.subjects, fallback: true });
          }
          return workerError(`Query failed: ${e}`, 500);
        }
      }

      // ADD
      if (path === "/rag/add" && request.method === "POST") {
        const body = await request.json() as { text?: string; subject?: string; path?: string };
        if (!body.text) return workerError("text required");
        if (!body.subject) return workerError("subject required");
        return json(await addContent(body.text, body.subject, body.path || "manual", env));
      }

      // ADD-BATCH
      if (path === "/rag/add-batch" && request.method === "POST") {
        const body = await request.json() as { chunks?: { text: string; subject: string; path: string }[] };
        if (!body.chunks?.length) return workerError("chunks array required");
        const { AI, UPSTASH_URL, UPSTASH_TOKEN } = env;
        const vectors: UpstashVector[] = [];
        const batchSize = 25;

        const toEmbed = body.chunks.map((c, i) => ({ text: c.text, subject: c.subject, path: c.path, chunkIndex: i }));
        if (toEmbed.length === 0) return json({ chunksAdded: 0, message: "No content" });

        for (let i = 0; i < toEmbed.length; i += batchSize) {
          const batch = toEmbed.slice(i, i + batchSize);
          try {
            const r = await AI.run("@cf/baai/bge-m3", { text: batch.map(b => b.text) }) as WorkersAIEmbedding;
            for (let j = 0; j < batch.length; j++) {
              vectors.push({
                id: `${batch[j].subject}-${batch[j].path.replace(/\//g, "-")}-${batch[j].chunkIndex}`,
                vector: extractEmbedding(r.data[j]),
                metadata: { subject: batch[j].subject, path: batch[j].path, chunkIndex: batch[j].chunkIndex, text: batch[j].text },
              });
            }
          } catch (e) { return json({ chunksAdded: 0, message: `Embedding failed: ${String(e)}` }); }
        }
        if (vectors.length === 0) return json({ chunksAdded: 0, message: "No vectors generated" });
        try {
          const result = await upstashUpsert(UPSTASH_URL, UPSTASH_TOKEN, vectors) as { status?: string; error?: string };
          if ("error" in result) {
            console.error(`Upsert error: ${result.error}`);
            return json({ chunksAdded: 0, message: `Upsert error: ${result.error}` });
          }
        } catch (e) { return json({ chunksAdded: 0, message: `Upsert failed: ${e}` }); }
        return json({ chunksAdded: vectors.length, message: `Indexed ${vectors.length} chunks` });
      }

      // PRACTICE — get a question
      if (path === "/rag/practice" && request.method === "GET") {
        const userId = 1;
        const subject = url.searchParams.get("subject") || undefined;
        const question = await getPracticeQuestion(userId, subject, env.DB);
        if (!question) return json({ question: null, message: "No questions available" });
        return json({ question: { id: question.id, text: question.question_text, difficulty: question.difficulty, topic_tags: question.topic_tags }, solution_preview: question.solution_text ? "✅ Solution available" : "❌ No solution yet" });
      }

      // GRADE — grade user answer
      if (path === "/rag/grade" && request.method === "POST") {
        const body = await request.json() as { question_id?: number; user_answer?: string; preview?: boolean };
        if (!body.question_id) return workerError("question_id required");
        if (!body.user_answer && !body.preview) return workerError("user_answer required");

        const question = await env.DB.prepare("SELECT * FROM questions WHERE id = ?").bind(body.question_id).first<Question>();
        if (!question) return workerError("Question not found", 404);

        const existing = await env.DB.prepare("SELECT * FROM user_progress WHERE user_id = 1 AND question_id = ?").bind(body.question_id).first<CardProgress>();

        if (body.preview) {
          return json({ preview: true, solution: question.solution_text, question: question.question_text });
        }

        const { score, feedback, partial } = await gradeAnswer(body.question_id, body.user_answer!, question, env);
        const isCorrect = score >= 3;
        const rating = smScoreToRating(score);
        const now = new Date();

        const r = existing
          ? f.next({ ...existing, elapsed_days: existing.elapsed_days ?? 0, scheduled_days: existing.scheduled_days ?? existing.interval, state: sm2StateToFSRS(existing.state), due: new Date(existing.due_date || now) } as any, now, rating)
          : f.next(createEmptyCard(now), now, rating);

        const sm = applyFSRSResult(r, now);

        if (!existing) {
          await env.DB.prepare(
            `INSERT INTO user_progress (user_id, question_id, attempted, correct, last_attempt, stability, difficulty, state, reps, lapses, interval, due_date, elapsed_days, scheduled_days, preview_answer)
             VALUES (1, ?, 1, ?, datetime('now'), ?, ?, ?, ?, ?, ?, datetime(?), ?, ?, 0)`
          ).bind(body.question_id, isCorrect && !partial ? 1 : 0, sm.stability, sm.difficulty, sm.state, sm.reps, sm.lapses, sm.interval, sm.due_date, sm.elapsed_days, sm.scheduled_days).run();
        } else {
          await env.DB.prepare(
            `UPDATE user_progress SET attempted = attempted + 1, correct = correct + ?, last_attempt = datetime('now'),
             stability = ?, difficulty = ?, state = ?, reps = ?, lapses = ?, interval = ?, due_date = datetime(?), elapsed_days = ?, scheduled_days = ?
             WHERE user_id = 1 AND question_id = ?`
          ).bind(isCorrect && !partial ? 1 : 0, sm.stability, sm.difficulty, sm.state, sm.reps, sm.lapses, sm.interval, sm.due_date, sm.elapsed_days, sm.scheduled_days, body.question_id).run();
        }

        return json({ score, feedback, correct: isCorrect, solution: question.solution_text, next_review: sm.due_date, state: sm.state, rating: Rating[rating] });
      }

      // RATE — user self-rates after seeing solution, FSRS runs with their chosen rating
      if (path === "/rag/rate" && request.method === "POST") {
        const body = await request.json() as { question_id?: number; rating?: string; correct?: boolean };
        if (!body.question_id) return workerError("question_id required");
        if (!body.rating) return workerError("rating required (Again|Hard|Good|Easy)");
        if (!["Again", "Hard", "Good", "Easy"].includes(body.rating)) {
          return workerError("rating must be one of: Again, Hard, Good, Easy");
        }

        const rating = Rating[body.rating as keyof typeof Rating] as unknown as number;
        const now = new Date();

        const existing = await env.DB.prepare(
          "SELECT * FROM user_progress WHERE user_id = 1 AND question_id = ?"
        ).bind(body.question_id).first<CardProgress>();

        const r = existing
          ? f.next({ ...existing, elapsed_days: existing.elapsed_days ?? 0, scheduled_days: existing.scheduled_days ?? existing.interval, state: sm2StateToFSRS(existing.state), due: new Date(existing.due_date || now) } as any, now, rating)
          : f.next(createEmptyCard(now), now, rating);

        const sm = applyFSRSResult(r, now);

        if (!existing) {
          await env.DB.prepare(
            `INSERT INTO user_progress (user_id, question_id, attempted, correct, last_attempt, stability, difficulty, state, reps, lapses, interval, due_date, elapsed_days, scheduled_days, preview_answer)
             VALUES (1, ?, 1, ?, datetime('now'), ?, ?, ?, ?, ?, ?, datetime(?), ?, ?, 0)`
          ).bind(body.question_id, body.correct ? 1 : 0, sm.stability, sm.difficulty, sm.state, sm.reps, sm.lapses || 0, sm.interval, sm.due_date, sm.elapsed_days, sm.scheduled_days).run();
        } else {
          await env.DB.prepare(
            `UPDATE user_progress SET attempted = attempted + 1, correct = correct + ?, last_attempt = datetime('now'),
             stability = ?, difficulty = ?, state = ?, reps = ?, lapses = ?, interval = ?, due_date = datetime(?), elapsed_days = ?, scheduled_days = ?
             WHERE user_id = 1 AND question_id = ?`
          ).bind(body.correct ? 1 : 0, sm.stability, sm.difficulty, sm.state, sm.reps, sm.lapses, sm.interval, sm.due_date, sm.elapsed_days, sm.scheduled_days, body.question_id).run();
        }

        return json({
          rating: body.rating,
          state: sm.state,
          next_review: sm.due_date,
          interval: sm.interval,
          message: `${body.rating} → Next review: ${new Date(sm.due_date).toLocaleDateString("fr-MA")} (${sm.interval} day${sm.interval !== 1 ? "s" : ""})`
        });
      }

      // QUESTIONS — list all questions
      if (path === "/rag/questions" && request.method === "GET") {
        const url = new URL(request.url);
        const subject = url.searchParams.get("subject");
        const difficulty = url.searchParams.get("difficulty");
        const limit = Math.min(parseInt(url.searchParams.get("limit") || "50"), 200);
        const offset = parseInt(url.searchParams.get("offset") || "0");

        let sql = "SELECT id, question_text, solution_text, topic_tags, difficulty, created_at FROM questions WHERE 1=1";
        let countSql = "SELECT COUNT(*) as c FROM questions WHERE 1=1";
        const bindings: (string | number)[] = [];
        if (subject) {
          sql += " AND topic_tags LIKE ?";
          countSql += " AND topic_tags LIKE ?";
          bindings.push(`%${subject}%`);
        }
        if (difficulty) {
          sql += " AND difficulty = ?";
          countSql += " AND difficulty = ?";
          bindings.push(parseInt(difficulty));
        }
        sql += " ORDER BY id LIMIT ? OFFSET ?";
        countSql += " LIMIT 1";

        const rows = await env.DB.prepare(sql).bind(...bindings, limit, offset).all<Question>();
        const totalBind = subject || difficulty ? bindings : [];
        const total = await env.DB.prepare(countSql).bind(...totalBind).first<{ c: number }>();

        return json({ questions: rows.results, total: total?.c || 0, limit, offset });
      }

      // SESSION — record and get session stats
      if (path === "/rag/session" && request.method === "POST") {
        const body = await request.json() as { questions_reviewed?: number; correct?: number; time_spent?: number };
        const today = new Date().toISOString().slice(0, 10);
        const existing = await env.DB.prepare(
          "SELECT id FROM study_sessions WHERE user_id = 1 AND date = ?"
        ).bind(today).first<{ id: number }>();

        if (existing) {
          await env.DB.prepare(
            "UPDATE study_sessions SET questions_reviewed = questions_reviewed + ?, correct = correct + ?, time_spent_seconds = time_spent_seconds + ? WHERE id = ?"
          ).bind(body.questions_reviewed || 0, body.correct || 0, body.time_spent || 0, existing.id).run();
        } else {
          await env.DB.prepare(
            "INSERT INTO study_sessions (user_id, date, questions_reviewed, correct, time_spent_seconds) VALUES (1, ?, ?, ?, ?)"
          ).bind(today, body.questions_reviewed || 0, body.correct || 0, body.time_spent || 0).run();
        }
        return json({ ok: true, date: today });
      }

      if (path === "/rag/session" && request.method === "GET") {
        const sessions = await env.DB.prepare(
          "SELECT date, questions_reviewed, correct, time_spent_seconds FROM study_sessions WHERE user_id = 1 ORDER BY date DESC LIMIT 30"
        ).all<{ date: string; questions_reviewed: number; correct: number; time_spent_seconds: number }>();

        const streak = await (async () => {
          let s = 0;
          const today = new Date();
          for (let i = 0; i < 365; i++) {
            const d = new Date(today);
            d.setDate(d.getDate() - i);
            const dateStr = d.toISOString().slice(0, 10);
            const found = sessions.results?.find(r => r.date === dateStr);
            if (found && found.questions_reviewed > 0) s++;
            else if (i > 0) break;
          }
          return s;
        })();

        const today = new Date().toISOString().slice(0, 10);
        const todaySession = sessions.results?.find(r => r.date === today);
        const totalReviews = sessions.results?.reduce((sum, s) => sum + s.questions_reviewed, 0) || 0;
        const totalCorrect = sessions.results?.reduce((sum, s) => sum + s.correct, 0) || 0;

        return json({
          streak,
          today: todaySession ? { reviews: todaySession.questions_reviewed, correct: todaySession.correct, time: todaySession.time_spent_seconds } : null,
          total_reviews: totalReviews,
          total_correct: totalCorrect,
          overall_accuracy: totalReviews > 0 ? Math.round(totalCorrect / totalReviews * 100) : 0,
          recent: sessions.results?.slice(0, 7).map(s => ({ date: s.date, reviews: s.questions_reviewed, correct: s.correct })),
        });
      }

      // TRACK — show progress
      if (path === "/rag/track" && request.method === "GET") {
        const stats = await env.DB.prepare(
          `SELECT
             COUNT(*) as total,
             SUM(attempted) as attempted,
             SUM(correct) as correct,
             SUM(CASE WHEN due_date IS NOT NULL AND due_date <= datetime('now') THEN 1 ELSE 0 END) as due_now,
             SUM(lapses) as lapses,
             SUM(reps) as total_reviews
           FROM user_progress WHERE user_id = 1`
        ).first<{ total: number; attempted: number; correct: number; due_now: number; lapses: number; total_reviews: number }>();

        const totalQ = await env.DB.prepare("SELECT COUNT(*) as c FROM questions").first<{ c: number }>();
        const dueQuestions = await env.DB.prepare(
          `SELECT q.id, q.question_text, q.difficulty, up.due_date FROM user_progress up
           JOIN questions q ON q.id = up.question_id WHERE up.user_id = 1 AND up.due_date IS NOT NULL AND up.due_date <= datetime('now')
           ORDER BY up.due_date ASC LIMIT 10`
        ).all<{ id: number; question_text: string; difficulty: number; due_date: string }>();

        return json({
          stats: {
            total_questions: totalQ?.c || 0,
            attempted: stats?.attempted || 0,
            correct: stats?.correct || 0,
            accuracy: stats?.attempted ? Math.round((stats!.correct! / stats!.attempted!) * 100) : 0,
            due_now: stats?.due_now || 0,
            lapses: stats?.lapses || 0,
            total_reviews: stats?.total_reviews || 0,
          },
          due_questions: dueQuestions.results || [],
        });
      }

      // DEBUG
      if (path === "/rag/debug" && request.method === "GET") {
        try {
          const r = await env.run("@cf/baai/bge-m3", { text: ["hello world test"] }) as unknown;
          const r2 = r as Record<string, unknown>;
          return json({ keys: Object.keys(r2), shape: r2.shape, dataLen: Array.isArray(r2.data) ? r2.data.length : "n/a", data0type: typeof r2.data });
        } catch (e) { return json({ error: String(e) }); }
      }

      // ADD-QUESTION — add a question
      if (path === "/rag/add-question" && request.method === "POST") {
        const body = await request.json() as { question?: string; subject?: string; solution?: string; difficulty?: number; solve?: boolean };
        if (!body.question) return workerError("question required");

        let solution = body.solution;
        if (!solution && (body.solve || !body.solution)) {
          try {
            solution = (await solveQuestion(body.question, body.subject || "general", env)).solution;
          } catch { solution = null; }
        }

        const tags = body.subject || "general";
        const diff = body.difficulty || 2;

        const result = await env.DB.prepare(
          "INSERT INTO questions (question_text, solution_text, topic_tags, difficulty) VALUES (?, ?, ?, ?)"
        ).bind(body.question, solution, tags, diff).run();

        return json({ id: result.meta?.last_row_id, question: body.question, solution, message: "Question added" });
      }

      // OCR - Multi-provider OCR chain with fallback
      if (path === "/rag/ocr" && request.method === "POST") {
        const contentType = request.headers.get("content-type") || "";

        if (contentType.includes("application/json")) {
          const body = await request.json() as { image?: string; url?: string; provider?: string };

          let imageData = body.image;

          if (!imageData && body.url) {
            try {
              const imgResp = await fetch(body.url);
              const arr = await imgResp.arrayBuffer();
              const bytes = Array.from(new Uint8Array(arr));
              imageData = `data:image/jpeg;base64,${btoa(String.fromCharCode(...bytes))}`;
            } catch (e) {
              return workerError(`Failed to fetch image: ${e}`);
            }
          }

          if (!imageData) {
            return workerError("image (base64) or url required");
          }

          // Single provider mode
          if (body.provider) {
            const start = Date.now();
            switch (body.provider) {
              case "mistral":
                if (!env.MISTRAL_API_KEY) return workerError("MISTRAL_API_KEY not configured", 503);
                return json(await ocrMistral(imageData, env.MISTRAL_API_KEY));
              case "ocrspace":
                return json(await ocrOCRspace(imageData));
              case "yandex":
                if (!env.YANDEX_IAM_TOKEN || !env.YANDEX_FOLDER_ID) return workerError("YANDEX credentials not configured", 503);
                return json(await ocrYandex(imageData, env.YANDEX_IAM_TOKEN, env.YANDEX_FOLDER_ID));
              case "llama-vision":
                return json(await ocrLlamaVision(imageData, env.AI));
              case "uform":
                return json(await ocrUform(imageData, env.AI));
              default:
                return workerError(`Unknown provider: ${body.provider}`);
            }
          }

          // Full chain mode
          const { results, best } = await runOCRChain(imageData, env);

          if (best) {
            return json({
              text: best.text,
              source: best.source,
              chars: best.chars,
              provider: best.provider,
              latency_ms: best.latency_ms,
              all_results: results.map(r => ({
                provider: r.provider,
                chars: r.chars,
                error: r.error,
                latency_ms: r.latency_ms,
              })),
            });
          }

          return json({
            text: "",
            source: "none",
            chars: 0,
            provider: "none",
            error: "All OCR providers failed",
            all_results: results.map(r => ({
              provider: r.provider,
              chars: r.chars,
              error: r.error,
              latency_ms: r.latency_ms,
            })),
          });
        }

        return workerError("Content-Type must be application/json");
      }

      // OCR providers list
      if (path === "/rag/ocr/providers" && request.method === "GET") {
        return json({
          providers: [
            { name: "mistral", tier: "cloud", free: false, languages: "100+", api_key: !!env.MISTRAL_API_KEY },
            { name: "ocrspace", tier: "cloud", free: true, languages: "100+", api_key: false },
            { name: "yandex", tier: "cloud", free: true, languages: "95+", api_key: !!(env.YANDEX_IAM_TOKEN && env.YANDEX_FOLDER_ID) },
            { name: "llama-vision", tier: "edge", free: true, languages: "multilingual", api_key: true },
            { name: "uform", tier: "edge", free: true, languages: "multilingual", api_key: true },
          ],
          chain: ["mistral", "ocrspace", "yandex", "llama-vision", "uform"],
        });
      }

      // SOLVE
      if (path === "/rag/solve" && request.method === "POST") {
        const body = await request.json() as { question?: string; subject?: string };
        if (!body.question) return workerError("question required");
        return json(await solveQuestion(body.question, body.subject || "general", env));
      }

      // SOLVE-EXAM
      if (path === "/rag/solve-exam" && request.method === "POST") {
        const body = await request.json() as { questions?: string[]; subject?: string };
        if (!body.questions?.length) return workerError("questions array required");
        const results = await Promise.all(
          body.questions.map((q) => solveQuestion(q, body.subject || "general", env))
        );
        return json({ results });
      }

      // ROOT
      if (path === "/" || path === "") {
        return json({
          name: "BAC Study RAG API", version: "1.0.0",
          endpoints: ["/rag/health", "/rag/status", "/rag/questions", "/rag/route", "/rag/expand", "/rag/search", "/rag/query", "/rag/add", "/rag/add-batch", "/rag/add-question", "/rag/solve", "/rag/solve-exam", "/rag/practice", "/rag/grade", "/rag/track", "/rag/ocr", "/rag/ocr/providers"],
        });
      }

      return workerError("Not found", 404);
    } catch (e) {
      console.error(`Worker error: ${e}`);
      return workerError(e instanceof Error ? e.message : "Internal error", 500);
    }
  },
};
