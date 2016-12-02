package models

type tool struct {
	name     string
	alias    string
	args     map[int32]string
	endpoint string
}

const (
	GoogleToolName     = "goolge"
	GoogleToolAlias    = "google"
	GoogleToolArgs1    = "key"
	GoogleToolArgs2    = "n"
	GoogleToolEndpoint = "/search/google"
)

var (
	GoogleTool tool
	//GoogleNewsTool tool
)
