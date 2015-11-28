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

/*func (c Post) loadPostById(id int) *models.Post {
  h, err := c.Txn.Get(models.Post{}, id)
  if err != nil {
    panic(err)
  }
  if h == nil {
    return nil
  }
  return h.(*models.Post)
}
*/

