package data

import (
	"errors"
	"time"

	"github.com/zrotrasukha/MOVAPI/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

type password struct {
	plaintext *string // to identify: if the a plaintext password is present at all
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func validateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func validatePassword(v *validator.Validator, password string) {
	v.Check(password != "", "email", "must be provided")
	v.Check(len(password) < 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) >= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be longer than 500 bytes long")

	validateEmail(v, user.Email)
	if user.Password.plaintext != nil {
		validatePassword(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		// A panic error will occur because, the has can
		// only be missing because of developers fault; logic error

		panic("missing password hash for user")
	}
}
