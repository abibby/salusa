package main

import (
	"context"
	"net/http"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/dialects/sqlite"
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
	"github.com/jmoiron/sqlx"
)

type Foo struct {
	model.BaseModel
	ID   int    `db:"id,autoincrement"`
	Name string `db:"name"`
}

type ListRequest struct {
	Tx *sqlx.Tx `inject:"r"`
}

var list = request.Handler(func(r *ListRequest) ([]*Foo, error) {
	foos, err := builder.From[*Foo]().Get(r.Tx)
	if err != nil {
		return nil, err
	}
	return foos, nil
})

type AddRequest struct {
	Name   string          `query:"name"`
	Update database.Update `inject:""`
}

var add = request.Handler(func(r *AddRequest) (*Foo, error) {
	foo := &Foo{Name: r.Name}
	err := r.Update(func(tx *sqlx.Tx) error {
		return model.Save(tx, foo)
	})
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

	r.Get("/foo", list)
	r.Get("/foo/create", add)
	r.Get("/foo/{foo}", get)

	err = http.ListenAndServe(":8087", r)
	if err != nil {
		panic(err)
	}
}
