package data

import (
	"Week04/internal/pkg/conf"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

func NewDB() (db *gorm.DB, cf func(), err error) {
	err = conf.LoadConf("../../configs/db.yaml")
	if err != nil {
		return
	}

	mysqlDns, err := conf.Get("mysql_dns")
	maxOpen, err := conf.Get("maxopen")
	maxIdle, err := conf.Get("maxidle")
	if err != nil {
		return
	}

	d, err := conf.ToString(mysqlDns)
	mo, err := conf.ToInt(maxOpen)
	mi, err := conf.ToInt(maxIdle)
	if err != nil {
		return
	}

	db, err = gorm.Open("mysql", d)
	if err != nil {
		return
	}
	db.DB().SetMaxIdleConns(mi)
	db.DB().SetMaxOpenConns(mo)
	db.DB().SetConnMaxLifetime(time.Second * 30)
	cf = func() {
		db.Close()
	}
	return
}
