package serve

import (
	"bytes"

	"github.com/yuin/goldmark"
)

func RenderMarkdown(src []byte) ([]byte, error) {
	md := goldmark.New()
	var buf bytes.Buffer
	if err := md.Convert(src, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
