package database

import (
	. "github.com/fishedee/assert"
	"os"
	"testing"
	"time"
)

type User struct {
	UserId     int       `xorm:"not null pk autoincr"`
	Name       string    `xorm:"not null varchar(32)"`
	CreateTime time.Time `xorm:"created not null"`
	ModifyTime time.Time `xorm:"updated not null"`
}

type Users struct {
	Count int
	Data  []User
}

func testDatabaseSingle(t *testing.T, database Database) {
	var err error
	//创建表
	engine := database.(*databaseImplement).Engine
	engine.DropTables(User{})
	err = engine.CreateTables(User{})
	if err != nil {
		panic(err)
	}

	//查询空
	count, err := database.Count(&User{})
	if err != nil {
		panic(err)
	}
	AssertEqual(t, count, int64(0))

	//插入
	user1 := User{
		Name: "fish",
	}
	_, err = database.Insert(&user1)
	if err != nil {
		panic(err)
	}
	user2 := User{
		Name: "jk",
	}
	_, err = database.Insert(&user2)
	if err != nil {
		panic(err)
	}
	t.Logf("%v %v", user1, user2)

	//查询
	var users []User
	err = database.Find(&users)
	if err != nil {
		panic(err)
	}
	AssertEqual(t, len(users), 2)
	AssertEqual(t, users[0].UserId, user1.UserId)
	AssertEqual(t, users[0].Name, user1.Name)
	AssertEqual(t, users[1].UserId, user2.UserId)
	AssertEqual(t, users[1].Name, user2.Name)

	//删除
	_, err = database.Where("userId = ?", user1.UserId).Delete(&User{})
	if err != nil {
		panic(err)
	}
	users = []User{}
	err = database.Find(&users)
	if err != nil {
		panic(err)
	}
	AssertEqual(t, len(users), 1)
	AssertEqual(t, users[0].UserId, user2.UserId)
	AssertEqual(t, users[0].Name, user2.Name)

	//修改
	time.Sleep(time.Second * 2)
	user2.Name = "mc"
	_, err = database.Where("userId = ?", user2.UserId).Update(user2)
	if err != nil {
		panic(err)
	}
	users = []User{}
	err = database.Find(&users)
	if err != nil {
		panic(err)
	}
	AssertEqual(t, len(users), 1)
	AssertEqual(t, users[0].UserId, user2.UserId)
	AssertEqual(t, users[0].Name, user2.Name)
	AssertEqual(t, users[0].CreateTime.Unix(), user2.CreateTime.Unix())
	AssertEqual(t, users[0].ModifyTime.Unix()-user2.ModifyTime.Unix() >= 2, true)
	t.Logf("%v", users)
}

func TestDatabase(t *testing.T) {
	var err error

	database, err := NewDatabase(DatabaseConfig{
		Driver:   "mysql",
		User:     "root",
		Passowrd: "1",
		Host:     "127.0.0.1",
		Port:     3306,
		Database: "test",
	})
	if err != nil {
		panic(err)
	}
	testDatabaseSingle(t, database)

	defer os.Remove("./test.db")
	database2, err := NewDatabase(DatabaseConfig{
		Driver:   "sqlite3",
		Database: "./test.db",
	})
	if err != nil {
		panic(err)
	}
	testDatabaseSingle(t, database2)
}

func TestDabaseTest(t *testing.T) {
	database := NewDatabaseTest()
	database.MustExec(`
		create table t_user(
			userId integer primary key autoincrement,
			name varchar(128) not null,
			createTime datetime not null,
			modifyTime datetime not null
		);
		create table t_category(
			categoryId integer primary key autoincrement,
			name varchar(128) not null,
			createTime datetime not null,
			modifyTime datetime not null
		);
	`)
	//查询空
	count := database.MustCount(&User{})
	AssertEqual(t, count, int64(0))

	//插入
	user1 := User{
		Name: "fish",
	}
	database.MustInsert(user1)
	user2 := User{
		Name: "jk",
	}
	database.MustInsert(user2)

	//查询
	var users []User
	database.MustFind(&users)
	for key, _ := range users {
		users[key].UserId = 0
		users[key].CreateTime = time.Time{}
		users[key].ModifyTime = time.Time{}
	}
	t.Logf("%v,%v", users, []User{
		user1,
		user2,
	})
	AssertEqual(t, users, []User{
		user1,
		user2,
	})
}
