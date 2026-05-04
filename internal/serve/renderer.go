package serve

import (
	"bytes"

	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark"
)

func RenderMarkdown(src []byte) ([]byte, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
			),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert(src, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
