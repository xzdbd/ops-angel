package models

import (
	"crypto/tls"
	"strconv"

	"github.com/astaxie/beego/httplib"
)

const (
	// Google Tool
	GoogleToolName     = "google"
	GoogleToolAlias    = "g"
	GoogleToolEndpoint = "/search/google"
	GoogleHelpMsg      = `google is a google search tool. Enjoy!

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

	newsResp.ArticleCount = len(GoogleResultList)
	for i := 0; i < newsResp.ArticleCount; i++ {
		//picUrl := getFavicons(GoogleResultList[i].URL)
		picUrl := "https://upload.wikimedia.org/wikipedia/commons/thumb/5/53/Google_%22G%22_Logo.svg/200px-Google_%22G%22_Logo.svg.png"
		item := Item{Title: GoogleResultList[i].Title, Description: GoogleResultList[i].Abstract, Url: GoogleResultList[i].URL, PicUrl: picUrl}
		newsResp.Articles = append(newsResp.Articles, &item)
	}

	return newsResp, nil
}

func getFavicons(domain string) string {
	//return "https://www.google.com/s2/favicons?domain=" + domain
	return "https://api.byi.pw/favicon/?url=" + domain
}
