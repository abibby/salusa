package main

import (
	"context"
	"net/http"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/dialects/sqlite"
	"github.com/abibby/salusa/database/insert"
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/database/models"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

type Foo struct {
	models.BaseModel
	ID   int    `db:"id,autoincrement"`
	Name string `db:"name"`
}

type ListRequest struct {
	Request *http.Request
}

var list = request.Handler(func(r *ListRequest) ([]*Foo, error) {
	tx := request.UseTx(r.Request)

	foos, err := builder.From[*Foo]().Get(tx)
	if err != nil {
		return nil, err
	}
	return foos, nil
})

type AddRequest struct {
	Request *http.Request
	Name    string `query:"name"`
}

var add = request.Handler(func(r *AddRequest) (*Foo, error) {
	tx := request.UseTx(r.Request)
	foo := &Foo{Name: r.Name}
	err := insert.Save(tx, foo)
	if err != nil {
		return nil, err
	}
	return foo, nil
})

type GetRequest struct {
	Foo *Foo `di:"foo"`
}

var get = request.Handler(func(r *GetRequest) (*Foo, error) {
	return r.Foo, nil
})

func main() {
	sqlite.UseSQLite()
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	r := router.New()

	c, err := migrate.CreateFromModel(&Foo{})
	if err != nil {
		panic(err)
	}

	err = c.Run(context.Background(), db)
	if err != nil {
		panic(err)
	}

	r.Use(request.WithDB(db))

	r.Get("/foo", list)
	r.Get("/foo/create", add)
	r.Get("/foo/{foo}", get)

	http.ListenAndServe(":8087", r)
}
