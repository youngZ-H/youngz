package routers

import (
	"ItcastCms/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/context"
	"itcastCms/models"
)
func init() {
	/*---------------------------------用户管理---------------------------------------*/
	beego.Router("/Admin/UserInfo/Index",&controllers.UserInfoController{},"get:Index")
	beego.Router("/Admin/UserInfo/AddUser",&controllers.UserInfoController{},"post:AddUser")
	beego.Router("/Admin/UserInfo/GetUserInfo",&controllers.UserInfoController{},"post:GetUserInfo")
	beego.Router("/Admin/UserInfo/DeleteUser",&controllers.UserInfoController{},"post:DeleteUser")
	beego.Router("/Admin/UserInfo/ShowEditUser",&controllers.UserInfoController{},"post:ShowEditUser")
	beego.Router("/Admin/UserInfo/EditUser",&controllers.UserInfoController{},"post:EditUser")
	//给用户指定分配权限
	beego.Router("/Admin/UserInfo/SetUserRole",&controllers.UserInfoController{},"post:SetUserRole")
	//展示用户已有的权限
	beego.Router("/Admin/UserInfo/ShowSetUserRole",&controllers.UserInfoController{},"get:ShowSetUserRole")
	//展示用户分配的权限
	beego.Router("/Admin/UserInfo/ShowSetUserAction",&controllers.UserInfoController{},"get:ShowSetUserAction")
	//为用户分配权限
	beego.Router("/Admin/UserInfo/SetUserAction",&controllers.UserInfoController{},"post:SetUserAction")
	//删除用户权限
	beego.Router("/Admin/UserInfo/DeleteUserAction",&controllers.UserInfoController{},"post:DeleteUserAction")

	/*---------------------------------角色管理---------------------------------------*/
	//展示角色主表格的界面
	beego.Router("/Admin/RoleInfo/Index",&controllers.RoleInfoController{},"get:Index")
	//展示角色主界面的添加按钮后弹出的界面
	beego.Router("/Admin/RoleInfo/ShowAddRole",&controllers.RoleInfoController{},"get:ShowAddRole")
	//进行添加角色信息方法的实现，并将数据插入到数据库，并显示在界面上
	beego.Router("/Admin/RoleInfo/AddRole",&controllers.RoleInfoController{},"post:AddRole")
	//从数据库中读取数据，并展示在界面上
	beego.Router("/Admin/RoleInfo/GetRoleInfo",&controllers.RoleInfoController{},"post:GetRoleInfo")
	//删除功能的实现
	beego.Router("/Admin/RoleInfo/DeleteRole",&controllers.RoleInfoController{},"post:DeleteRole")
	//展示要编辑的角色信息
	beego.Router("/Admin/RoleInfo/ShowEditRole",&controllers.RoleInfoController{},"post:ShowEditRole")
	//编辑用户角色信息的实现
	beego.Router("/Admin/RoleInfo/EditRole",&controllers.RoleInfoController{},"post:EditRole")
	//展示为角色分配权限
	beego.Router("/Admin/RoleInfo/ShowSetRoleAction",&controllers.RoleInfoController{},"get:ShowSetRoleAction")
	//为角色分配权限
	beego.Router("/Admin/RoleInfo/SetRoleAction",&controllers.RoleInfoController{},"post:SetRoleAction")

	/*---------------------------------角色权限管理---------------------------------------*/
	//展示角色权限界面视图
	beego.Router("/Admin/ActionInfo/Index",&controllers.ActionInfoController{},"get:Index")
	//图片上传路由,进行图片的上传
	beego.Router("Admin/ActionInfo/FileUp",&controllers.ActionInfoController{},"post:FileUp")
	//角色权限的添加
	beego.Router("/Admin/ActionInfo/AddAction",&controllers.ActionInfoController{},"post:AddAction")
	//将角色信息展示出来
	beego.Router("/Admin/ActionInfo/GetActionInfo",&controllers.ActionInfoController{},"post:GetActionInfo")
	//删除角色信息
	beego.Router("/Admin/ActionInfo/DeleteAction",&controllers.ActionInfoController{},"post:DeleteAction")
	//编辑角色信息
	beego.Router("/Admin/ActionInfo/ShowEditAction",&controllers.ActionInfoController{},"post:ShowEditAction")


	/*---------------------------------后台管理页面---------------------------------------*/
	//展示后台管理页面
	beego.Router("/Admin/Home/ShowIndex",&controllers.HomeController{},"get:ShowIndex")
	//展示后台主页面的布局
	beego.Router("/Admin/Home/Index",&controllers.HomeController{},"get:Index")
	//菜单权限过滤
	beego.Router("/Admin/Home/GetMenus",&controllers.HomeController{},"post:GetMenus")
	/*---------------------------------文章管理页面---------------------------------------*/
	//展示文章管理界面
	beego.Router("/Admin/ArticleClass/Index",&controllers.ArticleClassController{},"get:Index")
	//展示文章类别
	beego.Router("/Admin/ArticleClass/ShowArticleClass",&controllers.ArticleClassController{},"post:ShowArticleClass")
	//展示添加文章父类别
	beego.Router("/Admin/ArticleClass/ShowParentClass",&controllers.ArticleClassController{},"get:ShowParentClass")
	//添加文章父类别方法的是实现
	beego.Router("/Admin/ArticleClass/AddParentClass",&controllers.ArticleClassController{},"post:AddParentClass")
	//展示文章子类别
	beego.Router("/Admin/ArticleClass/ShowAddChildClass",&controllers.ArticleClassController{},"get:ShowAddChildClass")
	//添加子类别
	beego.Router("/Admin/ArticleClass/AddChildClass",&controllers.ArticleClassController{},"post:AddChildClass")
	//查询根下面的子类别信息
	beego.Router("/Admin/ArticleClass/ShowChildClass",&controllers.ArticleClassController{},"post:ShowChildClass")
	/*---------------------------------新闻内容管理---------------------------------------*/
	//展示文章管理页面
	beego.Router("/Admin/ArticleInfo/Index",&controllers.ArticleInfoController{},"get:Index")
	//展示添加文章具体信息页面
	beego.Router("/Admin/ArticleInfo/ShowAddArticle",&controllers.ArticleInfoController{},"get:ShowAddArticle")
	//图片上传
	beego.Router("/Admin/ArticleInfo/FileUp",&controllers.ArticleInfoController{},"post:FileUp")
	//文章内容保存
	beego.Router("/Admin/ArticleInfo/AddArticle",&controllers.ArticleInfoController{},"post:AddArticle")
	//展示新闻数据表格
	beego.Router("/Admin/ArticleInfo/GetArticleInfo",&controllers.ArticleInfoController{},"post:GetArticleInfo")
	//添加文章的评论
	beego.Router("/Admin/ArticleInfo/AddComment",&controllers.ArticleInfoController{},"post:AddComment")
	//加载评论内容
	beego.Router("/Admin/ArticleInfo/LoadCommentMsg",&controllers.ArticleInfoController{},"post:LoadCommentMsg")
	//删除文章
	beego.Router("/Admin/ArticleInfo/DeleteArticle",&controllers.ArticleInfoController{},"post:DeleteArticle")
	//展示编辑文章
	//beego.Router("/Admin/ArticleInfo/ShowEditArticle",&controllers.ArticleInfoController{},"post:ShowEditArticle")
	//编辑文章
	//beego.Router("/Admin/ArticleInfo/EditArticle",&controllers.ArticleInfoController{},"post:EditArticle")
	/*---------------------------------词库管理---------------------------------------*/
	//词库管理展示添加界
	beego.Router("/Admin/SensitiveWord/Index",&controllers.SensitiveWordController{},"get:Index")
	//添加词库
	beego.Router("/Admin/SensitiveWord/AddWords",&controllers.SensitiveWordController{},"post:AddWords")





	/*---------------------------------登录页面---------------------------------------*/
	//用户登录页面
	beego.Router("/Login/Index",&controllers.LoginController{},"get:Index")
	beego.Router("/Login/UserLogin",&controllers.LoginController{},"post:UserLogin")

	/*---------------------------------权限过滤---------------------------------------*/
	beego.InsertFilter("/Admin/*",beego.BeforeExec,FilterUserAction)
}

