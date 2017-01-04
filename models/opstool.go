package models

import "github.com/astaxie/beego"

type Tool interface {
	NewTool()
}

type toolBase struct {
	name     string
	alias    string
	endpoint string
}

const (
	APIADDRESS = "https://api.xzdbd.com/"
	//APIADDRESS = "http://11.11.1.6:8098/"
	APIVERSION = "v1"
)

var (
	apiuser     = beego.AppConfig.String("apiuser")
	apipassword = beego.AppConfig.String("apipassword")
)
