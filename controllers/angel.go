package controllers

import (
	"encoding/xml"
	"time"

	"strings"

	"strconv"

	"errors"

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
	beego.Trace("RequestBody:", string(c.Ctx.Input.RequestBody))
	var req models.Request
	xml.Unmarshal(c.Ctx.Input.RequestBody, &req)
	beego.Trace("req:", req)

	if req.MsgType == models.MsgTypeText {
		content := req.Content
		cmd := strings.Split(content, " ")
		toolname := cmd[0]

		switch toolname {
		case models.GoogleToolName:
			beego.Trace("It is a google tool.")
			resp, err := googleToolHandler(cmd)
			if err != nil {

			}
			resp.ToUserName = req.FromUserName
			resp.FromUserName = req.ToUserName
			resp.CreateTime = time.Duration(time.Now().Unix())

			beego.Trace("resp:", resp)
			c.Data["xml"] = resp
			c.ServeXML()
		}
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

func googleToolHandler(cmd []string) (models.NewsResponse, error) {
	var googleTool models.GoogleTool
	var resp models.NewsResponse

	googleTool.NewTool()
	length := len(cmd)
	if length == 2 {
		googleTool.Key = cmd[1]
		// default = 3
		googleTool.N = 3
	} else if length == 3 {
		googleTool.Key = cmd[1]
		var err error
		googleTool.N, err = strconv.Atoi(cmd[2])
		if err != nil {
			return resp, errors.New("参数错误。")
		}
	} else {
		return resp, errors.New("参数错误。")
	}

	beego.Trace("googleTool:", googleTool)
	resp, err := googleTool.Run()
	if err != nil {
		return resp, err
	}
	return resp, nil
}
