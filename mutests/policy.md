# Mutation Test — Policy 단독

### MUT-POLICY-001: Rego 파싱 오류
- 대상: `specs/gigbridge/policy/authz.rego`
- 변경: `allow {` → `allow (` (중괄호 → 소괄호)
- 기대: ERROR — Rego 파싱 실패
- 결과: PASS — policy.ParseDir에서 Rego 구문 오류 검출

### MUT-POLICY-002: Rego 파일 확장자 오류
- 대상: `specs/gigbridge/policy/authz.rego` → `authz.txt` 리네임
- 변경: .rego 확장자 제거
- 기대: SKIP — policy/ 디렉토리에 .rego 파일 없음으로 인식
- 결과: PASS — DetectSSOTs에서 *.rego glob 매칭 실패, KindPolicy 미등록
