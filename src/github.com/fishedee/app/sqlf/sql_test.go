//FIXME 需要加入mysql的测试
package sqlf

import (
	. "github.com/fishedee/app/log"
	. "github.com/fishedee/assert"
	. "github.com/fishedee/language"
	"testing"
	"time"
)

func initSqliteDatabase() SqlfDB {
	log, err := NewLog(LogConfig{
		Driver: "console",
	})
	if err != nil {
		panic(err)
	}
	db, err := NewSqlfDB(log, nil, SqlfDBConfig{
		Driver:     "sqlite3",
		SourceName: ":memory:",
		Debug:      true,
	})
	if err != nil {
		panic(err)
	}
	db.MustExec(`
	create table t_user(
		userId integer primary key autoincrement,
		name char(32) not null,
		age integer not null,
		money decimal(14,2) not null,
		loginTime timestamp not null,
		createTime timestamp not null default 0,
		modifyTime timestamp not null default 0
	);
	`)
	return db
}

func initMySqlDatabase() SqlfDB {
	log, err := NewLog(LogConfig{
		Driver: "console",
	})
	if err != nil {
		panic(err)
	}
	db, err := NewSqlfDB(log, nil, SqlfDBConfig{
		Driver:     "mysql",
		SourceName: "root:1@tcp(localhost:3306)/test?parseTime=true&loc=Local",
		Debug:      true,
	})
	if err != nil {
		panic(err)
	}
	db.MustExec(`
	drop table if exists t_user;
	`)
	db.MustExec(`
	create table t_user(
		userId int not null auto_increment,
		name char(32) not null,
		age integer not null,
		money decimal(14,2) not null,
		loginTime datetime not null,
		createTime datetime not null default '1970-01-01 08:00:00',
		modifyTime datetime not null default '1970-01-01 08:00:00',
		primary key(userId)
	)engine=innodb default charset=utf8mb4;`)
	return db
}

type User struct {
	UserId     int
	Name       string
	Age        int
	Money      Decimal
	LoginTime  time.Time
	CreateTime time.Time
	ModifyTime time.Time
}

type UserAdd struct {
	Name      string
	Age       int
	Money     Decimal
	LoginTime time.Time
}

type UserMod UserAdd

func testStructType(t *testing.T, db SqlfCommon) {
	//初始化，查询为空
	users := []User{}
	db.MustQuery(&users, "select * from t_user")

	AssertEqual(t, users, []User{})

	//第一次插入数据，批量插入
	userAdds := []UserAdd{
		UserAdd{Name: "fish", Age: 12, Money: "", LoginTime: time.Unix(1, 0)},
		UserAdd{Name: "cat", Age: 34, Money: "102.35", LoginTime: time.Unix(2, 0)},
	}
	db.MustExec("insert into t_user(?.column) values ?", userAdds, userAdds)

	db.MustQuery(&users, "select ?.column from t_user", users)

	AssertEqual(t, users, []User{
		User{UserId: 1, Name: "fish", Age: 12, Money: "0", LoginTime: time.Unix(1, 0), CreateTime: time.Unix(0, 0), ModifyTime: time.Unix(0, 0)},
		User{UserId: 2, Name: "cat", Age: 34, Money: "102.35", LoginTime: time.Unix(2, 0), CreateTime: time.Unix(0, 0), ModifyTime: time.Unix(0, 0)},
	})

	//删除一个数据
	db.MustExec("delete from t_user where userId = ?", 1)

	db.MustQuery(&users, "select ?.column from t_user", users)

	AssertEqual(t, users, []User{
		User{UserId: 2, Name: "cat", Age: 34, Money: "102.35", LoginTime: time.Unix(2, 0), CreateTime: time.Unix(0, 0), ModifyTime: time.Unix(0, 0)},
	})

	//更新一个数据
	userMod := UserMod{
		Name:      "cat2",
		Age:       789,
		Money:     "91.23",
		LoginTime: time.Unix(3, 0),
	}
	db.MustExec("update t_user set ?.setValue where userId = ?", userMod, 2)

	db.MustQuery(&users, "select ?.column from t_user", users)

	AssertEqual(t, users, []User{
		User{UserId: 2, Name: "cat2", Age: 789, Money: "91.23", LoginTime: time.Unix(3, 0), CreateTime: time.Unix(0, 0), ModifyTime: time.Unix(0, 0)},
	})

	//添加一个数据
	userAdd := UserAdd{
		Name:      "bird",
		Age:       56,
		Money:     "33",
		LoginTime: time.Unix(4, 0),
	}
	//这里的参数&符号不是必要的，省略后也可以正常运行，仅作测试使用
	db.MustExec("insert into t_user(?.column) values (?)", &userAdd, &userAdd)

	db.MustQuery(&users, "select ?.column from t_user", users)

	AssertEqual(t, users, []User{
		User{UserId: 2, Name: "cat2", Age: 789, Money: "91.23", LoginTime: time.Unix(3, 0), CreateTime: time.Unix(0, 0), ModifyTime: time.Unix(0, 0)},
		User{UserId: 3, Name: "bird", Age: 56, Money: "33", LoginTime: time.Unix(4, 0), CreateTime: time.Unix(0, 0), ModifyTime: time.Unix(0, 0)},
	})
}

