package models

import (
  "fmt"
  "github.com/revel/revel"
  "github.com/go-gorp/gorp"
  "time"
)

type Post struct {
  PostId          int
  Message         string
  Date            time.Time
  Likes           int
  UserId          int

  //Transient
  User            *User
}

//Called before post is inserted
func (p *Post) PreInsert(_ gorp.SqlExecutor) error {
  p.UserId = p.User.UserId
  return nil
}

func (post Post) Validate(v *revel.Validation) {
  v.Required(post.User)
}

//Called after select statement
func (b *Post) PostGet(s gorp.SqlExecutor) error {
  var (
    obj interface{}
    err error
  )
  obj, err = s.Get(User{}, b.UserId)
  if err != nil {
    return fmt.Errorf("Error post's user does not exist (%d): %s", b.UserId, err)
  }
  b.User = obj.(*User)
  return nil
}