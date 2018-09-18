package controllers

import (
	"github.com/astaxie/beego"
	"path"
	"time"
	"class3/models"
	"github.com/astaxie/beego/orm"
	"strconv"
	"math"
)

//定义展示界面的结构体

type ArticleController struct {
	beego.Controller
}
//绑定处理展示列表下拉框的post方法
func (this*ArticleController)HandleSelect()  {
	//1.获取下拉列表的数据
	typeName := this.GetString("select")
	//将数据打印出来进行测试
	//beego.Info(typeName)
	//2.判断数据
	if typeName == ""{
		beego.Info("获取下拉数据为空")
		return
	}
	//3.进行数据的处理
		//1.获取orm对象
		 o := orm.NewOrm()
		 //创建一个结构体切片进行数据的存储
		 var articles []models.Article
		 //2.获取要查询的对象
		 o.QueryTable("Article").RelatedSel("ArticleStype").Filter("ArticleStype__TypeName",typeName).All(&articles)
		// o.QueryTable("Article").RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)
		 //3.测试打印一下数据，是否插入数据成功
		//beego.Info(articles)

}





//为展示界面绑定get方法，加载展示界面
func (this*ArticleController)ShowArticle()  {
	////判断用户是否登陆，决定是否允许访问
	//userName := this.GetSession("userName")
	////判断用户名是否为空
	//if userName == nil{
	//	this.Redirect("/",302)
	//	return
	//}






	//展示列表页面
	//1.获取orm对象
	o := orm.NewOrm()
	//2.获取查询对象
	qs := o.QueryTable("Article")
	//3.定义一个结构体切片，用来存储查询的数据
	var articles []models.Article
	qs.All(&articles)
	//beego.Info(articles[1])
	//进行页数的显示

	count,_ := qs.RelatedSel("ArticleStype").Count()  //.Filter("ArticleStype__TypeName",typeName)
	//获取pageIndex

	pageIndex := this.GetString("pageIndex")
	pageIndex1,err:= strconv.Atoi(pageIndex)
	if err != nil{
		//beego.Info("获取页码失败",err)
		pageIndex1 = 1
	}
	//获取总页数
	pageSize := 2

	start := pageSize*(pageIndex1-1)
	qs.Limit(pageSize, start).RelatedSel("ArticleStype").All(&articles)

	pageCount := math.Ceil(float64(count)/float64(pageSize))
	var FirstPage = false //是否是首页
	var EndPage = false //是否是末页
	if pageIndex1 == 1{
		FirstPage = true
	}
	if float64(pageIndex1) == pageCount{
		EndPage = true
	}

	//文章列表类 文章类型选择的添加
	//1.获取orm对象，在前面获取文章展示列表时，已经获取
	//2.获取要操作的文章类型对象
	//定义一个切片，进行文章类型数据的存储
	var types []models.ArticleStype
	o.QueryTable("ArticleStype").All(&types)

		//根据文章的类型来显示数据
		//1.获取下拉列表的数据
		typeName := this.GetString("select")
		//将数据打印出来进行测试
		//beego.Info(typeName)
		//2.判断数据
		var articlewithtypes []models.Article

		if typeName == ""{
			qs.Limit(pageSize, start).RelatedSel("ArticleStype").All(&articlewithtypes)
		}else {
			o.QueryTable("Article").RelatedSel("ArticleStype").Filter("ArticleStype__TypeName",typeName).All(&articlewithtypes)

		}
		//3.进行数据的处理
		//1.获取orm对象
		//o := orm.NewOrm()
		////创建一个结构体切片进行数据的存储
		//var articles []models.Article
		////2.获取要查询的对象
		//o.QueryTable("Article").RelatedSel("ArticleStype").Filter("ArticleStype__TypeName",typeName).All(&articles)
		//// o.QueryTable("Article").RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)
		////3.测试打印一下数据，是否插入数据成功
		////beego.Info(articles)
	userName := this.GetSession("userName")
	this.Data["userName"] = userName

	this.Data["typeName"] = typeName
	//3.将数据传递给视图
	this.Data["types"] = types
	//将数据传递给视图
	this.Data["count"] = count
	this.Data["pageCount"] = pageCount
	this.Data["pageIndex"] = pageIndex1
	//4.将展示的数据传到前端的显示界面
	this.Data["articles"] = articlewithtypes
	this.Data["FirstPage"] = FirstPage
	this.Data["EndPage"] = EndPage

	this.Layout = "layout.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["contentHead"] = "head1.html"
	this.TplName = "index.html"

}

