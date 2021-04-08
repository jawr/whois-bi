package user

import (
	"time"

	"github.com/dchest/uniuri"
	"github.com/go-pg/pg/v10/orm"
	"github.com/pkg/errors"
)

type Recover struct {
	ID int `pg:",pk"`

	UserID int  `pg:",notnull,unique:user_id"`
	User   User `pg:"fk:user_id,rel:has-one"`

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

func (r *Recover) Insert(db orm.DB) error {
	_, err := db.Model(r).Returning("*").Insert()
	return err
}

func RecoverPassword(db orm.DB, code string, newPassword []byte) error {
	var rec Recover

	if err := db.Model(&rec).Where("code = ?", code).Select(); err != nil {
		return errors.New("invalid code")
	}

	if time.Now().After(rec.ValidUntil) {
		// delete expired recover
		db.Model(&rec).Delete()
		return errors.New("expired code")
	}

	_, err := db.Model((*User)(nil)).
		Set("password = ?", newPassword).
		Where("id = ?", rec.UserID).
		Update()
	if err != nil {
		return errors.WithMessage(err, "Update")
	}

	db.Model(&rec).Delete()

	return nil
}
