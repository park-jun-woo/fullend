✅ 완료

# Phase050 — @subscribe codegen + @exists validator 수정

## 목표
BUG027(@subscribe message struct 미출력 + route 등록 오류)과 BUG028(@exists와 @empty 강제 규칙 충돌)을 수정하여 fullend의 11개 시퀀스 타입 전체가 실전에서 동작하도록 한다.

## 배경
zenflow-try06 더미 프로젝트에서 웹훅 알림(@subscribe)과 템플릿 중복 방지(@exists)를 구현하면서 발견.
두 버그 모두 우회(동기 호출, UNIQUE 제약)로 기능 완수는 했으나, @subscribe와 @exists가 정상 동작하지 않으면 fullend의 시퀀스 커버리지에 구멍이 남는다.

---

## BUG027: @subscribe codegen 3건

### 현상
1. `.ssac` 파일에 message struct 선언 + `@subscribe "topic"` + `func OnXxx(message StructName) {}`
2. `fullend gen` 실행
3. 생성된 `.go` 파일에 struct 정의 없이 `message WorkflowExecutedMessage` 참조 → `go build` 실패
4. glue-gen의 `write_central_routes.go`가 @subscribe 핸들러를 `r.Handle()` HTTP route로 등록 → `gin.HandlerFunc` 타입 불일치
5. crosscheck `check_ssac_openapi.go`가 @subscribe 함수에도 operationId 매칭을 강제 → ERROR

### 원인 분석

**message struct 미출력**: `internal/ssac/generator/build_subscribe_func_body.go:17`에서 `msgType`을 참조하지만, struct 정의를 출력하는 코드가 없음. `.ssac` 파일의 `type ... struct` 블록을 파싱은 하지만 생성 `.go` 파일에 출력하지 않음.

**HTTP route 등록**: `internal/gen/gogin/write_central_routes.go`에서 모든 ServiceFunc를 `r.Handle()`로 등록. `fn.Subscribe != nil` 체크가 없어서 @subscribe 핸들러도 HTTP handler로 등록됨. `internal/gen/gogin/collect_subscribers.go`는 존재하지만, `generate_main_with_domains.go`에서 route 등록 시 제외 로직이 미구현.

**operationId 강제**: `internal/crosscheck/check_ssac_openapi.go:25-33`에서 모든 `funcNames`를 순회하며 `st.Operations`에 없으면 ERROR. `fn.Subscribe != nil` 조건 체크 없음.

### 수정 방안

#### 1. ssac-gen: message struct 출력
- **파일**: `internal/ssac/generator/go_target_generate_subscribe_func.go`
- **변경**: `.ssac` 파일에서 파싱한 message struct를 생성 `.go` 파일 상단에 `type StructName struct { ... }` 출력
- 파싱은 `internal/ssac/parser/parse_func_decl.go`에서 이미 수행 중 → struct 정보를 generator에 전달

#### 2. glue-gen: @subscribe 핸들러를 HTTP route에서 제외
- **파일**: `internal/gen/gogin/write_central_routes.go`
- **변경**: route 등록 루프에서 `fn.Subscribe != nil`이면 skip
- **파일**: `internal/gen/gogin/generate_main_with_domains.go`
- **변경**: `collectSubscribers()` 결과를 main.go에 `queue.Subscribe("topic", handler)` 호출로 출력 (template은 `internal/gen/gogin/main_with_domains_template.go`에 추가)

#### 3. crosscheck: @subscribe 함수는 operationId 매칭 면제
- **파일**: `internal/crosscheck/check_ssac_openapi.go`
- **변경**: L21-23의 funcNames 수집 시 `fn.Subscribe != nil`이면 제외
```go
for _, fn := range funcs {
    if fn.Subscribe != nil {
        continue  // @subscribe는 HTTP endpoint가 아니므로 operationId 불필요
    }
    funcNames[fn.Name] = fn.FileName
}
```

---

## BUG028: @empty 강제 규칙과 @exists 충돌

### 현상
FK 참조 컬럼으로 `@get` 조회 후 `@exists` 가드를 사용하려 하면, validator가 "FK 참조 조회 후 @empty 가드 필요" ERROR를 발생.

### 원인 분석
`internal/ssac/validator/validate_fk_reference_guard.go:60`에서 `hasEmptyGuardFor()`만 체크. `@exists` 가드는 `@empty`의 반대 — "있으면 탈출(409)"이므로 통과 후 해당 변수는 nil 확정. nil dereference 위험 없음.

`has_empty_guard_for.go`는 `parser.SeqEmpty`만 체크하고 `parser.SeqExists`는 무시.

### 수정 방안
`@exists`도 해당 변수의 nil 안전성을 보장하는 가드로 인정.

- **파일**: `internal/ssac/validator/has_empty_guard_for.go`
- **변경**: `@exists` 가드도 체크에 포함
```go
func hasEmptyGuardFor(seqs []parser.Sequence, varName string) bool {
    for _, s := range seqs {
        if s.Type == parser.SeqEmpty && rootVar(s.Target) == varName {
            return true
        }
        // @exists는 "있으면 탈출" — 통과 후 변수는 nil 확정이므로 이후 필드 접근 없음
        if s.Type == parser.SeqExists && rootVar(s.Target) == varName {
            return true
        }
    }
    return false
}
```

**근거**: `@exists existing "msg" 409`는 existing이 not nil이면 409 반환(함수 종료). 통과했으면 existing은 nil이 확정. 이후 시퀀스에서 `existing.Field` 접근이 없으므로 nil dereference 위험 없음.

---

## 변경 파일 요약

| 파일 | 버그 | 변경 내용 |
|---|---|---|
| `internal/ssac/generator/go_target_generate_subscribe_func.go` | BUG027 | message struct 출력 |
| `internal/gen/gogin/write_central_routes.go` | BUG027 | @subscribe route 제외 |
| `internal/gen/gogin/generate_main_with_domains.go` | BUG027 | queue.Subscribe 호출 생성 |
| `internal/gen/gogin/main_with_domains_template.go` | BUG027 | subscribe 템플릿 추가 |
| `internal/crosscheck/check_ssac_openapi.go` | BUG027 | @subscribe operationId 면제 |
| `internal/ssac/validator/has_empty_guard_for.go` | BUG028 | @exists 가드 인정 |

## 검증 방법

### 단위 테스트
- `internal/ssac/generator/go_subscribe_test.go` — message struct 출력 테스트
- `internal/ssac/validator/validator_test.go` — @exists 후 @empty 면제 테스트
- `internal/crosscheck/check_ssac_openapi.go` 관련 테스트 — @subscribe 면제 테스트

### 더미 프로젝트 검증
zenflow-try06에서:
1. `service/webhook/on_workflow_executed.ssac` 복원 (@subscribe + message struct)
2. `service/template/publish_template.ssac`에 @exists 복원
3. OpenAPI에서 dummy `/internal/on-workflow-executed` 경로 제거
4. `fullend validate` → ERROR 0
5. `fullend gen` → `go build` 성공 (미사용 import 없음)
6. hurl 7개 테스트 전부 통과

## 의존성
없음 (Phase049 이후 독립 작업)
