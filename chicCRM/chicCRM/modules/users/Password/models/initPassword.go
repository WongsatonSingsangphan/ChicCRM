package models

type InitPassword struct {
	Username        string `json:"username"`
	Newpassword     string `json:"newpassword"`
	Requires_action string `json:"requires_action"`
}

type ChangePassword struct {
	Username        string `json:"username"`
	Oldpassword     string `json:"oldpassword"`
	Newpassword     string `json:"newpassword"`
	Requires_action string `json:"requires_action"`
}

type RequestResetPassword struct {
	Username string `json:"username"`
}

type ResetPassword struct {
	Username    string `json:"username"`
	NewPassword string `json:"newpassword"`
}
