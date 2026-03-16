//ff:type feature=projectconfig type=model
//ff:what 백엔드 설정 구조체
package projectconfig

type Backend struct {
	Lang       string   `yaml:"lang"`
	Framework  string   `yaml:"framework"`
	Module     string   `yaml:"module"`
	Middleware []string `yaml:"middleware"`
	Auth       *Auth    `yaml:"auth"`
}
