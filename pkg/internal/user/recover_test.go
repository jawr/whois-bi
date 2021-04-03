package user

import (
	"testing"
	"time"

	"github.com/go-pg/pg/v10/orm"
)

func createRecover(t *testing.T, db orm.DB, user User) Recover {
	t.Helper()
	r := NewRecover(user)
	if err := r.Insert(db); err != nil {
		t.Fatalf("Insert expected nil got %q", err)
	}
	return r
}

func Test_RecoverPassword(t *testing.T) {
	conn := createConnection(t)
	defer conn.Close()

	tx := createTx(t, conn)
	defer tx.Rollback()

	user := createUser(t, tx)
	rec := createRecover(t, tx, user)

	newPassword, err := CreatePassword("NewStrongerPassword1")
	if err != nil {
		t.Fatalf("CreatePassword expected nil got %q", err)
	}

	err = RecoverPassword(tx, rec.Code, newPassword)
	if err != nil {
		t.Fatalf("RecoverPassword expected nil got %q", err)
	}
}

func Test_RecoverPasswordInvalidCode(t *testing.T) {
	conn := createConnection(t)
	defer conn.Close()

	tx := createTx(t, conn)
	defer tx.Rollback()

	newPassword, err := CreatePassword("NewStrongerPassword1")
	if err != nil {
		t.Fatalf("CreatePassword expected nil got %q", err)
	}

	err = RecoverPassword(tx, "madeupcode", newPassword)
	if err == nil {
		t.Fatal("RecoverPassword expected error got nil")
	}

	expected := `invalid code`
	if err.Error() != expected {
		t.Fatalf("RecoverPassword expected %q got %q", expected, err)
	}
}

func Test_RecoverPasswordExpired(t *testing.T) {
	conn := createConnection(t)
	defer conn.Close()

	tx := createTx(t, conn)
	defer tx.Rollback()

	user := createUser(t, tx)
	rec := createRecover(t, tx, user)
	rec.ValidUntil = time.Now().Add(-time.Hour * 24)
	if _, err := tx.Model(&rec).WherePK().Update(); err != nil {
		t.Fatalf("Update ValidUntil expected nil got %q", err)
	}

	newPassword, err := CreatePassword("NewStrongerPassword1")
	if err != nil {
		t.Fatalf("CreatePassword expected nil got %q", err)
	}

	err = RecoverPassword(tx, rec.Code, newPassword)
	if err == nil {
		t.Fatal("RecoverPassword expected error got nil")
	}

	expected := `expired code`
	if err.Error() != expected {
		t.Fatalf("RecoverPassword expected %q got %q", expected, err)
	}
}
