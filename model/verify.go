package model

type RBACEnforceReq struct {
	AccessToken string `json:"accessToken" binding:"required,len=32"`
	Path        string `json:"path" binding:"required"`
	Method      string `json:"method" binding:"required"`
}

type RBACEnforceResp struct {
	Verify     bool   `json:"verify"`
	Message    string `json:"message"`
	AppId      string `json:"appId"`
	Nickname   string `json:"nickname"`
	Username   string `json:"username"`
	ActionPass bool   `json:"actionPass"` // false-无权限
}
