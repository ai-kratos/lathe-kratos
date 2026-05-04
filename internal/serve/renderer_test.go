package serve_test

import (
	"strings"
	"testing"

	"github.com/devenjarvis/lathe/internal/serve"
)

func TestRenderMarkdown(t *testing.T) {
	src := []byte("# Hello World\n\nThis is a `test`.\n\n```go\nfmt.Println(\"hello\")\n```")

	out, err := serve.RenderMarkdown(src)
	if err != nil {
		t.Fatalf("RenderMarkdown() error = %v", err)
	}

	html := string(out)
	if !strings.Contains(html, "<h1>Hello World</h1>") {
		t.Errorf("RenderMarkdown() missing <h1>, got:\n%s", html)
	}
	if !strings.Contains(html, "<code>test</code>") {
		t.Errorf("RenderMarkdown() missing inline <code>, got:\n%s", html)
	}
	if !strings.Contains(html, "<pre") {
		t.Errorf("RenderMarkdown() code block not rendered as <pre>, got:\n%s", html)
	}
	if !strings.Contains(html, "Println") {
		t.Errorf("RenderMarkdown() code block content missing from output, got:\n%s", html)
	}
}
