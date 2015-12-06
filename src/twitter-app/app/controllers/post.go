package controllers

import (
 // "golang.org/x/crypto/bcrypt"
  "github.com/revel/revel"
  "twitter-app/app/models"
  "twitter-app/app/routes"
  "fmt"
  "time"
)

type Post struct {
  App
}

func (c App) SavePost(post models.Post) revel.Result {
  post.User = c.connected()
  user := post.User
  post.UserId = user.UserId
  post.Date = time.Now()
  post.Likes = 20
  post.Validate(c.Validation)

  fmt.Println("Inserted into ID ", post)
  fmt.Println("Post user = ", post.UserId)

  err := c.Txn.Insert(&post)

  if err != nil {
    panic(err)
  } 
  return c.Redirect(routes.App.Show(user.UserId))
}