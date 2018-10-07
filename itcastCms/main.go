package main

import (
	_ "itcastCms/routers"
	"github.com/astaxie/beego"
	_"itcastCms/models"
	"itcastCms/models"
)

func main() {
	beego.BConfig.WebConfig.Session.SessionOn = true//启用SESSION
	beego.AddFuncMap("checkRoleAction",CheckRoleAction)
	beego.AddFuncMap("checkUserRole",CheckUserRole)
	beego.AddFuncMap("checkExtAction",CheckExtAction)
	beego.AddFuncMap("checkPassAction",CheckPassAction)

	beego.Run()
}
//视图函数，判断角色是否已经存在
func CheckUserRole(userExtRoles []*models.RoleInfo,roleId int )(b bool)  {
	//判断传递过来的roleId是否在userExRoles 中
	 b = false
	 for i := 0; i < len(userExtRoles); i++{
	 	if userExtRoles[i].Id == roleId{
	 		b = true
	 		break
		}
	 }
	return b
}
func CheckRoleAction(roleExtActions []*models.ActionInfo, roleId int)(b bool)  {
	for _,actions := range roleExtActions{
		if actions.Id == roleId{
			 	b = true
			break
		}
	}
	return b
}
//视图函数的实现，判断用户已经存在的权限
func CheckExtAction(userExtActions []models.UserAction,actionId int )(b bool)  {
	//判断传递过来的actionId是否在userExRoles 中
	b = false
	for i := 0; i < len(userExtActions); i++{
		if userExtActions[i].Actions.Id == actionId{
			b = true
			break
		}
	}
	return b
}
//判断用户存在的权限是允许还是禁止
func CheckPassAction(userExtActions []models.UserAction,actionId int  )(b bool)  {
	//判断传递过来的actionId是否在userExRoles 中
	b = false
	for i := 0; i < len(userExtActions); i++{
		if userExtActions[i].Actions.Id == actionId{
			if userExtActions[i].IsPass == 1{
				b = true
				break
				}
			}
	}
	return b
}