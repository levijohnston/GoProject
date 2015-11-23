package controllers

import (
  "golang.org/x/crypto/bcrypt"
  "github.com/revel/revel"
  "twitter-app/app/models"
  "twitter-app/app/routes"
  "fmt"
)

type App struct {
  GorpController
}

func (c App) AddUser() revel.Result {
  if user := c.connected(); user != nil {
    c.RenderArgs["user"] = user
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
  if c.connected() != nil {
    return c.Redirect(routes.App.Index())
  }
  c.Flash.Error("Please log in first")
  return c.Render()
}


/*func (c App) Show() revel.Result {
  if c.connected() != nil {
    return c.Redirect(routes.App.Show())
  }
  c.Flash.Error("Please log in first")
  return c.Render()
}*/

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
  return c.Render(title, user)
}


func (c App) Register() revel.Result {
  return c.Render()
}

func (c App) SaveUser(user models.User, verifyPassword string) revel.Result {
  c.Validation.Required(verifyPassword)
  c.Validation.Required(verifyPassword == user.Password).
    Message("Password does not match")
  user.Validate(c.Validation)

  if c.Validation.HasErrors() {
    c.Validation.Keep()
    c.FlashParams()
    return c.Redirect(routes.App.Show(user.UserId))
  }

  user.HashedPassword, _ = bcrypt.GenerateFromPassword(
    []byte(user.Password), bcrypt.DefaultCost)
  err := c.Txn.Insert(&user)
  if err != nil {
    panic(err)
  }

  c.Session["user"] = user.Username
  c.Flash.Success("Welcome, " + user.Name)
  fmt.Println("User  %d", user.UserId)
  return c.Redirect(routes.App.Show(user.UserId))
}

func (c App) Login(username, password string, remember bool) revel.Result {
  user := c.getUser(username)
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
      return c.Redirect("http://google.com")
    }
  }

  c.Flash.Out["username"] = username
  c.Flash.Error("Login failed")
  return c.Redirect("http://google.com")
}

func (c App) Logout() revel.Result {
  for k := range c.Session {
    delete(c.Session, k)
  }
  return c.Redirect("http://google.com")
}

