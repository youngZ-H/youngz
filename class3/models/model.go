package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

//创建用户信息的结构体
type User struct {
	Id       int
	UserName string
	PassWord string
	Article  []*Article `orm:"rel(m2m)"`
}

//创建内容详情的结构体
type Article struct {
	Id           int           `orm:"pk;atuto"`
	Title        string        `orm:"size(20)""`
	Content      string        `orm:"size(500)"`
	Img          string        `orm:"size(50);null"`
	Time         time.Time     `orm:"type(datetime);auto_now_add"`
	Count        int           `"orm:default(0)"`
	ArticleStype *ArticleStype `orm:"rel(fk)"`
	Users        []*User       `orm:"reverse(many)"`
}

//创建文章类型的结构体,进行一对多的关系的创建
type ArticleStype struct {
	Id       int	 	//`orm:"pk;atuto"`
	TypeName string     `orm:"size(20)"`
	Articles []*Article `orm:"reverse(many)"`
}

func init() {
	/*
	进行数据库的操作
	1.连接数据库
	2.注册表
	3.插入数据
	*/
	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/newsWeb?charset=utf8")
	orm.RegisterModel(new(User), new(Article), new(ArticleStype))
	orm.RunSyncdb("default", false, true)
}
