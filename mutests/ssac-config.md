# Mutation Test — SSaC ↔ Config

### MUT-SSAC-CONFIG-001: currentUser 필드가 claims에 없음
- 대상: `specs/gigbridge/fullend.yaml`
- 변경: claims에서 `user_id` 키 삭제
- 기대: ERROR — SSaC에서 `currentUser.UserID` 참조하지만 claims에 user_id 정의 없음
- 결과: PASS — crosscheck/claims에서 currentUser 필드 → claims 매핑 검증

### MUT-SSAC-CONFIG-002: @auth 사용하지만 claims 설정 없음
- 대상: `specs/gigbridge/fullend.yaml`
- 변경: `backend.auth` 섹션 전체 삭제
- 기대: ERROR — SSaC에 @auth 시퀀스가 있지만 fullend.yaml에 claims 미설정
- 결과: PASS — crosscheck/claims에서 @auth 존재 + claims 미설정 검출
