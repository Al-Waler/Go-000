package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"sync/atomic"
	"time"
)

func main() {
	ctx,cancel:=context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)
	listener, err := net.Listen("tcp", ":5588")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	var num int32
	for {
		atomic.AddInt32(&num, 1)
		fmt.Println("accept.... ", num)
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		select {
		case <-ctx.Done():
			fmt.Println("server done")
			return
		default:
			go func() {
				m := make(chan string, 1)
				g.Go(func() error {
					return RConn(ctx, conn, m, num)
				})
				g.Go(func() error {
					return WConn(ctx, conn, m, num)
				})
				if err := g.Wait(); err != nil {
					log.Printf("client num : %d error:%T %v\n", num, errors.Cause(err), errors.Cause(err))
					log.Printf("client num : %d stack: %+v", num, err)
				} else {
					log.Println("done....")
				}
			}()
		}
	}
}

func RConn(ctx context.Context, conn net.Conn, m chan<- string, n int32) error {
	defer conn.Close()
	r := bufio.NewScanner(conn)
	for r.Scan() {
		select {
		case <-ctx.Done():
			fmt.Printf("client read %d close\n",n)
			close(m)
			return nil
		default:
			m <- r.Text()
		}
	}
	if err := r.Err(); err != nil {
		return errors.Wrap(err, "read err")
	}
	return nil
}

func WConn(ctx context.Context, conn net.Conn, m chan string, n int32) error {
	defer conn.Close()
	w := bufio.NewWriter(conn)
	select {
	case <-ctx.Done():
		fmt.Printf("client read %d close\n",n)
		close(m)
		return nil
	default:
		for c := range m {
			var buf bytes.Buffer
			buf.WriteString("reply msg client send:")
			buf.WriteString(c)
			buf.WriteString("\n")
			_, err := w.WriteString(buf.String())
			err = w.Flush()
			if err != nil {
				return errors.Wrap(err, "write err")
			}
		}
	}

	fmt.Printf("client %d write close\n",n)
	return nil
}
