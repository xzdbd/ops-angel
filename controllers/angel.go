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

	if req.MsgType == models.MsgTypeText {
		content := req.Content
		toolname := strings.Split(content, " ")[0]

		switch toolname {
		case models.GoogleToolName:
			beego.Trace("It is a google tool.")
			resp, respHelp := googleToolHandler(content, req)
			beego.Trace("resp:", resp, "respHelp", respHelp.Content)
			if respHelp.Content != "" {
				c.Data["xml"] = respHelp
				c.ServeXML()
			}

			c.Data["xml"] = resp
			c.ServeXML()
		}
	} else if req.MsgType == models.MsgTypeEvent && req.Event == models.MsgTypeEventSubscribe {
		subscribeHandler(req)
	} else {
		var resp models.TextResponse
		resp.ToUserName = req.FromUserName
		resp.FromUserName = req.ToUserName
		resp.Content = "Hello, it works!"
		resp.CreateTime = time.Duration(time.Now().Unix())
		resp.MsgType = models.MsgTypeText
		c.Data["xml"] = resp
		c.ServeXML()
	}

}

func googleToolHandler(content string, req models.Request) (models.NewsResponse, models.TextResponse) {
	var googleTool models.GoogleTool
	var resp models.NewsResponse

	googleTool.NewTool()

	cmd := strings.SplitN(content, " ", 2)
	beego.Trace("1:", cmd[0], "2:", cmd[1])
	length := len(cmd)

	if length == 2 {
		googleTool.Key = cmd[1]
		googleTool.N = 4
	} else {
		return resp, googleToolHelpHandler(req, googleTool)
	}

	resp, err := googleTool.Run()
	if err != nil {
		return resp, googleToolHelpHandler(req, googleTool)
	}
	resp.ToUserName = req.FromUserName
	resp.FromUserName = req.ToUserName
	resp.CreateTime = time.Duration(time.Now().Unix())
	return resp, models.TextResponse{}
}

func googleToolHelpHandler(req models.Request, g models.GoogleTool) models.TextResponse {
	var respHelp models.TextResponse
	respHelp.ToUserName = req.FromUserName
	respHelp.FromUserName = req.ToUserName
	respHelp.Content = g.HelpMsg
	respHelp.CreateTime = time.Duration(time.Now().Unix())
	respHelp.MsgType = models.MsgTypeText
	return respHelp
}

func subscribeHandler(req models.Request) models.TextResponse {
	var resp models.TextResponse
	resp.ToUserName = req.FromUserName
	resp.FromUserName = req.ToUserName
	resp.Content = `感谢订阅运维小天使官方微信\U0001f606，目前支持的工具：
		1.google
	输入工具名获取使用帮助。`
	resp.CreateTime = time.Duration(time.Now().Unix())
	resp.MsgType = models.MsgTypeText
	return resp
}