//为展示界面绑定get方法，用来展示加载文章的界面的
func (this*ArticleController)ShowAddArticle()  {
	//进行添加文章界面的文章类型，下拉框的从操作
	//1.获取orm对象，前面已经获取
	o := orm.NewOrm()
	//2.获取文章类型对象，并定义一个文章类型的结构体切片
	var types []models.ArticleStype
	o.QueryTable("ArticleStype").All(&types)
	//3.将数据传递给视图
	this.Data["types"] = types
	this.Layout = "layout.html"
	this.TplName = "add.html"
}
//为i添加文章界面绑定post方法，用来进行数据的存储
func (this*ArticleController)HandleAddArticle()  {
	//1.获取用户输入的文章标题和文章内容
	articleName := this.GetString("articleName")
	articleContent := this.GetString("content")

	//2.获取文章上传的文件
	f, h, err := this.GetFile("uploadname")
	if err != nil{
		beego.Info("文件上传错误1")
		return
	}
	defer f.Close()
	//获取上传文件的信息
		//1.判断文件格式是否争取，提取后追信息
	ext := path.Ext(h.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jepg"{
		beego.Info("长传文件格式不正确")
		return
	}
		//2.判断文件的大小
	if h.Size > 500000{
		beego.Info("文件太大，无法上传")
		return
	}
		//3.进行文件的存储
	fileName := time.Now().Format("2006-01-02 15:04:05")
	err = this.SaveToFile("uploadname",fileName + ext)
	if err != nil{
		beego.Info("文件上传错误2")
		return
	}
	//this.Ctx.WriteString("上传文件成功")
	//beego.Info(articleName,articleContent)
	//beego.Info(ext)
	//beego.Info(h.Filename)
	//向数据库中存储数据的信息
	//1.创建orm对象
	o := orm.NewOrm()
	//2.获取插入对象
	article := models.Article{}
	//3.将从浏览器读取到的数据存储到数据库对应的变量上
	article.Title = articleName
	article.Content = articleContent
	article.Img = fileName+ext

	//进行文章类型的表单提交处理
	//1.先从视图中获取数据
	typeName := this.GetString("select")
	//2.判断数据是否为空
	if typeName == ""{
		beego.Info("获取添加文章界面中的文章类型出错", err)
		return
	}
	//3.获取要插入数据的对象
	var articleStype models.ArticleStype
	articleStype.TypeName = typeName
	err = o.Read(&articleStype,"TypeName")
	if err != nil{
		beego.Info("读取数据类型出错", err)
		return
	}
	article.ArticleStype = &articleStype





	//4.插入数据
	_, err = o.Insert(&article)
	if err != nil {
		beego.Info("向article表中插入数据时发生错误")
		beego.Info(err)
		return
	}

	//5.数据插入成功后，跳转到add.html界面
	this.Layout = "layout.html"
	this.Redirect("/Article/ShowArticle",302)

}
func (this*ArticleController)ShowContent()  {
	//1.获取数据的id
	id := this.GetString("id")
	beego.Info(id)
	//2.查询数据
	id2,_:= strconv.Atoi(id)
		//1.获取orm对象
		o :=  orm.NewOrm()
		//2.获取查询对象
		article := models.Article{Id:id2}
		err := o.Read(&article)
		if err != nil{
			beego.Info("查询数据为空")
			beego.Info(err)
			return
		}

		article.Count++
		_, err = o.Update(&article)
		if err != nil{
			beego.Info("更新数据失败")
			beego.Info(err)
			return
		}
		//多表操作
		//1.获取操作的表对象
		//article := models.Article{Id:id2}
		//获取多对多操作对象
		m2m := o.QueryM2M(&article,"Users")
		//获取插入对象
		userName := this.GetSession("userName")
		user := models.User{}
		user.UserName = userName.(string)
		o.Read(&user,"UserName")
		//插入多对多对象
		_, err = m2m.Add(&user)
		if err != nil{
			beego.Info("插入多对多数据是出现错误", err)
			return
		}
		var users []models.User
		//多对多对象查询
		//o.LoadRelated(&article, "Users")
		//多表的去重查询
		o.QueryTable("User").Filter("Article__Article__Id", id2).Distinct().All(&users)
		//o.QueryTable("User").Filter("Article__Article__Id2",id2).Distinct().All(&users)
	this.Data["users"] = users

	//3.将查询数据显示在前端界面
	this.Data["article"] = article


	this.Layout = "layout.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["contentHead"] = "head.html"
	this.TplName = "content.html"
	//this.Redirect("/Article/ShowContent",302)
}

