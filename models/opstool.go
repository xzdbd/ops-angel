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
	Key string
	N   int
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
	GoogleToolAlias    = "google"
	GoogleToolEndpoint = "/search/google"
)

func (g *GoogleTool) NewTool() {
	g.name = GoogleToolName
	g.alias = GoogleToolAlias
	g.endpoint = GoogleToolEndpoint
}

func (g *GoogleTool) Run() (NewsResponse, error) {
	var GoogleResultList []*GoogleResult
	GoogleResultList = make([]*GoogleResult, g.N)
	var newsResp NewsResponse
	newsResp.MsgType = MsgTypeNews

	req := httplib.Get(APIADDRESS + APIVERSION + g.endpoint)
	req.Param("key", g.Key)
	req.Param("n", strconv.Itoa(g.N))
	req.SetBasicAuth("xzdbd", "xzdbd1989")
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	err := req.ToJSON(&GoogleResultList)
	if err != nil {
		return newsResp, err
	}

	newsResp.ArticleCount = g.N
	for i := 0; i < g.N; i++ {
		item := Item{Title: GoogleResultList[i].Title, Description: GoogleResultList[i].Abstract, Url: GoogleResultList[i].URL}
		newsResp.Articles = append(newsResp.Articles, &item)
	}
	beego.Trace("newsResp:", newsResp)

	return newsResp, nil
}
