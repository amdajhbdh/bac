#!/bin/bash
# BAC Unified CLI Wrapper
# Usage: ./bac solve "problem"

cd "$(dirname "$0")"

if [ "$1" = "solve" ]; then
    shift
    nu -c "source bac.nu; cmd-solve [$(printf '%s\n' "$@" | jq -Rs 'split("\n")[:-1] | map(if . != "" then . else null end) | compact | join(", ")')]"
else
    nu -c "source bac.nu; main $@"
fi
