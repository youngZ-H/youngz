package controllers

import (
	"github.com/astaxie/beego"
	"itcastCms/models"
	"time"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

type UserInfoController struct {
	beego.Controller
}
//存储搜索的条件
type UserData struct {
	PageIndex int
	PageSize int
	Name string
	Remark string
	TotalCount int64
}

func(this *UserInfoController) Index()  {
	this.TplName="UserInfo/Index.html"
}
func(this* UserInfoController) AddUser()  {
	var userInfo=models.UserInfo{}
	userInfo.UserName=this.GetString("userName")//接收用户名
	userInfo.UserPwd=this.GetString("userPwd")
	userInfo.Remark=this.GetString("userRemark")
	userInfo.ModifDate=time.Now()
	userInfo.AddDate=time.Now()
	userInfo.DelFlag=0//表示正常，1表示表示软删除。
	o:=orm.NewOrm()
	_,err:=o.Insert(&userInfo)
	if err==nil{
		//Data中的key必须为"json"
		this.Data["json"]=map[string]interface{}{"flag":"ok"}

	}else{
		this.Data["json"]=map[string]interface{}{"flag":"no"}
	}
	//怎样将数据生成JSON.
	this.ServeJSON()
}
func (this *UserInfoController)GetUserInfo()  {
	//1.接收前端发送过来的当前页码，与每页显示的记录数
	pageIndex,_:=strconv.Atoi(this.GetString("page"))//当前页码
	pageSize,_:=strconv.Atoi(this.GetString("rows")) //每页显示记录数
	//接收前端传递过来的搜索数据
	name := this.GetString("name")
	remark := this.GetString("remark")
	//封装一个函数用来进行数据的搜索查询和分页显示
	//封装一个类，进行方法的绑定，将对象传递给方法
	//顶一个变量，将从前端获取的数据进行复制给对象
	var userSearchData=UserData{}
	userSearchData.Remark=remark
	userSearchData.PageSize=pageSize
	userSearchData.PageIndex=pageIndex
	userSearchData.Name=name
	serverData:=userSearchData.SearchUserData()

	//将搜索的数据已json的格式解析发给前端
	this.Data["json"]=map[string]interface{}{"rows":serverData,"total":userSearchData.TotalCount}
	this.ServeJSON()

	//2.实现分页查询数据
	//Limit(第一个参数，第二个参数）
	//第一个参数:获取多少条数据
	//第二个参数:从哪里开始取
	/*start := (pageIndex-1)*pageSize
	var users []models.UserInfo
	o:=orm.NewOrm()
	o.QueryTable("user_info").Filter("del_flag",0).OrderBy("id").Limit(pageSize,start).All(&users)
	count,_:=o.QueryTable("user_info").Filter("del_flag",0).Count()
	this.Data["json"]=map[string]interface{}{"rows":users,"total":count}
	this.ServeJSON()*/
}
func (this *UserInfoController)DeleteUser()  {
	ids := this.GetString("strId")
	strIds:=strings.Split(ids,",")
	o:=orm.NewOrm()
	var userInfo=models.UserInfo{}
	for i:=0;i<len(strIds);i++{
		id,_:=strconv.Atoi(strIds[i])
		userInfo.Id=id
		o.Delete(&userInfo)
	}
	this.Data["json"]=map[string]interface{}{"flag":"ok"}
	this.ServeJSON()
}
//为UserData绑定搜索的方法
func (this *UserData)SearchUserData()[]models.UserInfo  {
	//构建搜索条件
	//创建orm对象
	o := orm.NewOrm()
	//创建搜索的对象
	temp := o.QueryTable("user_info")
	//对搜索的对象的值进行判断
	if this.Name != ""{
		temp = temp.Filter("user_name__icontains",this.Name)
	}
	if this.Remark != ""{
		temp = temp.Filter("remark__icontains",this.Remark)
	}
	//根据数据是否被删除的条件，进行数据的过滤
	temp = temp.Filter("del_flag",0)
	//进行分页
	start := (this.PageIndex-1)*this.PageSize
	var users []models.UserInfo
	temp.OrderBy("id").Limit(this.PageSize,start).All(&users)
	this.TotalCount,_=temp.Count()
	//将数据返回函数
	return users
}
//绑定方法，处理要编辑的数据
func (this *UserInfoController)ShowEditUser()  {
	//从客户端接收要编辑数据的id值
	userId,err := this.GetInt("userId")
	beego.Info(userId)
	if err != nil{
		beego.Info("获取编辑数据的id发生错误",err)
		return
	}
	//声明一个用户信息对象
	var userInfo models.UserInfo

	//获取编辑对象
	o := orm.NewOrm()
	o.QueryTable("user_info").Filter("id",userId).One(&userInfo)
	//将要查询的数据转为json格式，发送给前端
	this.Data["json"] = map[string]interface{}{"userInfo":userInfo}
	this.ServeJSON()

}
//绑定方法，将编辑好的数据存储到数据库中，并返回前端一个标志位
func(this* UserInfoController) EditUser()  {
	var userInfo=models.UserInfo{}
	id, err := this.GetInt("id")
	if err != nil{
		beego.Info("获取编辑数据的id发生错误",err)
		return
	}
	userInfo.Id = id
	beego.Info("获取id发生错误", err)
	userInfo.UserName=this.GetString("userName")//接收用户名
	userInfo.UserPwd=this.GetString("userPwd")
	userInfo.Remark=this.GetString("userRemark")
	userInfo.ModifDate=time.Now()
	userInfo.AddDate=time.Now()
	userInfo.DelFlag=0//表示正常，1表示表示软删除。
	//打印测试修改好的数据
	beego.Info("打印测试修改好的数据",userInfo)
	o:=orm.NewOrm()
	_,err =o.Update(&userInfo)

	if err==nil{
		//Data中的key必须为"json"
		this.Data["json"]=map[string]interface{}{"flag":"ok"}

	}else{
		beego.Info("更新失败错误", err)
		this.Data["json"]=map[string]interface{}{"flag":"no"}
	}
	//怎样将数据生成JSON.
	this.ServeJSON()
}
//绑定方法，展示用户现在已有的角色
func (this *UserInfoController)ShowSetUserRole()  {
//1：接收传递过来的用户编号
userId,_:=this.GetInt("userId")
//2:查询出用户已经有的角色。
beego.Info("查询用户已有id", userId)

var userInfo models.UserInfo
o:=orm.NewOrm()
o.QueryTable("user_info").Filter("id",userId).One(&userInfo)
var userExtRoles []*models.RoleInfo//表示用户已有的角色。
o.LoadRelated(&userInfo,"Roles")
for _,role:=range userInfo.Roles{//延迟加载 /懒加载
userExtRoles=append(userExtRoles,role)
}
//3:查询出所有角色。
var allRoles[]models.RoleInfo
o.QueryTable("role_info").Filter("del_flag",0).All(&allRoles)
this.Data["allRoles"]=allRoles
this.Data["userExtRoles"]=userExtRoles
this.Data["userInfo"]=userInfo

this.TplName="UserInfo/ShowSetUserRole.html"
}
//绑定方法，进行用户角色的重新分配
func (this *UserInfoController)SetUserRole()  {
	//从前端获取form表单的提交信息
	allKeys := this.Ctx.Request.PostForm
	//allKEYS 的类型是：map[string][]string,我们只需要key的值
	//循环遍历allKeys的值，取出k值,并把id解析出来
	//定义一个切片，用来存储已有角色的id
	var list []int
	for k,_ := range allKeys{
		//进行字符串的判断，判断是否是否含有"cba_",如果含有，表明是已有权限
		if strings.Contains(k,"cba_"){
			//进行字符串的替换
			strId := strings.Replace(k,"cba_","",-1)
			//进行字符串的转换，字符型，转换为整型
			id, err := strconv.Atoi(strId)
			if err != nil{
				beego.Info("判断是否含有角色id整型发生错误",err)
				return
			}
			//将所有取得的id存入切片中
			list = append(list,id)
		}
		//获取用户id
		userId,err := this.GetInt("userId")
		if err != nil {
			if err != nil{
				beego.Info("用户id转整型发生错误",err)
				return
			}
		}
		//获取用户所有的角色信息，并将其清空
		//创建orm对象
		o := orm.NewOrm()
		var userInfo models.UserInfo
		o.QueryTable("user_info").Filter("id",userId).One(&userInfo)
		//建立多对多的关系
		o.LoadRelated(&userInfo,"Roles")
		m2m := o.QueryM2M(&userInfo,"Roles")
		//循环遍历中间表的信息，并将其数据清空
		for _,role := range userInfo.Roles{
			m2m.Remove(role)
		}
		//将新的角色信息添加到中间表中
		var roleInfo models.RoleInfo
		for _,id := range list{
			o.QueryTable("role_info").Filter("id",id).One(&roleInfo)
			m2m.Add(roleInfo)
		}
		//通过json格式返回个前端添加成功标志位
		this.Data["json"] = map[string]interface{}{"flag":"ok"}
		this.ServeJSON()
	}
}
//展示用户分配的权限
func (this *UserInfoController)ShowSetUserAction()  {
	//从前端获取要分配角色的id
	userId,err := this.GetInt("userId")
	if err != nil{
		beego.Info("为用户分配角色信息获取id值时发生错误",err)
		return
	}
	//创建orm对象
	o := orm.NewOrm()
	//查询当前的用户信息
	var userInfo models.UserInfo
	o.QueryTable("user_info").Filter("id",userId).One(&userInfo)
	//查询用户拥有的权限信息
	var userExtActions []models.UserAction
	o.QueryTable("user_action").Filter("users_id",userId).All(&userExtActions)
	//查询列表中所有权限信息
	var allActions []models.ActionInfo
	o.QueryTable("action_info").Filter("del_flag",0).All(&allActions)
	//将查询到的所有的用户信息传递给前端
	this.Data["userInfo"] = userInfo
	this.Data["userExtActions"] = userExtActions
	this.Data["allActions"] = allActions
	//在视图中展示对应的信息
	this.TplName = "UserInfo/ShowSetUserAction.html"
}
//用户权限的分配
func (this *UserInfoController)SetUserAction()  {
	//从前端获取发送过来的数据
	userId,_ := this.GetInt("userId")
	actionId,_ := this.GetInt("actionId")
	isPass,_ := this.GetInt("isPass")
	//判断传递过来的权限，用户是否存在，如果存在，更新权限的状态，如果没有重新插入新的数据
	var userAction models.UserAction
	//查询传递过来的用户的权限信息
	//创建orm对象
	o := orm.NewOrm()
	o.QueryTable("user_action").Filter("users_id",userId).Filter("actions_id",actionId).One(&userAction)
	//判断传递过来的权限，用户是否存在
	if userAction.Id > 0{
		//用户权限存在，重新获取用户的权限的状态
		userAction.IsPass = isPass
		//将数据更新到数据库中
		_, err := o.Update(&userAction)
		if err != nil{
			beego.Info("更新用户的权限信息时，发生错误",err)
			return
		}
	}else {
		//用户无此权限，需要添加此权限
		userAction.IsPass = isPass
		//根据传递过来的id查询对应的权限
		var actionInfo models.ActionInfo
		o.QueryTable("action_info").Filter("id",actionId).One(&actionInfo)
		//根据传递过来的id查询对应的用户信息
		var userInfo models.UserInfo
		o.QueryTable("user_info").Filter("id",userId).One(&userInfo)
		userAction.Actions = &actionInfo
		userAction.Users = &userInfo
		//将数据插入数据库
		_, err := o.Insert(&userAction)
		if err != nil{
			beego.Info("向数据库中插入用户权限信息发生错误",err)
			return
		}
	}
	this.Data["json"] = map[string]interface{}{"flag":"ok"}
	this.ServeJSON()
}
//删除用户方法的实现
func (this *UserInfoController)DeleteUserAction()  {
	//从前端获取用户的id和所拥有的权限的id
	userId,_ := this.GetInt("userId")
	actionId,_ := this.GetInt("actionId")
	//创建orm对象
	o := orm.NewOrm()
	//创建要查询的对象
	var userAction models.UserAction
	//查询中间表，用户的权限信息
	o.QueryTable("user_action").Filter("users_id",userId).Filter("actions_id",actionId).One(&userAction)
	//进行数据的删除
	o.Delete(&userAction)
	//向前端返回数据
	this.Data["json"] = map[string]interface{}{"flag":"ok"}
	this.ServeJSON()
}











