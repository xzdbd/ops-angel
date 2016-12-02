package controllers

import (
	"encoding/xml"
	"time"

	"strings"

	"github.com/astaxie/beego"
	"github.com/xzdbd/ops-angel/models"
)

type AngelController struct {
	beego.Controller
}

func (c *AngelController) Get() {
	signature := c.GetString("signature")
	timestamp := c.GetString("timestamp")
	nonce := c.GetString("nonce")
	echostr := c.GetString("echostr")

	if models.CheckSignature(timestamp, nonce) == signature {
		c.Ctx.WriteString(echostr)
	} else {
		c.Ctx.WriteString("")
	}
}

func (c *AngelController) Post() {
	var req models.Request
	xml.Unmarshal(c.Ctx.Input.RequestBody, &req)
	beego.Trace("req:", req)

	content := req.Content
	beego.Trace("content:", content)
	toolname := strings.Split(content, " ")[0]
	beego.Trace("toolname:", toolname)
	switch toolname {
	case models.GoogleToolName:
		beego.Trace("It is a google tool.")
	}

	var resp models.TextResponse
	resp.ToUserName = req.FromUserName
	resp.FromUserName = req.ToUserName
	resp.Content = "Hello, it works!"
	resp.CreateTime = time.Duration(time.Now().Unix())
	resp.MsgType = models.MsgTypeText

	c.Data["xml"] = resp
	c.ServeXML()
}
