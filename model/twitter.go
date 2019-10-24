package model

import "time"

type TwitterImageModel struct {
	ID       string    `db:"twitter_id"`
	Twitter  string    `db:"twitter_icon_url"`
	Updateat time.Time `db:"update_at"`
}
