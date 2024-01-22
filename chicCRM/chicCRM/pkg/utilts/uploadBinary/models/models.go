package models

type RequestOrganizeBinary struct {
	OrganizeID string `form:"organizeID"`
	FileBytes  [][]byte
}

type EditLogoRequest struct {
	OrganizeID string `form:"organizeID"`
	FileBytes  [][]byte
}

type EditPersonalProfileRequest struct {
	PersonalID string `form:"personalID"`
	FileBytes  [][]byte
}
