# Mutation Test — SSaC ↔ States

### MUT-SSAC-STATES-001: 상태 전이 누락
- 대상: `specs/gigbridge/states/gig.md`
- 변경: `under_review --> disputed : RaiseDispute` 행 삭제
- 기대: ERROR — SSaC RaiseDispute의 @state 전이가 States에 없음
- 결과: PASS — States↔SSaC 전이 누락 검출
