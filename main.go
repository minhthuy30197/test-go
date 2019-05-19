package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"unicode"

	"git.hocngay.com/test/model"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func main() {
	connect := model.CreateConnect("postgres", "123456", "postgres", "localhost:5432")
	setupDatabase(connect.DB)
	go hello()
	log.Println(model.Messages)
	go SendEmail()
	log.Println(model.Messages)
	time.Sleep(10 * time.Second)
}

func hello() {
	fmt.Println("Hello world goroutine")
}

func SendEmail() {
	log.Println("jojo")
	msg := <- model.Messages
	fmt.Println("------ ", msg)
}

func setupDatabase(db *pg.DB) {
	argsWithProg := os.Args
	if len(argsWithProg) > 1 && os.Args[1] == "release" {
	} else {
		LogQueryToConsole(db)
	}

	err := MigrationDb(db)
	if err != nil {
		panic(err)
	}
}

func MigrationDb(db *pg.DB) error {
	// Tạo schema
	var schemas = []string{"book"}
	for _, schema := range schemas {
		_, err := db.Exec("CREATE SCHEMA IF NOT EXISTS " + schema + ";")
		if err != nil {
			return err
		}
	}

	var book model.Book
	err := createTable(&book, "book", "books", db)
	if err != nil {
		return err
	}

	return nil
}

// TableIsExists kiểm tra bảng đã tồn tại
func TableIsExists(schema, tableName string, db *pg.DB) (bool, error) {
	var exist bool
	_, err := db.Query(&exist, `
		SELECT EXISTS (
			SELECT 1
			FROM   information_schema.tables 
			WHERE  table_schema = ?
			AND    table_name = ?
			)`, schema, tableName)
	if err != nil {
		return exist, err
	}
	return exist, err
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(q *pg.QueryEvent) {}

func (d dbLogger) AfterQuery(q *pg.QueryEvent) {
	fmt.Println(q.FormattedQuery())
}

// LogQueryToConsole Log câu lệnh query
func LogQueryToConsole(db *pg.DB) {
	db.AddQueryHook(dbLogger{})
}

func createTable(model interface{}, schema, tableName string, db *pg.DB) error {
	exist, err := TableIsExists(schema, tableName, db)
	if err != nil {
		return err
	}
	if !exist {
		err = db.CreateTable(model, &orm.CreateTableOptions{
			Temp:          false,
			FKConstraints: true,
			IfNotExists:   true,
		})

		if err != nil {
			return err
		}
	}

	return err
}

// ToSnake Change word to Snake Case
func ToSnake(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}
