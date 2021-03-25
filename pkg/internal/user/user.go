package user

import (
	"time"

	"github.com/dchest/uniuri"
	"github.com/go-pg/pg/v10"
	"github.com/hesahesa/pwdbro"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID int `pg:",pk"`

	Email    string `pg:",unique,notnull"`
	Password []byte `pg:",notnull"`

	// meta data
	AddedAt   time.Time `pg:",notnull,default:now()"`
	DeletedAt time.Time `pg:",soft_delete"`

	VerifiedAt   time.Time
	VerifiedCode string `pg:",notnull"`

	LastLoginAt time.Time
}

func ValidatePassword(password string) error {
	checker := pwdbro.NewDefaultPwdBro()
	status, err := checker.RunChecks(password)
	if err != nil {
		return err
	}

	for _, s := range status {
		if !s.Safe {
			return errors.New(s.Message)
		}
	}

	return nil
}

func CreatePassword(password string) ([]byte, error) {
	if err := ValidatePassword(password); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "bcrypt")
	}

	return hash, nil
}

// create a new user and bcrypt the password
func NewUser(email, password string) (User, error) {

	passwordHash, err := CreatePassword(password)
	if err != nil {
		return User{}, err
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
