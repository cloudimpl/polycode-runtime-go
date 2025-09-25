package runtime

type ClientEnv struct {
	AppName    string `json:"appName"`
	AppPort    uint   `json:"appPort"`
	SidecarApi string `json:"sidecarApi"`
}
