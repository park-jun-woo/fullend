package image

import (
	"bytes"

	"github.com/disintegration/imaging"
)

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
	size := in.Size
	if size <= 0 {
		size = 200
	}
	src, err := imaging.Decode(bytes.NewReader(in.Data))
	if err != nil {
		return GenerateThumbnailOutput{}, err
	}
	thumb := imaging.Fill(src, size, size, imaging.Center, imaging.Lanczos)
	var buf bytes.Buffer
	if err := imaging.Encode(&buf, thumb, imaging.JPEG); err != nil {
		return GenerateThumbnailOutput{}, err
	}
	return GenerateThumbnailOutput{Data: buf.Bytes()}, nil
}
