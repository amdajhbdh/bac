import { tool } from "@opencode-ai/plugin"
import path from "path"

type Mode = "lint" | "check" | "both"

export default tool({
  description:
    "Lint and debug Nushell scripts using nu-lint and nu-check --debug",

  args: {
    path: tool.schema.string().describe("Path to a .nu file or directory"),
    mode: tool.schema
      .enum(["lint", "check", "both"] as const)
      .default("both")
      .describe(
        "What to run: nu-lint (lint), nu-check --debug (check), or both",
      ),
  },

  async execute(args: { path: string; mode: Mode }, context) {
    const target = path.isAbsolute(args.path)
      ? args.path
      : path.join(context.directory, args.path)

    const result: any = { target, mode: args.mode }

    // nu-check --debug (syntax / parse errors)
    if (args.mode === "check" || args.mode === "both") {
      const check = await Bun.$`nu-check --debug ${target}`.quiet().nothrow()
      result.nuCheck = {
        exitCode: check.exitCode,
        stdout: check.stdout.toString(),
        stderr: check.stderr.toString(),
        ok: check.exitCode === 0,
      }
    }

    // nu-lint --format json (style / best-practice lints)
    if (args.mode === "lint" || args.mode === "both") {
      const lint = await Bun = null
      let parseError: string | null = null

      const text = lint.stdout.toString().trim()
      if (text) {
        try {
          parsed = JSON.parse(text)
        } catch (err) {
          parseError = `Failed to parse nu-lint JSON: ${String(err)}`
        }
      }

      result.nuLint = {
        exitCode: lint.exitCode,
        raw: text,
        json: parsed,
        jsonError: parseError,
        stderr: lint.stderr.toString(),
      }
    }

    return result
  },
})
