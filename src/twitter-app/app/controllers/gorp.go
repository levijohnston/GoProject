package controllers
 
import (
    "database/sql"
    "github.com/coopernurse/gorp"
    _ "github.com/go-sql-driver/mysql"
    r "github.com/revel/revel"
    "social-paster/app/models"
)
 
var (
    Dbm *gorp.DbMap
)
 
func InitDB() {
    db, err := sql.Open("mysql", "db_username:db_password@/paster")
    if(err != nil){
      panic("Unable to connect to the database")
    }
    Dbm = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
 
    setColumnSizes := func(t *gorp.TableMap, colSizes map[string]int) {
        for col, size := range colSizes {
            t.ColMap(col).MaxSize = size
        }
    }
 
    t := Dbm.AddTable(models.User{}).SetKeys(true, "Id")
    t.ColMap("ConfirmPassword").Transient = true
        
        t.ColMap("Email").SetNotNull(true)
    t.ColMap("FirstName").SetNotNull(true)
    t.ColMap("LastName").SetNotNull(true)
    t.ColMap("Password").SetNotNull(true)
 
    t.ColMap("Email").Unique = true
 
    setColumnSizes(t, map[string]int{
        "FirstName": 30,
        "LastName":  30,
    })
 
    Dbm.TraceOn("[gorp]", r.INFO)
    Dbm.CreateTablesIfNotExists()
 
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
