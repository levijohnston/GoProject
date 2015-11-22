package models

import (
  "fmt"
  "github.com/revel/revel"
  "regexp"
  "https://github.com/go-gorp/gorp"
)

type User struct {
  Id                 int64
  FirstName          string
  LastName           string
  Email              string
  Username, Password string
  HashedPassword     []byte

  Version      int64
  DateCreated  int64
  LastUpdated  int64

  ConfirmPassword string
}

func (u *User) String() string {
  return fmt.Sprintf("User(%s)", u.Username)
}