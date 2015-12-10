package controllers
import (
  "golang.org/x/crypto/bcrypt"
  "database/sql"
  "github.com/go-gorp/gorp"
  _ "github.com/mattn/go-sqlite3"
  r "github.com/revel/revel"
  "github.com/revel/modules/db/app"
  "twitter-app/app/models"
  "time"
  "fmt"
)

var (
  Dbm *gorp.DbMap
)

func InitDB() {
  db.Init()
  Dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.SqliteDialect{}}

  setColumnSizes := func(t *gorp.TableMap, colSizes map[string]int) {
      for col, size := range colSizes {
          t.ColMap(col).MaxSize = size
      }
  }

  t := Dbm.AddTable(models.User{}).SetKeys(true, "UserId")
  t.ColMap("Password").Transient = true
  setColumnSizes(t, map[string]int{
      "Username": 20,
      "Name":     100,
  })

 t = Dbm.AddTable(models.Post{}).SetKeys(true, "PostId")
  t.ColMap("User").Transient = true
  setColumnSizes(t, map[string]int{
      "Message": 160,
      "Date": 10,
      "Likes": 20,
  })

  y := Dbm.AddTable(models.Friend{}).SetKeys(true, "FriendId")
  setColumnSizes(y, map[string]int{
      "UserIdOne": 10,
      "UserIdTwo": 10,
      "AreFriends": 1,
  })

  Dbm.TraceOn("[gorp]", r.INFO)
  Dbm.CreateTables()

  bcryptPassword, _ := bcrypt.GenerateFromPassword(
      []byte("demo"), bcrypt.DefaultCost)
  demoUser := &models.User{0, "Levi Johnston", "levigene123", "asdf", bcryptPassword}
  demoUser2 := &models.User{0, "Zack Peters", "pistol123", "asdf", bcryptPassword}
  demoUser3 := &models.User{0, "Tim Courtney", "timmy", "asdf", bcryptPassword}



  if err := Dbm.Insert(demoUser); err != nil {
      panic(err)
  }
  if err := Dbm.Insert(demoUser2); err != nil {
      panic(err)
  }
  if err := Dbm.Insert(demoUser3); err != nil {
      panic(err)
  }

  posts := []*models.Post{
    &models.Post{0, "Hello World this post is by levi", time.Now(), 3, 1, demoUser},
    &models.Post{0, "Levi's status update", time.Now(), 0, 1,demoUser},
    &models.Post{0, "Hello World", time.Now(), 0, 1, demoUser},
    &models.Post{0, "Post by Zack", time.Now(), 0, 1, demoUser2},
    &models.Post{0, "Post by Tim only visible to him and Zack", time.Now(), 0, 1, demoUser3},

  }
  for _, post := range posts {
    if err := Dbm.Insert(post); err != nil {
      panic(err)
    }
    fmt.Println("User from post = ", post.User)
  }

    friends := &models.Friend{0, 1, 2, true}
    if err := Dbm.Insert(friends); err != nil {
      panic(err)
    }

    friends2 := &models.Friend{0, 2, 3, true}
    if err := Dbm.Insert(friends2); err != nil {
      panic(err)
    }

}

type GorpController struct {
  *r.Controller
  Txn *gorp.Transaction
}

func (c *GorpController) Begin() r.Result {
  txn, err := Dbm.Begin()
  if err != nil {
      panic(err)
  }
  c.Txn = txn
  return nil
}

func (c *GorpController) Commit() r.Result {
  if c.Txn == nil {
      return nil
  }
  if err := c.Txn.Commit(); err != nil && err != sql.ErrTxDone {
      panic(err)
  }
  c.Txn = nil
  return nil
}

func (c *GorpController) Rollback() r.Result {
  if c.Txn == nil {
      return nil
  }
  if err := c.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
      panic(err)
  }
  c.Txn = nil
  return nil
}