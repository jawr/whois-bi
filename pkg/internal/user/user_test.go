package user

import (
	"testing"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/jawr/whois-bi/pkg/internal/db"
)

const (
	email    = "testuser@place.com"
	password = "MyFancyPassword1"
)

func createConnection(t *testing.T) *pg.DB {
	t.Helper()
	conn, err := db.SetupDatabase()
	if err != nil {
		t.Fatalf("SetupDatabase unexpected error %q", err)
	}
	return conn
}

func createTx(t *testing.T, conn *pg.DB) *pg.Tx {
	t.Helper()
	tx, err := conn.Begin()
	if err != nil {
		t.Fatalf("Begin expected nil got %q", err)
	}

	return tx
}

func createUser(t *testing.T, db orm.DB) User {
	t.Helper()

	user, err := NewUser(email, password)
	if err != nil {
		t.Fatalf("NewUser expected nil got %q", err)
	}

	if err := user.Insert(db); err != nil {
		t.Fatalf("Insert expected nil got %q", err)
	}

	return user
}

func Test_ValidatePassword(t *testing.T) {
	t.Parallel()

	type tcase struct {
		password string
		err      string
	}

	cases := []tcase{
		tcase{"pass", "password must be at least 8 characters long"},
		tcase{"password123", "password must have at least one upper case character"},
		tcase{"PASSwordsARE", "password must have at least one numeric character"},
		tcase{"PASSWORDSARECOOL1", "password must have at least one lower case character"},
	}

	for _, tc := range cases {
		t.Run(tc.password, func(tt *testing.T) {
			err := ValidatePassword(tc.password)
			if err == nil {
				tt.Fatal("expected an error")
			}
			if err.Error() != tc.err {
				tt.Errorf("expected %q got %q", tc.err, err.Error())
			}
		})
	}
}

func Test_NewUser(t *testing.T) {
	t.Parallel()

	conn := createConnection(t)
	defer conn.Close()

	tx := createTx(t, conn)
	defer tx.Rollback()

	user := createUser(t, tx)

	if user.ID == 0 {
		t.Fatal("Expected user ID to not be 0")
	}

	if !user.VerifiedAt.IsZero() {
		t.Fatal("Expected VerifiedAt to be zero")
	}

	if err := VerifyUser(tx, user.VerifiedCode); err != nil {
		t.Fatalf("VerifyUser expected err to be nil, got %q", err)
	}

	user2, err := GetUser(tx, email)
	if err != nil {
		t.Fatalf("GetUser expected nil, but got %q", err)
	}

	if user2.VerifiedAt.IsZero() {
		t.Fatalf("Expected VerifiedAt to not be zero")
	}

	if user.ID != user2.ID {
		t.Fatalf("Expected user ID to match, %d vs %d", user.ID, user2.ID)
	}
}

func Test_DuplicateUser(t *testing.T) {
	t.Parallel()

	email := "testuser@place.com"
	password := "MyFancyPassword1"

	conn := createConnection(t)
	defer conn.Close()

	tx := createTx(t, conn)
	defer tx.Rollback()

	createUser(t, tx)

	dupe, err := NewUser(email, password)
	if err != nil {
		t.Fatalf("NewUser expected nil got %q", err)
	}

	err = dupe.Insert(tx)
	if err == nil {
		t.Fatal("dupe Insert expected error got nil")
	}

	expected := `ERROR #23505 duplicate key value violates unique constraint "users_email_key"`
	if err.Error() != expected {
		t.Fatalf("dupe Insert expected error to be %q, got %q", expected, err)
	}
}

func Test_BadPassword(t *testing.T) {
	t.Parallel()

	email := "testuser@place.com"
	password := "mybadpassword1"

	_, err := NewUser(email, password)
	if err == nil {
		t.Fatal("NewUser expected error got nil")
	}

	expected := `password must have at least one upper case character`
	if err.Error() != expected {
		t.Fatalf("NewUser expected %q got %q", expected, err)
	}
}

func Test_VerifyNoUser(t *testing.T) {
	t.Parallel()

	conn := createConnection(t)
	defer conn.Close()

	tx := createTx(t, conn)
	defer tx.Rollback()

	err := VerifyUser(tx, "noneexistint")
	if err == nil {
		t.Fatal("VerifyUser expected error got nil")
	}

	expected := `Select: pg: no rows in result set`
	if err.Error() != expected {
		t.Fatalf("VerifyUser expected %q to be nil, got %q", expected, err)
	}
}

func Test_VerifyAlreadyVerified(t *testing.T) {
	t.Parallel()

	conn := createConnection(t)
	defer conn.Close()

	tx := createTx(t, conn)
	defer tx.Rollback()

	user := createUser(t, tx)

	err := VerifyUser(tx, user.VerifiedCode)
	if err != nil {
		t.Fatalf("VerifyUser expected nil got %q", err)
	}

	err = VerifyUser(tx, user.VerifiedCode)
	if err == nil {
		t.Fatal("VerifyUser dupe expected error got nil")
	}

	expected := `Select: pg: no rows in result set`
	if err.Error() != expected {
		t.Fatalf("VerifyUser dupe expected %q got %q", expected, err)
	}
}

func Test_GetInvalidUser(t *testing.T) {
	t.Parallel()

	conn := createConnection(t)
	defer conn.Close()

	tx := createTx(t, conn)
	defer tx.Rollback()

	_, err := GetUser(tx, "doesnotexist")

	expected := `pg: no rows in result set`
	if err.Error() != expected {
		t.Fatalf("GetUser() expected %q got %q", expected, err)
	}
}
