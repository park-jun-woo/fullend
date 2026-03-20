# Mutation Test — Config 단독

### MUT-CONFIG-001: fullend.yaml 문법 오류
- 대상: `specs/gigbridge/fullend.yaml`
- 변경: `backend:` → `backend` (콜론 제거)
- 기대: ERROR — YAML 파싱 실패
- 결과: PASS — Config 파서 단계에서 YAML 구조 오류 검출

### MUT-CONFIG-002: metadata name 누락
- 대상: `specs/gigbridge/fullend.yaml`
- 변경: `name: gigbridge` 라인 삭제
- 기대: WARNING — 프로젝트 이름 누락
- 결과: PASS — validateConfig에서 metadata 필수 필드 검사
