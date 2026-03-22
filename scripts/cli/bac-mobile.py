#!/usr/bin/env python3
"""
BAC Study - Mobile/Web Interface (Termux-friendly)
Single file, works offline, uses Worker API
"""

import http.server
import socketserver
import json
import urllib.request
import urllib.parse

PORT = 8080
WORKER = "https://bac-api.amdajhbdh.workers.dev"

HTML = """<!DOCTYPE html>
<html lang="fr">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=yes">
  <meta name="theme-color" content="#6366f1">
  <meta name="apple-mobile-web-app-capable" content="yes">
  <title>BAC Study</title>
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; 
           background: #0f172a; color: #e2e8f0; min-height: 100vh; padding: 16px; 
           touch-action: manipulation; }
    .container { max-width: 500px; margin: 0 auto; }
    h1 { color: #6366f1; text-align: center; margin-bottom: 20px; font-size: 24px; }
    .card { background: #1e293b; border-radius: 12px; padding: 16px; margin-bottom: 12px; }
    input, textarea { width: 100%; background: #0f172a; border: 1px solid #334155; 
                      border-radius: 8px; color: #fff; padding: 12px; font-size: 16px; 
                      -webkit-appearance: none; }
    input:focus, textarea:focus { outline: none; border-color: #6366f1; }
    textarea { min-height: 80px; resize: none; }
    button { width: 100%; background: #6366f1; color: white; border: none; 
             padding: 14px; border-radius: 8px; font-size: 16px; font-weight: 600;
             cursor: pointer; touch-action: manipulation; }
    button:active { opacity: 0.8; }
    button:disabled { opacity: 0.5; }
    .btn-row { display: flex; gap: 8px; }
    .btn-row button { flex: 1; }
    .result { background: #0f172a; border-radius: 8px; padding: 12px; margin-top: 12px; 
              white-space: pre-wrap; word-wrap: break-word; font-size: 14px; line-height: 1.5; }
    .sources { font-size: 12px; color: #94a3b8; margin-top: 8px; }
    .sources span { display: block; margin: 4px 0; }
    .loading { text-align: center; color: #94a3b8; padding: 20px; }
    .error { color: #ef4444; padding: 12px; }
    .tabs { display: flex; gap: 8px; margin-bottom: 16px; }
    .tab { flex: 1; background: #1e293b; border: 1px solid #334155; color: #94a3b8; 
           padding: 10px; border-radius: 8px; text-align: center; cursor: pointer; font-size: 14px; }
    .tab.active { background: #6366f1; color: white; border-color: #6366f1; }
    .tab-content { display: none; }
    .tab-content.active { display: block; }
    .nav { display: flex; gap: 8px; margin-bottom: 16px; }
    .nav button { padding: 10px; font-size: 14px; }
    pre { overflow-x: auto; }
  </style>
</head>
<body>
  <div class="container">
    <h1>📚 BAC Study</h1>
    
    <div class="nav">
      <button onclick="showTab('query')" class="tab active" id="tab-query">Query</button>
      <button onclick="showTab('solve')" class="tab" id="tab-solve">Solve</button>
      <button onclick="showTab('practice')" class="tab" id="tab-practice">Practice</button>
      <button onclick="showTab('status')" class="tab" id="tab-status">Status</button>
    </div>

    <!-- Query Tab -->
    <div class="tab-content active" id="content-query">
      <div class="card">
        <textarea id="query-input" placeholder="Ask a question... (e.g., 'explique la mitose')"></textarea>
        <button onclick="doQuery()" id="query-btn" style="margin-top: 12px;">Chercher</button>
      </div>
      <div id="query-result"></div>
    </div>

    <!-- Solve Tab -->
    <div class="tab-content" id="content-solve">
      <div class="card">
        <textarea id="solve-input" placeholder="Problem to solve... (e.g., 'résous x² + 5x + 6 = 0')"></textarea>
        <button onclick="doSolve()" id="solve-btn" style="margin-top: 12px;">Résoudre</button>
      </div>
      <div id="solve-result"></div>
    </div>

    <!-- Practice Tab -->
    <div class="tab-content" id="content-practice">
      <div class="card">
        <button onclick="loadPractice()">Nouvelle Question</button>
      </div>
      <div id="practice-question"></div>
      <div id="practice-answer" style="display:none;">
        <textarea id="user-answer" placeholder="Ta réponse..."></textarea>
        <button onclick="checkAnswer()" style="margin-top: 8px;">Vérifier</button>
      </div>
      <div id="practice-result"></div>
    </div>

    <!-- Status Tab -->
    <div class="tab-content" id="content-status">
      <div class="card">
        <button onclick="loadStatus()">Actualiser</button>
      </div>
      <div id="status-result"></div>
    </div>
  </div>

  <script>
    const WORKER = 'https://bac-api.amdajhbdh.workers.dev';
    let currentQuestion = null;

    function showTab(tab) {
      document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
      document.querySelectorAll('.tab-content').forEach(c => c.classList.remove('active'));
      document.getElementById('tab-' + tab).classList.add('active');
      document.getElementById('content-' + tab).classList.add('active');
    }

    function setLoading(id, msg) {
      document.getElementById(id).innerHTML = '<div class="loading">' + msg + '</div>';
    }

    async function apiCall(endpoint, data) {
      try {
        const resp = await fetch(WORKER + endpoint, {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          body: JSON.stringify(data)
        });
        return await resp.json();
      } catch (e) {
        return { error: e.message };
      }
    }

    // Query
    async function doQuery() {
      const input = document.getElementById('query-input');
      const btn = document.getElementById('query-btn');
      const result = document.getElementById('query-result');
      const q = input.value.trim();
      if (!q) return;
      
      btn.disabled = true;
      btn.textContent = '...';
      setLoading('query-result', 'Recherche...');
      
      const r = await apiCall('/rag/query', { query: q });
      
      btn.disabled = false;
      btn.textContent = 'Chercher';
      
      if (r.error) {
        result.innerHTML = '<div class="error">Erreur: ' + r.error + '</div>';
      } else {
        let srcs = r.sources?.map(s => '<span>' + s.subject + ': ' + s.source.split('/').pop() + '</span>').join('') || '';
        result.innerHTML = '<div class="result">' + r.answer + '</div><div class="sources">' + srcs + '</div>';
      }
    }

    // Solve
    async function doSolve() {
      const input = document.getElementById('solve-input');
      const btn = document.getElementById('solve-btn');
      const result = document.getElementById('solve-result');
      const q = input.value.trim();
      if (!q) return;
      
      btn.disabled = true;
      btn.textContent = '...';
      setLoading('solve-result', 'Calcul en cours...');
      
      const r = await apiCall('/rag/solve', { question: q });
      
      btn.disabled = false;
      btn.textContent = 'Résoudre';
      
      if (r.error) {
        result.innerHTML = '<div class="error">Erreur: ' + r.error + '</div>';
      } else {
        result.innerHTML = '<div class="result">' + (r.solution || 'Pas de solution') + '</div>';
      }
    }

    // Practice
    async function loadPractice() {
      setLoading('practice-question', 'Chargement...');
      const r = await fetch(WORKER + '/rag/practice').then(x => x.json()).catch(x => ({error: 'Erreur'}));
      
      if (r.question) {
        currentQuestion = r.question;
        document.getElementById('practice-question').innerHTML = '<div class="card"><strong>Q' + r.question.id + '</strong><br>' + r.question.text + '</div>';
        document.getElementById('practice-answer').style.display = 'block';
        document.getElementById('user-answer').value = '';
        document.getElementById('practice-result').innerHTML = '';
      } else {
        document.getElementById('practice-question').innerHTML = '<div class="card">Pas de question disponible</div>';
      }
    }

    async function checkAnswer() {
      const ans = document.getElementById('user-answer').value;
      if (!currentQuestion || !ans) return;
      
      const r = await apiCall('/rag/grade', { question_id: currentQuestion.id, user_answer: ans });
      const res = r.result || {};
      document.getElementById('practice-result').innerHTML = '<div class="card"><strong>Score: ' + (res.score || '?') + '/5</strong><br>' + (res.feedback || '') + '</div>';
    }

    // Status
    async function loadStatus() {
      setLoading('status-result', 'Chargement...');
      const s = await fetch(WORKER + '/rag/subjects').then(x => x.json()).catch(x => ({}));
      const st = await fetch(WORKER + '/rag/status').then(x => x.json()).catch(x => ({}));
      
      document.getElementById('status-result').innerHTML = '<div class="card"><strong>Système</strong><br>Vecteurs: ' + (st.total_vectors || '?') + '<br><br><strong>Matières</strong><br>Biology: ' + (s.subject_biology || 0) + '<br>Chemistry: ' + (s.subject_chemistry || 0) + '<br>Math: ' + (s.subject_math || 0) + '<br>Physics: ' + (s.subject_physics || 0) + '</div>';
    }

    // Init
    loadStatus();
    document.getElementById('query-input').addEventListener('keydown', e => { if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); doQuery(); }});
  </script>
</body>
</html>"""


class Handler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        if self.path == "/":
            self.send_response(200)
            self.send_header("Content-type", "text/html")
            self.end_headers()
            self.wfile.write(HTML.encode())
        else:
            super().do_GET()


print(f"BAC Study Mobile Interface")
print(f"==============================")
print(f"Starting server on port {PORT}...")
print(f"Open http://localhost:{PORT} in browser")
print(f"Or on Termux: http://127.0.0.1:{PORT}")
print(f"")
print(f"Press Ctrl+C to stop")

with socketserver.TCPServer(("", PORT), Handler) as httpd:
    httpd.serve_forever()
