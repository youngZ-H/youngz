package routers

import (
	"class3/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	//beego.Router("/", &controllers.MainController{})
	//注册过滤器函数路由
	beego.InsertFilter("/Article/*",beego.BeforeRouter,FilterFunc)

	//注册界面路由
	beego.Router("/register", &controllers.RegController{}, "get:ShowReg;post:HandleReg")
	//登陆界面路由
	beego.Router("/", &controllers.LoginController{}, "get:ShowLogin;post:HandleLogin")
	//展示列表界面路由
	beego.Router("/Article/ShowArticle", &controllers.ArticleController{}, "get:ShowArticle;post:HandleSelect")
	//添加文章界面路由
	beego.Router("/Article/AddArticle", &controllers.ArticleController{}, "get:ShowAddArticle;post:HandleAddArticle")
	//展示文章详情界面路由
	beego.Router("/Article/ShowContent",&controllers.ArticleController{},"get:ShowContent")
	//创建删除的路由界面
	beego.Router("/Article/DeleteArticle",&controllers.ArticleController{},"get:ShowDelete")
	//进行更新操作
	beego.Router("/Article/UpdateArticle",&controllers.ArticleController{},"get:ShowUpdate;post:HandleUpdate")
	//进行文章类型的添加操作
	beego.Router("/Article/AddArticleType",&controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")
	//进行文章类型的删除操作
	beego.Router("/Article/DeleteArticleType",&controllers.ArticleController{},"get:ShowDeleteType")
	//退出界面的操作
	beego.Router("/Logout",&controllers.ArticleController{},"get:ShowLogout")
}
//定义过滤器函数
var FilterFunc = func(ctx *context.Context){
	userName := ctx.Input.Session("userName")
	if userName == nil{
		ctx.Redirect(302,"/")
	}
	//如果不为空，继续往下执行
}