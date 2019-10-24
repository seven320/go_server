package service

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"../model"
	"../repository"
	"../twitter"
	"github.com/jmoiron/sqlx"
)

type TwitterDB struct {
	db *sqlx.DB
}

func GetTwitterImage(db *sqlx.DB, id string) (*model.TwitterImageModel, error) {
	t, err := repository.GetTwitterImage(db, id)
	if err != nil && err == sql.ErrNoRows {
		log.Printf("検索")
		imgurl, err := twitter.GetUserImage(id)
		if err != nil {
			log.Printf("twitter error:%s", err)
			return nil, err
		}
		t.ID = id
		t.Twitter = imgurl
		repository.CreateTwitterImage(db, t)
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
