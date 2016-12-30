package routers

import (
	"hellobee/controllers"

	"github.com/astaxie/beego"
)

func init() {
	//beego.Router("/", &controllers.WeightController{})
	beego.Router("/weight", &controllers.WeightController{})
	beego.Router("/online", &controllers.OnlineController{})
}
