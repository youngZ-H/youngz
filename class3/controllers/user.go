package controllers

import (
	"github.com/astaxie/beego"
	"class3/models"

	"github.com/astaxie/beego/orm"
	"time"
)

//声明注册的类，并继承beego的父类
type RegController struct {
	beego.Controller
}

//为注册类绑定get方法,获取页面的信息
func (this *RegController) ShowReg() {
	this.TplName = "register.html"
}

//为注册函数绑定post方法，获取注册的信息
func (this *RegController) HandleReg() {
	//进行数据的处理
	/*
	1.从浏览器获取到数据
	2.进行数据的处理
	3.将数据插入数据库
	4.返回视图，创建成功返回登陆视图，创建失败返回注册视图
	*/
	//1.从浏览器获取到数据
	name := this.GetString("userName")
	passwd := this.GetString("password")
	//2.进行数据的处理,判断用户名和密码是否为空
	if name == "" || passwd == "" {
		//如果为空，返回注册页面重新进行注册
		this.TplName = "register.html"
		return
	}
	//3.将数据插入数据库
	//1.获取orm对象
	o := orm.NewOrm()
	//2.获取插入对象
	user := models.User{}
	//3.插入数据
	user.UserName = name
	user.PassWord = passwd
	_, err := o.Insert(&user)
	if err != nil {
		beego.Info("插入数据错误")
		return
	}

	//4.返回视图
	//this.TplName = "login.html"
	this.Redirect("/", 302)
}

//定义登陆的类，并继承父类
type LoginController struct {
	beego.Controller
}

//绑定登陆类的get方法，用来加载登陆界面
func (this *LoginController) ShowLogin() {
	//从视图中获取用户名
	//name := this.GetString("userName")
	name := this.Ctx.GetCookie("userName")

	if name != ""{
		this.Data["check"] = "checked"
		//将用户名传递给视图
		this.Data["userName"] = name

	}

	this.TplName = "login.html"
}

//绑定登陆类的post方法
func (this *LoginController) HandleLogin() {
	/*
	1.从浏览器获取用户的登陆信息
	2.数据处理，判断数据是否为空
	3.数据的查询，判断用户名和密码是否正确
	4.返回视图
	*/
	//1.从浏览器获取用户的登陆信息
	name := this.GetString("userName")
	passwd := this.GetString("password")
	//2.数据处理，判断数据是否为空
	if name == "" || passwd == "" {
		beego.Info("用户名或密码不能为空")
		this.TplName = "login.html"
		return
	}
	//3.数据的查询，判断用户名和密码是否正确
	//1.创建orm对象
	o := orm.NewOrm()
	//2.获取查询对象
	user := models.User{}
	user.UserName = name
	err := o.Read(&user, "UserName")
	if err != nil {
		beego.Info("用户名错误")
		this.TplName = "login.html"
		return
	}
	if user.PassWord != passwd {
		beego.Info("密码错误")
		this.TplName = "login.html"
		return
	}
	//记住用户名
	//获取remember的状态值
	check := this.GetString("remember")
	//打印测试所获取的值
	//beego.Info(check)
	//对获取的数据进行判断处理
	if check == "on"{
		this.Ctx.SetCookie("userName",name,time.Second*3600)
	}else{
		this.Ctx.SetCookie("userName",name,-1)
	}

	this.SetSession("userName",name)






	//this.Ctx.WriteString("登陆成功")
	//this.TplName = "index.html"
	this.Redirect("/Article/ShowArticle",302)

}
