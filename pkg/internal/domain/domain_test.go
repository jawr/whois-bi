package domain

import (
	"testing"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/jawr/whois-bi/pkg/internal/db"
	"github.com/jawr/whois-bi/pkg/internal/user"
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

func createOwner(t *testing.T, db orm.DB) user.User {
	t.Helper()
	u, err := user.NewUser("hi@whois.bi", "SuperStrongPassword1")
	if err != nil {
		t.Fatalf("NewUser() expected nil got %q", err)
	}
	if err := u.Insert(db); err != nil {
		t.Fatalf("User.Insert() expected nil got %q", err)
	}
	if err := user.VerifyUser(db, u.VerifiedCode); err != nil {
		t.Fatalf("VerifyUser() expected nil got %q", err)
	}
	return u
}

func createDomain(t *testing.T, db orm.DB, o user.User, name string) Domain {
	t.Helper()
	d := NewDomain(name, o)
	if err := d.Insert(db); err != nil {
		t.Fatalf("Domain.Insert() expected nil got %q", err)
	}
	return d
}

func Test_CreateDomain(t *testing.T) {
	t.Parallel()

	conn := createConnection(t)
	defer conn.Close()

	tx := createTx(t, conn)
	defer tx.Rollback()

	o := createOwner(t, tx)

	d := createDomain(t, tx, o, "testdomain.com")
	d2, err := GetDomain(tx, "testdomain.com")
	if err != nil {
		t.Fatalf("GetDomain() expected nil got %q", err)
	}

	if d.ID != d2.ID {
		t.Errorf("Expected IDs to match %d vs %d", d.ID, d2.ID)
	}
}

func Test_DomainStringer(t *testing.T) {
	t.Parallel()

	const name string = "testdomain.com"
	dom := Domain{Domain: name}
	if dom.String() != name {
		t.Errorf("Expected %q got %q", name, dom.String())
	}
}

func Test_GetRecords(t *testing.T) {
	t.Parallel()

	conn := createConnection(t)
	defer conn.Close()

	tx := createTx(t, conn)
	defer tx.Rollback()
	o := createOwner(t, tx)

	d := createDomain(t, tx, o, "testdomain.com")

	records, err := d.GetRecords(tx)
	if err != nil {
		t.Fatalf("GetRecords() expected nil got %q", err)
	}

	if len(records) != 0 {
		t.Fatalf("expected 0 records got %d", len(records))
	}
}

func Test_GetRecordsError(t *testing.T) {
	t.Parallel()

	conn := createConnection(t)
	defer conn.Close()

	tx := createTx(t, conn)
	defer tx.Rollback()

	d := Domain{Domain: "testdomain.com"}

	records, err := d.GetRecords(tx)
	if err != nil {
		t.Fatalf("GetRecords() expected nil error got %q", err)
	}

	if len(records) != 0 {
		t.Fatalf("GetRecords() expected 0 got %d", len(records))
	}
}

func Test_GetDomainsWhereLastJobBefore(t *testing.T) {
	t.Parallel()

	conn := createConnection(t)
	defer conn.Close()

	tx := createTx(t, conn)
	defer tx.Rollback()

	o := createOwner(t, tx)

	d1 := createDomain(t, tx, o, "testdomain1.com")
	d2 := createDomain(t, tx, o, "testdomain2.com")
	createDomain(t, tx, o, "testdomain3.com")
	createDomain(t, tx, o, "testdomain4.com")

	jobNeeded, err := GetDomainsWhereLastJobBefore(tx, time.Minute)
	if err != nil {
		t.Fatalf("GetDomainsWhereLastJobBefore() expected nil got %q", err)
	}

	if len(jobNeeded) != 4 {
		t.Fatalf("GetDomainsWhereLastJobBefore() expected 4 got %d", len(jobNeeded))
	}

	d1.LastJobAt = time.Now()
	if _, err := tx.Model(&d1).WherePK().Update(); err != nil {
		t.Fatalf("Update() d1 expected nil got %q", err)
	}
	d2.LastJobAt = time.Now()
	if _, err := tx.Model(&d2).WherePK().Update(); err != nil {
		t.Fatalf("Update() d2 expected nil got %q", err)
	}

	jobNeeded, err = GetDomainsWhereLastJobBefore(tx, time.Minute)
	if err != nil {
		t.Fatalf("GetDomainsWhereLastJobBefore() expected nil got %q", err)
	}

	if len(jobNeeded) != 2 {
		t.Fatalf("GetDomainsWhereLastJobBefore() expected 2 got %d", len(jobNeeded))
	}
}
