package service

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	// "github.com/seven320/go_server/dbutil"
	// "github.com/seven320/go_server/model"
	// "github.com/seven320/go_server/repository"
	// "github.com/seven320/go_server/twitter"

	"../dbutil"
	"../model"
	"../repository"
	"../twitter"
)

type TwitterImage struct {
	db *sqlx.DB
}

func NewTwitterImage(db *sqlx.DB) *TwitterImage {
	return &TwitterImage{db}
}

func (ti *TwitterImage) GetTwitterImage(db *sqlx.DB, id string) (*model.TwitterImageModel, error) {
	t, err := repository.GetTwitterImage(ti.db, id)
	if err != nil && err == sql.ErrNoRows {
		log.Printf("検索")
		imgurl, err := twitter.GetUserImage(id)
		if err != nil {
			log.Printf("twitter error:%s", err)
			return nil, err
		}
		t.Twitter = imgurl
		if err := dbutil.TXHandler(ti.db, func(tx *sqlx.Tx) error {
			_, err := repository.CreateTwitterImage(tx, t)
			if err != nil {
				return err
			}
			if err := tx.Commit(); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return nil, errors.Wrap(err, "failed twitter image insert transaction")
		}
		return t, nil
	} else if err != nil {
		log.Printf("%s", err)
		return nil, err
	}
	elapsed := int(time.Since(t.Updateat).Hours())
	fmt.Printf("elapsed, %d", elapsed)
	if elapsed > 24 {
		imgurl, err := twitter.GetUserImage(id)
		if err != nil {
			log.Printf("twitter error:%s", err)
			return nil, err
		}
		if t.Twitter != imgurl { //imagedate更新
			t.Twitter = imgurl
			_, err := repository.UpdateTwitterImage(db, t)
			if err != nil {
				log.Printf("update error :%s", err)
				return nil, err
			}
		}
	}
	_, err = repository.UpdateAccessCount(db, t)
	if err != nil {
		log.Printf("update error: %s", err)
		return nil, err
	}
	return t, nil
}
