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
	beego.Info("User request, User:", req.FromUserName, "Message Type:", req.MsgType, "Content:", req.Content)
	if req.MsgType == models.MsgTypeText {
		content := req.Content
		toolname := strings.Split(content, " ")[0]

		switch toolname {
		case models.GoogleToolName, models.GoogleToolAlias:
			beego.Info("User request google tool, User:", req.FromUserName, "Command:", req.Content)
			resp, respHelp := googleToolHandler(content, req)
			if respHelp.Content != "" {
				beego.Info("Response to the user with google tool help. User:", req.FromUserName)
				c.Data["xml"] = respHelp
				c.ServeXML()
				break
			}
			beego.Info("Response to the user with google result. User:", req.FromUserName, "Count:", resp.ArticleCount)
			c.Data["xml"] = resp
			c.ServeXML()

		case models.DockerCloudToolName, models.DockerCloudToolAlias:
			beego.Info("User request dockercloud tool, User:", req.FromUserName, "Command:", req.Content)
			resp := dockerCloudToolHandler(content, req)
			beego.Info("Response to the user with dockercloud result. User:", req.FromUserName, "content:", req.Content)
			c.Data["xml"] = resp
			c.ServeXML()

		default:
			c.Data["xml"] = descriptionHandler(req)
			c.ServeXML()
		}
	} else if req.MsgType == models.MsgTypeEvent && req.Event == models.MsgTypeEventSubscribe {
		c.Data["xml"] = subscribeHandler(req)
		c.ServeXML()
	} else {
		c.Data["xml"] = descriptionHandler(req)
		c.ServeXML()
	}

}

func googleToolHandler(content string, req models.Request) (models.NewsResponse, models.TextResponse) {
	var googleTool models.GoogleTool
	var resp models.NewsResponse

	googleTool.NewTool()

	cmd := strings.SplitN(content, " ", 2)
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

func dockerCloudToolHandler(content string, req models.Request) models.TextResponse {
	var dcTool models.DockerCloudTool
	var resp models.TextResponse

	dcTool.NewTool()

	cmd := strings.Split(content, " ")
	length := len(cmd)

	if length == 2 {
		if cmd[1] == "service" {
			dcTool.Action = "status"
			dcTool.ServiceName = "all"
		} else {
			return dockerCloudToolHelpHandler(req, dcTool)
		}
	} else if length == 3 {
		if cmd[1] == "service" {
			dcTool.Action = "status"
			dcTool.ServiceName = cmd[2]
		} else {
			return dockerCloudToolHelpHandler(req, dcTool)
		}
	} else if length == 4 {
		if cmd[1] == "service" {
			if cmd[3] == "start" || cmd[3] == "stop" {
				dcTool.Action = cmd[3]
				var valid bool
				valid, resp = validatePrivilegedAction(req)
				if !valid {
					return resp
				}
			} else if cmd[3] == "status" {
				dcTool.Action = cmd[3]
			} else {
				return dockerCloudToolHelpHandler(req, dcTool)
			}
			dcTool.ServiceName = cmd[2]
		} else {
			return dockerCloudToolHelpHandler(req, dcTool)
		}
	}

	resp, err := dcTool.Run()
	if err != nil {
		return dockerCloudToolHelpHandler(req, dcTool)
	}
	resp.ToUserName = req.FromUserName
	resp.FromUserName = req.ToUserName
	resp.CreateTime = time.Duration(time.Now().Unix())
	return resp
}

func dockerCloudToolHelpHandler(req models.Request, dc models.DockerCloudTool) models.TextResponse {
	var respHelp models.TextResponse
	respHelp.ToUserName = req.FromUserName
	respHelp.FromUserName = req.ToUserName
	respHelp.Content = dc.HelpMsg
	respHelp.CreateTime = time.Duration(time.Now().Unix())
	respHelp.MsgType = models.MsgTypeText
	return respHelp
}

func subscribeHandler(req models.Request) models.TextResponse {
	var resp models.TextResponse
	resp.ToUserName = req.FromUserName
	resp.FromUserName = req.ToUserName
	resp.Content = `感谢订阅运维小天使官方微信，目前支持的工具：
		1.google
输入工具名获取使用帮助。`
	resp.CreateTime = time.Duration(time.Now().Unix())
	resp.MsgType = models.MsgTypeText
	return resp
}

func descriptionHandler(req models.Request) models.TextResponse {
	var resp models.TextResponse
	resp.ToUserName = req.FromUserName
	resp.FromUserName = req.ToUserName
	resp.Content = `运维小天使官方微信，目前支持的工具：
		1.google
输入工具名获取使用帮助。`
	resp.CreateTime = time.Duration(time.Now().Unix())
	resp.MsgType = models.MsgTypeText
	return resp
}

func validatePrivilegedAction(req models.Request) (bool, models.TextResponse) {
	var textResp models.TextResponse
	users := beego.AppConfig.Strings("privilegeduser")
	for i := 0; i < len(users); i++ {
		if req.FromUserName == users[i] {
			return true, textResp
		}
	}
	textResp.Content = "您没有权限执行该操作。"
	return false, textResp
}
