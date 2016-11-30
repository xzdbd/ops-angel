package routers

import (
	"github.com/xzdbd/ops-angel/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
