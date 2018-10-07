package controllers

import (
	"github.com/astaxie/beego"
	"itcastCms/models"
	"time"
	"github.com/astaxie/beego/orm"
)

type ArticleClassController struct {
	beego.Controller
}
//展示添加文章类别的视图界面
func (this *ArticleClassController)Index()  {
	this.TplName = "ArticleClass/Index.html"
}
//展示添加文章父类别方法的实现
func (this *ArticleClassController)ShowParentClass()  {
	this.TplName = "ArticleClass/ShowAddParentClass.html"
}
//添加文章父类别方法的实现
func (this *ArticleClassController)AddParentClass()  {
	//从前端读取到数据，存储到数据库中
	//定义一个对象
	var ParentClass models.ArticleClass
	//进行数据的赋值
	ParentClass.DelFlag = 0
	ParentClass.Remark = this.GetString("remark")
	ParentClass.ClassName = this.GetString("className")
	ParentClass.CreateUserId = this.GetSession("userId").(int)
	ParentClass.CreateDate = time.Now()
	ParentClass.ParentId = 0
	//创建orm对象
	o := orm.NewOrm()
	_,err := o.Insert(&ParentClass)
	if err == nil{
		this.Data["json"] = map[string]interface{}{"flag":"ok"}
	}else {
		this.Data["json"] = map[string]interface{}{"flag":"no"}
	}
	this.ServeJSON()
}
//展示添加的新闻类别方法的实现
func (this *ArticleClassController)ShowArticleClass()  {
	//创建orm对象
	o := orm.NewOrm()
	//定义一个展示文章类别的对象
	var articleClasses []models.ArticleClass
	//查询数据
	o.QueryTable("article_class").Filter("parent_id",0).Filter("del_flag",0).All(&articleClasses)
	//将数据通过ajax的形式发送到前端进行数据的展示
	this.Data["json"] = map[string]interface{}{"rows":articleClasses}
	this.ServeJSON()
}
//展示添加子类别的方法
func (this *ArticleClassController)ShowAddChildClass()  {
	//从前端获取所选中的父类别的id号
	id,err:=this.GetInt("cId")
	if err != nil{
		beego.Info("展示新闻子类别获取父类行号时发生错误",err)
		return
	}
	//创建orm对象，获取对应的父类别的信息
	o := orm.NewOrm()
	//定义一个新闻类别的对象，用来存储父类的信息
	var classInfo models.ArticleClass
	o.QueryTable("article_class").Filter("id",id).One(&classInfo)
	//将获取到的父类信息发送到添加子类别的展示界面
	this.Data["classInfo"] = classInfo
	this.TplName = "ArticleClass/ShowAddChildClass.html"
}
//完成子类别方法的添加
func (this *ArticleClassController)AddChildClass()  {
	var classInfo=models.ArticleClass{}
	classInfo.CreateUserId=this.GetSession("userId").(int)
	classInfo.CreateDate=time.Now()
	classInfo.ParentId,_=this.GetInt("classId")
	classInfo.Remark=this.GetString("remark")
	classInfo.ClassName=this.GetString("className")
	o:=orm.NewOrm()
	_,err:=o.Insert(&classInfo)
	if err==nil{
		this.Data["json"]=map[string]interface{}{"flag":"ok"}
	}else{
		this.Data["json"]=map[string]interface{}{"flag":"no"}
	}
	this.ServeJSON()
}
//查询根下面的子类别信息。
func (this *ArticleClassController)ShowChildClass()  {
	cid,_:=this.GetInt("id")//根的编号。
	o:=orm.NewOrm()
	var articleClasses[]models.ArticleClass
	o.QueryTable("article_class").Filter("parent_id",cid).All(&articleClasses)
	this.Data["json"]=map[string]interface{}{"rows":articleClasses}
	this.ServeJSON()
}






















