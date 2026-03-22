#!/usr/bin/env python3
"""
Telegram Bot for BAC Study RAG System.
Receives messages → calls Worker API → responds.
"""

import os, json, subprocess, logging
from pathlib import Path

logging.basicConfig(level=logging.INFO)
log = logging.getLogger(__name__)

TOKEN = os.environ.get(
    "TELEGRAM_BOT_TOKEN", "8234814718:AAGYJLgBfLOMqWXuXWoYAMAfNF2CMBqKDAc"
)
WORKER = "https://bac-api.amdajhbdh.workers.dev"
quiz_state = {}
session_state = {}


def curl(method, path, data=None, timeout=120):
    cmd = ["curl", "-s", "-X", method, f"{WORKER}{path}"]
    if data:
        cmd.extend(["-H", "Content-Type: application/json", "-d", json.dumps(data)])
    result = subprocess.run(cmd, capture_output=True, text=True, timeout=timeout + 10)
    try:
        return json.loads(result.stdout)
    except:
        return {}


def record_session(chat_id, correct=None):
    import time

    if chat_id not in session_state:
        session_state[chat_id] = {"start": time.time(), "reviews": 0, "correct": 0}
    if correct is not None:
        session_state[chat_id]["reviews"] += 1
        if correct:
            session_state[chat_id]["correct"] += 1
    elapsed = int(time.time() - session_state[chat_id]["start"])
    curl(
        "POST",
        "/rag/session",
        {
            "questions_reviewed": session_state[chat_id]["reviews"],
            "correct": session_state[chat_id]["correct"],
            "time_spent": elapsed,
        },
        timeout=10,
    )


def send_message(chat_id, text, parse_mode="Markdown", reply_markup=None):
    payload = {"chat_id": chat_id, "text": text, "parse_mode": parse_mode}
    if reply_markup:
        payload["reply_markup"] = reply_markup
    cmd = [
        "curl",
        "-s",
        "-X",
        "POST",
        f"https://api.telegram.org/bot{TOKEN}/sendMessage",
        "-H",
        "Content-Type: application/json",
        "-d",
        json.dumps(payload),
    ]
    subprocess.run(cmd, capture_output=True)


def edit_message(chat_id, message_id, text, parse_mode="Markdown", reply_markup=None):
    payload = {
        "chat_id": chat_id,
        "message_id": message_id,
        "text": text,
        "parse_mode": parse_mode,
    }
    if reply_markup:
        payload["reply_markup"] = reply_markup
    cmd = [
        "curl",
        "-s",
        "-X",
        "POST",
        f"https://api.telegram.org/bot{TOKEN}/editMessageText",
        "-H",
        "Content-Type: application/json",
        "-d",
        json.dumps(payload),
    ]
    subprocess.run(cmd, capture_output=True)


def send_action(chat_id, action="typing"):
    cmd = [
        "curl",
        "-s",
        "-X",
        "POST",
        f"https://api.telegram.org/bot{TOKEN}/sendChatAction",
        "-H",
        "Content-Type: application/json",
        "-d",
        json.dumps({"chat_id": chat_id, "action": action}),
    ]
    subprocess.run(cmd, capture_output=True)


