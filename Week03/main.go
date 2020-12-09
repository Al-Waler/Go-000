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
var stop chan struct{}
func init()  {
	stop = make(chan struct{})
	c = make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt,syscall.SIGQUIT)
}

func main() {
	g:=new(errgroup.Group)
	g.Go(serveApp)
	g.Go(serveDebug)

	if err := g.Wait(); err != nil {
		log.Printf("error:%T %v\n", errors.Cause(err), errors.Cause(err))
		log.Printf("stack: %+v", err)
	}
}


func serveApp() error {
	defer func() {
		if r := recover(); r != nil {
			errors.Wrap(fmt.Errorf("v",r), "serve app panic")
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
			close(stop)
		case <-stop:
			s.Shutdown(context.Background())
		}
	}()
	return errors.Wrap(s.ListenAndServe(), "app err")
}

func serveDebug() error {
	defer func() {
		if r := recover(); r != nil {
			errors.Wrap(fmt.Errorf("%v",r), "debug app panic")
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
			close(stop)
		case <-stop:
			s.Shutdown(context.Background())
		}
	}()
	return errors.Wrap(s.ListenAndServe(), "debug err")
}

