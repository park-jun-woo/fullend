# Phase011: Terraform SSOT 제거 ✅ 완료

## 목표

Terraform을 fullend 관리 SSOT 목록에서 제거한다. 10개 → 9개 SSOT.

## 배경

Terraform은 인프라 선언(HCL)으로, 애플리케이션 레벨 SSOT와 교차 검증 규칙이 없다:

- DDL ↔ Terraform: DDL에 DB 엔진 정보 없음 → 검증 불가
- fullend.yaml ↔ Terraform: `deploy.*` 필드 미구현
- 기타 SSOT와 연결점 없음

현재 fullend가 Terraform에 대해 하는 일:
- **validate**: `terraform/*.tf` 파일 존재 확인 + 파일 수 카운트
- **gen**: `terraform fmt` exec 호출 (포맷팅만)
- **status**: 파일 수 출력
- **crosscheck**: 없음

파일 존재 확인과 `terraform fmt`는 CI/Makefile 영역이지 SSOT 오케스트레이터의 역할이 아니다.

## 변경 내용

### 1. `internal/orchestrator/detect.go`

- `KindTerraform` 상수 삭제
- `checks` 배열에서 `{KindTerraform, "terraform/*.tf"}` 제거
- `AllSSOTKinds()`에서 `KindTerraform` 제거
- `kindNames` 맵에서 `"terraform"` 제거

### 2. `internal/orchestrator/validate.go`

- `allKinds`에서 `KindTerraform` 제거
- `case KindTerraform:` 분기 삭제
- `validateTerraform` 함수 삭제

### 3. `internal/orchestrator/status.go`

- `case KindTerraform:` 분기 삭제
- `statusTerraform` 함수 삭제

### 4. `internal/orchestrator/gen.go`

- `terraformAvailable` 변수 + `exec.LookPath("terraform")` 체크 삭제
- `genTerraform` 호출부 삭제
- `genTerraform` 함수 삭제

### 5. `cmd/fullend/main.go`

- skip kinds 목록에서 `terraform` 제거

### 6. `artifacts/manual-for-ai.md`

- 디렉토리 구조에서 `terraform/*.tf` 행 삭제
- skip kinds 목록에서 `terraform` 제거

### 7. `CLAUDE.md`

- 프로젝트 개요 "10개 SSOT" → "9개 SSOT", Terraform 제거

### 8. `README.md`

- 프로젝트 설명 "10 SSOT sources" → "9 SSOT sources", Terraform 제거
- skip kinds, status 출력 예시에서 Terraform 행 삭제

### 9. `artifacts/AGENTS.md`

- SSOT 테이블에서 Terraform 행 삭제

## 영향 없는 범위

- `files/`, `specs/plans-v1/` — 과거 기록/아이디어 문서. 수정하지 않음
- `specs/gigbridge/`, `specs/zenflow/` — 프로젝트 specs에 `terraform/` 디렉토리 없으면 영향 없음

## 검증

```
go build ./...
go test ./...
```

전체 통과 확인.
