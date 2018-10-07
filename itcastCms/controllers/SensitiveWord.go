package controllers

import (
	"github.com/astaxie/beego"
	"strings"
	"itcastCms/models"
	"github.com/astaxie/beego/orm"
)
//创建词库结构体
type SensitiveWordController struct {
	beego.Controller
}
//绑定展示添加词库的界面
func (this *SensitiveWordController)Index()  {
	this.TplName = "SensitiveWord/Index.html"
}
//添加词库
func (this *SensitiveWordController)AddWords()  {
	//从前端获取添加到词库的所有词
	contentMsg := this.GetString("contentMsg")
	//将词库中的词按照行进行分割
	strs := strings.Split(contentMsg,"\r\n")
	//循环遍历strs中的数据，取出每一行，再进行分割
	//创建一个切片用来存储敏感词
	var wordMsgs []models.SensitiveWord
	for i:=0;i<len(strs);i++{
		strMsg := strings.Split(strs[i],"=")
		//创建一个对象
		var SensitiveWord models.SensitiveWord
		//取出敏感词的部分
		SensitiveWord.WordPattern = strMsg[0]
		//对敏感词的类别进行判断，判断是什么类别的敏感词
		if strMsg[1] == "{BANNED}"{
			SensitiveWord.IsForbid = 1
		}else if strMsg[1] == "{MOD}"{
			SensitiveWord.IsMod = 1
		}else {
			SensitiveWord.ReplaceWord = strMsg[1]
		}
		//将查询出来的词存入切片中
		wordMsgs = append(wordMsgs,SensitiveWord)
	}
	o := orm.NewOrm()
	//进行多数据的插入数据库操作
	o.InsertMulti(len(wordMsgs),wordMsgs)
	this.Data["json"] = map[string]interface{}{"flag":"yes"}
	this.ServeJSON()
}