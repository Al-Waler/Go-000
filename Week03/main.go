package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var c chan os.Signal

func init() {
	c = make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGQUIT)
}

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		return serveApp(ctx)
	})
	g.Go(func() error {
		return serveDebug(ctx)
	})

	if err := g.Wait(); err != nil {
		log.Printf("error:%T %v\n", errors.Cause(err), errors.Cause(err))
		log.Printf("stack: %+v", err)
	}
}

func serveApp(ctx context.Context) error {
	defer func() {
		if r := recover(); r != nil {
			errors.Wrap(fmt.Errorf("v", r), "serve app panic")
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello world")
	})

	s := &http.Server{
		Addr:    ":8888",
		Handler: mux,
	}

	go func() {
		select {
		case <-c:
			s.Shutdown(context.Background())
		case <-ctx.Done():
			s.Shutdown(context.Background())
		}
	}()
	return errors.Wrap(s.ListenAndServe(), "app err")
}

func serveDebug(ctx context.Context) error {
	defer func() {
		if r := recover(); r != nil {
			errors.Wrap(fmt.Errorf("%v", r), "debug app panic")
		}
	}()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "debug")
	})

	s := &http.Server{
		Addr:    ":9999",
		Handler: mux,
	}

	go func() {
		select {
		case <-c:
			s.Shutdown(context.Background())
		case <-ctx.Done():
			s.Shutdown(context.Background())
		}
	}()
	return errors.Wrap(s.ListenAndServe(), "debug err")
}
