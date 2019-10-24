package repository

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/seven320/go_server/model"
	// "../model"
)

func GetTwitterImage(db *sqlx.DB, id string) (*model.TwitterImageModel, error) {
	t := model.TwitterImageModel{}
	if err := db.Get(&t,
		`SELECT twitter_id, twitter_icon_url, update_at FROM twitter_user WHERE twitter_id = ?`,
		id); err != nil {
		return nil, err
	}
	return &t, nil
}

// create なので間にトランザクションを挟む，そのために受け取るDBは.Tx
func CreateTwitterImage(db *sqlx.Tx, ti *model.TwitterImageModel) (sql.Result, error) {
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
