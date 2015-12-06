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
  User            *User
}

func (p *Post) String() string {
  return fmt.Sprintf("Post(%s)", p.User)
}

func (p Post) getUser() string {
  return p.User.Name
}

//func (b Booking) Total() int {
 // return b.Hotel.Price * b.Nights()
//}


func (p *Post) PreInsert(_ gorp.SqlExecutor) error {
  p.UserId = p.User.UserId
  return nil
}

func (post Post) Validate(v *revel.Validation) {
  v.Required(post.User)
}