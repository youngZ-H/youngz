package controllers

import (
	"github.com/astaxie/beego"
	"strconv"
	"time"
	"os"
	"path"
	"github.com/astaxie/beego/orm"
	"itcastCms/models"
	"io/ioutil"
	"strings"
	"math"
	"html"
	"regexp"
)
//创建文章内容类
type ArticleInfoController struct {
	beego.Controller
}
//为文章内容类绑定展示视图的方法
func (this *ArticleInfoController)Index()  {
	this.TplName = "ArticleInfo/Index.html"
}
//展示添加文章的详细信息的界面
func (this *ArticleInfoController)ShowAddArticle()  {
	//展示新闻类别的单选框
	o:=orm.NewOrm()
	var articleClass []models.ArticleClass
	o.QueryTable("article_class").Filter("parent_id__gte",0).All(&articleClass)
	this.Data["classInfo"]=articleClass
	this.TplName = "ArticleInfo/ShowAddArticle.html"
}
//文章里面图片的上传，将图片加载到编辑器里面
func (this *ArticleInfoController)FileUp()  {
	//从前端获取文件
	f,h,err := this.GetFile("fileUp")
	defer f.Close()
	//进行错误的判断
	if err != nil{
		//先前端返回错误信息
		this.Data["json"] = map[string]interface{}{"flag":"no","msg":"读取文件时发生错误"}
	}
	//对文件的类型进行判断
	//获取文件的后缀名
	fileExt := path.Ext(h.Filename)
	if fileExt == ".jpg" || fileExt == ".png" || fileExt == ".jpeg"{
		//文件类型满足条件，进行文件的大小判断
		if h.Size < 5000000{
			//都满足条件，进行文件的存储
			//定义存储的路径名
			dir := "./static/fileUp"+"/"+strconv.Itoa(time.Now().Year())+"/"+time.Now().Month().String()+"/"+strconv.Itoa(time.Now().Day())+"/"
			//判断存储之前，文件夹是否存在
			_, err := os.Stat(dir)
			//进行判断
			if err != nil{
				//如果有错误，表示文件夹不存在，进行文件夹的重新创建
				os.MkdirAll(dir, os.ModePerm) //用此方式创建的文件夹具有读写和执行的全部权限
			}
			//如果文件夹已经存在，进行文件的重命名
			newFileName:=strconv.Itoa(time.Now().Year())+time.Now().Month().String()+strconv.Itoa(time.Now().Day())+strconv.Itoa(time.Now().Hour())+strconv.Itoa(time.Now().Minute())+strconv.Itoa(time.Now().Nanosecond())
			//构建一个完整的文件路径
			fullDir := dir + newFileName + fileExt

			//进行文件的存储
			err = this.SaveToFile("fileUp",fullDir)
			if  err != nil{
				this.Data["json"] = map[string]interface{}{"flag":"no","msg":"文件存储失败"}
			}else {
				this.Data["json"] = map[string]interface{}{"flag":"ok","msg":fullDir}
			}
		}else {
			this.Data["json"] = map[string]interface{}{"flag":"no","msg":"上传的文件太大"}
		}
	}else {
		//如果不是以上的文件类型，代表上传的文件类型错误，返回错误
		this.Data["json"] = map[string]interface{}{"flag":"no","msg":"上传的文件类型错误"}
	}
	this.ServeJSON()
}
//文章详细信息的内容保存，即将数据插入到数据库中
func (this *ArticleInfoController)AddArticle()  {
	var articleInfo=models.ArticleInfo{}
	articleInfo.DelFlag=0
	articleInfo.AddDate=time.Now()
	articleInfo.ArticleContent=this.GetString("ArticleContent")
	articleInfo.Origin=this.GetString("Origin")
	articleInfo.Title=this.GetString("title")
	articleInfo.ModifyDate=time.Now()
	articleInfo.FullTitle=this.GetString("Fulltitle")
	articleInfo.Author=this.GetString("Author")
	articleInfo.Intro=this.GetString("Intro")
	articleInfo.KeyWords=this.GetString("KeyWords")
	articleInfo.PhotoUrl=this.GetString("PhotoUrl")
	o:=orm.NewOrm()
	num, err := o.Insert(&articleInfo)//插入文章信息,返回的整型表示的含义是插入的记录的主键编号。
	if err != nil{
		beego.Info("向数据库中插入文章数据时发生错误",err)
	}
	//查询类别信息
	var classInfo models.ArticleClass
	id,_ := this.GetInt("classId")
	o.QueryTable("article_class").Filter("id",id).One(&classInfo)
	//创建M2M对象
	m2m := o.QueryM2M(&articleInfo,"ArticleClasses")
	m2m.Add(classInfo)
	//创建静态页面
	CreateStaticPage(int(num))
	this.Data["json"]=map[string]interface{}{"flag":"ok"}
	this.ServeJSON()
}
//生成静态文章页面
func CreateStaticPage(cId int)  {
	//从存储文章的数据库中读取数据
	o := orm.NewOrm()
	var articleInfo models.ArticleInfo
	o.QueryTable("article_info").Filter("id",cId).One(&articleInfo)
	//读取模板文件，替换其中的占位符，生成静态文件
	dir := "./static/ArticelTemplate/ArticelTemplateInfo.html"  //指定模板文件的路径
	fs, err := os.Open(dir)
	defer fs.Close()
	if err != nil{
		beego.Info("打开模板文件时发生错误",err)
	}else {
		//读取模板文件，替换掉其中的占位符
		content, err := ioutil.ReadAll(fs)
		if err != nil{
			beego.Info("读取模板文件发生错误",err)
		}else {
			//替换掉其中的占位符
			fileContent := string(content)//将字节转换成字符串
			fileContent=strings.Replace(fileContent,"$Title",articleInfo.Title,-1)
			fileContent=strings.Replace(fileContent,"$Origin",articleInfo.Origin,-1)
			fileContent=strings.Replace(fileContent,"$AddDate",articleInfo.AddDate.Format("2006-01-02"),-1)
			fileContent=strings.Replace(fileContent,"$ArticleContent",articleInfo.ArticleContent,-1)
			fileContent=strings.Replace(fileContent,"$articleId",strconv.Itoa(articleInfo.Id),-1)

			//3：创建目录(根据添加新闻的日期时间，来创建文件夹。)
			month:=articleInfo.AddDate.Month().String()//获取到了添加日期中月份。默认是英文
			//将英文的月份转换成数字月份。
			var m int
			for i:=0;i<len(months) ;i++  {
				if months[i]==month{
					m=i+1
					break
				}
			}
			var dirPath string  //构建路径
			day := articleInfo.AddDate.Day()
			var d string
			if day < 10{
				dd := strconv.Itoa(day)
				d = "0" + dd
			}else {
				d = strconv.Itoa(day)
			}
			if m<10 {
				n:="0"+strconv.Itoa(m)
				dirPath="./static/Article/"+strconv.Itoa(articleInfo.AddDate.Year())+"/"+n+"/"+d+"/"
			}else{
				dirPath="./static/Article/"+strconv.Itoa(articleInfo.AddDate.Year())+"/"+strconv.Itoa(m)+"/"+d+"/"
			}
			//判断文件夹是否存在
			_,err := os.Stat(dirPath)
			if err != nil{
				//如果文件夹不存在，就创建新的文件夹
				os.MkdirAll(dirPath,os.ModePerm)
			}
			//创建完整的文件路径
			fullDir := dirPath + strconv.Itoa(cId) + ".html"
			//写入文件
			err = ioutil.WriteFile(fullDir,[]byte(fileContent),0644)
			if err != nil{
				beego.Info("写入文件失败",err)
			}
		}
	}
}
var months = []string{
	"January",
	"February",
	"March",
	"April",
	"May",
	"June",
	"July",
	"August",
	"September",
	"October",
	"November",
	"December",
}
//展示新闻数据表格
func (this *ArticleInfoController)GetArticleInfo()  {
	PageIndex,_ := this.GetInt("page")
	PageSize,_ := this.GetInt("rows")
	start := (PageIndex-1)*PageSize
	var articleInfo []models.ArticleInfo
	o := orm.NewOrm()
	o.QueryTable("article_info").Filter("del_flag",0).OrderBy("id").Limit(PageSize,start).All(&articleInfo)
	count,_:= o.QueryTable("article_info").Filter("del_flag",0).Count()
	//将数据传递给前端界面进行展示
	this.Data["json"] = map[string]interface{}{"rows":articleInfo,"total":count}
	this.ServeJSON()
}
//添加文章评论
func (this *ArticleInfoController)AddComment(){
	msg := html.EscapeString(this.GetString("msg"))
	//调用禁用词的函数，进行匹配输入的评论是否含有禁用词
	if ForbindWord(msg){
		this.Data["json"] = map[string]interface{}{"flag":"no","msg":"输入的评论中含有禁用词"}
	}else if CheckModWord(msg) {
		var articleComment models.ArticleComment
		articleComment.AddDate = time.Now()
		//articleComment.Msg = this.GetString("msg")
		articleComment.Msg = msg
		articleComment.IsPass = 0
		articleId,err := this.GetInt("articleId")
		if err != nil{
			beego.Info("存储文章评论时，获取文章的id时发生错误",err)
		}

		var article models.ArticleInfo
		o := orm.NewOrm()
		o.QueryTable("article_info").Filter("id",articleId).One(&article)
		articleComment.Article = &article
		_,err = o.Insert(&articleComment)
		if err == nil{
			this.Data["json"] = map[string]interface{}{"flag":"ok"}
		}else {
			this.Data["json"] = map[string]interface{}{"flag":"no"}
		}
		this.Data["json"] = map[string]interface{}{"flag":"no","msg":"输入的评论中含有审查词"}
	}else if CheckReplaceWord(msg){
		comments := ReplaceWord(msg)
		var articleComment models.ArticleComment
		articleComment.AddDate = time.Now()
		//articleComment.Msg = this.GetString("msg")
		articleComment.Msg = comments
		articleComment.IsPass = 1
		articleId,err := this.GetInt("articleId")
		if err != nil{
			beego.Info("存储文章评论时，获取文章的id时发生错误",err)
		}

		var article models.ArticleInfo
		o := orm.NewOrm()
		o.QueryTable("article_info").Filter("id",articleId).One(&article)
		articleComment.Article = &article
		_,err = o.Insert(&articleComment)
		if err == nil{
			this.Data["json"] = map[string]interface{}{"flag":"ok"}
		}else {
			this.Data["json"] = map[string]interface{}{"flag":"no"}
		}
	}else {
		var articleComment models.ArticleComment
		articleComment.AddDate = time.Now()
		//articleComment.Msg = this.GetString("msg")
		articleComment.Msg = msg
		articleComment.IsPass = 1
		articleId,err := this.GetInt("articleId")
		if err != nil{
			beego.Info("存储文章评论时，获取文章的id时发生错误",err)
		}

		var article models.ArticleInfo
		o := orm.NewOrm()
		o.QueryTable("article_info").Filter("id",articleId).One(&article)
		articleComment.Article = &article
		_,err = o.Insert(&articleComment)
		if err == nil{
			this.Data["json"] = map[string]interface{}{"flag":"ok"}
		}else {
			this.Data["json"] = map[string]interface{}{"flag":"no"}
		}
	}

	this.ServeJSON()

}
//加载文章评论内容的方法的实现
func (this *ArticleInfoController)LoadCommentMsg()  {
	//从前端端读取要加载的文章内容
	articleId,_ := this.GetInt("articleId")
	//进行评论的分页
	//从前端获取当前的页数
	PageIndex,_ := this.GetInt("PageIndex")
	//指定每页显示的数据的大小
	PageSize := 3
	//求出文章评论的总的记录数
	o := orm.NewOrm()
	recordCount,_ := o.QueryTable("article_comment").Filter("article_id",articleId).Filter("is_pass",1).Count()
	//求出总的页数
	pageCount := int(math.Ceil(float64(recordCount)/float64(PageSize)))
	//对传递过来的当前页数进行校验
	if PageIndex < 1 {
		PageIndex = 1
	}
	if PageIndex > pageCount{
		PageIndex = pageCount
	}
	//求出起始页
	start := (PageIndex-1)*PageSize
	//从数据库中查询出文章的每一评论页的评论内容
	var articleComments []models.ArticleComment
	o.QueryTable("article_comment").Filter("article_id",articleId).Filter("is_pass",1).Limit(PageSize,start).All(&articleComments)
	pageBar := CreatePageBar(PageIndex,pageCount)
	this.Data["json"] = map[string]interface{}{"msg":articleComments,"pageBar":pageBar}
	this.ServeJSON()
}
//创建页码条
func CreatePageBar(PageIndex,PageCount int)(strHtml string)  {
	//判断
	if PageCount == 1{
		return ""
	}
	start := PageIndex-5
	if start < 1 {
		start = 1
	}
	end := start + 9
	if end >PageCount {
		end = PageCount
	}
	if PageIndex > 1{
		strHtml = "<a class='pageBarLink' href='/Admin/ArticleInfo/LoadCommentMsg?PageIndex="+strconv.Itoa(PageIndex-1)+"'>上一页</a>"
	}
	for i:=start;i<=end;i++{
		if PageIndex == i{
			strHtml = strHtml + strconv.Itoa(i)
		}else {
			strHtml = strHtml + "<a class='pageBarLink' href='/Admin/ArticleInfo/LoadCommentMsg?PageIndex="+strconv.Itoa(i)+"'>"+strconv.Itoa(i)+"</a>"

		}

	}
	if PageIndex < PageCount{
		strHtml = strHtml + "<a class='pageBarLink' href='/Admin/ArticleInfo/LoadCommentMsg?PageIndex="+strconv.Itoa(PageIndex+1)+"'>下一页</a>"
	}
	return strHtml
}
//禁用词
func ForbindWord(msg string)(b bool)  {
	//1:查询表中的禁用词。
	o:=orm.NewOrm()
	var forWords[]models.SensitiveWord
	o.QueryTable("sensitive_word").Filter("is_forbid",1).All(&forWords)
	//2:进行过滤
	// "词的名称1|词名称2|词名称3"
	var words[]string
	for i:=0;i<len(forWords);i++{
		words=append(words,forWords[i].WordPattern)
	}
	str:=strings.Join(words,"|")
	if str == ""{
		return false
	}
	reg:=regexp.MustCompile(str)
	return reg.MatchString(msg)
}
//审查词
func CheckModWord(msg string)(b bool)  {
	//将数据库中所有的禁用词全部查询出来，存储到一个切片中
	o := orm.NewOrm()
	//声明一个禁用词的对象，用来存储从数据库中查询出来的词
	var checkWords []models.SensitiveWord
	o.QueryTable("sensitive_word").Filter("is_mod",1).All(&checkWords)
	//循环遍历切片，取出所有的禁用词
	//定义一个切片用来存储禁用词
	var words []string
	for i := 0; i < len(checkWords); i++{
		words = append(words,checkWords[i].WordPattern)
	}
	//再将切片转换为字符串
	strs := strings.Join(words,"|")
	if strs == ""{
		return false
	}

	//创建正则表达式的规则
	reg := regexp.MustCompile(strs)
	//进行正则表达式规则的匹配

	return reg.MatchString(msg)
}
//替换词判断
func CheckReplaceWord(msg string)(b bool)  {
	//将数据库中所有的禁用词全部查询出来，存储到一个切片中
	o := orm.NewOrm()
	//声明一个禁用词的对象，用来存储从数据库中查询出来的词
	var replaceWords []models.SensitiveWord
	o.QueryTable("sensitive_word").Filter("is_mod",0).Filter("is_forbid",0).All(&replaceWords)
	//循环遍历切片，取出所有的禁用词
	//定义一个切片用来存储禁用词
	var words []string
	for i := 0; i < len(replaceWords); i++{
		words = append(words,replaceWords[i].WordPattern)
	}
	//再将切片转换为字符串
	strs := strings.Join(words,"|")
	//创建正则表达式的规则
	reg := regexp.MustCompile(strs)

	return reg.MatchString(msg)
}
//进行替换词的替换
func ReplaceWord(msg string)(comments string)  {
	//将数据库中所有的禁用词全部查询出来，存储到一个切片中
	o := orm.NewOrm()
	//声明一个禁用词的对象，用来存储从数据库中查询出来的词
	var replaceWords []models.SensitiveWord
	o.QueryTable("sensitive_word").Filter("is_mod",0).Filter("is_forbid",0).All(&replaceWords)
	//循环遍历切片，取出所有的禁用词
	//定义一个切片用来存储禁用词
	var words []string
	for i := 0; i < len(replaceWords); i++{
		words = append(words,replaceWords[i].WordPattern)
	}
	//再将切片转换为字符串
	strs := strings.Join(words,"|")
	//创建正则表达式的规则
	reg := regexp.MustCompile(strs)
	//进行正则表达式规则的匹配
	b := reg.MatchString(msg)
	comments = msg
	if b {
		for i := 0; i < len(replaceWords); i++{
			comments = strings.Replace(comments,replaceWords[i].WordPattern,replaceWords[i].ReplaceWord,-1)
		}
	}
	return comments
}
//删除文章
func (this *ArticleInfoController)DeleteArticle()  {
	ids := this.GetString("strId")
	strIds:=strings.Split(ids,",")
	o:=orm.NewOrm()
	var articleInfo models.ArticleInfo
	for i:=0;i<len(strIds);i++{
		id,_:=strconv.Atoi(strIds[i])
		articleInfo.Id=id
		o.Delete(&articleInfo)
	}
	this.Data["json"]=map[string]interface{}{"flag":"ok"}
	this.ServeJSON()
}
















