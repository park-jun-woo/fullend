//ff:type feature=projectconfig type=model
//ff:what 프론트엔드 설정 구조체
package manifest

type Frontend struct {
	Lang      string `yaml:"lang"`
	Framework string `yaml:"framework"`
	Bundler   string `yaml:"bundler"`
	Name      string `yaml:"name"`
}
