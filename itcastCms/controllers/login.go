package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"itcastCms/models"
)

type LoginController struct {
	beego.Controller
}
func (this *LoginController)Index(){
	this.TplName = "Login/Index.html"
}
func (this * LoginController)UserLogin()  {
	//从前端获取用户名和用户的密码
	userName := this.GetString("LoginCode")
	userPwd := this.GetString("LoginPwd")
	//创建orm对象
	o := orm.NewOrm()
	//创建一个用户信息的对象，用来存储查询的信息
	var userInfo models.UserInfo
	o.QueryTable("user_info").Filter("user_name",userName).Filter("user_pwd",userPwd).One(&userInfo)
	//对查询到的信息进行判断，是否在数据库中用这个用户名
	if userInfo.Id > 0 {
		//因为主键id是至少大于1的数，说明有此用户名的存在
		//进行session的值得设置
		this.SetSession("userId", userInfo.Id)
		this.SetSession("userName", userInfo.UserName)
		//将从数据库中取出来的数据传递给前端
		this.Data["json"] = map[string]interface{}{"flag": "ok"}

	}else {
		this.Data["json"] = map[string]interface{}{"flag":"no"}
	}
	//将数据转换为ajax的形式传递给前端
	this.ServeJSON()
}