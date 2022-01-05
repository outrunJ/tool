package sqler

import (
	"fmt"
	"testing"
	"time"
)

func TestClause(t *testing.T) {

	type Depot struct {
		Cash int
		Date int `dbName:"date"`
	}
	type User struct {
		UserName  string
		PWD       int `dbName:"pwd"`
		CreatedAt *time.Time
		Depot     *Depot
	}

	//now := time.Now()
	s := &User{
		UserName: "a",
		PWD:      1,
		//CreatedAt: &now,
		Depot: &Depot{
			Cash: 2,
		},
	}

	sqle := Sqler(s)
	fmt.Println(sqle.Table())
	fmt.Println(sqle.Columns("UserName"))
	fmt.Println(sqle.ColumnsEx("UserName"))
	fmt.Println(sqle.AllColumns("UserName"))
	fmt.Println(sqle.AllColumnsEx("UserName"))
	fmt.Println(sqle.Values("UserName"))
	fmt.Println(sqle.ValuesEx("UserName"))
	fmt.Println(sqle.ClauseColumns("UserName"))
	fmt.Println(sqle.ClauseColumnsEx("UserName"))
	fmt.Println(sqle.ClauseAllColumns("UserName"))
	fmt.Println(sqle.ClauseAllColumnsEx("UserName"))
	fmt.Println(sqle.ClausePatterns("UserName"))
	fmt.Println(sqle.ClausePatternsEx("UserName"))
	fmt.Println(sqle.ClauseEqualComma("UserName"))
	fmt.Println(sqle.ClauseEqualCommaEx("UserName"))
	fmt.Println(sqle.ClauseEqualAnd("UserName"))
	fmt.Println(sqle.ClauseEqualAndEx("UserName"))
	fmt.Println(sqle.SqlSelect())
	fmt.Println(sqle.SqlInsert())
	fmt.Println(sqle.SqlDelete())
}
