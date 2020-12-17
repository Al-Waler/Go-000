package main

import (
	"Week04/internal/acg/data"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main()  {
	fmt.Println("start")
	c:=make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGQUIT)

	g:=new(errgroup.Group)
	g.Go(func() error {
		_,cf,err:=data.InitApp()
		if err!=nil {
			panic(err)
		}
		<-c
		cf()
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Printf("error:%T %v\n", errors.Cause(err), errors.Cause(err))
		log.Printf("stack: %+v", err)
	}
}
