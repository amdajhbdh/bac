#!/bin/bash
# BAC Study CLI shortcuts — add to ~/.bashrc or ~/.zshrc
# source ~/Documents/bac/scripts/bac-aliases.sh

export BAC_WORKER="https://bac-api.amdajhbdh.workers.dev"

bac() {
	python3 ~/Documents/bac/scripts/bac-cli.py "$@"
}

# Shortcuts
alias bq="bac query"
alias bs="bac solve"
alias bp="bac practice"
alias bt="bac track"
alias ba="bac add"
alias bst="bac status"
alias bsrch="bac search"
alias bg="bac grade"
alias bdue="bac due"
alias bl="bac loop"
alias blloop="bac loop -n 20"
alias bqz="bac quiz -n 5"
alias bretry="bac retry"

# Practice loop — keeps giving questions until Ctrl+C
bac-practice-loop() {
	while true; do
		echo ""
		echo "=== BAC Practice ==="
		result=$(python3 ~/Documents/bac/scripts/bac-cli.py practice)
		echo "$result"
		echo ""
		qid=$(echo "$result" | grep -oP 'Question #\K\d+')
		echo -n "Ta réponse: "
		read -r answer
		if [ -n "$answer" ]; then
			echo ""
			python3 ~/Documents/bac/scripts/bac-cli.py grade "$qid" "$answer"
		fi
		echo ""
		echo -n "Continue? (Enter for next, q to quit): "
		read -r cont
		[ "$cont" = "q" ] && break
	done
}
alias bstats='bac stats'

# Quick RAG query from CLI (default: human-readable, --json: JSON output)
rag() {
	local JSON_FLAG=""
	local QUERY=""

	# Parse args
	while [[ $# -gt 0 ]]; do
		case $1 in
		--json | -j)
			JSON_FLAG=1
			shift
			;;
		*)
			QUERY="$1"
			shift
			;;
		esac
	done

	if [[ -z "$QUERY" ]]; then
		echo "Usage: rag \"question\" [--json]" >&2
		return 1
	fi

	local RESPONSE=$(curl -s -X POST "https://bac-api.amdajhbdh.workers.dev/rag/query" \
		-H "Content-Type: application/json" \
		-d "{\"query\": \"$QUERY\"}")

	if [[ -n "$JSON_FLAG" ]]; then
		echo "$RESPONSE"
	else
		echo "$RESPONSE" | jq -r '
			"=== Réponse ===",
			.answer,
			"",
			"=== Sources ===",
			(.sources[:3] | .[] | "\(.subject): \(.source)")
		'
	fi
}

# Quick solve from CLI (default: solution only, --json: full JSON)
solve() {
	local JSON_FLAG=""
	local QUESTION=""

	while [[ $# -gt 0 ]]; do
		case $1 in
		--json | -j)
			JSON_FLAG=1
			shift
			;;
		*)
			QUESTION="$1"
			shift
			;;
		esac
	done

	if [[ -z "$QUESTION" ]]; then
		echo "Usage: solve \"question\" [--json]" >&2
		return 1
	fi

	local RESPONSE=$(curl -s -X POST "https://bac-api.amdajhbdh.workers.dev/rag/solve" \
		-H "Content-Type: application/json" \
		-d "{\"question\": \"$QUESTION\"}")

	if [[ -n "$JSON_FLAG" ]]; then
		echo "$RESPONSE"
	else
		echo "$RESPONSE" | jq -r '.solution // "No solution"'
	fi
}

# Quick search from CLI
search() {
	local JSON_FLAG=""
	local QUERY=""

	while [[ $# -gt 0 ]]; do
		case $1 in
		--json | -j)
			JSON_FLAG=1
			shift
			;;
		*)
			QUERY="$1"
			shift
			;;
		esac
	done

	if [[ -z "$QUERY" ]]; then
		echo "Usage: search \"query\" [--json]" >&2
		return 1
	fi

	local RESPONSE=$(curl -s -X POST "https://bac-api.amdajhbdh.workers.dev/rag/search" \
		-H "Content-Type: application/json" \
		-d "{\"query\": \"$QUERY\", \"topK\": 5}")

	if [[ -n "$JSON_FLAG" ]]; then
		echo "$RESPONSE"
	else
		echo "$RESPONSE" | jq -r '
			.results = (.chunks | length) |
			(.chunks[] | "\(.metadata.subject) [\(.score | . * 100 | floor)%]: \(.metadata.text[:100])...")
		'
	fi
}

# Status (always JSON for scripting)
status() {
	curl -s "https://bac-api.amdajhbdh.workers.dev/rag/status"
}

# Subject counts (always JSON)
subjects() {
	curl -s "https://bac-api.amdajhbdh.workers.dev/rag/subjects"
}
