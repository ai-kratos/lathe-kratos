package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/devenjarvis/lathe/internal/config"
	"github.com/devenjarvis/lathe/internal/serve"
	"github.com/spf13/cobra"
)

var servePort int
var serveNoOpen bool

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the tutorial web server and open the browser",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runServe(cmd.Context(), servePort, serveNoOpen, openBrowser, listenAndServeGraceful)
	},
}

func runServe(ctx context.Context, port int, noOpen bool, open func(string), listenAndServe func(context.Context, string, http.Handler) error) error {
	dir, err := config.TutorialsDir()
	if err != nil {
		return err
	}
	srv := serve.NewServer(dir)
	url := fmt.Sprintf("http://localhost:%d", port)

	// Record the running server so the worker CLI (`lathe work ...`) can find
	// its URL, and clean it up on shutdown. Best-effort: a failed write only
	// means the worker can't auto-discover the server, not that serving fails.
	rt := &config.ServeRuntime{URL: url, PID: os.Getpid(), Started: time.Now().UTC()}
	if werr := config.WriteServeRuntime(rt); werr != nil {
		_, _ = fmt.Fprintf(os.Stderr, "warning: could not write serve runtime file: %v\n", werr)
	}
	defer func() { _ = config.RemoveServeRuntime() }()

	fmt.Printf("Serving tutorials at %s\n", url)
	// Nudge toward live mode without spawning anything: starting the loop is
	// the user's call (it can't be agent-agnostic or non-metered otherwise —
	// see the worker-bridge note in AGENTS.md).
	fmt.Println("Live mode: run /lathe-work in your coding agent to drive Ask/Verify/Extend here (otherwise the buttons hand you a command to paste).")
	if !noOpen {
		open(url)
	}

	// Bind all available interfaces/protocol families so tutorials can be
	// read from other devices on the local network or tailnet. Mutating
	// endpoints still reject foreign Origin/Referer headers in the serve
	// package.
	return listenAndServe(ctx, fmt.Sprintf(":%d", port), srv.Handler())
}

// listenAndServeGraceful serves until ctx is cancelled (Ctrl-C / SIGTERM),
// then shuts down gracefully so deferred cleanup (the serve runtime file) runs
// — a plain ListenAndServe never returns on a signal.
func listenAndServeGraceful(ctx context.Context, addr string, handler http.Handler) error {
	httpSrv := &http.Server{Addr: addr, Handler: handler}

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	errCh := make(chan error, 1)
	go func() {
		err := httpSrv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		errCh <- err
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return httpSrv.Shutdown(shutdownCtx)
	}
}

func openBrowser(url string) {
	var bin string
	switch runtime.GOOS {
	case "darwin":
		bin = "open"
	case "linux":
		bin = "xdg-open"
	default:
		return
	}
	if err := exec.Command(bin, url).Start(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "warning: could not open browser: %v\n", err)
	}
}

func init() {
	serveCmd.Flags().IntVar(&servePort, "port", 4242, "port to listen on")
	serveCmd.Flags().BoolVar(&serveNoOpen, "no-open", false, "do not open a browser")
	rootCmd.AddCommand(serveCmd)
}
