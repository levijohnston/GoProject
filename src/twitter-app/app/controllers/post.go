package controllers

import (
  "github.com/revel/revel"
  "twitter-app/app/models"
  "twitter-app/app/routes"
  "fmt"
  "time"
)

type Post struct {
  App
}

func (c App) loadPostById(id int) *models.Post {
  post, err := c.Txn.Get(models.Post{}, id)
  if err != nil {
    panic(err)
  }
  if post == nil {
    return nil
  }
  return post.(*models.Post)
}

func (c App) SavePost(post models.Post) revel.Result {
  post.User = c.connected()
  user := post.User
  post.UserId = user.UserId
  post.Date = time.Now()
  post.Likes = 0
  post.Validate(c.Validation)

  fmt.Println("Inserted into ID ", post)
  fmt.Println("Post user = ", post.UserId)

  err := c.Txn.Insert(&post)

  if err != nil {
    panic(err)
  } 
  return c.Redirect(routes.App.Show(user.UserId))
}

func (c App) LikePost(postId int) revel.Result {
  post := c.loadPostById(postId)

  likes := post.Likes

  likes = likes + 1
  post.Likes = likes
    fmt.Println("Number of likes ", post.Likes)

  _, err := c.Txn.Update(post)
  if err != nil {
    panic(err)
  }
  return c.Redirect(routes.App.Index())
}