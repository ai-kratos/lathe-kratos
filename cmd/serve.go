package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"

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
		return runServe(servePort, serveNoOpen, openBrowser, http.ListenAndServe)
	},
}

func runServe(port int, noOpen bool, open func(string), listenAndServe func(string, http.Handler) error) error {
	dir, err := config.TutorialsDir()
	if err != nil {
		return err
	}
	srv := serve.NewServer(dir)
	url := fmt.Sprintf("http://localhost:%d", port)
	fmt.Printf("Serving tutorials at %s\n", url)
	if !noOpen {
		open(url)
	}
	// Bind all available interfaces/protocol families so tutorials can be
	// read from other devices on the local network or tailnet. Mutating
	// endpoints still reject foreign Origin/Referer headers in the serve
	// package.
	return listenAndServe(fmt.Sprintf(":%d", port), srv.Handler())
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
