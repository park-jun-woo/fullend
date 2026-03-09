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
	suffix := in.Suffix
	if suffix == "" {
		suffix = "..."
	}
	runes := []rune(in.Text)
	if len(runes) <= in.MaxLength {
		return TruncateTextOutput{Truncated: in.Text}, nil
	}
	return TruncateTextOutput{Truncated: string(runes[:in.MaxLength]) + suffix}, nil
}
