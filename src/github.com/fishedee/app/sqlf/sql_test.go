package sqlf

import (
	. "github.com/fishedee/assert"
	. "github.com/fishedee/language"
	"testing"
	"time"
)

func initDatabase() SqlfDB {
	db := NewSqlfDbTest()
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

func TestStructType(t *testing.T) {
	db := initDatabase()

	//初始化，查询为空
	users := []User{}
	db.MustQuery(&users, "select * from t_user")

	AssertEqual(t, users, []User{})

	//第一次插入数据，批量插入
	userAdds := []UserAdd{
		UserAdd{Name: "fish", Age: 12, Money: "100.1", LoginTime: time.Unix(1, 0)},
		UserAdd{Name: "cat", Age: 34, Money: "102.35", LoginTime: time.Unix(2, 0)},
	}
	db.MustExec("insert into t_user(?.column) values ?", userAdds, userAdds)

	db.MustQuery(&users, "select ?.column from t_user", users)

	AssertEqual(t, users, []User{
		User{UserId: 1, Name: "fish", Age: 12, Money: "100.1", LoginTime: time.Unix(1, 0), CreateTime: time.Unix(0, 0), ModifyTime: time.Unix(0, 0)},
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

func TestBuildInType(t *testing.T) {
	db := initDatabase()
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
