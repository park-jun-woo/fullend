//ff:func feature=gen-react type=generator control=iteration dimension=1
//ff:what package.json 파일을 생성한다

package react

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func writePackageJSON(dir string, stmlDeps map[string]string) error {
	deps := map[string]string{
		"react":            "^18",
		"react-dom":        "^18",
		"react-router-dom": "^6",
	}
	// Merge stml dependencies.
	for k, v := range stmlDeps {
		deps[k] = v
	}

	pkg := map[string]interface{}{
		"private": true,
		"type":    "module",
		"scripts": map[string]string{
			"dev":   "vite",
			"build": "tsc && vite build",
		},
		"dependencies": deps,
		"devDependencies": map[string]string{
			"@types/react":         "^18",
			"@types/react-dom":     "^18",
			"@vitejs/plugin-react": "^4",
			"typescript":           "^5",
			"vite":                 "^5",
			"tailwindcss":          "^3",
			"postcss":              "^8",
			"autoprefixer":         "^10",
		},
	}
	b, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "package.json"), b, 0644)
}