def handle_command(chat_id, text):
    parts = text.strip().split(maxsplit=1)
    cmd = parts[0].lower()
    args = parts[1] if len(parts) > 1 else ""

    if cmd == "/start":
        return """*Bienvenue au BAC Study Bot !* 📚

*Pratique*
`/practice` — Question avec spaced repetition
`/quiz [matière]` — Quiz choix multiple (5 questions)
`/due` — Questions en retard
`/grade [id] [réponse]` — Noter ta réponse

*Apprentissage*
`/query [question]` — Poser une question (RAG)
`/solve [question]` — Résoudre un exercice
`/solve-exam` — Résoudre un examen complet

*OCR (Envoyer une image)*
Le bot détecte automatiquement le texte des images envoyées et répond avec le contenu RAG.

*Questions*
`/add [question]` — Ajouter une question
`/track` — Voir ta progression
`/status` — État du système
`/session` — Stats de streak
`/ocr-providers` — Voir les providers OCR disponibles"""

    if cmd == "/ocr-providers":
        r = curl("GET", "/rag/ocr/providers")
        providers = r.get("providers", [])
        chain = r.get("chain", [])
        lines = ["*🔍 OCR Providers (18 disponibles):*\n"]
        lines.append("_Tier Cloud (API):_")
        for p in providers:
            status = "✅" if p.get("api_key") else "🔑" if p.get("free") else "❌"
            lines.append(
                f"  {status} {p['name']} ({p['tier']}, {'gratuit' if p.get('free') else 'payant'})"
            )
        lines.append("\n_Chaîne actuelle:_")
        lines.append(" → ".join(chain))
        lines.append("\n_💡 Envoie une image au bot pour tester l'OCR!_")
        return "\n".join(lines)

    if cmd == "/ocr-test":
        return "*📸 Pour tester l'OCR:*\n\nEnvoie simplement une image au bot !\n\nLe bot essaie automatiquement:\n1. Mistral OCR (le plus précis)\n2. Workers AI (Llama Vision)\n3. OCR.space (gratuit)\n4. Yandex Vision (gratuit)\n5. uform (local)"

    if cmd == "/help":
        return """*Guide du Bot*

*Pratique:* `/practice` → question, puis note (Again/Hard/Good/Easy)
*Quiz:* `/quiz math` → quiz 5 questions
*Examen:* `/solve-exam` → paste un exam complet, résout chaque question
*RAG:* pose ta question directement en français"""

    if cmd == "/help":
        return """*Guide du Bot*

*Pratique:* `/practice` → question, puis note avec les boutons (Again/Hard/Good/Easy)
*Quiz:* `/quiz math` → quiz 5 questions
*Grade:* `/grade 42 Ma réponse` → score + boutons pour noter
*RAG:* pose ta question directement en français"""

    if cmd == "/status":
        r = curl("GET", "/rag/status")
        vectors = r.get("total_vectors", "?")
        r2 = curl("GET", "/rag/track")
        s = r2.get("stats", {})
        return f"*État du système*\n\n📚 Vecteurs: {vectors}\n❓ Questions: {s.get('total_questions', 0)}\n📝 Attempted: {s.get('attempted', 0)}\n✅ Accuracy: {s.get('accuracy', 0)}%"

    if cmd == "/session":
        r = curl("GET", "/rag/session")
        streak = r.get("streak", 0)
        r2 = curl("GET", "/rag/track")
        s = r2.get("stats", {})
        today = r.get("today") or {}
        lines = [f"*🔥 Streak: {streak} jour{'s' if streak != 1 else ''}*"]
        if today:
            lines.append(
                f"Aujourd'hui: {today.get('reviews', 0)} revues, {today.get('correct', 0)} correctes ({round((today.get('time') or 0) / 60)}min)"
            )
        week = r.get("recent", [])
        if week:
            lines.append("\n*7 derniers jours:*")
            for d in week[-7:]:
                date = d.get("date", "")[-5:]
                reviews = d.get("reviews", 0)
                correct = d.get("correct", 0)
                bar = "▓" * reviews + "░" * max(0, 10 - reviews)
                lines.append(f"{date} {bar} {reviews}r/{correct}c")
        lines.extend(
            [
                f"\n*Total:* {s.get('attempted', 0)} tentatives, {s.get('accuracy', 0)}% accuracy",
                f"*Questions dues:* {s.get('due_now', 0)}",
            ]
        )
        return "\n".join(lines)

    if cmd == "/track":
        r = curl("GET", "/rag/track")
        s = r.get("stats", {})
        r2 = curl("GET", "/rag/session")
        streak = r2.get("streak", 0)
        today = r2.get("today")
        lines = [
            f"*📊 Progression* 🔥{streak} jour{'s' if streak != 1 else ''}",
        ]
        if today:
            lines.append(
                f"Aujourd'hui: {today.get('reviews', 0)} questions, {today.get('correct', 0)} correctes"
            )
        lines.extend(
            [
                f"Questions: {s.get('attempted', 0)}/{s.get('total_questions', 0)}",
                f"Accuracy: {s.get('accuracy', 0)}%",
                f"Due now: {s.get('due_now', 0)}",
            ]
        )
        if r.get("due_questions"):
            lines.append("\n*⏰ Questions dues:*")
            for q in r["due_questions"][:5]:
                lines.append(f"• [{q['id']}] {q['question_text'][:50]}")
        return "\n".join(lines)

    if cmd == "/due":
        r = curl("GET", "/rag/track")
        due = r.get("due_questions", [])
        if not due:
            return "🎉 Pas de questions en retard! Lance `/practice` pour continuer."
        lines = ["*⏰ Questions en retard:*"]
        for q in due[:10]:
            lines.append(f"[{q['id']}] {q['question_text'][:60]}")
        lines.append("\n_Réponds: /grade [id] Ta réponse_")
        return "\n".join(lines)

    if cmd == "/practice":
        r = curl("GET", "/rag/practice")
        q = r.get("question")
        if not q:
            return "❌ Pas de questions disponibles."
        text = f"*❓ Question #{q['id']}* [{'⭐' * q['difficulty']}]\n\n{q['text']}\n\n_Tape ta réponse en dessous_"
        reply_markup = {
            "inline_keyboard": [
                [
                    {"text": "📝 Voir solution", "callback_data": f"preview:{q['id']}"},
                    {"text": "⏭️ Question suivante", "callback_data": f"next:{q['id']}"},
                ]
            ]
        }
        send_message(chat_id, text, reply_markup=reply_markup)
        return None

    if cmd == "/quiz":
        subject = args.split()[0] if args else ""
        params = "limit=100"
        if subject in ("math", "biology", "chemistry", "physics"):
            params += f"&subject={subject}"
        qs = curl("GET", f"/rag/questions?{params}")
        questions = qs.get("questions", [])
        if len(questions) < 4:
            return "Pas assez de questions pour un quiz."
        import random

        random.shuffle(questions)
        quiz_qs = []
        for q in questions:
            if len(quiz_qs) >= 5:
                break
            correct = q.get("solution_text", "").strip()
            if not correct or len(correct) < 3:
                continue
            tag = q.get("topic_tags", "") or ""
            candidates = [
                r
                for r in questions
                if r["id"] != q["id"]
                and r.get("topic_tags", "") == tag
                and r.get("solution_text")
            ]
            if len(candidates) < 3:
                candidates = [
                    r
                    for r in questions
                    if r["id"] != q["id"] and r.get("solution_text")
                ]
            if len(candidates) < 3:
                continue
            others = random.sample(candidates, 3)
            options = [(correct, q["id"])] + [
                (r.get("solution_text", "").strip(), r["id"]) for r in others
            ]
            random.shuffle(options)
            quiz_qs.append(
                {"id": q["id"], "text": q["question_text"], "options": options}
            )

        if not quiz_qs:
            return "Pas assez de questions valides."

        quiz_data = {
            "questions": quiz_qs,
            "current": 0,
            "correct": 0,
            "total": len(quiz_qs),
            "subject": subject,
        }
        quiz_state[chat_id] = quiz_data

        q = quiz_qs[0]
        lines = [f"*Quiz — Question 1/5*" + (f" ({subject})" if subject else "")]
        lines.append(f"\n{q['text']}\n")
        for i, (opt, _) in enumerate(q["options"]):
            lines.append(f"{i + 1}. {opt[:70]}")
        lines.append("\n_Tappe 1-4 pour répondre_")

        keyboard = [
            [{"text": str(i + 1), "callback_data": f"qz:{chat_id}:{i}"}]
            for i in range(4)
        ]
        send_message(
            chat_id, "\n".join(lines), reply_markup={"inline_keyboard": keyboard}
        )
        return None

    if cmd == "/retry":
        r = curl("GET", "/rag/track")
        due = r.get("due_questions", [])
        if not due:
            return "🎉 Pas de questions en retard!"
        import random

        random.shuffle(due)
        qi = due[0]
        pv = curl("POST", "/rag/grade", {"question_id": qi["id"], "preview": True})
        text = f"*⚡ Retry: Question #{qi['id']}*\n\n{qi['question_text']}\n\n*Solution:*\n{pv.get('solution', 'N/A')[:300]}"
        send_message(chat_id, text)
        return f"_Réponds: /grade {qi['id']} Ta réponse_"

    if cmd == "/query":
        if not args:
            return "Usage: `/query Ta question`"
        send_action(chat_id, "typing")
        r = curl("POST", "/rag/query", {"query": args}, timeout=120)
        if r.get("answer"):
            answer = r["answer"][:2000]
            sources = ""
            if r.get("sources"):
                sources = "\n\n*Sources:*\n" + "\n".join(
                    f"• [{s['subject']}] {s['source'].rsplit('/', 1)[-1][:50]}"
                    for s in r["sources"][:3]
                )
            return f"*{args}*\n\n{answer}{sources}"
        return "❌ Pas de réponse trouvée."

    if cmd == "/solve":
        if not args:
            return "Usage: `/solve Ta question`"
        send_action(chat_id, "typing")
        r = curl("POST", "/rag/solve", {"question": args}, timeout=120)
        if r.get("solution"):
            sol = r["solution"][:2000]
            return f"*❓ Question:*\n{args}\n\n*✅ Solution:*\n{sol}"
        return "❌ Pas de solution trouvée."

    if cmd == "/solve-exam":
        if not args:
            return "Usage: `/solve-exam Colle ton examen ici`"
        lines = args.strip().split("\n")
        questions = [
            l.strip()
            for l in lines
            if l.strip() and not l.startswith("#") and len(l.strip()) > 10
        ]
        if not questions:
            return "❌ Pas de questions détectées. Colle un examen avec des questions."
        send_message(
            chat_id, f"📝 Détecté {len(questions)} questions, résolution en cours..."
        )
        results = []
        for i, q in enumerate(questions[:10]):
            send_action(chat_id, "typing")
            r = curl("POST", "/rag/solve", {"question": q}, timeout=120)
            sol = r.get("solution", "Erreur")[:500]
            results.append(f"*{i + 1}.* {q[:60]}...\n{sol}")
            if (i + 1) % 3 == 0:
                send_message(
                    chat_id, f"Progress: {i + 1}/{len(questions)} questions résolues"
                )
        send_message(chat_id, f"✅ Terminé! {len(results)} réponses:")
        for r in results[:5]:
            send_message(chat_id, r[:1000])
        if len(results) > 5:
            send_message(
                chat_id,
                f"... et {len(results) - 5} réponses supplémentaires. Use /query pour plus de détails.",
            )
        return None

    if cmd == "/grade":
        parts = args.split(maxsplit=1)
        if not parts:
            return "Usage: `/grade [id] [réponse]` ou `/grade [id] --preview`"
        qid = parts[0]
        if len(parts) < 2:
            return f"Usage: `/grade {qid} Ta réponse`"
        answer = parts[1]
        if answer == "--preview":
            r = curl("POST", "/rag/grade", {"question_id": int(qid), "preview": True})
            return f"*Question:* {r.get('question', '')}\n\n*Solution:*\n{r.get('solution', 'N/A')}"
        send_action(chat_id, "typing")
        r = curl(
            "POST",
            "/rag/grade",
            {"question_id": int(qid), "user_answer": answer},
            timeout=120,
        )
        emoji = "✅" if r.get("correct") else "❌"
        score = r.get("score", "?")
        feedback = r.get("feedback", "")[:500]
        solution = r.get("solution", "")[:500]
        reply_markup = {
            "inline_keyboard": [
                [
                    {"text": "Again ❌", "callback_data": f"rate:{qid}:Again"},
                    {"text": "Hard 🟡", "callback_data": f"rate:{qid}:Hard"},
                    {"text": "Good ✅", "callback_data": f"rate:{qid}:Good"},
                    {"text": "Easy 🔥", "callback_data": f"rate:{qid}:Easy"},
                ],
                [
                    {"text": "⏭️ Question suivante", "callback_data": "next:0"},
                ],
            ]
        }
        send_message(
            chat_id,
            f"{emoji} *Score: {score}/5*\n\n{feedback}\n\n*Solution:*\n{solution}",
            reply_markup=reply_markup,
        )
        return None

    if cmd == "/preview":
        parts = args.split(maxsplit=1)
        if not parts:
            return "Usage: `/preview [id]`"
        qid = parts[0]
        r = curl("POST", "/rag/grade", {"question_id": int(qid), "preview": True})
        return f"*Question:* {r.get('question', '')}\n\n*Solution:*\n{r.get('solution', 'N/A')}"

    if cmd == "/add":
        if not args:
            return "Usage: `/add [question]`"
        send_action(chat_id, "typing")
        r = curl(
            "POST",
            "/rag/add-question",
            {"question": args, "solve": True, "difficulty": 2},
            timeout=120,
        )
        if r.get("id"):
            sol = r.get("solution", "")[:300] if r.get("solution") else "En cours..."
            return f"✅ Question #{r['id']} ajoutée!\n\n*Q:* {r['question'][:200]}\n\n*Solution:*\n{sol}"
        return f"❌ Erreur: {r.get('error', 'inconnue')}"

    return f"Commande inconnue: {cmd}. Tapez /help pour les commandes disponibles."


