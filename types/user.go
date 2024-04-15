package types

import (
	"fmt"
	"net/mail"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      = 12
	minFirstNameLen = 3
	minLastNameLen  = 3
	minPasswordLen  = 7
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Age       uint8  `json:"age"`
	Password  string `json:"password"`
}

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       uint   `json:"age"`
}

func (p CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}

	if len(p.FirstName) < minFirstNameLen {
		errors["firstName"] = fmt.Sprintf("firstName length should be at least %d characters", minFirstNameLen)
	}

	if len(p.LastName) < minLastNameLen {
		errors["lastName"] = fmt.Sprintf("lastName length should be at least %d characters", minLastNameLen)
	}

	if len(p.Password) < minPasswordLen {
		errors["password"] = fmt.Sprintf("password length should be at least %d characters", minPasswordLen)
	}

	if !isEmailValid(p.Email) {
		errors["email"] = "email is invalid"
	}

	if p.Age > 200 || p.Age <= 0 {
		errors["age"] = "age should be betweeb 1 and 200"
	}

	return errors
}

func isEmailValid(e string) bool {
	_, err := mail.ParseAddress(e)
	return err == nil
}

func (p UpdateUserParams) ToBson() bson.M {
	m := bson.M{}
	if len(p.FirstName) > 0 {
		m["firstName"] = p.FirstName
	}

	if len(p.LastName) > 0 {
		m["lastName"] = p.LastName
	}

	if p.Age > 0 && p.Age <= 200 {
		m["age"] = p.Age
	}
	return m
}

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName    string             `bson:"firstName" json:"firstName"`
	LastName     string             `bson:"lastName" json:"lastName"`
	Email        string             `bson:"email" json:"email"`
	Age          uint8              `bson:"age" json:"age"`
	Passwordhash string             `bson:"passwordHash" json:"-"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpwd, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:    params.FirstName,
		LastName:     params.LastName,
		Email:        params.Email,
		Age:          params.Age,
		Passwordhash: string(encpwd),
	}, nil
}
