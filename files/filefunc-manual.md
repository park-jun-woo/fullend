# filefunc — Manual for AI Agents

## How to Navigate

1. Read `codebook.yaml` — project vocabulary (required/optional keys and allowed values)
2. `filefunc chain func <target> --chon 2` — trace call relationships before modifying
3. `rg '//ff:func feature=validate'` — grep with codebook values to find files
4. Read `//ff:what` to narrow down — skip body if what is sufficient
5. Full read only the files you need, then work

---

## Rules

### File structure

| Rule | Description | Severity |
|---|---|---|
| F1 | One func per file (filename = func name, snake_case) | ERROR |
| F2 | One type per file (filename = type name, snake_case) | ERROR |
| F3 | One method per file (`{receiver}_{method}.go`) | ERROR |
| F4 | init() must not exist alone (requires var or func) | ERROR |
| F5 | _test.go files may have multiple funcs | exception |
| F6 | Semantically grouped consts allowed in one file | exception |

### Code quality

| Rule | Description | Severity |
|---|---|---|
| Q1 | Nesting depth ≤ 2 (use early return, merge conditions, extract func) | ERROR |
| Q2 | Func max 1000 lines | ERROR |
| Q3 | Func recommended max 100 lines | WARNING |

### Annotation

| Rule | Description | Severity |
|---|---|---|
| A1 | Func files require `//ff:func`, type files require `//ff:type` | ERROR |
| A2 | Annotation values must exist in codebook | ERROR |
| A3 | Func/type files require `//ff:what` | ERROR |
| A6 | Annotations must be at the top of the file (above package) | ERROR |
| A7 | `//ff:checked` hash mismatch → body changed after LLM verification | ERROR |
| A8 | Required codebook keys must be present in annotation | ERROR |

### Codebook format

| Rule | Description | Severity |
|---|---|---|
| C1 | `required` section must have at least one key with at least one value | ERROR |
| C2 | No duplicate values within the same key | ERROR |
| C3 | All values lowercase + hyphens only (`[a-z][a-z0-9-]*`) | ERROR |

Codebook is validated first. If codebook fails, code validation does not run.

### Exceptions (not violations)

- const-only and var-only files do not require annotations
- If no `//ff:checked` exists in the project, A7 is skipped entirely

---

## Annotations

Write at the **very top** of every func/type file (above package declaration):

```go
//ff:func feature=validate type=rule
//ff:what F1: validates one func per file
//ff:checked llm=gpt-oss:20b hash=a3f8c1d2    (auto by llmc)
package validate
```

| Annotation | Required | Description |
|---|---|---|
| `//ff:func` | func files | Metadata (feature, type). Values from codebook.yaml |
| `//ff:type` | type files | Metadata (feature, type). Values from codebook.yaml |
| `//ff:what` | func/type files | One-line description. What does this do? |
| `//ff:why` | optional | Why designed this way? User decisions only |
| `//ff:checked` | auto (llmc) | LLM verification signature. Do not write manually |

### Naming

- Filenames: `snake_case`
- Variables/functions: `camelCase`
- Types: `PascalCase`
- gofmt compliance, early return pattern

---

## Codebook

`codebook.yaml` must exist in the project root (next to `go.mod`). `required` keys must be in every annotation (A8). `optional` keys are used when relevant.

```yaml
required:
  feature: [validate, annotate, chain, parse, codebook, report, cli]
  type: [command, rule, parser, walker, model, formatter, loader, util]

optional:
  pattern: [error-collection, file-visitor, rule-registry]
  level: [error, warning, info]
```

Amend codebook.yaml when new values are needed.

---

## Commands

```bash
filefunc validate ./internal/                        # validate (codebook auto-detected)
filefunc validate --format json ./internal/          # JSON output
filefunc chain func RunAll --chon 2                  # call relationships
filefunc chain feature validate                      # feature-wide chain
filefunc llmc ./internal/                            # LLM what-body verification
filefunc llmc --model qwen3:8b --threshold 0.9 ./internal/
```

Exit code 1 on violations. Zero violations required before committing.

### .ffignore

Place in project root. Same syntax as `.gitignore`. Excludes paths from all commands.

```
vendor/
*.pb.go
*_gen.go
```

---

## Common Mistakes

| Mistake | Fix |
|---|---|
| Two funcs in one file | Extract helper functions into separate files |
| depth 3 (for→switch→if, for→if→if) | Type assertions + early continue, merge conditions, or extract func |
| Missing //ff:what | Write annotations first when creating a file |
| Value not in codebook | Check codebook.yaml first. Amend if absent |
| //ff:checked hash mismatch | Run `filefunc llmc` to re-verify |
| "codebook.yaml required" | Create codebook.yaml next to go.mod |
