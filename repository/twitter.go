package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"../model"
	"../twitter"
	"github.com/jmoiron/sqlx"
)

// package repository
func GetTwitterImage(db *sqlx.DB, id string) (*model.TwitterImageModel, error) {
	t := model.TwitterImageModel{}
	if err := db.Get(&t,
		`SELECT twitter_id, twitter_icon_url, update_at FROM twitter_user WHERE twitter_id = ?`,
		id); err != nil {
		return nil, err
	}
	return &t, nil
}

func FindTwitterImage(db *sqlx.DB, id string) (*model.TwitterImageModel, error) {
	t := model.TwitterImageModel{}
	if err := db.Get(&t,
		`SELECT twitter_id, twitter_icon_url, update_at FROM twitter_user WHERE twitter_id = ?`,
		id); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("検索")
			imgurl, err := twitter.GetUserImage(id)
			if err != nil {
				log.Printf("twitter error:%s", err)
				return nil, err
			}
			t.ID = id
			t.Twitter = imgurl
			CreateTwitterImage(db, &t)
			return &t, nil
		} else {
			log.Printf("%s", err)
		}
		return nil, err
	}
	elapsed := int(time.Since(t.Updateat).Hours())
	fmt.Printf("elapsed, %d", elapsed)
	if elapsed > 24 {
		UpdateTwitterImage(db, &t)
	}
	UpdateAccessCount(db, &t)
	return &t, nil
}

func CreateTwitterImage(db *sqlx.DB, ti *model.TwitterImageModel) (sql.Result, error) {
	stmt, err := db.Prepare("INSERT INTO twitter_user SET twitter_id = ?, twitter_icon_url = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	log.Printf("insert")
	return stmt.Exec(ti.ID, ti.Twitter)
}

func UpdateTwitterImage(db *sqlx.DB, ti *model.TwitterImageModel) (sql.Result, error) {
	stmt, err := db.Prepare("UPDATE twitter_user SET twitter_icon_url = ? WHERE twitter_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	log.Printf("update")
	return stmt.Exec(ti.Twitter, ti.ID)
}

func UpdateAccessCount(db *sqlx.DB, ti *model.TwitterImageModel) (sql.Result, error) {
	stmt, err := db.Prepare("UPDATE twitter_user SET access_cnt = access_cnt + 1 WHERE twitter_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	log.Printf("count up")
	return stmt.Exec(ti.ID)
}
