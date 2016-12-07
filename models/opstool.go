package models

import (
	"crypto/tls"

	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
)

type Tool interface {
	NewTool()
}

type toolBase struct {
	name     string
	alias    string
	endpoint string
}

type GoogleTool struct {
	toolBase
	Key     string
	N       int
	HelpMsg string
}

type GoogleResult struct {
	Abstract  string
	Title     string
	URL       string
	Sitelinks []Sitelinks
}

type Sitelinks struct {
	Index    string
	Abstract string
	Title    string
	URL      string
}

const (
	APIADDRESS = "https://api.xzdbd.com/"
	APIVERSION = "v1"

	// Google Tool
	GoogleToolName     = "google"
	GoogleToolAlias    = "g"
	GoogleToolEndpoint = "/search/google"
	GoogleHelpMsg      = `google is a google search tool in wechat. It will return 4 results by default.

	Usage:
		google KEY 
	or
		g KEY 

	KEY:
		search key words.

	Example:
		google happy day	
	`
)

var (
	apiuser     = beego.AppConfig.String("apiuser")
	apipassword = beego.AppConfig.String("apipassword")
)

func (g *GoogleTool) NewTool() {
	g.name = GoogleToolName
	g.alias = GoogleToolAlias
	g.endpoint = GoogleToolEndpoint
	g.HelpMsg = GoogleHelpMsg
}

func (g *GoogleTool) Run() (NewsResponse, error) {
	var GoogleResultList []*GoogleResult
	GoogleResultList = make([]*GoogleResult, g.N)
	var newsResp NewsResponse
	newsResp.MsgType = MsgTypeNews

	req := httplib.Get(APIADDRESS + APIVERSION + g.endpoint)
	req.Param("key", g.Key)
	req.Param("n", strconv.Itoa(g.N))
	req.SetBasicAuth(apiuser, apipassword)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	err := req.ToJSON(&GoogleResultList)
	if err != nil {
		return newsResp, err
	}

	newsResp.ArticleCount = g.N
	for i := 0; i < g.N; i++ {
		item := Item{Title: GoogleResultList[i].Title, Description: GoogleResultList[i].Abstract, Url: GoogleResultList[i].URL, PicUrl: "https://lh3.googleusercontent.com/0-BzaWtxoAnsBjQ_wzUcKxyF07XE7v2Kkg1ogPVUdzmQpvaz118uHQEGU6BdtzJuzfo=h1264"}
		newsResp.Articles = append(newsResp.Articles, &item)
	}
	beego.Trace("newsResp:", newsResp)

	return newsResp, nil
}
