package main

import (
	"flag"
	"fmt"
	"github.com/mshayler/accountapi/internal/accountmanager"
	"github.com/mshayler/accountapi/internal/endpoints"
	"github.com/mshayler/accountapi/internal/service"
	httpTransport "github.com/mshayler/accountapi/internal/transport/http"
	"github.com/oklog/run"
	"github.com/op/go-logging"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	httpAddr := flag.String("http-addr", ":8081", "HTTP listen address")

	// Create Logger
	lgr, err := logging.GetLogger("Accounts API")
	if err != nil {
		panic("Failed to create logger")
	}

	// Create account client
	ac, err := accountmanager.NewManager(lgr)
	if err != nil {
		lgr.Fatal("Failed to create accounts client")
		panic("exiting...")
	}

	// Create service
	svc, err := service.New(lgr, ac)
	if err != nil {
		lgr.Fatal("Failed to create service")
		panic("exiting...")
	}

	// Make go-kit endpoints
	eps := endpoints.MakeEndpoints(svc, time.Minute)

	// Create HTTP handlers
	httpHandlers := createHTTPHandlers(eps, lgr)

	// Create HTTP listener
	httpListener := createHTTPListener(httpAddr, lgr)

	var g run.Group

	startHTTPServer(&g, lgr, httpAddr, httpHandlers, httpListener)

	// Start interrupt listener
	g = withInterrupt(g, lgr)

	lgr.Error("exit", g.Run())
}

// Create the HTTP server
func startHTTPServer(
	g *run.Group,
	lgr *logging.Logger,
	httpAddr *string,
	httpHandlers http.Handler,
	httpListener net.Listener,
) {
	g.Add(func() error {
		lgr.Info(fmt.Sprintf("http transport started on addr: %s", *httpAddr))

		httpServer := &http.Server{
			Addr:    *httpAddr,
			Handler: httpHandlers,
		}

		// Serve via HTTP
		return httpServer.Serve(httpListener)
	}, func(err error) {
		lgr.Info("http server stopping...")

		cerr := httpListener.Close()
		if cerr != nil {
			lgr.Error("http server error", cerr)
		}
	})
}

// Create HTTP handlers
func createHTTPHandlers(eps endpoints.Endpoints, lgr *logging.Logger) http.Handler {
	httpHandlers, err := httpTransport.NewHandlers(eps, nil, lgr)
	if err != nil {
		lgr.Fatal("failed to create http handlers", err)
	}
	return httpHandlers
}

// Create HTTP listener
func createHTTPListener(httpAddr *string, lgr *logging.Logger) net.Listener {
	httpListener, err := net.Listen("tcp", *httpAddr)
	if nil != err {
		lgr.Fatal("failed to create http listener", err)
	}
	return httpListener
}

func withInterrupt(g run.Group, lgr *logging.Logger) run.Group {
	interruptCh := make(chan struct{})
	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("os signal %s", sig)
		case <-interruptCh:
			return nil
		}
	}, func(err error) {
		lgr.Info("interrupt listener stopping...")
		close(interruptCh)
	})
	return g
}
