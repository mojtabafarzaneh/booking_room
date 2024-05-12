package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      = 12
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 8
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (params CreateUserParams) Validate() []string {
	errors := []string{}
	if len(params.FirstName) < minFirstNameLen {
		errors = append(errors, fmt.Sprintf("first name should be atleast %d charecter", minFirstNameLen))
	}
	if len(params.LastName) < minLastNameLen {
		errors = append(errors, fmt.Sprintf("last name should be atleast %d charecter", minLastNameLen))
	}
	if len(params.Password) < minPasswordLen {
		errors = append(errors, fmt.Sprintf("password atleast should have atleast %d charecter", minPasswordLen))
	}
	if !isEmailValid(params.Email) {
		errors = append(errors, fmt.Sprintf("email is invalid"))
	}

	return errors
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"FirstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"EncryptedPassword" json:"-"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)

	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
