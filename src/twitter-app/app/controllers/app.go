package controllers

import (
  "golang.org/x/crypto/bcrypt"
  "github.com/revel/revel"
  "twitter-app/app/models"
  "twitter-app/app/routes"
  "fmt"
  //"time"
)

type App struct {
  *revel.Controller
  GorpController
}

func (c App) AddUser() revel.Result {
  if user := c.connected(); user != nil {
    c.RenderArgs["user"] = user
  }
  return nil
}

func (c App) checkUser() revel.Result {
  if user := c.connected(); user == nil {
    c.Flash.Error("Please log in first")
    return c.Redirect(routes.App.Index())
  }
  return nil
}

func (c App) connected() *models.User {
  if c.RenderArgs["user"] != nil {
    return c.RenderArgs["user"].(*models.User)
  }
  if username, ok := c.Session["user"]; ok {
    return c.getUser(username)
  }
  return nil
}

func (c App) getUser(username string) *models.User {
  users, err := c.Txn.Select(models.User{}, `select * from User where Username = ?`, username)
  if err != nil {
    panic(err)
  }
  if len(users) == 0 {
    return nil
  }
  return users[0].(*models.User)
}

func (c App) Index() revel.Result {
  user := c.connected()
  if user != nil {
    query := "SELECT p. * " +
    "FROM Post p, Friend f " +
    "WHERE f.UserIdOne = ? AND p.UserId = f.UserIdTwo " + " AND f.AreFriends = ? " +
    "OR f.UserIdTwo = ? AND p.UserId = f.UserIdOne " +  "AND f.AreFriends = ? " + 
    "OR p.UserId = ? " +
    "ORDER BY PostId DESC"
    results, err := c.Txn.Select(models.Post{}, query, user.UserId, true, user.UserId, true, user.UserId)
    if err != nil {
        panic(err)
    }

    var posts []*models.Post
    for _, r := range results {
      b := r.(*models.Post)
      posts = append(posts, b)
    }
    return c.Render(user, posts)
  }
  return c.Render()
}

func (c App) loadUserById(id int) *models.User {
  h, err := c.Txn.Get(models.User{}, id)
  if err != nil {
    panic(err)
  }
  if h == nil {
    return nil
  }
  return h.(*models.User)
}


func (c App) Show(id int) revel.Result {
  user := c.loadUserById(id)
  if user == nil {
    return c.NotFound("User %d does not exist", id)
  }
  title := user.Name
  results, err := c.Txn.Select(models.Post{}, `select * from Post where UserId = ? ORDER BY PostId DESC `, c.connected().UserId)
  post1, err := c.Txn.Select(models.Post{}, `select * from Post where PostId = 1`)

  fmt.Println("Post user = ", post1)

  if err != nil {
    panic(err)
  }

  var posts []*models.Post
  for _, r := range results {
    b := r.(*models.Post)
    posts = append(posts, b)
    fmt.Println("Post user = ", b.User)
  }

  return c.Render(title, user, posts)
}


func (c App) Register() revel.Result {
  return c.Render()
}

func (c App) SaveUser(user models.User, passwordConfirmation string) revel.Result {

  c.Validation.Required(passwordConfirmation)
  c.Validation.Required(passwordConfirmation == user.Password).
    Message("Password does not match")
  user.Validate(c.Validation)

  if c.Validation.HasErrors() {
      fmt.Println("User errors ")

    c.Validation.Keep()
    c.FlashParams()
    return c.Redirect(routes.App.Index())
  }

  user.HashedPassword, _ = bcrypt.GenerateFromPassword(
    []byte(user.Password), bcrypt.DefaultCost)
  err := c.Txn.Insert(&user)
  if err != nil {
    panic(err)
  }

  c.Session["user"] = user.Username
  c.Flash.Success("Welcome, " + user.Name)
  return c.Redirect(routes.App.Show(user.UserId))
}

func (c App) Login(username, password string, remember bool) revel.Result {
  fmt.Println("Input: ", username)
  user := c.getUser(username)
  fmt.Println("User: ", user)

  users, _ := c.Txn.Select(models.User{}, `select * from User`)



for _, r := range users {
    b := r.(*models.User)
    //posts = append(posts, b)
    fmt.Println("Getting ", b.Username)
  }


  if user != nil {
    err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
    if err == nil {
      c.Session["user"] = username
      if remember {
        c.Session.SetDefaultExpiration()
      } else {
        c.Session.SetNoExpiration()
      }
      c.Flash.Success("Welcome, " + username)
      return c.Redirect(routes.App.Show(user.UserId))
    }
  }

  c.Flash.Out["username"] = username
  c.Flash.Error("Invalid username and password")
  return c.Redirect(routes.App.Index())
}

func (c App) Logout() revel.Result {
  for k := range c.Session {
    delete(c.Session, k)
  }
  return c.Redirect(routes.App.Register())
}

