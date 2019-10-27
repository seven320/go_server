package controller

// queryの解析などを行う

import (
	"database/sql"
	"net/http"
	"net/url"

	"github.com/jmoiron/sqlx"

	// "github.com/seven320/go_server/service"
	"../service"
)

type TwitterDB struct {
	db *sqlx.DB
}

func NewTwitterImage(db *sqlx.DB) *TwitterDB {
	return &TwitterDB{db: db}
}

func (t *TwitterDB) Show(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	// vars := mux.Vars(r)

	u, _ := url.Parse(r.URL.String())
	query := u.Query()
	id := query.Get("id")
	twitterService := service.NewTwitterImage(t.db)
	twitterimage, err := twitterService.GetTwitterImage(t.db, id)
	if err != nil && err == sql.ErrNoRows {
		return http.StatusBadRequest, nil, err
	} else if err != nil {
		// return http.StatusBadRequest, nil, err
		return http.StatusBadRequest, nil, err
	}

	return http.StatusCreated, twitterimage, nil
}
