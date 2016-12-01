package routers

import (
	"github.com/astaxie/beego"
	"github.com/xzdbd/ops-angel/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/weixin", &controllers.AngelController{})
}
