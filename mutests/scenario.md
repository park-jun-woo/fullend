# Mutation Test — Scenario 단독

### MUT-SCENARIO-001: .feature 파일 사용 (deprecated)
- 대상: `specs/gigbridge/tests/` 디렉토리
- 변경: `scenario-smoke.hurl` → `scenario-smoke.feature` 리네임
- 기대: ERROR — .feature 파일은 deprecated, .hurl 사용 필요
- 결과: PASS — validateScenarioHurl에서 .feature 파일 존재 시 ERROR 반환

### MUT-SCENARIO-002: 빈 tests 디렉토리
- 대상: `specs/gigbridge/tests/` 디렉토리
- 변경: 모든 .hurl 파일 삭제
- 기대: WARNING — 시나리오 테스트 미작성 경고
- 결과: PASS — scenario-*.hurl + invariant-*.hurl 0건 시 WARNING 출력
