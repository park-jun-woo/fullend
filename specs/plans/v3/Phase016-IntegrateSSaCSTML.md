# Phase016: SSaC·STML fullend 통합

## 목표

SSaC(parser/validator/generator)와 STML(parser/validator/generator)을 fullend `internal/`로 이동한다. 향후 원천은 fullend가 통합 관리하며, SSaC/STML repo에는 단순 파일 복사로 내려준다.

## 결정 사항

| 항목 | 결정 |
|------|------|
| 원천 | fullend 단일 관리 |
| SSaC/STML repo | fullend에서 해당 디렉토리 단순 복사 (독립 빌드 보장 안 함) |
| SSaC v1/ | 레거시 — 복사 대상 아님 |
| SSaC/STML cmd/ | fullend CLI로 통합 — 복사 대상 아님 |

## 현황

### fullend → SSaC 의존 (33개 파일, 51회 참조)

```
github.com/geul-org/ssac/parser       → ssacparser.ServiceFunc 등
github.com/geul-org/ssac/validator     → ssacvalidator.DDLTable, Index, SymbolTable 등
github.com/geul-org/ssac/generator     → ssacgenerator.Generate 등
```

### fullend → STML 의존 (4개 파일, 6회 참조)

```
github.com/geul-org/stml/parser        → stmlparser.Page 등
github.com/geul-org/stml/validator     → stmlvalidator.Validate 등
github.com/geul-org/stml/generator     → stmlgenerator.Generate 등
```

### SSaC 코드 구조 (이동 대상)

| 디렉토리 | 파일 수 | 크기 | 역할 |
|----------|--------|------|------|
| parser/ | 3 | ~19KB | .ssac 파서 + 타입 |
| validator/ | 5 | ~60KB | DDL 파서, 심볼 테이블, 검증 |
| generator/ | 14 | ~84KB | Go handler/interface/args 코드젠 |
| testdata/ | 1+ | — | 테스트 데이터 |

외부 의존: `github.com/ettle/strcase`, `gopkg.in/yaml.v3` (fullend에 이미 존재)

### STML 코드 구조 (이동 대상)

| 디렉토리 | 파일 수 | 크기 | 역할 |
|----------|--------|------|------|
| parser/ | 3 | ~21KB | HTML 파서 + 타입 |
| validator/ | 4 | ~14KB | STML ↔ OpenAPI 검증 |
| generator/ | 6 | ~44KB | React TSX 코드젠 |

외부 의존: `golang.org/x/net` (HTML parser)

## 변경 계획

### 1단계: 파일 이동

```
ssac/parser/     → fullend/internal/ssac/parser/
ssac/validator/  → fullend/internal/ssac/validator/
ssac/generator/  → fullend/internal/ssac/generator/
ssac/testdata/   → fullend/internal/ssac/testdata/

stml/parser/     → fullend/internal/stml/parser/
stml/validator/  → fullend/internal/stml/validator/
stml/generator/  → fullend/internal/stml/generator/
```

### 2단계: import 경로 변경

모든 fullend Go 파일에서:

```
github.com/geul-org/ssac/parser     → github.com/geul-org/fullend/internal/ssac/parser
github.com/geul-org/ssac/validator  → github.com/geul-org/fullend/internal/ssac/validator
github.com/geul-org/ssac/generator  → github.com/geul-org/fullend/internal/ssac/generator
github.com/geul-org/stml/parser     → github.com/geul-org/fullend/internal/stml/parser
github.com/geul-org/stml/validator  → github.com/geul-org/fullend/internal/stml/validator
github.com/geul-org/stml/generator  → github.com/geul-org/fullend/internal/stml/generator
```

SSaC/STML 내부 파일의 상호 import도 변경:

```
github.com/geul-org/ssac/parser     → github.com/geul-org/fullend/internal/ssac/parser
(ssac generator → ssac parser 참조 등)
```

### 3단계: go.mod 정리

- `github.com/geul-org/ssac` require + replace 제거
- `github.com/geul-org/stml` require + replace 제거
- `golang.org/x/net` 의존 추가 (STML용)
- `go mod tidy`

### 4단계: 빌드 + 테스트

```bash
go build ./cmd/fullend/
go test ./...
./fullend validate specs/gigbridge
./fullend gen specs/gigbridge artifacts/gigbridge
```

### 5단계: SSaC/STML repo 동기화

```bash
# SSaC repo에 복사
cp -r fullend/internal/ssac/parser/*    ssac/parser/
cp -r fullend/internal/ssac/validator/* ssac/validator/
cp -r fullend/internal/ssac/generator/* ssac/generator/

# STML repo에 복사
cp -r fullend/internal/stml/parser/*    stml/parser/
cp -r fullend/internal/stml/validator/* stml/validator/
cp -r fullend/internal/stml/generator/* stml/generator/
```

## 변경 파일 요약

| 범위 | 파일 수 | 변경 |
|------|--------|------|
| SSaC 파일 복사 | ~22 | 신규 (internal/ssac/) |
| STML 파일 복사 | ~13 | 신규 (internal/stml/) |
| fullend import 변경 | 33 | ssac import 경로 |
| fullend import 변경 | 4 | stml import 경로 |
| SSaC 내부 import 변경 | ~14 | generator → parser 등 내부 참조 |
| STML 내부 import 변경 | ~6 | generator → parser 등 내부 참조 |
| go.mod | 1 | require/replace 정리 |

## 검증

```bash
go build ./cmd/fullend/
go test ./...
./fullend validate specs/gigbridge
./fullend gen specs/gigbridge artifacts/gigbridge
cd artifacts/gigbridge/backend && go build -o server ./cmd/
```

## 리스크

| 리스크 | 대응 |
|--------|------|
| SSaC/STML 내부 import 누락 | `go build`로 즉시 발견 |
| testdata 경로 변경 | 상대 경로 사용 중이면 수정 |
| SSaC/STML repo 독립 빌드 깨짐 | 설계상 허용 (단순 복사 미러) |
