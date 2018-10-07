package models

import (
	"time"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init(){
	var dbhost string
	var dbport string
	var dbuser string
	var dbpassword string
	var db string
	//获取配置文件中对应的配置信息
	dbhost = beego.AppConfig.String("dbhost")
	dbport = beego.AppConfig.String("dbport")
	dbuser = beego.AppConfig.String("dbuser")
	dbpassword = beego.AppConfig.String("dbpassword")
	db = beego.AppConfig.String("db")
	orm.RegisterDriver("mysql", orm.DRMySQL) //注册mysql Driver
	//构造conn连接
	conn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + db + "?charset=utf8"
	//注册数据库连接
	orm.RegisterDataBase("default", "mysql", conn)

	orm.RegisterModel(new(UserInfo),new(RoleInfo),new(ActionInfo),new(UserAction),new(ArticleClass),new(ArticleInfo),new(ArticleComment),new(SensitiveWord))//注册模型
	orm.RunSyncdb("default", false, true)

}
type UserInfo struct {
	Id int
	UserName string //用户名
	UserPwd string//密码
	Remark string//备注
	AddDate time.Time//添加日期
	ModifDate time.Time//修改日期
	DelFlag int// 删除标记 软删除。
	Roles []*RoleInfo`orm:"rel(m2m)"`
	UserActions []*UserAction `orm:"reverse(many)"`
}
//创建角色信息表
type RoleInfo struct {
	Id int
	RoleName string `orm:"size(32)"`
	Remark string
	DelFlag int
	AddDate time.Time
	ModifDate time.Time
	Users []*UserInfo`orm:"reverse(many)"`
	Actions []*ActionInfo`orm:"rel(m2m)"`
}
//创建角色权限信息表
type ActionInfo struct {
	Id int
	Remark string
	DelFlag int
	AddDate time.Time
	ModifDate time.Time
	Url string
	HttpMethod string
	ActionInfoName string
	ActionTypeEnum int //权限类型。
	MenuIcon string //图片地址
	IconWidth int
	IconHeight int
	Roles[]*RoleInfo `orm:"reverse(many)"`
	UserActions []*UserAction `orm:"reverse(many)"`

}
//中间表模型
type UserAction struct {
	Id int
	IsPass int
	Users *UserInfo `orm:"rel(fk)"`
	Actions *ActionInfo `orm:"rel(fk)"`
}
//文章类别
type ArticleClass struct {
	Id int//主键
	ClassName string//类别名称
	ParentId int //父类别的编号
	CreateUserId int //创建类别的用户编号
	CreateDate time.Time //创建时间
	DelFlag int //删除标记
	Remark string //备注
	Articles []*ArticleInfo`orm:"reverse(many)"`
}
//文章内容管理
type ArticleInfo struct {
	Id int
	KeyWords string  //关键词
	Title string    //标题
	FullTitle string  //全标题
	Intro string  //导读
	ArticleContent string `orm:"type(text)"` //新闻内容
	Author string//作者
	Origin string//来源
	AddDate time.Time//添加日期
	ModifyDate time.Time//修改日期
	DelFlag int//删除标记
	PhotoUrl string//图片地址
	ArticleClasses []*ArticleClass`orm:"rel(m2m)"`
	Comments[]*ArticleComment `orm:"reverse(many)"`
}
//文章评论
type ArticleComment struct {
	Id int `from:"-"`
	Msg string
	AddDate time.Time
	IsPass int
	Article *ArticleInfo `orm:"rel(fk)"`
}
//添加词库
type SensitiveWord struct {
	Id int
	WordPattern string
	IsForbid int
	IsMod int
	ReplaceWord string `orm:"null"`
}