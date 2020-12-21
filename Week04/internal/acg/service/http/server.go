package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Msg   string      `json:"msg"`
	Error string      `json:"error,omitempty"`
}

func New() (s *http.Server, err error) {
	r := gin.Default()
	initRoute(r)

	s = &http.Server{
		Addr:    ":3939",
		Handler: r,
	}
	go func() {
		if err = s.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				err = errors.Wrap(err, "http server error")
				panic(err)
			}
			panic(err)
		}
	}()

	return
}

func initRoute(r *gin.Engine)  {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, Response{
			Code: 2000,
			Msg:  "Hello World",
		})
	})
}
