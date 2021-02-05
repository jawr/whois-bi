package user

import (
	"time"

	"github.com/dchest/uniuri"
	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID int `sql:",pk"`

	Email    string `sql:",unique,notnull"`
	Password []byte `sql:",notnull"`

	// meta data
	AddedAt   time.Time `sql:",notnull,default:now()"`
	DeletedAt time.Time `pg:",soft_delete"`

	VerifiedAt   time.Time
	VerifiedCode string `sql:",notnull"`

	LastLoginAt time.Time
}

// create a new user and bcrypt the password
func NewUser(email, password string) (User, error) {

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, errors.Wrap(err, "bcrypt")
	}

	verifiedCode := uniuri.NewLen(32)

	user := User{
		Email:        email,
		Password:     passwordHash,
		VerifiedCode: verifiedCode,
	}

	return user, nil
}

// insert user in to database
func (u *User) Insert(db *pg.DB) error {
	if _, err := db.Model(u).Returning("*").Insert(); err != nil {
		return err
	}
	return nil
}

func GetUser(db *pg.DB, email string) (User, error) {
	var user User
	if err := db.Model(&user).Where("email = ?", email).Select(); err != nil {
		return User{}, err
	}

	return user, nil
}

func VerifyUser(db *pg.DB, code string) error {
	var user User
	err := db.Model(&user).Where("verified_code = ? AND verified_at IS NULL", code).Select()
	if err != nil {
		return errors.Wrap(err, "Count")
	}

	if user.ID == 0 {
		return errors.New("Not found")
	}

	if _, err := db.Model(&user).Set("verified_at = now()").WherePK().Update(); err != nil {
		return errors.Wrap(err, "Update")
	}

	return nil
}
