package dber

import (
	"fmt"
	"git.meiqia.com/business_platform/tool"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"context"
)

var dberTest DBer

func init() {
	dberTest = tool.PanicErrWith(NewSqlDBer(
		"mysql",
		"root:asdf@/tool_test?charset=utf8mb4,utf8&parseTime=true",
		1,
		1,
	)).(DBer)
}

func beforeDBerTest() {
	dberTest.Exec("truncate user")
}
func afterDBerTest() {
	dberTest.Exec("truncate user")
}

type TestData struct {
	Cash int
	Desc      string
}
type TestUser struct {
	Id        int
	Name      string
	NilStr    string
	CreatedAt *time.Time
	*TestData
}

func TestConn(t *testing.T) {
	beforeDBerTest()
	defer afterDBerTest()
	var err error
	assert := assert.New(t)

	_, err = dberTest.Exec("insert into user(id, name) values(1, \"a\")")
	assert.NoError(err)

	upd := func(db DBer) {
		_, err := db.Exec("update user set name = 'b' where id = 1")
		assert.NoError(err)
	}

	get := func(db DBer) {
		_, err := db.Query("select * from user")
		assert.NoError(err)
	}

	tx := func(db DBer) error {
		fmt.Println("33")
		upd(db)
		get(db)
		go func() {
			fmt.Println("77")
			upd(dberTest)
			fmt.Println("88")
		}()
		time.Sleep(5 * time.Microsecond)
		fmt.Println("44")
		return err
	}

	go func() {
		fmt.Println("11")
		err, _, _ := dberTest.Tx(tx, nil)
		fmt.Println("22")
		assert.NoError(err)
	}()

	go func() {
		time.Sleep(time.Microsecond)
		fmt.Println("55")
		upd(dberTest)
		fmt.Println("66")
	}()

	time.Sleep(1000 * time.Microsecond)

}

// TODO: fully automatic testing
func TestQueryScan(t *testing.T) {
	beforeDBerTest()
	defer afterDBerTest()
	assert := assert.New(t)

	userSlice := &[]*TestUser{}

	assert.NoError(dberTest.SQLQuery("select * from user").Each(func(scan func(...interface{}) error) error {
		user := new(TestUser)
		err := scan(&user.Id, &user.Name)
		if err != nil {
			return err
		}
		*userSlice = append(*userSlice, user)
		return nil
	}))

	assert.NotEmpty(&userSlice)

}

// TODO: fully automatic testing, not using Exec
func TestTx(t *testing.T) {
	beforeDBerTest()
	defer afterDBerTest()
	assert := assert.New(t)

	// err
	err, _, _ := dberTest.Tx(func(db DBer) error {
		_, err := db.Exec("insert into user(name) values(\"a1\")")
		assert.NoError(err)

		err, _, _ = db.Tx(func(db DBer) error {
			_, err = db.Exec("insert into user(name) values(\"a2\")")
			assert.NoError(err)

			err, _, _ := db.Tx(func(db DBer) error {
				_, err := db.Exec("insert into user(name) values(\"a3\")")
				return err
			}, func() error {
				return fmt.Errorf("rb")
			})

			return err
		}, nil)

		return err
	}, nil)
	assert.NoError(err)

	err, _, _ = dberTest.Tx(func(db DBer) error {
		_, err = db.Exec("insert into user(name) values(\"b1\")")
		assert.NoError(err)

		err, _, rbErr := db.Tx(func(db DBer) error {
			_, err := db.Exec("insert wrong into user(name) values(\"b2\")")
			return err
		}, func() error {
			return fmt.Errorf("rb")
		})
		assert.Error(rbErr)
		return err
	}, nil)
	assert.Error(err)

	// panic
	err, _, _ = dberTest.Tx(func(d DBer) error {
		panic("a")
		return nil
	}, func() error {
		return nil
	})
	assert.Error(err)
}

func TestModel(t *testing.T) {
	beforeDBerTest()
	defer afterDBerTest()
	assert := assert.New(t)

	users := &[]*TestUser{}
	tim := time.Unix(123, 0).UTC()
	timeString := tim.String()
	assert.NoError(tool.Err(dberTest.Exec("insert into user(name, `desc`, created_at, cash) values(\"a\", \"b\", '" + timeString + "' , 2)")))
	assert.NoError(dberTest.SQLQuery("select * from user").Model(users))

	assert.Equal(&TestUser{
		Id:        1,
		Name:      "a",
		CreatedAt: &tim,
		TestData: &TestData{
			Cash: 2,
			Desc:      "b",
		},
	}, (*users)[0])
}

//func TestModelCache(t *testing.T)  {
//	beforeDBerTest()
//	defer afterDBerTest()
//	assert := assert.New(t)
//
//	dberTest.AddModel
//}

func TestContext(t *testing.T) {
	beforeDBerTest()
	defer afterDBerTest()
	assert := assert.New(t)
	var err error
	dbe := dberTest.Context(context.Background())

	_, err = dbe.Exec("insert into user(name) values(\"a\")")
	assert.NoError(err)

	rows, err := dbe.Query("select * from user")
	assert.NoError(err)
	rows.Close()

	u1 := &TestUser{}
	err = dbe.SQL("select * from user").Query().Model(u1)
	assert.NoError(err)
	assert.Equal(u1.Name, "a")

	u2 := &TestUser{}
	err = dbe.SQLQuery("select name from user").Each(func(scan func(...interface{}) error) error {
		return scan(&u2.Name)
	})
	assert.NoError(err)
	assert.Equal(u2.Name, "a")

	row := dbe.QueryRow("select * from user")
	row.Scan()

	err, _, _ = dbe.Tx(func(d DBer) error {
		err, _, _ := d.Tx(func(d DBer) error {
			_, err := d.Exec("update user set name=\"b\"")
			return err
		}, nil)
		return err
	}, nil)
	assert.NoError(err)
}

func TestNewID(t *testing.T) {
	beforeDBerTest()
	defer afterDBerTest()
	assert := assert.New(t)
	var err error

	dberTest.SetNewIDContext(func(ctx context.Context) (string, error) {
		return "1", nil
	})

	id, err := dberTest.Context(context.Background()).NewID()
	assert.NoError(err)
	assert.Equal(id, "1")
}
