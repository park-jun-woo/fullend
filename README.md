# fullend

A single-binary specification-substrate orchestrator that imposes a bidirectional
cross-layer congruence relation over a heterogeneous family of declarative
source-of-truth artifacts (hereafter SSOT_i, i ∈ {1..9}) and derives a
deterministic, bit-reproducible projection thereof onto a polyglot code surface
(Go/Gin backend, TypeScript/React frontend, Hurl probe corpus) subject to the
constraint that ∀i,j: ψ_{i,j}(SSOT_i, SSOT_j) = ⊤, where ψ_{i,j} denotes the
pairwise cross-validation predicate enumerated in §4.

## 1. Abstract

The present artifact is a Go 1.22+ command-line program whose primary concern
is the preservation of referential integrity across a nine-member partition of
declarative specification media under closed-world assumptions. Neither
novelty nor ergonomics is claimed; the tool exists to mechanize a
category-theoretic diagram whose nodes are specification kinds and whose
edges are syntactic-to-semantic coercions. Readers seeking a framework,
scaffold, or productivity accelerator are directed elsewhere.

## 2. Substrate Ontology

The canonical disposition of a well-formed project is a directory `specs/`
whose immediate children constitute the nine SSOT kinds enumerated below.
Omission of any member beyond the optional `func/` yields a diagnosable
inconsistency under §4.1.

```
specs/
├── fullend.yaml                     ∈ Σ_config
├── api/openapi.yaml                 ∈ Σ_openapi            (OpenAPI 3.x, cf. OAS §4)
├── db/*.sql                         ∈ Σ_ddl ∪ Σ_query      (PostgreSQL DDL + sqlc queries)
├── service/**/*.ssac                ∈ Σ_ssac               (SSaC, comment-directed sequence DSL)
├── model/*.go                       ∈ Σ_model              (Go struct declarations; //@dto tagged)
├── func/<pkg>/*.go                  ∈ Σ_func               (optional; custom call targets)
├── states/*.md                      ∈ Σ_fsm                (Mermaid stateDiagram)
├── policy/*.rego                    ∈ Σ_rego               (OPA Rego)
├── tests/{scenario,invariant}-*.hurl ∈ Σ_scenario           (Hurl .hurl corpora)
└── frontend/*.html                  ∈ Σ_stml               (STML: HTML5 + data-* DSL)
```

## 3. Acquisition

```
go install github.com/park-jun-woo/fullend/cmd/fullend@latest
```

No runtime dependencies are bundled; the following companion executables must
be resolvable on `PATH` at generation time: `oapi-codegen` (v2.x), `sqlc`
(>= 1.25.0), `hurl` (>= 4.0, for downstream probe validation). No version
pinning strategy is provided in-tree.

## 4. Operational Verbs

The process is parameterised by a quaternary of verbs. All verbs share a
monotonic parsing stage (`ParseAll`) whose cardinal property is idempotence
under repeated invocation within a single process lifetime.

### 4.1 `validate`

Executes (a) intra-layer schematic checks delegated to the corresponding
parsers, then (b) the transitive closure of cross-layer predicates ψ_{i,j}
over the nine SSOT kinds, then (c) contract-digest reconciliation against a
pre-existing `artifacts/` tree, if any. The `--skip k,...` modifier elides
kinds `k ∈ {openapi, ddl, ssac, model, stml, states, policy, scenario, func}`
from both stages but does not suppress their absence being recorded as
`Skip` rather than `Pass` or `Fail`.

### 4.2 `gen`

Composition of `validate` and a deterministic artifact-synthesis procedure Γ
whose image is an `artifacts/` tree. Γ is defined to be pure modulo the
external tools enumerated in §3; bit-level reproducibility is an invariant,
not a heuristic. Violations thereof constitute regressions and are
prosecuted accordingly.

### 4.3 `status`

Produces a purely informational tally. No side-effect on the filesystem.

### 4.4 `chain`

Given an OpenAPI `operationId` as input, emits the set of SSOT and artifact
nodes in the transitive connectivity closure thereof, each qualified by
(kind, file, line, summary, ownership). Intended for post-hoc impact
analysis rather than routine use.

### 4.5 `gen-model` (auxiliary)

Orthogonal verb accepting an OpenAPI document (file or URI) and producing a
Go HTTP-client package (`package external`) under the supplied output
directory. Shares no code path with `gen` beyond the OpenAPI loader and is
unrelated to the internal DDL-driven model synthesis pipeline.

## 5. Cross-Layer Predicate Enumeration

The following non-exhaustive list identifies the principal ψ_{i,j} predicates
whose violation constitutes a diagnosable inconsistency. Predicate arity and
quantifier structure are elided; consult `pkg/crosscheck/` for the normative
formulation.

