package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"log"
	"os"
)

// you should only handle errors once
func main() {
	user, err := Biz(12)
	if err != nil {
		log.Printf("%T %v\n", errors.Cause(err), errors.Cause(err))
		log.Printf("stack trace: \n%+v\n", err)
		os.Exit(1)
	}
	log.Println(user)
}

func Biz(id uint) (user User, err error) {
	return Dao(id)
}

type User struct {
	Id             uint
	CreatedAt      string
	UpdatedAt      string
	DeletedAt      string
	Username       string
	PasswordDigest string
	Nickname       string
	Status         string
	Avatar         string
}

func Dao(id uint) (user User, err error) {
	db, err := sql.Open("mysql", "root:uUiknmbGFDBIu9801827654@tcp(127.0.0.1:3307)/acg?parseTime=true")
	if err != nil {
		return user, errors.New("数据库链接失败")
	}
	sqlstr := "SELECT id, username, nickname, avatar FROM users WHERE id=?"
	err = db.QueryRow(sqlstr, id).Scan(&user.Id, &user.Username, &user.Nickname, &user.Avatar)
	defer db.Close()
	// 直接向上抛
	if errors.Is(err, sql.ErrNoRows) {
		return
	}
	return user,errors.Wrap(err, "出错了")
}
