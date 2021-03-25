package user

import (
	"time"

	"github.com/dchest/uniuri"
)

type Recover struct {
	ID int `pg:",pk"`

	UserID int  `pg:",notnull,unique:user_id"`
	User   User `pg:"fk:owner_id,rel:has-one"`

	Code string `pg:",notnull"`

	AddedAt    time.Time `pg:",notnull,default:now()"`
	ValidUntil time.Time `pg:",notnull"`
}

func NewRecover(user User) Recover {
	r := Recover{
		UserID:     user.ID,
		User:       user,
		Code:       uniuri.NewLen(32),
		ValidUntil: time.Now().Add(time.Hour * 24),
	}

	return r
}
