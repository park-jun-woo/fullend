✅ 완료

# Phase 025: pkg/ 기본 제공 함수 확장 — 6개 패키지 16개 함수

## 목표

fullend가 기본 제공하는 패키지 함수를 확장한다.
프로젝트마다 반복 구현하는 공통 기능을 `pkg/`에 래핑하여 `@func pkg.name` 한 줄로 사용 가능하게 한다.
재발명 없이 검증된 Go 라이브러리를 `func(In) (Out, error)` 계약으로 감싼다.

---

## 함수 목록

### pkg/auth (기존 3 + 신규 3 = 6)

| 함수 | 설명 | 래핑 대상 |
|---|---|---|
| ~~hashPassword~~ | ~~bcrypt 해싱~~ | ~~이미 구현~~ |
| ~~verifyPassword~~ | ~~bcrypt 검증~~ | ~~이미 구현~~ |
| ~~issueToken~~ | ~~JWT 발급~~ | ~~이미 구현~~ |
| `verifyToken` | JWT 검증 → claims 추출 | `golang-jwt/jwt/v5` |
| `refreshToken` | 리프레시 토큰 발급 (7일 만료) | `golang-jwt/jwt/v5` |
| `generateResetToken` | 비밀번호 리셋용 랜덤 토큰 (32바이트 hex) | `crypto/rand` + `encoding/hex` |

### pkg/crypto (신규 4)

| 함수 | 설명 | 래핑 대상 |
|---|---|---|
| `encrypt` | AES-256-GCM 대칭 암호화 | `crypto/aes` + `crypto/cipher` |
| `decrypt` | AES-256-GCM 복호화 | `crypto/aes` + `crypto/cipher` |
| `generateOTP` | TOTP 시크릿 + QR URL 생성 | `github.com/pquerna/otp/totp` |
| `verifyOTP` | TOTP 코드 검증 | `github.com/pquerna/otp/totp` |

### pkg/storage (신규 3)

| 함수 | 설명 | 래핑 대상 |
|---|---|---|
| `uploadFile` | S3 호환 파일 업로드 | `aws-sdk-go-v2/service/s3` |
| `deleteFile` | S3 호환 파일 삭제 | `aws-sdk-go-v2/service/s3` |
| `presignURL` | 서명된 다운로드 URL 생성 | `aws-sdk-go-v2/s3/presign` |

### pkg/mail (신규 2)

| 함수 | 설명 | 래핑 대상 |
|---|---|---|
| `sendEmail` | SMTP 이메일 발송 | `net/smtp` (표준 라이브러리) |
| `sendTemplateEmail` | 템플릿 기반 이메일 (Go template) | `html/template` + `net/smtp` |

### pkg/text (신규 3)

| 함수 | 설명 | 래핑 대상 |
|---|---|---|
| `generateSlug` | 유니코드 → URL-safe slug | `github.com/gosimple/slug` |
| `sanitizeHTML` | XSS 방지 HTML 정제 | `github.com/microcosm-cc/bluemonday` |
| `truncateText` | 유니코드 안전 텍스트 자르기 (rune 기반) | 표준 라이브러리 |

### pkg/image (신규 2)

| 함수 | 설명 | 래핑 대상 |
|---|---|---|
| `resizeImage` | 이미지 리사이즈 + 포맷 변환 | `github.com/disintegration/imaging` |
| `generateThumbnail` | 고정 크기 썸네일 생성 | `github.com/disintegration/imaging` |

---

## 변경 파일 목록

### 1. pkg/ 신규 파일 (16개)

