package user

import (
	"testing"

	"github.com/jawr/whois-bi/pkg/internal/db"
)

func Test_ValidatePassword(t *testing.T) {
	type tcase struct {
		password string
		err      string
	}

	cases := []tcase{
		tcase{"pass", "password must be at least 8 characters long"},
		tcase{"password123", "password must have at least one upper case character"},
		tcase{"PASSwordsARE", "password must have at least one numeric character"},
		tcase{"PASSWORDSARECOOL", "password must have at least one lower case character"},
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
	email := "testuser@place.com"
	password := "MyFancyPassword1"

	dbConn, err := db.SetupDatabase()
	if err != nil {
		t.Fatalf("SetupDatabase unexpected error %q", err)
	}
	defer dbConn.Close()

	user, err := NewUser(email, password)
	if err != nil {
		t.Fatalf("NewUser expected nil got %q", err)
	}

	if err := user.Insert(dbConn); err != nil {
		t.Fatalf("Insert expected nil got %q", err)
	}

	if user.ID == 0 {
		t.Fatal("Expected user ID to not be 0")
	}

	if !user.VerifiedAt.IsZero() {
		t.Fatal("Expected VerifiedAt to be zero")
	}

	if err := VerifyUser(dbConn, user.VerifiedCode); err != nil {
		t.Fatalf("VerifyUser expected err to be nil, got %q", err)
	}

	user2, err := GetUser(dbConn, email)
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
