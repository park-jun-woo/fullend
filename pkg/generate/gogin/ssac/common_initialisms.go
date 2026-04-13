//ff:func feature=ssac-gen type=util control=sequence topic=string-convert
//ff:what Go 컨벤션 공통 이니셜리즘 맵 정의
package ssac

// commonInitialisms는 Go 컨벤션에서 대소문자를 통일하는 공통 이니셜리즘이다.
// https://github.com/golang/lint/blob/master/lint.go#L770
var commonInitialisms = map[string]bool{
	"ACL": true, "API": true, "ASCII": true, "CPU": true, "CSS": true,
	"DNS": true, "EOF": true, "HTML": true, "HTTP": true, "HTTPS": true,
	"ID": true, "IP": true, "JSON": true, "QPS": true, "RAM": true,
	"RPC": true, "SLA": true, "SMTP": true, "SQL": true, "SSH": true,
	"TCP": true, "TLS": true, "TTL": true, "UDP": true, "UI": true,
	"UID": true, "UUID": true, "URI": true, "URL": true, "XML": true,
}