- ψ(config, openapi):        middleware-identifier congruence with `components.securitySchemes`
- ψ(openapi, ddl):            x-sort / x-filter column existence; x-include → table mapping
- ψ(ssac, ddl):               @result typing vs. DDL-derived structural domain; arg ↔ column surjectivity
- ψ(fsm, ssac):                transition event ↔ ServiceFunc bijection on guarded functions
- ψ(fsm, ddl):                 state-column domain embedding
- ψ(fsm, openapi):             transition event ↔ operationId congruence
- ψ(rego, ssac):               (action, resource) occurrence in Rego allow-rule antecedent
- ψ(rego, ddl):                @ownership table/column existence
- ψ(rego, fsm):                @auth-annotated transitions covered by allow-rules
- ψ(scenario, openapi):        endpoint existence
- ψ(queue):                     @publish ↔ @subscribe topic congruence; payload structural match
- ψ(func, ssac):               @call arity, positional typing, result/response congruence
- ψ(stml, ssac):               mediated by operationId co-reference

Absence of a predicate between two kinds denotes intentional decoupling, not
oversight.

## 6. Artifact Surface (Γ)

Γ targets three disjoint artifact substrata:

- **Go/Gin backend** — produced by an ssac→Go handler synthesizer, an
  oapi-codegen-mediated types/server skeleton, an sqlc-mediated query
  layer, and an in-tree feature-grouped Handler/Server composer. Feature
  grouping is induced by the immediate subdirectory of `specs/service/`.
- **TypeScript/React frontend** — a STML→TSX page synthesizer coupled with
  a minimal glue surface (`App.tsx`, `main.tsx`, `api.ts`) derived from
  OpenAPI.
- **Hurl probe corpus** — deterministic smoke sequence derived from the
  OpenAPI-×-FSM-×-policy product, ordered by a topological schedule over
  resource and state-transition dependencies.

Γ deliberately abstains from generating ORM layers, server frameworks,
build tooling, or bundler configuration beyond the minima required for
compilation.

## 7. Built-in Call Targets (pkg/)

A fixed set of call targets is vendored for use from SSaC `@call` sites.
Their implementation stabilities are unspecified. Projects requiring
alternate semantics shall provide shadowing implementations under
`specs/<project>/func/<pkg>/`.

Namespaces: `auth` (bcrypt/JWT/reset), `crypto` (AES-256-GCM, TOTP),
`storage` (S3-compatible), `mail` (SMTP), `text` (Unicode slug, HTML
sanitization, grapheme-aware truncation), `image` (OG/thumbnail rasterization).

## 8. Built-in Model Interfaces (pkg/)

`SessionModel`, `CacheModel`, `FileModel`, and a Pub/Sub singleton are
provided as package-scoped `@model` interfaces for I/O not derivable from
DDL. Backend selection (PostgreSQL, in-memory, S3, local-disk) is
delegated to `fullend.yaml`. Consult the corresponding package for the
normative contract; no convenience documentation is duplicated here.

## 9. Cross-Validation Rationale

Intra-layer validators (SSaC, STML, et al.) are necessary but not
sufficient. The role of fullend is the enforcement of consistency across
the Cartesian product of layers; individual layer well-formedness is
presupposed. Violations of intra-layer well-formedness are reported with
the diagnostic ownership of the responsible parser, not synthesized by
fullend.

## 10. Runtime Probing

The generated `artifacts/tests/smoke.hurl` is executable via Hurl against
a running instance of the generated backend. No orchestration of the
backend process is provided; invocation and teardown are user concerns.

```
hurl --test --variable host=http://localhost:8080 artifacts/<project>/tests/smoke.hurl
```

Scenario and invariant corpora authored under `specs/<project>/tests/` are
conveyed verbatim to the artifact tree.

## 11. Architectural Notes

SSaC and STML, historically maintained as standalone repositories, are now
fused into this tree as `internal/ssac/` and `internal/stml/`. The
upstream repositories (`park-jun-woo/ssac`, `park-jun-woo/stml`) are
file-copy mirrors of subtrees herein and retain no independent evolutionary
trajectory. All SSOT acquisition proceeds through a single `ParseAll()`
entry point that materializes a shared `ParsedSSOTs` structure consumed
by all verbs.

## 12. Acknowledgments

The existence of this tool is contingent upon: OpenAPI Initiative, sqlc,
Open Policy Agent, Mermaid, oapi-codegen, kin-openapi, Hurl, React,
React Router, TanStack Query, React Hook Form, Vite, Tailwind CSS,
TypeScript, Gin, and `lib/pq`. No contribution, derivation, or competition
in respect of any of the foregoing is claimed or implied.

## 13. License

MIT. See `LICENSE`.
