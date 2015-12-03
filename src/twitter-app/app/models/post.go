package models

import (
  "fmt"
  "github.com/revel/revel"
  "time"
)

type Post struct {
  PostId          int
  Message         string
  Date            time.Time
  Likes           int
  UserId          int
  User            *User
}

func (p *Post) String() string {
  return fmt.Sprintf("Post(%s)", p.Message)
}

func (post Post) Validate(v *revel.Validation) {
 // v.Required(post.User)
}

