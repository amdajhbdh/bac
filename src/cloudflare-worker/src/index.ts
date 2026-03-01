export interface Env {
	DB: D1Database;
	QUESTIONS: Vectorize;
	AI: Ai;
}

export default {
	async fetch(request: Request, env: Env, ctx: ExecutionContext): Promise<Response> {
		const url = new URL(request.url);
		const path = url.pathname;

		// CORS headers
		const corsHeaders = {
			"Access-Control-Allow-Origin": "*",
			"Access-Control-Allow-Methods": "GET, POST, OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type",
		};

		// Handle CORS preflight
		if (request.method === "OPTIONS") {
			return new Response(null, { headers: corsHeaders });
		}

		try {
			// Health check
			if (path === "/health") {
				return Response.json({ status: "ok", timestamp: Date.now() }, { headers: corsHeaders });
			}

			// Solve math problem with AI
			if (path === "/solve" && request.method === "POST") {
				const { problem } = await request.json();

				if (!problem) {
					return Response.json({ error: "problem is required" }, { status: 400, headers: corsHeaders });
				}

				// Use Workers AI to solve
				const result = await env.AI.run("@cf/meta/llama-3.2-3b-instruct", {
					messages: [
						{
							role: "system",
							content: "Tu es un assistant pédagogique pour les étudiants du BAC C en Mauritanie. Explique les solutions de manière claire et détaillée en français. Utilise des étapes claires.",
						},
						{ role: "user", content: problem }
					],
					max_tokens: 1000,
				});

				return Response.json(
					{
						problem,
						solution: result.response,
						model: "llama-3.2-3b-instruct",
					},
					{ headers: corsHeaders }
				);
			}

			// Get all questions
			if (path === "/questions" && request.method === "GET") {
				const { results } = await env.DB.prepare(
					"SELECT id, question_text, solution_text, difficulty FROM questions LIMIT 20"
				).all();

				return Response.json({ questions: results }, { headers: corsHeaders });
			}

			// Add new question
			if (path === "/questions" && request.method === "POST") {
				const { question_text, solution_text, difficulty } = await request.json();

				if (!question_text) {
					return Response.json({ error: "question_text is required" }, { status: 400, headers: corsHeaders });
				}

				const { success } = await env.DB.prepare(
					"INSERT INTO questions (question_text, solution_text, difficulty) VALUES (?, ?, ?)"
				).bind(question_text, solution_text || null, difficulty || 1).run();

				return Response.json({ success: !!success }, { headers: corsHeaders });
			}

			// Search similar questions using Vectorize
			if (path === "/search" && request.method === "POST") {
				const { query } = await request.json();

				if (!query) {
					return Response.json({ error: "query is required" }, { status: 400, headers: corsHeaders });
				}

				// Generate embedding
				const embedding = await env.AI.run("@cf/baai/bge-small-en-v1.5", { text: query });

				// Search Vectorize
				const results = await env.QUESTIONS.query(embedding.data[0], { topK: 5, returnMetadata: true });

				return Response.json({ results: results.matches }, { headers: corsHeaders });
			}

			// ===== USER ENDPOINTS =====

			// Register new user
			if (path === "/auth/register" && request.method === "POST") {
				const { name, email } = await request.json();

				if (!name || !email) {
					return Response.json({ error: "name and email are required" }, { status: 400, headers: corsHeaders });
				}

				try {
					const result = await env.DB.prepare(
						"INSERT INTO users (name, email, points) VALUES (?, ?, 0)"
					).bind(name, email).run();

					return Response.json({ 
						success: true, 
						user_id: result.lastInsertRowid,
						message: "User registered successfully" 
					}, { headers: corsHeaders });
				} catch (err: any) {
					if (err.message?.includes("UNIQUE constraint")) {
						return Response.json({ error: "Email already exists" }, { status: 409, headers: corsHeaders });
					}
					throw err;
				}
			}

			// Login user (simple email lookup)
			if (path === "/auth/login" && request.method === "POST") {
				const { email } = await request.json();

				if (!email) {
					return Response.json({ error: "email is required" }, { status: 400, headers: corsHeaders });
				}

				const { results } = await env.DB.prepare(
					"SELECT id, name, email, points FROM users WHERE email = ?"
				).bind(email).all();

				if (results.length === 0) {
					return Response.json({ error: "User not found" }, { status: 404, headers: corsHeaders });
				}

				return Response.json({ user: results[0] }, { headers: corsHeaders });
			}

			// Get user profile
			if (path === "/user/profile" && request.method === "POST") {
				const { user_id } = await request.json();

				if (!user_id) {
					return Response.json({ error: "user_id is required" }, { status: 400, headers: corsHeaders });
				}

				const { results } = await env.DB.prepare(`
					SELECT u.id, u.name, u.email, u.points, 
						   COALESCE(SUM(up.attempted), 0) as total_attempted,
						   COALESCE(SUM(up.correct), 0) as total_correct
					FROM users u
					LEFT JOIN user_progress up ON u.id = up.user_id
					WHERE u.id = ?
					GROUP BY u.id
				`).bind(user_id).all();

				if (results.length === 0) {
					return Response.json({ error: "User not found" }, { status: 404, headers: corsHeaders });
				}

				return Response.json({ user: results[0] }, { headers: corsHeaders });
			}

			// ===== PREDICTION ENDPOINTS =====

			// Get predictions
			if (path === "/predictions" && request.method === "GET") {
				const { results } = await env.DB.prepare(`
					SELECT p.*, q.question_text 
					FROM predictions p
					LEFT JOIN questions q ON p.question_id = q.id
					ORDER BY p.predicted_date DESC
					LIMIT 20
				`).all();

				return Response.json({ predictions: results }, { headers: corsHeaders });
			}

			// Generate predictions (AI-powered)
			if (path === "/predictions/generate" && request.method === "POST") {
				// Get all questions
				const { results: questions } = await env.DB.prepare(
					"SELECT id, question_text FROM questions"
				).all();

				// Use AI to predict likely exam questions
				const prompt = `Voici une liste de questions potentielles pour le BAC C Mauritanie. 
				Analyse et donne une probabilité (0-1) pour chaque question d'apparaître à l'examen.
				Questions: ${questions.map((q: any) => q.question_text).join(", ")}
				
				Réponds en JSON comme: [{"question_id": 1, "probability": 0.8}, {"question_id": 2, "probability": 0.3}, ...]`;

				const result = await env.AI.run("@cf/meta/llama-3.2-1b-instruct", {
					messages: [{ role: "user", content: prompt }],
					max_tokens: 500,
				});

				// Parse and save predictions
				try {
					const predictions = JSON.parse(result.response);
					for (const p of predictions) {
						await env.DB.prepare(
							"INSERT INTO predictions (question_id, probability, confidence) VALUES (?, ?, ?)"
						).bind(p.question_id, p.probability, p.probability * 0.8).run();
					}
					return Response.json({ 
						success: true, 
						count: predictions.length,
						message: "Predictions generated successfully" 
					}, { headers: corsHeaders });
				} catch {
					return Response.json({ 
						success: true, 
						ai_response: result.response,
						message: "Predictions generated (manual review needed)" 
					}, { headers: corsHeaders });
				}
			}

			// ===== LEADERBOARD ENDPOINTS =====

			// Get leaderboard
			if (path === "/leaderboard" && request.method === "GET") {
				const { results } = await env.DB.prepare(`
					SELECT u.id, u.name, u.points,
						   COALESCE(SUM(up.attempted), 0) as total_attempted,
						   COALESCE(SUM(up.correct), 0) as total_correct
					FROM users u
					LEFT JOIN user_progress up ON u.id = up.user_id
					GROUP BY u.id
					ORDER BY u.points DESC
					LIMIT 10
				`).all();

				return Response.json({ leaderboard: results }, { headers: corsHeaders });
			}

			// Update user points
			if (path === "/user/points" && request.method === "POST") {
				const { user_id, points } = await request.json();

				if (!user_id || points === undefined) {
					return Response.json({ error: "user_id and points are required" }, { status: 400, headers: corsHeaders });
				}

				await env.DB.prepare(
					"UPDATE users SET points = points + ? WHERE id = ?"
				).bind(points, user_id).run();

				return Response.json({ success: true, message: "Points updated" }, { headers: corsHeaders });
			}

			// ===== PRACTICE ENDPOINTS =====

			// Submit answer
			if (path === "/practice/answer" && request.method === "POST") {
				const { user_id, question_id, correct, time_spent } = await request.json();

				if (!user_id || !question_id) {
					return Response.json({ error: "user_id and question_id are required" }, { status: 400, headers: corsHeaders });
				}

				// Check if progress exists
				const { results } = await env.DB.prepare(
					"SELECT id FROM user_progress WHERE user_id = ? AND question_id = ?"
				).bind(user_id, question_id).all();

				if (results.length > 0) {
					await env.DB.prepare(`
						UPDATE user_progress 
						SET attempted = attempted + 1, 
						    correct = correct + ?, 
						    last_attempt = datetime('now'),
						    time_spent_seconds = COALESCE(time_spent_seconds, 0) + ?
						WHERE user_id = ? AND question_id = ?
					`).bind(correct ? 1 : 0, time_spent || 0, user_id, question_id).run();
				} else {
					await env.DB.prepare(`
						INSERT INTO user_progress (user_id, question_id, attempted, correct, last_attempt, time_spent_seconds)
						VALUES (?, ?, 1, ?, datetime('now'), ?)
					`).bind(user_id, question_id, correct ? 1 : 0, time_spent || 0).run();
				}

				// Update user points
				const pointsEarned = correct ? 10 : 1;
				await env.DB.prepare(
					"UPDATE users SET points = points + ? WHERE id = ?"
				).bind(pointsEarned, user_id).run();

				return Response.json({ 
					success: true, 
					correct,
					points_earned: pointsEarned 
				}, { headers: corsHeaders });
			}

			// Get practice questions (random)
			if (path === "/practice/questions" && request.method === "GET") {
				const url = new URL(request.url);
				const difficulty = url.searchParams.get("difficulty");
				const limit = parseInt(url.searchParams.get("limit") || "10");

				let query = "SELECT id, question_text, difficulty FROM questions";
				const params: any[] = [];

				if (difficulty) {
					query += " WHERE difficulty = ?";
					params.push(parseInt(difficulty));
				}

				query += " ORDER BY RANDOM() LIMIT ?";
				params.push(limit);

				const { results } = await env.DB.prepare(query).bind(...params).all();

				return Response.json({ questions: results }, { headers: corsHeaders });
			}

			// Root path - API info
			if (path === "/") {
				return Response.json(
					{
						name: "BAC Unified API",
						version: "1.0.0",
						endpoints: {
							health: "GET /health",
							solve: "POST /solve",
							questions: "GET /questions, POST /questions",
							search: "POST /search",
							auth: "POST /auth/register, POST /auth/login",
							user: "POST /user/profile, POST /user/points",
							predictions: "GET /predictions, POST /predictions/generate",
							leaderboard: "GET /leaderboard",
							practice: "POST /practice/answer, GET /practice/questions",
						},
					},
					{ headers: corsHeaders }
				);
			}

			return Response.json({ error: "Not found" }, { status: 404, headers: corsHeaders });
		} catch (error) {
			return Response.json(
				{ error: error.message || "Internal server error" },
				{ status: 500, headers: corsHeaders }
			);
		}
	},
};
