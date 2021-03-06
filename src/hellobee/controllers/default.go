package controllers

import (
	"encoding/json"
	"fmt"
	"hellobee/czserver"
	"hellobee/models"
	con "strconv"

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
type ParamController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Redirect("static/html/index.html", 301)
}

type JsonData struct {
	Error    string
	Result   uint32
	Total    uint32
	PageSize uint32
	Weights  []models.OneWeight
}

/**
 * 分页函数，适用任何表
 * 返回 总记录条数,总页数,以及当前请求的数据RawSeter,调用中需要"rs.QueryRows(&tblog)"就行了  --tblog是一个Tb_log对象
 * 参数：表名，当前页数，页面大小，条件（查询条件,格式为 " and name='zhifeiya' and age=12 "）
 */
func GetPagesInfo(tableName string, currentpage int, pagesize int, conditions string) (int, int, orm.RawSeter) {
	if currentpage <= 1 {
		currentpage = 1
	}
	if pagesize == 0 {
		pagesize = 20
	}
	var rs orm.RawSeter
	o := orm.NewOrm()
	var totalItem, totalpages int = 0, 0 //总条数,总页数
	sql := "SELECT count(*) FROM " + tableName + "  where 1>0 " + conditions
	o.Raw(sql).QueryRow(&totalItem) //获取总条数
	fmt.Println(sql, totalItem)
	if totalItem <= pagesize {
		totalpages = 1
	} else if totalItem > pagesize {
		temp := totalItem / pagesize
		if (totalItem % pagesize) != 0 {
			temp = temp + 1
		}
		totalpages = temp
	}
	rs = o.Raw("select *  from  " + tableName + " order by id desc " + conditions + " LIMIT " + con.Itoa((currentpage-1)*pagesize) + "," + con.Itoa(pagesize))
	return totalItem, totalpages, rs
}

func (c *WeightController) Get() {
	fmt.Println("weight reqeust")
	var cond = ""
	page, err := c.GetInt("pages")
	if err != nil {
		fmt.Println(err)
		page = 1
	}
	id, err := c.GetInt("id")
	if err == nil {
		cond = fmt.Sprintf(" and dev_id=%d", id)
	}
	fmt.Println("page=", page, "id=", id)

	all, pagesize, rs := GetPagesInfo("one_weight", page, 10, cond)

	//o := orm.NewOrm()
	ws := make([]models.OneWeight, 0)

	//sql := "select * from one_weight order by id desc limit 10"
	//fmt.Println(sql)
	//_, er := o.Raw(sql).QueryRows(&ws)
	rs.QueryRows(&ws)
	res := JsonData{}

	res.Error = "ok"
	res.Result = 0
	res.Total = uint32(all)
	res.PageSize = uint32(pagesize)
	res.Weights = ws
	fmt.Println(all, pagesize, ws)
	c.Data["json"] = &res

	c.ServeJSON()
	return
}

func (c *OnlineController) Get() {

	clients := czserver.GetClient()
	//fmt.Println("client=", clients)
	c.Data["json"] = &clients
	c.ServeJSON()

}

func (c *ParamController) Post() {
	fmt.Println("params post")
	u := czserver.DevInfo{}
	res := make(map[string]interface{})
	result := 0
	message := "ok"
	defer func() {
		res["result"] = result
		res["message"] = message
		c.Data["json"] = res
		c.ServeJSON()
	}()
	fmt.Println(c.Ctx.Input.RequestBody)
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &u)
	if err != nil {
		fmt.Println(err)
		result = 1
		message = "json convert failed"
		return
	}

	fmt.Println(u)

}
func (c *ParamController) Get() {
	fmt.Println("params get")
	clients := czserver.GetClient()
	c.Data["json"] = &clients
	c.ServeJSON()
}
