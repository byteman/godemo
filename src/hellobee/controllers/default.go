package controllers

import (
	"fmt"
	"hellobee/models"

	"hellobee/czserver"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type MainController struct {
	beego.Controller
}
type WeightController struct {
	beego.Controller
}

type OnlineController struct {
	beego.Controller
}

func (c *MainController) Get() {

}

type JsonData struct {
	Error   string
	Result  uint32
	Weights []models.OneWeight
}

func (c *WeightController) Get() {
	fmt.Println("weight reqeust")
	o := orm.NewOrm()
	ws := make([]models.OneWeight, 20)

	sql := "select * from one_weight order by id desc limit 10"
	fmt.Println(sql)
	_, er := o.Raw(sql).QueryRows(&ws)
	res := JsonData{}
	if er != nil {
		fmt.Println("查询出错")
		res.Error = "error"
		res.Result = 1
	} else {

		res.Error = "ok"
		res.Result = 0
		res.Weights = ws

	}

	c.Data["json"] = &res

	c.ServeJSON()
	return
}

func (c *OnlineController) Get() {
	fmt.Println("online reqeust")
	page, err := c.GetInt("pages")
	if err != nil {
		fmt.Println(err)
		//return
	}
	fmt.Println("page=", page)

	clients := czserver.GetClient()
	fmt.Println("client=", clients)
	c.Data["json"] = &clients
	c.ServeJSON()

}
