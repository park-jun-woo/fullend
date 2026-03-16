//ff:type feature=projectconfig type=model
//ff:what 인가 패키지 설정 구조체
package projectconfig

// AuthzConfig configures the authorization package.
type AuthzConfig struct {
	Package string `yaml:"package"` // custom authz package path, default: github.com/geul-org/fullend/pkg/authz
}
