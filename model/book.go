package model

import (
	"errors"
	"log"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type Connect struct {
	DB *pg.DB
}

var Messages = make(chan string)

func CreateConnect(user string, password string, database string, address string) Connect {
	db := pg.Connect(&pg.Options{
		User:     user,
		Password: password,
		Database: database,
		Addr:     address,
	})

	connect := Connect{
		DB: db,
	}

	return connect
}

type Book struct {
	TableName []byte `json:"table_name" sql:"book.books"`
	Id        int32  `json:"id" sql:",pk"`
	Name      string `json:"name"`
	Author    string `json:"author"`
	Category  string `json:"category"`
}

func (book Book) InsertBook(DB orm.DB) error {
	log.Println("hihi")
	// Validate thông tin
	if len(book.Name) == 0 {
		return errors.New("Tiêu đề sách không được rỗng")
	}
	if len(book.Author) == 0 {
		return errors.New("Tác giả sách không được rỗng")
	}

	err := DB.Insert(&book)
	if err != nil {
		return errors.New("Không thể lưu sách vào cơ sở dữ liệu")
	}

	Messages <- "send mail"

	return nil
}