def handle_callback(update):
    cb = update.get("callback_query", {})
    if not cb:
        return
    chat_id = cb.get("message", {}).get("chat", {}).get("id")
    msg_id = cb.get("message", {}).get("message_id")
    data = cb.get("data", "")
    log.info(f"Callback from {chat_id}: {data}")

    try:
        if data.startswith("preview:"):
            qid = data.split(":", 1)[1]
            r = curl("POST", "/rag/grade", {"question_id": int(qid), "preview": True})
            reply_markup = {
                "inline_keyboard": [
                    [
                        {"text": "Again ❌", "callback_data": f"rate:{qid}:Again"},
                        {"text": "Hard 🟡", "callback_data": f"rate:{qid}:Hard"},
                        {"text": "Good ✅", "callback_data": f"rate:{qid}:Good"},
                        {"text": "Easy 🔥", "callback_data": f"rate:{qid}:Easy"},
                    ],
                    [
                        {"text": "⏭️ Question suivante", "callback_data": "next:0"},
                    ],
                ]
            }
            send_message(
                chat_id,
                f"*Solution:*\n{r.get('solution', 'N/A')[:500]}\n\n_Quel était ton niveau de confiance?_",
                reply_markup=reply_markup,
            )
        elif data.startswith("rate:"):
            parts = data.split(":")
            if len(parts) < 3:
                return
            qid = parts[1]
            rating = parts[2]
            is_correct = None
            if rating in ("Good", "Easy"):
                is_correct = True
            elif rating == "Again":
                is_correct = False
            r = curl(
                "POST",
                "/rag/rate",
                {"question_id": int(qid), "rating": rating, "correct": is_correct},
            )
            msg = r.get("message", "")
            state_emoji = {"learning": "🔵", "review": "🟢", "new": "⚪"}.get(
                r.get("state", ""), "🔵"
            )
            lines = [
                f"{state_emoji} *{rating}* — {msg}",
            ]
            keyboard = [[{"text": "⏭️ Question suivante", "callback_data": "next:0"}]]
            send_message(
                chat_id, "\n".join(lines), reply_markup={"inline_keyboard": keyboard}
            )
            record_session(chat_id, is_correct)
        elif data.startswith("next:") and data != "next:0":
            qid = data.split(":")[1]
            r = curl("POST", "/rag/grade", {"question_id": int(qid), "preview": True})
            text = f"*Solution:*\n{r.get('solution', 'N/A')}\n\n_Tape ta réponse: /grade {qid} Ta réponse_"
            send_message(chat_id, text)
        elif data.startswith("next:0"):
            handle_command(chat_id, "/practice")
        elif data.startswith("qz:"):
            parts = data.split(":")
            if len(parts) < 3:
                return
            qz_chat_id = int(parts[1])
            choice = int(parts[2])
            if qz_chat_id not in quiz_state:
                send_message(chat_id, "Quiz expiré. Tappe /quiz pour recommencer.")
                return
            qz = quiz_state[qz_chat_id]
            q = qz["questions"][qz["current"]]
            selected_opt, selected_id = q["options"][choice]
            correct_opt, correct_id = next(
                (o, o[1]) for o in q["options"] if o[1] == q["id"]
            )
            is_correct = selected_id == q["id"]
            if is_correct:
                qz["correct"] += 1
            qz["total"] += 1
            record_session(chat_id, is_correct)

            lines = [
                f"{'✅' if is_correct else '❌'} Tu as choisi: {selected_opt[:60]}"
            ]
            if not is_correct:
                lines.append(f"→ Correct: {correct_opt[:60]}")
            lines.append(f"\nScore: {qz['correct']}/{qz['total']}")
            send_message(chat_id, "\n".join(lines))

            qz["current"] += 1
            if qz["current"] >= qz["total"]:
                total = qz["total"]
                lines = [f"*Quiz terminé! 🎉*"]
                lines.append(
                    f"Score final: {qz['correct']}/{total} ({round(qz['correct'] / total * 100)}%)"
                )
                send_message(chat_id, "\n".join(lines))
                record_session(chat_id, None)
                del quiz_state[qz_chat_id]
            else:
                nq = qz["questions"][qz["current"]]
                subj = qz.get("subject", "")
                qnum = qz["current"] + 1
                total = len(qz["questions"])
                lines = [
                    f"*Quiz — Question {qnum}/{total}*" + (f" ({subj})" if subj else "")
                ]
                lines.append(f"\n{nq['text']}\n")
                for i, (opt, _) in enumerate(nq["options"]):
                    lines.append(f"{i + 1}. {opt[:70]}")
                keyboard = [
                    [{"text": str(i + 1), "callback_data": f"qz:{qz_chat_id}:{i}"}]
                    for i in range(4)
                ]
                send_message(
                    chat_id,
                    "\n".join(lines),
                    reply_markup={"inline_keyboard": keyboard},
                )

    except Exception as e:
        log.error(f"Callback error: {e}")
        send_message(chat_id, f"❌ Erreur: {e}")