//绑定删除的方法
func (this*ArticleController)ShowDelete()  {
	//1.获取查询对象
	id := this.GetString("id")
	id2, _ := strconv.Atoi(id)
		//1.获取orm对象
		o := orm.NewOrm()
		//2.获取查询对象
		article := models.Article{Id:id2}
	//2.删除选中的对象
	_, err := o.Delete(&article)
	if err != nil{
		beego.Info("删除数据失败")
		beego.Info(err)
		return
	}
	//3.进行界面的重定向
	this.Redirect("/Article/ShowArticle",302)

}
//绑定展示函数更新的方法
func (this*ArticleController)ShowUpdate()  {
	//获取数据
	id := this.GetString("id")
	//对数据进行判断
	beego.Info("###"+id)
	if id == ""{
		beego.Info("连接错误")
		return
	}
	//获取orm对象
	o := orm.NewOrm()
	//获取查询对象
	//进行类型的转换
	id2,_ := strconv.Atoi(id)
	article := models.Article{Id:id2}
	//读取数据
	err := o.Read(&article)
	if err != nil{
		beego.Info("显示更新数据失败")
		beego.Info(err)
		return
	}
	//把数据传递给视图
	this.Data["article"] = article
	//this.Layout = "layout.html"
	this.TplName = "update.html"

}



//绑定函数更新的方法
func (this*ArticleController)HandleUpdate()  {
	//1.拿取数据
	id,_ := this.GetInt("id")
	name := this.GetString("articleName")
	content := this.GetString("content")
	//2.进行数据的判断
	if name == "" || content == ""{
		beego.Info("更新数据为空")
		return
	}
	f,h,err := this.GetFile("uploadname")
	defer f.Close()
	if err != nil{
		beego.Info("获取文件内容失败")
		beego.Info(err)
	}
	//判断文件的格式是否正确
	ext := path.Ext(h.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jepg"{
		beego.Info("上传的文件格式不对")
		return
	}
	//判断文件的大小
	if h.Size > 500000{
		beego.Info("上传的文件过大")
		return
	}
	//保证上传的文件不重名
	fileName := time.Now().Format("2006-01-02 15:04:05")
	//进行文件的保存
	this.SaveToFile("uploadname",fileName + ext)
	//获取orm对象
	o := orm.NewOrm()
	//获取要查询的对象
	article := models.Article{Id:id}
	//读取要更新的数据
	err = o.Read(&article)
	if err != nil{
		beego.Info("读取数据失败")
		beego.Info(err)
	}
	//对要更新的数据，进行重新的赋值
	article.Title = name
	article.Content = content
	article.Img = fileName + ext
	//更新到数据库
	_, err = o.Update(&article)
	if err != nil{
		beego.Info("更新数据失败")
		beego.Info(err)
	}
	//更新完成后，进行画面的跳转
	this.Redirect("/Article/ShowArticle",302)
}
//进行文章类型的展示方法绑定
func (this*ArticleController)ShowAddType()  {
	//1.查询数据
	//获取orm对象
	o := orm.NewOrm()
	//创建一个数组进行文章类型的查询
	var articleStypes []models.ArticleStype
	_, err:=o.QueryTable("ArticleStype").All(&articleStypes)
	if err != nil{
		beego.Info("查询数据类型错误", err)
	}
	this.Data["ArticleStype"] = articleStypes
	this.Layout = "layout.html"
	this.TplName = "addType.html"
}
//进行文章类型的添加的方法绑定
func (this*ArticleController)HandleAddType()  {
	//从前端获取输入的文章类型数据
	typeName := this.GetString("typeName")
	//进行数据的判断
	if typeName == ""{
		beego.Info("查询的数据类型为空")
	}
	//获取orm对象
	o := orm.NewOrm()
	//获取要插入的数据的对象
	//articleStype := models.ArticleStype{}
	var articleStype models.ArticleStype
	//进行数据的赋值
	articleStype.TypeName = typeName
	//将数据插入数据库
	_, err := o.Insert(&articleStype)
	if err != nil{
		beego.Info("向数据库中插入文章类型出错", err)
		return
	}
	//进行视图的重定向
	this.Redirect("/Article/AddArticleType",302)
}

//进行文章类型的删除绑定操作
func (this*ArticleController)ShowDeleteType()  {
	//获取要删除的数据信息
	id := this.GetString("id")
	//进行数据类型的转换
	id2,_ := strconv.Atoi(id)
	//获取orm操作对象
	o := orm.NewOrm()
	//获取要删除的对象
	articleStype := models.ArticleStype{Id:id2}
	//执行删除操作
	_, err := o.Delete(&articleStype)
	if err != nil{
		beego.Info("删除文章类型出错", err)
		return
	}
	//进行重定向
	this.Redirect("/Article/AddArticleType",302)

}

//推出界面的函数方法绑定
func (this*ArticleController)ShowLogout()  {
	//获取session
	this.DelSession("userName")
	//删除session的状态，进行界面的重定向
	this.Redirect("/",302)
}





















