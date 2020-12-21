package data

import (
	"Week04/internal/pkg/conf"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

func NewRedis() (redisDB *redis.Client, cf func(), err error) {
	err = conf.LoadConf("../../configs/redis.yaml")
	addr, err := conf.Get("addr")
	password, err := conf.Get("password")
	db, err := conf.Get("db")
	if err != nil {
		return nil, nil, errors.Wrap(err, "redis 配置文件读取错误")
	}

	address, err := conf.ToString(addr)
	pw, err := conf.ToString(password)
	dbCode, err := conf.ToInt(db)
	if err != nil {
		return nil, nil, errors.Wrap(err, "redis 配置读取错误")
	}
	redisDB = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: pw,     // no password set
		DB:       dbCode, // use default DB
	})

	cf= func(){
		redisDB.Close()
	}
	return
}