```
pkg/
├── auth/
│   ├── hash_password.go          # 기존
│   ├── verify_password.go        # 기존
│   ├── issue_token.go            # 기존
│   ├── verify_token.go           # 신규
│   ├── refresh_token.go          # 신규
│   └── generate_reset_token.go   # 신규
├── crypto/
│   ├── encrypt.go                # 신규
│   ├── decrypt.go                # 신규
│   ├── generate_otp.go           # 신규
│   └── verify_otp.go             # 신규
├── storage/
│   ├── upload_file.go            # 신규
│   ├── delete_file.go            # 신규
│   └── presign_url.go            # 신규
├── mail/
│   ├── send_email.go             # 신규
│   └── send_template_email.go    # 신규
├── text/
│   ├── generate_slug.go          # 신규
│   ├── sanitize_html.go          # 신규
│   └── truncate_text.go          # 신규
└── image/
    ├── resize_image.go           # 신규
    └── generate_thumbnail.go     # 신규
```

### 2. 문서 수정

| 파일 | 변경 |
|---|---|
| `artifacts/manual-for-ai.md` | 기본 제공 함수 목록 갱신 (3 → 19) |
| `README.md` | Default Functions 테이블 갱신 |

### 3. go.mod

신규 의존성 추가.

---

## 상세 설계

### auth.verifyToken

```go
package auth

// @func verifyToken
// @description JWT 토큰을 검증하고 claims를 추출한다

type VerifyTokenInput struct {
    Token  string
    Secret string
}

type VerifyTokenOutput struct {
    UserID int64
    Email  string
    Role   string
}

func VerifyToken(in VerifyTokenInput) (VerifyTokenOutput, error) {
    // jwt.Parse → MapClaims 추출
}
```

### auth.refreshToken

```go
package auth

// @func refreshToken
// @description 리프레시 토큰을 발급한다 (7일 만료)

type RefreshTokenInput struct {
    UserID int64
    Email  string
    Role   string
}

type RefreshTokenOutput struct {
    RefreshToken string
}

func RefreshToken(in RefreshTokenInput) (RefreshTokenOutput, error) {
    // issueToken과 동일하되 만료 7일
}
```

### auth.generateResetToken

```go
package auth

// @func generateResetToken
// @description 비밀번호 리셋용 32바이트 랜덤 hex 토큰을 생성한다

type GenerateResetTokenInput struct{}

type GenerateResetTokenOutput struct {
    Token string
}

func GenerateResetToken(in GenerateResetTokenInput) (GenerateResetTokenOutput, error) {
    // crypto/rand.Read(32) → hex.EncodeToString
}
```

### crypto.encrypt

```go
package crypto

// @func encrypt
// @description 평문을 AES-256-GCM으로 암호화한다

type EncryptInput struct {
    Plaintext string
    Key       string // 32바이트 hex
}

type EncryptOutput struct {
    Ciphertext string // base64 인코딩
}

func Encrypt(in EncryptInput) (EncryptOutput, error) {
    // hex.DecodeString(key) → aes.NewCipher → cipher.NewGCM → Seal
}
```

### crypto.decrypt

```go
package crypto

// @func decrypt
// @description AES-256-GCM 암호문을 복호화한다

type DecryptInput struct {
    Ciphertext string // base64 인코딩
    Key        string // 32바이트 hex
}

type DecryptOutput struct {
    Plaintext string
}

func Decrypt(in DecryptInput) (DecryptOutput, error) {
    // base64.Decode → aes.NewCipher → cipher.NewGCM → Open
}
```

### crypto.generateOTP

```go
package crypto

// @func generateOTP
// @description TOTP 시크릿과 QR 프로비저닝 URL을 생성한다

type GenerateOTPInput struct {
    Issuer      string
    AccountName string
}

type GenerateOTPOutput struct {
    Secret string
    URL    string // otpauth:// URL (QR 코드용)
}

func GenerateOTP(in GenerateOTPInput) (GenerateOTPOutput, error) {
    // totp.Generate → key.Secret(), key.URL()
}
```

### crypto.verifyOTP

