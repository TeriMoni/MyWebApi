package routers

import (
	"MyWebApi/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/perception", &controllers.ApiController{})
}
