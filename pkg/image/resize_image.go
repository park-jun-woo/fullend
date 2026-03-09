package image

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/disintegration/imaging"
)

// @func resizeImage
// @description 이미지를 지정 크기로 리사이즈한다

type ResizeImageInput struct {
	Data   []byte
	Width  int
	Height int    // 0이면 비율 유지
	Format string // "jpeg", "png" (빈 문자열이면 원본 포맷 유지)
}

type ResizeImageOutput struct {
	Data []byte
}

func ResizeImage(in ResizeImageInput) (ResizeImageOutput, error) {
	src, err := imaging.Decode(bytes.NewReader(in.Data))
	if err != nil {
		return ResizeImageOutput{}, err
	}
	resized := imaging.Resize(src, in.Width, in.Height, imaging.Lanczos)
	format, err := resolveFormat(in.Format)
	if err != nil {
		return ResizeImageOutput{}, err
	}
	var buf bytes.Buffer
	if err := imaging.Encode(&buf, resized, format); err != nil {
		return ResizeImageOutput{}, err
	}
	return ResizeImageOutput{Data: buf.Bytes()}, nil
}

func resolveFormat(format string) (imaging.Format, error) {
	switch strings.ToLower(format) {
	case "jpeg", "jpg", "":
		return imaging.JPEG, nil
	case "png":
		return imaging.PNG, nil
	case "gif":
		return imaging.GIF, nil
	case "bmp":
		return imaging.BMP, nil
	case "tiff":
		return imaging.TIFF, nil
	default:
		return 0, fmt.Errorf("unsupported image format: %s", format)
	}
}
