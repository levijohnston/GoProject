package controllers

import "github.com/revel/revel"

type User struct {
  *revel.Controller
}

func (c User) connected() *models.User {
  if c.RenderArgs["user"] != nil {
    return c.RenderArgs["user"].(*models.User)
  }
  if username, ok := c.Session["user"]; ok {
    return c.getUser(username)
  }
  return nil
}

func (c Application) SaveUser(user models.User, verifyPassword string) revel.Result {
  c.Validation.Required(verifyPassword)
  c.Validation.Required(verifyPassword == user.Password).
    Message("Password does not match")
  user.Validate(c.Validation)

  if c.Validation.HasErrors() {
    c.Validation.Keep()
    c.FlashParams()
    return c.Redirect(routes.Application.Register())
  }

  user.HashedPassword, _ = bcrypt.GenerateFromPassword(
    []byte(user.Password), bcrypt.DefaultCost)
  err := c.Txn.Insert(&user)
  if err != nil {
    panic(err)
  }

  c.Session["user"] = user.Username
  c.Flash.Success("Welcome, " + user.Name)
  return c.Redirect(routes.Hotels.Index())
}