```go
package crypto

// @func verifyOTP
// @description TOTP 코드가 시크릿과 일치하는지 검증한다

type VerifyOTPInput struct {
    Code   string
    Secret string
}

type VerifyOTPOutput struct{}

func VerifyOTP(in VerifyOTPInput) (VerifyOTPOutput, error) {
    // totp.Validate(code, secret) → false면 error 반환
}
```

### storage.uploadFile

```go
package storage

// @func uploadFile
// @description S3 호환 스토리지에 파일을 업로드한다

type UploadFileInput struct {
    Bucket      string
    Key         string
    Data        []byte
    ContentType string
    Endpoint    string // MinIO 등 커스텀 엔드포인트 (빈 문자열이면 AWS 기본)
    Region      string
}

type UploadFileOutput struct {
    URL string
}

func UploadFile(in UploadFileInput) (UploadFileOutput, error) {
    // s3.PutObject
}
```

### storage.deleteFile

```go
package storage

// @func deleteFile
// @description S3 호환 스토리지에서 파일을 삭제한다

type DeleteFileInput struct {
    Bucket   string
    Key      string
    Endpoint string
    Region   string
}

type DeleteFileOutput struct{}

func DeleteFile(in DeleteFileInput) (DeleteFileOutput, error) {
    // s3.DeleteObject
}
```

### storage.presignURL

```go
package storage

// @func presignURL
// @description 서명된 다운로드 URL을 생성한다

type PresignURLInput struct {
    Bucket    string
    Key       string
    ExpiresIn int // 초 단위 (기본 3600)
    Endpoint  string
    Region    string
}

type PresignURLOutput struct {
    URL string
}

func PresignURL(in PresignURLInput) (PresignURLOutput, error) {
    // s3.PresignClient → PresignGetObject
}
```

### mail.sendEmail

```go
package mail

// @func sendEmail
// @description SMTP를 통해 이메일을 발송한다

type SendEmailInput struct {
    Host     string // SMTP 호스트 (예: smtp.gmail.com)
    Port     int    // SMTP 포트 (예: 587)
    Username string
    Password string
    From     string
    To       string
    Subject  string
    Body     string // plain text
}

type SendEmailOutput struct{}

func SendEmail(in SendEmailInput) (SendEmailOutput, error) {
    // smtp.PlainAuth → smtp.SendMail
}
```

### mail.sendTemplateEmail

```go
package mail

// @func sendTemplateEmail
// @description Go 템플릿으로 HTML 이메일을 발송한다

type SendTemplateEmailInput struct {
    Host         string
    Port         int
    Username     string
    Password     string
    From         string
    To           string
    Subject      string
    TemplateName string            // 템플릿 파일 경로 또는 인라인 템플릿
    Data         map[string]string // 템플릿 변수
}

type SendTemplateEmailOutput struct{}

func SendTemplateEmail(in SendTemplateEmailInput) (SendTemplateEmailOutput, error) {
    // template.Execute → MIME multipart → smtp.SendMail
}
```

### text.generateSlug

```go
package text

// @func generateSlug
// @description 텍스트를 URL-safe slug로 변환한다

type GenerateSlugInput struct {
    Text string
}

type GenerateSlugOutput struct {
    Slug string
}

func GenerateSlug(in GenerateSlugInput) (GenerateSlugOutput, error) {
    // slug.Make(text)
}
```

### text.sanitizeHTML

```go
package text

// @func sanitizeHTML
// @description HTML에서 위험한 태그와 속성을 제거한다 (XSS 방지)

type SanitizeHTMLInput struct {
    HTML string
}

type SanitizeHTMLOutput struct {
    Sanitized string
}

func SanitizeHTML(in SanitizeHTMLInput) (SanitizeHTMLOutput, error) {
    // bluemonday.UGCPolicy().Sanitize(html)
}
```

### text.truncateText

