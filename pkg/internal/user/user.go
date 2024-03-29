package user

import (
	"fmt"
	"time"
	"unicode"

	"github.com/dchest/uniuri"
	"github.com/go-pg/pg/v10/orm"
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

var passwordValidation = map[string][]*unicode.RangeTable{
	"upper case": {unicode.Upper, unicode.Title},
	"lower case": {unicode.Lower},
	"numeric":    {unicode.Number, unicode.Digit},
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

next:
	for name, classes := range passwordValidation {
		for _, r := range password {
			if unicode.IsOneOf(classes, r) {
				continue next
			}
		}
		return fmt.Errorf("password must have at least one %s character", name)
	}
	return nil
}

func CreatePassword(password string) ([]byte, error) {
	if err := ValidatePassword(password); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.WithMessage(err, "bcrypt")
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
func (u *User) Insert(db orm.DB) error {
	if _, err := db.Model(u).Returning("*").Insert(); err != nil {
		return err
	}
	return nil
}

func GetUser(db orm.DB, email string) (User, error) {
	var user User
	if err := db.Model(&user).Where("email = ?", email).Select(); err != nil {
		return User{}, err
	}

	return user, nil
}

func VerifyUser(db orm.DB, code string) error {
	var user User
	err := db.Model(&user).Where("verified_code = ? AND verified_at IS NULL", code).Select()
	if err != nil {
		return errors.WithMessage(err, "Select")
	}

	if _, err := db.Model(&user).Set("verified_at = now()").WherePK().Update(); err != nil {
		return errors.WithMessage(err, "Update")
	}

	return nil
}
