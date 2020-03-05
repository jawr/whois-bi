package user

import (
	"time"

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

	LastLoginAt time.Time
}

// create a new user and bcrypt the password
func NewUser(email, password string) (User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, errors.Wrap(err, "bcrypt")
	}

	user := User{
		Email:    email,
		Password: passwordHash,
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
