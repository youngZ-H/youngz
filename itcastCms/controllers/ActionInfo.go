package controllers

import (
	"github.com/astaxie/beego"
	"path"
	"strconv"
	"time"
	"os"
	"itcastCms/models"
	"github.com/astaxie/beego/orm"
	"strings"
)

type ActionInfoController struct {
	beego.Controller
}
//展示权限表格的界面
func (this *ActionInfoController)Index()  {
	this.TplName = "ActionInfo/Index.html"
}
//为权限控制器绑定上传文件的方法
func (this *ActionInfoController)FileUp()  {
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
//进行性权限类的添加角色信息的方法绑定
func (this *ActionInfoController)AddAction()  {
	//定义一个角色的对象
	var actionInfo models.ActionInfo
	//从前端获取数据
	actionInfo.Remark = this.GetString("Remark")
	actionInfo.DelFlag = 0
	actionInfo.AddDate = time.Now()
	actionInfo.ModifDate = time.Now()
	actionInfo.Url = this.GetString("Url")
	actionInfo.HttpMethod = this.GetString("HttpMethod")
	actionInfo.ActionInfoName = this.GetString("ActionInfoName")
	actionInfo.ActionTypeEnum,_ = this.GetInt("ActionTypeEnum")
	actionInfo.MenuIcon = this.GetString("MenuIcon")
	actionInfo.IconWidth = 0
	actionInfo.IconHeight = 0
	//创建orm对象
	o := orm.NewOrm()
	//将从前端表格中发送回来的数据库插入到数据库中
	_, err := o.Insert(&actionInfo)
	if err == nil {
		this.Data["json"] = map[string]interface{}{"flag":"ok","msg":"文件保存到数据库成功"}
	}else {
		this.Data["json"] = map[string]interface{}{"flag":"no","msg":"文件保存到数据库失败"}
	}
	//将保存成功的数据通过ajax的方式发送回前端
	this.ServeJSON()
}
//将权限信息展示出来
func (this *ActionInfoController)GetActionInfo()  {
	//从前端获取数据
	PageIndex,_ := this.GetInt("page")
	PageSize,_ := this.GetInt("rows")
	//起始页为
	start := (PageIndex-1)*PageSize
	//创建orm对象
	o := orm.NewOrm()
	var actions []models.ActionInfo
	//查询数据库权限表中的所有数据
	o.QueryTable("action_info").Filter("del_flag",0).OrderBy("id").Limit(PageSize,start).All(&actions)
	count,_ := o.QueryTable("action_info").Filter("del_flag",0).Count()
	//将数据传递给前端
	this.Data["json"] = map[string]interface{}{"rows":actions,"total":count}
	this.ServeJSON()
}
//删除角色信息方法的实现
func (this *ActionInfoController)DeleteAction()  {
	//从前端接收数据
	strId := this.GetString("strId")
	//进行数据的处理，将字符串分割为切片
	id := strings.Split(strId,",")
	beego.Info(id)

	//定义一个角色对象
	var actionInfo models.ActionInfo

	//创建orm对象，获取要删除的对象
	o := orm.NewOrm()
	//遍历数组
	for i:=0;i<len(id);i++{
		strId,_ := strconv.Atoi(id[i])
		actionInfo.Id = strId

		//删除选中的数据
		o.Delete(&actionInfo)
	}
	//返回删除信号标志位
	this.Data["json"] = map[string]interface{}{"flag":"ok"}
	this.ServeJSON()
}
//展示要编辑的角色信息
func (this *ActionInfoController)ShowEditAction()  {
	//从前端读取要编辑的数据的id号
	id,err := this.GetInt("Id")
	if err != nil{
		beego.Info("获取要编辑信息的id发生错误", err)
		return
	}
	//创建orm对象
	o := orm.NewOrm()
	//获取要编辑的对象
	var actionInfo models.ActionInfo
	o.QueryTable("action_info").Filter("id",id).One(&actionInfo)
	//将数据编码为ajax形式，发送给前端要编辑的角色信息
	this.Data["json"] = map[string]interface{}{"actionInfo":actionInfo}
	this.ServeJSON()
}





















