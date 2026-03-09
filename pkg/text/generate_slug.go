package text

import "github.com/gosimple/slug"

// @func generateSlug
// @description 텍스트를 URL-safe slug로 변환한다

type GenerateSlugInput struct {
	Text string
}

type GenerateSlugOutput struct {
	Slug string
}

func GenerateSlug(in GenerateSlugInput) (GenerateSlugOutput, error) {
	return GenerateSlugOutput{Slug: slug.Make(in.Text)}, nil
}
