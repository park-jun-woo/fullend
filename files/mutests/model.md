# Mutation Test — Model 단독

### MUT-MODEL-001: model 디렉토리에 .go 파일 없음
- 대상: `specs/gigbridge/model/`
- 변경: 모든 *.go 파일 삭제
- 기대: FAIL — model 디렉토리에 Go 파일이 없으면 검증 실패
- 결과: PASS — validateModel에서 *.go glob 결과 0건 시 FAIL 반환

### MUT-MODEL-002: @dto 타입 선언 오류
- 대상: `specs/gigbridge/model/dto.go`
- 변경: `// @dto` 주석 삭제 (type 선언은 유지)
- 기대: WARNING — @dto 미부착 타입이 SSaC @result에서 사용되면 DDL 테이블 매칭 시도
- 결과: PASS — crosscheck/ssac_ddl에서 @dto 없는 타입의 DDL 매칭 실패 감지