def telegram_api(method, path, data=None, timeout=120):
    """Call Telegram Bot API."""
    cmd = ["curl", "-s", f"https://api.telegram.org/bot{TOKEN}/{path}"]
    if data:
        cmd.extend(["-H", "Content-Type: application/json", "-d", json.dumps(data)])
    result = subprocess.run(cmd, capture_output=True, text=True, timeout=timeout + 10)
    try:
        return json.loads(result.stdout)
    except:
        return {}


def handle_photo_message(chat_id, photo, file_id):
    """Handle photo message - download, OCR, and query RAG."""
    import tempfile
    import os
    import base64

    try:
        # Send "processing" message
        send_message(chat_id, "📸 Image reçue! Traitement en cours...")
        send_action(chat_id, "typing")

        # Get file path from Telegram
        file_info = telegram_api("getFile", f"getFile?file_id={file_id}")
        file_path = file_info.get("result", {}).get("file_path")

        if not file_path:
            log.error(f"getFile failed: {file_info}")
            send_message(chat_id, "❌ Impossible de télécharger l'image")
            return

        # Download the file to temp
        with tempfile.NamedTemporaryFile(suffix=".jpg", delete=False) as tmp:
            tmp_path = tmp.name

        download_result = subprocess.run(
            [
                "curl",
                "-s",
                f"https://api.telegram.org/bot{TOKEN}/file/{file_path}",
                "-o",
                tmp_path,
            ],
            capture_output=True,
        )

        # Check if file was downloaded
        import os

        file_size = os.path.getsize(tmp_path) if os.path.exists(tmp_path) else 0
        log.info(f"Downloaded image: {file_size} bytes")

        if file_size < 100:
            send_message(chat_id, f"❌ Image trop petite: {file_size} bytes")
            os.unlink(tmp_path)
            return

        # Try Workers AI OCR first (cloud), fallback to local
        ocr_resp = None
        extracted_text = ""
        ocr_source = "unknown"
        try:
            # Read image as base64
            img_base64 = base64.b64encode(open(tmp_path, "rb").read()).decode()
            log.info(f"Sending {len(img_base64)} bytes to OCR")

            ocr_cmd = [
                "curl",
                "-s",
                "-X",
                "POST",
                f"{WORKER}/rag/ocr",
                "-H",
                "Content-Type: application/json",
                "-d",
                json.dumps({"image": img_base64}),
                "--max-time",
                "60",
            ]
            ocr_result = subprocess.run(
                ocr_cmd, capture_output=True, text=True, timeout=65
            )
            log.info(f"OCR response: {ocr_result.stdout[:500]}")

            try:
                ocr_resp = json.loads(ocr_result.stdout)
                if ocr_resp.get("text"):
                    extracted_text = ocr_resp["text"]
                    ocr_source = ocr_resp.get("source", "Workers AI")
                    log.info(
                        f"OCR success: {len(extracted_text)} chars via {ocr_source}"
                    )
                else:
                    raise Exception(f"No text. Response: {ocr_result.stdout[:200]}")
            except Exception as e:
                log.error(f"OCR parse failed: {e}")
                raise Exception("OCR parse failed")
        except Exception as e:
            log.info(f"Workers AI OCR failed: {e}, trying local PaddleOCR")

            # Try PaddleOCR locally
            try:
                from paddleocr import PaddleOCR

                ocr = PaddleOCR(
                    use_angle_cls=True, lang="fr", use_gpu=False, show_log=False
                )
                result = ocr.ocr(tmp_path, cls=True)

                if result and result[0]:
                    text_lines = [
                        str(line[1][0]) for line in result[0] if line and len(line) >= 2
                    ]
                    extracted_text = "\n".join(text_lines)
                    ocr_source = "PaddleOCR"
                else:
                    extracted_text = ""
                    ocr_source = "PaddleOCR"
            except ImportError:
                # Fallback to Tesseract
                tesseract_result = subprocess.run(
                    ["tesseract", tmp_path, "stdout", "-l", "fra"],
                    capture_output=True,
                    text=True,
                    timeout=60,
                )
                extracted_text = (
                    tesseract_result.stdout.strip()
                    if tesseract_result.returncode == 0
                    else ""
                )
                ocr_source = "Tesseract"
            except Exception as e2:
                log.error(f"Local OCR failed: {e2}")
                send_message(chat_id, f"❌ OCR échoué: {str(e2)[:200]}")
                os.unlink(tmp_path)
                return

        # Clean up temp file
        os.unlink(tmp_path)

        if not extracted_text or len(extracted_text.strip()) < 5:
            send_message(chat_id, "❌ Aucun texte détecté dans l'image")
            return

        # Show extracted text preview with provider info
        preview = (
            extracted_text[:300] + "..."
            if len(extracted_text) > 300
            else extracted_text
        )

        # Get chain info from OCR response if available
        chain_info = ""
        if ocr_resp and "all_results" in ocr_resp:
            tried = ocr_resp.get("all_results", [])
            if len(tried) > 1:
                providers_tried = [r["provider"] for r in tried]
                chain_info = f"\n\n_🔗 Chaîne: {' → '.join(providers_tried)}_"

        send_message(
            chat_id, f"📝 *Texte détecté ({ocr_source})*\n_{preview}_{chain_info}"
        )

        # Query RAG with extracted text
        send_action(chat_id, "typing")
        r = curl("POST", "/rag/query", {"query": extracted_text}, timeout=120)

        if r.get("answer"):
            answer = r["answer"][:2000]
            sources = ""
            if r.get("sources"):
                subs = list(set(s.get("subject", "") for s in r["sources"][:3]))
                sources = f"\n\n_Sources: {', '.join(subs)}_"
            send_message(chat_id, f"💡 *Réponse:*\n{answer}{sources}")
        elif r.get("error"):
            send_message(
                chat_id,
                f"❌ Erreur RAG: {r['error'][:200]}\n\n📝 Texte: {extracted_text[:500]}",
            )
        else:
            send_message(chat_id, f"📝 Texte détecté: {extracted_text[:1000]}")

    except Exception as e:
        log.error(f"Photo handling error: {e}")
        send_message(chat_id, f"❌ Erreur: {str(e)[:200]}")