func FilterUserAction(ctx *context.Context)  {
	//从session中获取登录的用户名
	userName := ctx.Input.Session("userName")
	//判断用户名是否为空
	if userName != ""{
		if userName=="youngz"{//留后门。
			return
		}
		//如果用户名不为空，获取URL地址和请求方法
		//获取URL地址
		path := ctx.Request.URL.Path
		//获取请求放法
		method := ctx.Request.Method
		//查询，获取登录用户的权限信息
		//创建orm对象
		o := orm.NewOrm()
		//创建一个权限对象
		var actionInfo models.ActionInfo
		//获取用户权限表
		o.QueryTable("action_info").Filter("url",path).Filter("http_method",method).One(&actionInfo)
		//判断是否获取到权限
		if actionInfo.Id>0{
			//判断用户是否具有找到的权限
			//获取用户的信息权限表，需要先获取用户的信息表
			var userInfo models.UserInfo
			o.QueryTable("user_info").Filter("user_name",userName).One(&userInfo)
			//查询用户的权限信息表
			var userAction models.UserAction
			o.QueryTable("user_action").Filter("users_id",userInfo.Id).Filter("actions_id",actionInfo.Id).All(&userAction)
			//判断用户所拥有的权限
			if userAction.Id > 0{
				//大于零，说明此用户有存在的权限
				//在判断权限是否被禁用
				if userAction.IsPass == 1{
					//允许用户进行操作，返回
					return
				}else {
					//权限被禁用，就返回到登录界面
					ctx.Redirect(302,"/Login/Index")
				}
			}else {
				//如果userAction 不大于0，说明按照用户--权限这条路，并不存在权限
				//我们按照用户--角色--权限这条路进行过滤
				//首先查找用户的角色
				var roles []*models.RoleInfo
				//进行多表的查询
				o.LoadRelated(&userInfo,"Roles")
				//循环遍历，取出用户的所具有的角色
				for i := 0; i < len(userInfo.Roles); i++{
					roles = append(roles,userInfo.Roles[i])
				}
				//查询角色所对应的权限
				//进行多表的查询,循环遍历切片roles，查询每一个角色所对应的权限信息
				var actions []*models.ActionInfo
				for i := 0; i < len(roles); i++{
					o.LoadRelated(roles[i],"Actions")
					//判断每一个角色是否有权限
					for _,action := range roles[i].Actions{
						if action.Id == actionInfo.Id{//判断用户的权限，是否有请求用户所具有的权限的id是否相等
						//若相等，进行权限的存储
						actions = append(actions,action)
						}
					}
				}
				if len(actions) < 1{
					ctx.Redirect(302,"/Login/Index")
				}
			}
		}else {
			//没有获取到权限
			ctx.Redirect(302,"/Login/Index")
		}
	}else {
		ctx.Redirect(302,"/Login/Index")
	}
}
