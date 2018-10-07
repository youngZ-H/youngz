package controllers

import (
	"github.com/astaxie/beego"
	"itcastCms/models"
	"time"
	"github.com/astaxie/beego/orm"
	"strings"
	"strconv"
)

type RoleInfoController struct {
	beego.Controller
}
//展示用户的角色界面
func (this *RoleInfoController)Index()  {
	this.TplName="RoleInfo/Index.html"
}

func (this *RoleInfoController)ShowAddRole()  {

	this.TplName="RoleInfo/ShowAddRole.html"
}
//添加信息的方法实现
func (this *RoleInfoController)AddRole()  {
	var roleInfo=models.RoleInfo{}
	roleInfo.RoleName=this.GetString("roleName")
	roleInfo.Remark=this.GetString("roleRemark")
	roleInfo.DelFlag=0
	roleInfo.AddDate=time.Now()
	roleInfo.ModifDate=time.Now()
	beego.Info(roleInfo)
	o:=orm.NewOrm()
	_,err:=o.Insert(&roleInfo)
	if err==nil{
		this.Data["json"]=map[string]interface{}{"flag":"ok"}
	}else{
		this.Data["json"]=map[string]interface{}{"flag":"no"}
	}
	this.ServeJSON()
}


//展示现有的角色信息
func (this *RoleInfoController)GetRoleInfo()  {
	pageIndex,_:=this.GetInt("page")
	pageSize,_:=this.GetInt("rows")
	start:=(pageIndex-1)*pageSize
	o:=orm.NewOrm()
	var roles []models.RoleInfo
	o.QueryTable("role_info").Filter("del_flag",0).OrderBy("Id").Limit(pageSize,start).All(&roles)
	count,_:=o.QueryTable("role_info").Filter("del_flag",0).Count()
	this.Data["json"]=map[string]interface{}{"rows":roles,"total":count}
	this.ServeJSON()
}
//删除功能的实现
func (this *RoleInfoController)DeleteRole()  {
	//从前端接收数据
	strId := this.GetString("strId")
	//进行数据的处理，将字符串分割为切片
	id := strings.Split(strId,",")
	beego.Info(id)

	//定义一个角色对象
	var roleInfo models.RoleInfo

	//创建orm对象，获取要删除的对象
	o := orm.NewOrm()
	//遍历数组
	for i:=0;i<len(id);i++{
		strId,_ := strconv.Atoi(id[i])
		roleInfo.Id = strId

		//删除选中的数据
		o.Delete(&roleInfo)
	}
	//返回删除信号标志位
	this.Data["json"] = map[string]interface{}{"flag":"ok"}
	this.ServeJSON()
}
//展示要编辑的角色信息，即修改信息化时，弹出的表格，要含有原来的信息
func (this *RoleInfoController)ShowEditRole()  {
	//从前端读取要编辑的数据的id号
	id,err := this.GetInt("Id")
	if err != nil{
		beego.Info("获取要编辑信息的id发生错误", err)
		return
	}
	//创建orm对象
	o := orm.NewOrm()
	//获取要编辑的对象
	var roleInfo models.RoleInfo
	o.QueryTable("role_info").Filter("id",id).One(&roleInfo)
	//将数据编码为ajax形式，发送给前端要编辑的角色信息
	this.Data["json"] = map[string]interface{}{"roleInfo":roleInfo}
	this.ServeJSON()
	this.TplName="RoleInfo/EditRole.html"
}
//编辑角色信息，进行数据的修改
func (this *RoleInfoController)EditRole()  {
	//从前端读取获取的到数据
	//定义一个对象
	var roleInfo = models.RoleInfo{}
	//将从前端读取到的值，进行重新赋值
	roleInfo.Id,_ = this.GetInt("roleId")
	roleInfo.RoleName = this.GetString("roleName")
	roleInfo.Remark = this.GetString("roleRemark")
	roleInfo.ModifDate = time.Now()
	roleInfo.AddDate = time.Now()
	roleInfo.DelFlag = 0
	//创建orm对象
	o := orm.NewOrm()
	beego.Info(roleInfo)
	_,err := o.Update(&roleInfo)
	//判断标志位，并将标志位信息返回给前端
	if err == nil{
		this.Data["json"] = map[string]interface{}{"flag":"ok"}
	}else{
		this.Data["json"] = map[string]interface{}{"flag":"no"}
	}
	this.ServeJSON()
}
//展示位用户分配角色
func (this *RoleInfoController)ShowSetRoleAction()  {
//从前端获取选中的用户id
	id,_:=this.GetInt("roleId")
	//查询角色信息
	//创建orm对象
	o := orm.NewOrm()
	//定义一个角色对象
	var roleInfo models.RoleInfo
	o.QueryTable("role_info").Filter("Id",id).One(&roleInfo)
	//定义一个用来存储一个角色所拥有的所有权限信息
	var roleExtActions []*models.ActionInfo
	o.LoadRelated(&roleInfo,"Actions")
	//循环遍历角色所拥有的信息
	for _,action := range roleInfo.Actions{
		//将用户所拥有的信息存入切片中
		roleExtActions = append(roleExtActions,action)
	}
	//查询所有的权限信息
	var allActions []models.ActionInfo
	o.QueryTable("action_info").Filter("del_flag",0).All(&allActions)
	//将数据传递给前端
	this.Data["roleInfo"] = roleInfo
	this.Data["roleExtActions"] = roleExtActions
	this.Data["allActions"] = allActions
	//指定视图
	this.TplName = "RoleInfo/ShowSetRoleAction.html"
}
//为角色分配权限方法的实现
func (this *RoleInfoController)SetRoleAction()  {
	//1:接受角色id.
	roleId,_:=this.GetInt("roleId")
	//获取选中的权限的编号
	allKeys:=this.Ctx.Request.PostForm//获取所有的表单 map[string] []string
	var list[]int
	for key,_:=range allKeys {
		if strings.Contains(key,"cba_"){
			id:=strings.Replace(key,"cba_","",-1)
			strId,_:=strconv.Atoi(id)
			list=append(list,strId)
		}
	}
	//2:查询角色的信息
	o:=orm.NewOrm()
	var roleInfo models.RoleInfo
	o.QueryTable("role_info").Filter("id",roleId).One(&roleInfo)
	//3:根据查询出的角色信息，找出已有的权限信息看，并且全部的干掉。
	o.LoadRelated(&roleInfo,"Actions")
	m2m:=o.QueryM2M(&roleInfo,"Actions")
	for _,action:=range roleInfo.Actions {
		m2m.Remove(action)
	}
	//4：重新给角色分配权限。
	//根据list切片集合中存储的权限ID,查询出对应的权限信息，重新赋值给角色。
	var actionInfo models.ActionInfo
	for i:=0; i<len(list);i++  {
		o.QueryTable("action_info").Filter("id",list[i]).One(&actionInfo)
		m2m.Add(actionInfo)
	}
	this.Data["json"]=map[string]interface{}{"flag":"ok"}
	this.ServeJSON()

}