def handle_message(update):
    msg = update.get("message", {})
    chat = msg.get("chat", {})
    text = msg.get("text", "")
    photo = msg.get("photo", [])
    voice = msg.get("voice", {})

    chat_id = chat.get("id")

    # Handle photo
    if photo:
        file_id = photo[-1].get("file_id")
        if file_id:
            log.info(f"Handling photo from {chat_id}")
            handle_photo_message(chat_id, photo, file_id)
            return

    # Handle voice (future: whisper transcription)
    if voice:
        send_message(
            chat_id,
            "🎤 Messages vocaux non encore supportés. Envoyez une image à la place!",
        )
        return

    if not text:
        return

    is_group = chat.get("type") in ("group", "supergroup")

    if is_group:
        if text.startswith("/"):
            parts = text[1:].split(maxsplit=1)
            if not parts[0].lower().startswith("bac"):
                return
            text = parts[1] if len(parts) > 1 else ""

    log.info(f"Handling message from {chat_id}: {text[:50]}")

    try:
        if text.startswith("/"):
            response = handle_command(chat_id, text)
            if response:
                send_message(chat_id, response)
        else:
            send_action(chat_id, "typing")
            r = curl("POST", "/rag/query", {"query": text}, timeout=120)
            if r.get("answer"):
                answer = r["answer"][:2000]
                sources = ""
                if r.get("sources"):
                    subs = list(set(s.get("subject", "") for s in r["sources"][:3]))
                    sources = f"\n\n_Sources: {', '.join(subs)}_"
                send_message(chat_id, f"*{text[:100]}*\n\n{answer}{sources}")
            elif r.get("error"):
                send_message(chat_id, f"❌ {r['error'][:200]}")
            else:
                send_message(chat_id, "Pas de réponse trouvée.")
    except Exception as e:
        log.error(f"Error handling message: {e}")
        send_message(chat_id, f"❌ Erreur: {e}")


