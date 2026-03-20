# Mutation Test — SSaC ↔ Queue

### MUT-SSAC-QUEUE-001: @publish 사용하지만 queue 미설정
- 대상: `specs/gigbridge/fullend.yaml`
- 변경: `queue:` 섹션 삭제
- 기대: ERROR — SSaC에 @publish 시퀀스가 있지만 fullend.yaml에 queue.backend 미설정
- 결과: PASS — crosscheck/queue에서 @publish 존재 + queue 미설정 검출

### MUT-SSAC-QUEUE-002: @publish topic에 대응하는 @subscribe 없음
- 대상: `specs/gigbridge/service/gig/publish_gig.ssac`
- 변경: @publish topic을 `gig.published.TYPO`로 변경
- 기대: WARNING — 토픽에 대응하는 @subscribe 함수 없음
- 결과: PASS — crosscheck/queue에서 publish→subscribe 토픽 매칭 검증

### MUT-SSAC-QUEUE-003: @subscribe message 필드 불일치
- 대상: @subscribe 함수의 message struct 정의
- 변경: @subscribe message에 @publish payload에 없는 필드 추가
- 기대: WARNING — @subscribe message struct 필드가 @publish payload에 없음
- 결과: 미실행
