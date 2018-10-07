package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"itcastCms/models"
)
//定义home类
type HomeController struct {
	beego.Controller
}
//为home类绑定展示视图的方法
func (this *HomeController)ShowIndex(){
	this.TplName = "Home/ShowIndex.html"
}
func (this *HomeController)Index()  {
	this.TplName = "Home/Index.html"
}
//菜单权限过滤
func (this *HomeController)GetMenus()  {
	//首先根据的“用户”---“角色”--"权限"这条线进行过滤。
	//1:根据登录用户，获取对应的用户信息。
		//根据登录信息，获取用户的id
		userId := this.GetSession("userId")
		o := orm.NewOrm()
		var userInfo models.UserInfo
		o.QueryTable("user_info").Filter("id",userId).One(&userInfo)
	//2:根据登录用户信息，查找对应的角色。
		//进行多表查询
		o.LoadRelated(&userInfo,"Roles")
		var roles []*models.RoleInfo
		//查询出每个用户所拥有的角色
		for _, role := range userInfo.Roles{
			roles = append(roles,role)
		}
	//3:根据角色，查询对应的权限。
	//声明一个切片用来存储角色的权限信息的
	var actions[]*models.ActionInfo
	for i:=0;i<len(roles);i++{
		//进行角色和权限的多表查询
		o.LoadRelated(roles[i],"Actions")
		//查询出角色的权限存储到一个切片中
		for _,action := range roles[i].Actions{
			actions = append(actions,action)
		}
	}
	//4；判断这些查询出的权限，哪些菜单权限。
	var menuActions []*models.ActionInfo
	for i:=0;i<len(actions);i++{
		if actions[i].ActionTypeEnum == 1{
			menuActions = append(menuActions,actions[i])
		}
	}
	//按照“用户”--“权限”这条线，查询出登录的菜单权限。
	//查询出用户所拥有的全部权限
	var subActions []*models.UserAction
	o.QueryTable("user_action").Filter("users_id",userId).All(&subActions)
	//判断用户的权限那些是菜单权限
	var subMenuActions []*models.ActionInfo
	for _,subAction := range subActions{
		var actionInfo models.ActionInfo
		o.QueryTable("action_info").Filter("id",subAction.Actions.Id).Filter("action_type_enum",1).One(&actionInfo)
		//判断一下是否是菜单权限
		if actionInfo.Id>0{
			subMenuActions = append(subMenuActions,&actionInfo)
		}
	}
	//将两条线合并
	menuActions = append(menuActions,subMenuActions...)
	//进行去重操作
	temp := RemoveRepeatedElement(menuActions)
	//接下来判断temp中存储的当前登录用户的权限是否有被禁止，如果有，则去掉。
	//查询一下是否有被禁止的权限
	var userForbidActions []models.UserAction
	o.QueryTable("user_action").Filter("users_id",userId).Filter("is_pass",0).All(&userForbidActions)
	//如果forbidActions 中有元素，就说明有被禁止的权限
	var newTemp []*models.ActionInfo
	if len(userForbidActions)>0{
		//找到了禁用权限后，清除。
		for _,action := range temp{
			if CheckUserForAction(userForbidActions,action.Id) == false{
				newTemp = append(newTemp,action)
			}
		}
		this.Data["json"] = map[string]interface{}{"menus":newTemp}
	}else {
		//没有找到禁用的权限
		this.Data["json"] = map[string]interface{}{"menus":temp}
	}
	this.ServeJSON()

}
//判断用户是否有禁用权限函数的实现
func CheckUserForAction(userForbidActions []models.UserAction,actionId int)(b bool) {
	b = false
	for i:=0;i<len(userForbidActions);i++{
		if userForbidActions[i].Actions.Id == actionId{
			b = true
			break
		}
	}
	return
}

//去重操作
func RemoveRepeatedElement(arr []*models.ActionInfo) (newArr []*models.ActionInfo) {
	newArr = make([]*models.ActionInfo, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i].Id == arr[j].Id {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}