def main():
    if not TOKEN:
        log.error("TELEGRAM_BOT_TOKEN not set!")
        return

    cmd = [
        "curl",
        "-s",
        "-N",
        f"https://api.telegram.org/bot{TOKEN}/getUpdates",
        "-H",
        "Content-Type: application/json",
        "-d",
        json.dumps({"timeout": 30, "offset": -1}),
    ]
    result = subprocess.run(cmd, capture_output=True, text=True, timeout=35)
    try:
        data = json.loads(result.stdout)
        updates = data.get("result", [])
        offset = None
        for update in updates:
            offset = update.get("update_id", 0) + 1
            if update.get("callback_query"):
                handle_callback(update)
            else:
                handle_message(update)
        if offset:
            subprocess.run(
                [
                    "curl",
                    "-s",
                    "-X",
                    "POST",
                    f"https://api.telegram.org/bot{TOKEN}/getUpdates",
                    "-H",
                    "Content-Type: application/json",
                    "-d",
                    json.dumps({"offset": offset}),
                ],
                capture_output=True,
                timeout=5,
            )
    except Exception as e:
        log.error(f"Polling error: {e}")


if __name__ == "__main__":
    print("Starting BAC Telegram Bot...")
    print("Set TELEGRAM_BOT_TOKEN env var to run.")
    if TOKEN:
        print(f"Worker: {WORKER}")
        while True:
            main()
