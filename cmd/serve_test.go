package cmd

import (
	"context"
	"net/http"
	"testing"
)

func TestRunServeOpensBrowserByDefault(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	var opened []string
	var listenedAddr string
	err := runServe(context.Background(), 4242, false, func(url string) {
		opened = append(opened, url)
	}, func(_ context.Context, addr string, handler http.Handler) error {
		listenedAddr = addr
		if handler == nil {
			t.Fatal("listen handler = nil")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("runServe: %v", err)
	}

	if len(opened) != 1 || opened[0] != "http://localhost:4242" {
		t.Fatalf("opened = %v, want [http://localhost:4242]", opened)
	}
	if listenedAddr != ":4242" {
		t.Fatalf("listen addr = %q, want :4242", listenedAddr)
	}
}

func TestRunServeNoOpenSkipsBrowser(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	var opened []string
	var listenedAddr string
	err := runServe(context.Background(), 4243, true, func(url string) {
		opened = append(opened, url)
	}, func(_ context.Context, addr string, handler http.Handler) error {
		listenedAddr = addr
		if handler == nil {
			t.Fatal("listen handler = nil")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("runServe: %v", err)
	}

	if len(opened) != 0 {
		t.Fatalf("opened = %v, want none", opened)
	}
	if listenedAddr != ":4243" {
		t.Fatalf("listen addr = %q, want :4243", listenedAddr)
	}
}
