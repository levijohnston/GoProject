package controllers
import (
  "golang.org/x/crypto/bcrypt"
  "database/sql"
  "github.com/go-gorp/gorp"
  _ "github.com/mattn/go-sqlite3"
  r "github.com/revel/revel"
  "github.com/revel/modules/db/app"
  "twitter-app/app/models"
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

  Dbm.TraceOn("[gorp]", r.INFO)
  Dbm.CreateTables()

  bcryptPassword, _ := bcrypt.GenerateFromPassword(
      []byte("demo"), bcrypt.DefaultCost)
  demoUser := &models.User{0, "Levi Johnston", "levigene123", "asdf", bcryptPassword}

  if err := Dbm.Insert(demoUser); err != nil {
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