func testStructTypeAll(t *testing.T, initDatabase func() SqlfDB) {
	db := initDatabase()
	testStructType(t, db)

	db2 := initDatabase().MustBegin()
	defer db2.MustClose()
	testStructType(t, db2)
	db2.MustCommit()

}

func testBuildInType(t *testing.T, db SqlfCommon) {
	users := []User{}

	//测试单个type类型
	db.MustExec("insert into t_user(name,age,money,loginTime) values(?,?,?,?)", "fish", 123, Decimal("23"), time.Unix(1, 0))

	db.MustQuery(&users, "select ?.column from t_user", users)

	AssertEqual(t, users, []User{
		User{UserId: 1, Name: "fish", Age: 123, Money: "23", LoginTime: time.Unix(1, 0), CreateTime: time.Unix(0, 0), ModifyTime: time.Unix(0, 0)},
	})

	//测试[]type类型
	db.MustQuery(&users, "select * from t_user where name in (?) and age in (?) and money in (?) and loginTime in (?)",
		[]string{"12", "23"},
		[]int{1, 2, 3},
		[]Decimal{"123", "456", "789", "0ab"},
		[]time.Time{time.Unix(1, 0)},
	)
	AssertEqual(t, users, []User{})
}

func testBuildInTypeAll(t *testing.T, initDatabase func() SqlfDB) {
	db := initDatabase()
	testBuildInType(t, db)

	db2 := initDatabase().MustBegin()
	defer db2.MustClose()
	testBuildInType(t, db2)
	db2.MustCommit()
}

func testTxCommit(t *testing.T, initDatabase func() SqlfDB) {
	db := initDatabase()

	tx := db.MustBegin()

	//添加一个数据
	userAdd := UserAdd{
		Name:      "bird",
		Age:       56,
		Money:     "33",
		LoginTime: time.Unix(4, 0),
	}
	tx.MustExec("insert into t_user(?.column) values (?)", userAdd, userAdd)

	tx.MustCommit()

	var users []User

	db.MustQuery(&users, "select ?.column from t_user", users)

	AssertEqual(t, users, []User{
		User{UserId: 1, Name: "bird", Age: 56, Money: "33", LoginTime: time.Unix(4, 0), CreateTime: time.Unix(0, 0), ModifyTime: time.Unix(0, 0)},
	})
}

func testTxRollBack(t *testing.T, initDatabase func() SqlfDB) {
	db := initDatabase()

	tx := db.MustBegin()

	//添加一个数据
	userAdd := UserAdd{
		Name:      "bird",
		Age:       56,
		Money:     "33",
		LoginTime: time.Unix(4, 0),
	}
	tx.MustExec("insert into t_user(?.column) values (?)", userAdd, userAdd)

	tx.MustRollback()

	var users []User

	db.MustQuery(&users, "select ?.column from t_user", users)

	AssertEqual(t, users, []User{})
}

func testTxCloseCommit(t *testing.T, initDatabase func() SqlfDB) {
	db := initDatabase()

	tx := db.MustBegin()

	func() {
		defer tx.MustClose()

		//添加一个数据
		userAdd := UserAdd{
			Name:      "bird",
			Age:       56,
			Money:     "33",
			LoginTime: time.Unix(4, 0),
		}
		tx.MustExec("insert into t_user(?.column) values (?)", userAdd, userAdd)

		tx.MustCommit()
	}()

	var users []User

	db.MustQuery(&users, "select ?.column from t_user", users)

	AssertEqual(t, users, []User{
		User{UserId: 1, Name: "bird", Age: 56, Money: "33", LoginTime: time.Unix(4, 0), CreateTime: time.Unix(0, 0), ModifyTime: time.Unix(0, 0)},
	})
}

func testTxCloseRollback(t *testing.T, initDatabase func() SqlfDB) {
	db := initDatabase()

	tx := db.MustBegin()

	func() {
		defer CatchCrash(func(e Exception) {

		})
		defer tx.MustClose()

		//添加一个数据
		userAdd := UserAdd{
			Name:      "bird",
			Age:       56,
			Money:     "33",
			LoginTime: time.Unix(4, 0),
		}
		tx.MustExec("insert into t_user(?.column) values (?)", userAdd, userAdd)

		panic("ud")

		tx.MustCommit()
	}()

	var users []User

	db.MustQuery(&users, "select ?.column from t_user", users)

	AssertEqual(t, users, []User{})
}

func testAll(t *testing.T, initDatabase func() SqlfDB) {
	testStructTypeAll(t, initDatabase)
	testBuildInTypeAll(t, initDatabase)
	testTxCommit(t, initDatabase)
	testTxRollBack(t, initDatabase)
	testTxCloseCommit(t, initDatabase)
	testTxCloseRollback(t, initDatabase)
}

func TestAll(t *testing.T) {
	testAll(t, initSqliteDatabase)
	testAll(t, initMySqlDatabase)
}