```go
package text

// @func truncateText
// @description 유니코드 안전하게 텍스트를 자른다

type TruncateTextInput struct {
    Text      string
    MaxLength int
    Suffix    string // 말줄임 (기본 "...")
}

type TruncateTextOutput struct {
    Truncated string
}

func TruncateText(in TruncateTextInput) (TruncateTextOutput, error) {
    // []rune 변환 → 길이 비교 → suffix 추가
}
```

### image.resizeImage

```go
package image

// @func resizeImage
// @description 이미지를 지정 크기로 리사이즈한다

type ResizeImageInput struct {
    Data   []byte
    Width  int
    Height int    // 0이면 비율 유지
    Format string // "jpeg", "png", "webp" (빈 문자열이면 원본 포맷 유지)
}

type ResizeImageOutput struct {
    Data []byte
}

func ResizeImage(in ResizeImageInput) (ResizeImageOutput, error) {
    // imaging.Decode → imaging.Resize → imaging.Encode
}
```

### image.generateThumbnail

```go
package image

// @func generateThumbnail
// @description 이미지를 정사각형 크기로 크롭하여 썸네일을 생성한다

type GenerateThumbnailInput struct {
    Data []byte
    Size int // 한 변의 크기 (기본 200)
}

type GenerateThumbnailOutput struct {
    Data []byte
}

func GenerateThumbnail(in GenerateThumbnailInput) (GenerateThumbnailOutput, error) {
    // imaging.Decode → imaging.Fill(size, size, Center, Lanczos) → Encode(JPEG)
}
```

---

## 의존성

### 신규 외부 의존성

| 패키지 | 용도 | 라이선스 |
|---|---|---|
| `github.com/pquerna/otp` | TOTP 생성/검증 | Apache 2.0 |
| `github.com/aws/aws-sdk-go-v2/service/s3` | S3 파일 업로드/삭제 | Apache 2.0 |
| `github.com/aws/aws-sdk-go-v2/config` | AWS 설정 로드 | Apache 2.0 |
| `github.com/aws/aws-sdk-go-v2/credentials` | AWS 자격증명 | Apache 2.0 |
| `github.com/gosimple/slug` | URL slug 생성 | MPL 2.0 |
| `github.com/microcosm-cc/bluemonday` | HTML 정제 | BSD 3-Clause |
| `github.com/disintegration/imaging` | 이미지 리사이즈 | MIT |

### 표준 라이브러리만 사용

| 패키지 | 용도 |
|---|---|
| `crypto/rand`, `encoding/hex` | 랜덤 토큰 생성 |
| `crypto/aes`, `crypto/cipher` | AES-256-GCM |
| `encoding/base64` | 암호문 인코딩 |
| `net/smtp` | 이메일 발송 |
| `html/template` | 이메일 템플릿 |

---

## 검증 방법

```bash
# 1. fullend 테스트 (기존 + 신규 파서 테스트 통과)
cd ~/.clari/repos/fullend && go test ./...

# 2. validate — fullend pkg/ 기본 함수 19개 인식
go run ./artifacts/cmd/fullend/main.go validate specs/dummy-lesson

# 3. 개별 패키지 빌드 확인
go build ./pkg/auth/...
go build ./pkg/crypto/...
go build ./pkg/storage/...
go build ./pkg/mail/...
go build ./pkg/text/...
go build ./pkg/image/...
```

---

## SSaC 사용 예시

```go
// 비밀번호 리셋 플로우
// @sequence call
// @func auth.generateResetToken
// @result resetToken string
//
// @sequence call
// @func mail.sendTemplateEmail
// @param user.Email
// @param resetToken

// 파일 업로드 플로우
// @sequence call
// @func image.resizeImage
// @param FileData request
// @result resized bytes
//
// @sequence call
// @func storage.uploadFile
// @param resized
// @result fileURL string

// 게시글 작성
// @sequence call
// @func text.sanitizeHTML
// @param Content request
// @result sanitized string
//
// @sequence call
// @func text.generateSlug
// @param Title request
// @result slug string
```
