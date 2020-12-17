package data

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

type App struct {
	db      *gorm.DB
	redisDb *redis.Client
	http    *http.Server
}


func NewApp(db *gorm.DB, r *redis.Client, h *http.Server) (app *App, cf func() , err error) {
	app = &App{
		db:      db,
		redisDb: r,
		http:    h,
	}
	cf = func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		if err = h.Shutdown(ctx); err != nil {
			// 记录日志
		}
		cancel()
	}
	return
}
