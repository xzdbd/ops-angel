package models

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strings"

	"github.com/astaxie/beego/httplib"
	"github.com/docker/go-dockercloud/dockercloud"
)

const (

	// Docker Cloud Tool
	DockerCloudToolName     = "dockercloud"
	DockerCloudToolAlias    = "dc"
	DockerCloudToolEndpoint = "/dockercloud"
	DockerCloudHelpMsg      = `dockercloud is an operations tool.

Usage:
	dockercloud service NAME [status]|start|stop
or
	dc service NAME [status]|start|stop
		
Example 
	dc service test status`
)

type DockerCloudTool struct {
	toolBase
	ServiceName string
	Action      string
	Privileged  bool
	HelpMsg     string
}

func (dc *DockerCloudTool) NewTool() {
	dc.name = DockerCloudToolName
	dc.alias = DockerCloudToolAlias
	dc.endpoint = DockerCloudToolAlias
	dc.HelpMsg = DockerCloudHelpMsg
}

func (dc *DockerCloudTool) Run() (TextResponse, error) {
	var dcList dockercloud.SListResponse
	var textResp TextResponse
	textResp.MsgType = MsgTypeText

	switch dc.Action {
	case "status":
		if dc.ServiceName != "" {
			if strings.ToLower(dc.ServiceName) == "all" {
				var err error
				dcList, err = getAllDockerCloudService()
				if err != nil {
					return textResp, err
				}
				textResp.Content = fmt.Sprintf("共有%d个服务。\n", dcList.Meta.TotalCount)
				for i := 0; i < dcList.Meta.TotalCount; i++ {
					textResp.Content += fmt.Sprintf("%d. %s: %s\n", i+1, dcList.Objects[i].Name, dcList.Objects[i].State)
				}
			} else {
				var err error
				dcList, err = getDockerCloudServiceByName(dc.ServiceName)
				if err != nil {
					return textResp, err
				}
				if dcList.Meta.TotalCount < 1 {
					textResp.Content = fmt.Sprintf("没有找到名称为%s的服务。", dc.ServiceName)
				} else {
					textResp.Content = fmt.Sprintf("%s: %s\n", dcList.Objects[0].Name, dcList.Objects[0].State)
				}
			}
		}
	case "start":
		if dc.ServiceName != "" {
			_, err := actionDockerCloudService(dc.ServiceName, "start")
			if err != nil {
				textResp.Content = fmt.Sprintf("服务启动错误，错误信息：%s\n", err.Error())
			} else {
				textResp.Content = fmt.Sprintf("服务启动成功，请稍后查看该服务状态。")
			}
		}
	case "stop":
		if dc.ServiceName != "" {
			_, err := actionDockerCloudService(dc.ServiceName, "stop")
			if err != nil {
				textResp.Content = fmt.Sprintf("服务错误，错误信息：%s\n", err.Error())
			} else {
				textResp.Content = fmt.Sprintf("服务停止成功，请稍后查看该服务状态。")
			}
		}
	case "redeploy":
		if dc.ServiceName != "" {
			_, err := actionDockerCloudService(dc.ServiceName, "redeploy")
			if err != nil {
				textResp.Content = fmt.Sprintf("服务重新部署错误，错误信息：%s\n", err.Error())
			} else {
				textResp.Content = fmt.Sprintf("服务重新部署成功，请稍后查看该服务状态。")
			}
		}
	default:
		return textResp, errors.New("Invalid Action. Valid actions are 'start', 'stop' and 'status'")
	}
	return textResp, nil
}

func getAllDockerCloudService() (dockercloud.SListResponse, error) {
	var dcList dockercloud.SListResponse
	req := httplib.Get(APIADDRESS + APIVERSION + DockerCloudToolEndpoint + "/service")
	req.SetBasicAuth(apiuser, apipassword)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	err := req.ToJSON(&dcList)
	if err != nil {
		return dcList, err
	}
	return dcList, nil
}

func getDockerCloudServiceByName(name string) (dockercloud.SListResponse, error) {
	var dcList dockercloud.SListResponse
	req := httplib.Get(APIADDRESS + APIVERSION + DockerCloudToolEndpoint + "/service/" + name)
	req.SetBasicAuth(apiuser, apipassword)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	err := req.ToJSON(&dcList)
	if err != nil {
		return dcList, err
	}
	return dcList, nil
}

func getDockerCloudServiceUuid(name string) (string, error) {
	dcList, err := getDockerCloudServiceByName(name)
	if err != nil {
		return "", err
	}
	if dcList.Meta.TotalCount < 1 {
		return "", fmt.Errorf("没有找到名称为%s的服务。", name)
	}
	return dcList.Objects[0].Uuid, nil
}

// start, stop and redeploy service
func actionDockerCloudService(name string, action string) (dockercloud.Service, error) {
	var service dockercloud.Service
	var req *httplib.BeegoHTTPRequest
	uuid, err := getDockerCloudServiceUuid(name)
	if err != nil {
		return service, err
	}
	switch action {
	case "start":
		req = httplib.Post(APIADDRESS + APIVERSION + DockerCloudToolEndpoint + "/service/" + uuid + "/start")
	case "stop":
		req = httplib.Post(APIADDRESS + APIVERSION + DockerCloudToolEndpoint + "/service/" + uuid + "/stop")
	case "redeploy":
		req = httplib.Post(APIADDRESS + APIVERSION + DockerCloudToolEndpoint + "/service/" + uuid + "/redeploy")
	}
	req.SetBasicAuth(apiuser, apipassword)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	err = req.ToJSON(&service)
	if err != nil {
		return service, err
	}
	return service, nil
}
