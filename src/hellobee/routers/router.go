package routers

import (
	"hellobee/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/weight", &controllers.WeightController{})
	beego.Router("/online", &controllers.OnlineController{})
	beego.Router("/params", &controllers.ParamController{})
}
