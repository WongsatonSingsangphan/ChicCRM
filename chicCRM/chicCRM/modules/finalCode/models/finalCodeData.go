package models

type TenantFeatureInput struct {
	ID string `json:"id"`
}

type TenantFeature struct {
	TntftID         string  `json:"tntft_id"`
	TntftTntID      string  `json:"tntft_tnt_id"`
	TntftName       string  `json:"tntft_name"`
	TntftType       string  `json:"tntft_type"`
	TntftStart      string  `json:"tntft_start"`
	TntftExp        string  `json:"tntft_exp"`
	TntftStatus     string  `json:"tntft_status"`
	TntftLicense    string  `json:"tntft_license"`
	TntftUnit       string  `json:"tntft_unit"`
	TntftLimit      int     `json:"tntft_limit"`
	TntftEditor     *string `json:"tntft_editor"`
	TntftCreator    string  `json:"tntft_creator"`
	TntftCreateLoc  *string `json:"tntft_createlocation"`
	TntftUpdateLoc  *string `json:"tntft_updatelocation"`
	TntftCreateTime string  `json:"tntft_createtime"`
	TntftUpdateTime string  `json:"tntft_updatetime"`
	OrgmbatFeature  string  `json:"orgmbat_feature"`
	OrgmbatRight    string  `json:"orgmbat_right"`
}

type TenantFeatureMember struct {
	TntftID      string `json:"tntft_id"`
	TntftTntID   string `json:"tntft_tnt_id"`
	TntftName    string `json:"tntft_name"`
	TntftType    string `json:"tntft_type"`
	TntftStart   string `json:"tntft_start"`
	TntftExp     string `json:"tntft_exp"`
	TntftStatus  string `json:"tntft_status"`
	TntftLicense string `json:"tntft_license"`
	TntftUnit    string `json:"tntft_unit"`
	TntftLimit   int    `json:"tntft_limit"`
	// TntftEditor     string `json:"tntft_editor"`
	TntftEditor  *string `json:"tntft_editor"`
	TntftCreator string  `json:"tntft_creator"`
	// TntftCreateLoc  string         `json:"tntft_createlocation"`
	TntftCreateLoc  *string `json:"tntft_createlocation"`
	TntftUpdateLoc  *string `json:"tntft_updatelocation"`
	TntftCreateTime string  `json:"tntft_createtime"`
	TntftUpdateTime string  `json:"tntft_updatetime"`
}

type OrganizeMemberAuthorization struct {
	TntftID string `json:"orgmbat_tntft_id"`
	Feature string `json:"orgmbat_feature"`
	Right   string `json:"orgmbat_right"`
}

type PolicyData struct {
	ScdpolID             string  `json:"scdpol_id"`
	ScdusrOrgID          string  `json:"scdusr_org_id"`
	ScdpolFiletype       *string `json:"scdpol_filetype"`
	ScdpolName           *string `json:"scdpol_name"`
	ScdpolType           *string `json:"scdpol_type"`
	ScdpolRecipient      *string `json:"scdpol_reciepient"`
	ScdpolStartTime      *string `json:"scdpol_starttime"`
	ScdpolEndTime        *string `json:"scdpol_endtime"`
	ScdpolPeriodDay      int     `json:"scdpol_periodday"`
	ScdpolPeriodHour     int     `json:"scdpol_periodhour"`
	ScdpolNumberOpen     int     `json:"scdpol_numberopen"`
	ScdpolNoLimit        bool    `json:"scdpol_nolimit"`
	ScdpolCvtOriginal    bool    `json:"scdpol_cvtoriginal"`
	ScdpolEdit           bool    `json:"scdpol_edit"`
	ScdpolPrint          bool    `json:"scdpol_print"`
	ScdpolCopy           bool    `json:"scdpol_copy"`
	ScdpolScrWatermark   bool    `json:"scdpol_scrwatermark"`
	ScdpolWatermark      bool    `json:"scdpol_watermark"`
	ScdpolCvtHTML        bool    `json:"scdpol_cvthtml"`
	ScdpolCvtFCL         bool    `json:"scdpol_cvtfcl"`
	ScdpolMarcro         bool    `json:"scdpol_marcro"`
	ScdpolMsgText        *string `json:"scdpol_msgtext"`
	ScdpolCreateLocation *string `json:"scdpol_createlocation"`
	ScdpolUpdateLocation *string `json:"scdpol_updatelocation"`
	ScdpolCreateTime     string  `json:"scdpol_createtime"`
	ScdpolUpdateTime     string  `json:"scdpol_updatetime"`
}




type ResponseStatus struct {
	Status  string `json:"Status"`
	Message string `json:"Message"`
}

type RequestActivity struct {
	OrgmbatOrgmbid string
	// ค่าลูป
	ScdactCommand      string `json:"scdact_command"`
	ScdactBinary       []byte
	ScdactFilename     string `json:"scdact_filename"`
	ScdactFiletype     string `json:"scdact_filetype"`
	ScdactFilehash     string `json:"scdact_filehash"`
	ScdactFilesize     string `json:"scdact_filesize"`
	ScdactFilecreated  string `json:"scdact_filecreated"`
	ScdactFilemodified string `json:"scdact_filemodified"`
	// ไม่ค่าลูป
	ScdactAction         string `json:"scdact_action"`
	ScdactActor          string `json:"scdact_actor"`
	ScdactStatus         string `json:"scdact_status"`
	ScdactReqID          string `json:"scdact_reqid"`
	ScdactFilelocation   string `json:"scdact_filelocation"`
	ScdactName           string `json:"scdact_name"`
	ScdactType           string `json:"scdact_type"`
	ScdactReciepient     string `json:"scdact_reciepient"`
	ScdactStartTime      string `json:"scdact_starttime"`
	ScdactEndTime        string `json:"scdact_endtime"`
	ScdactPeriodDay      string `json:"scdact_periodday"`
	ScdactPeriodHour     string `json:"scdact_periodhour"`
	ScdactNumberOpen     string `json:"scdact_numberopen"`
	ScdactNoLimit        bool   `json:"scdact_nolimit"`
	ScdactCvtOriginal    bool   `json:"scdact_cvtoriginal"`
	ScdactEdit           bool   `json:"scdact_edit"`
	ScdactPrint          bool   `json:"scdact_print"`
	ScdactCopy           bool   `json:"scdact_copy"`
	ScdactScrWatermark   bool   `json:"scdact_scrwatermark"`
	ScdactWatermark      bool   `json:"scdact_watermark"`
	ScdactCvtHtml        bool   `json:"scdact_cvthtml"`
	ScdactCvtFcl         bool   `json:"scdact_cvtfcl"`
	ScdactMarcro         bool   `json:"scdact_marcro"`
	ScdactEnableCv       bool   `json:"scdact_enableconvertoriginal"`
	ScdactMsgText        string `json:"scdact_msgtext"`
	ScdactSubject        string `json:"scdact_subject"`
	ScdactSender         string `json:"scdact_sender"`
	ScdactCreateLocation string `json:"scdact_createlocation"`
	ScdactUpdateLocation string `json:"scdact_updatelocation"`
	ScdactActionTime     string
}


type ResponseStatusFC struct {
	Status   string            `json:"status"`
	Message  map[string]string `json:"message"`
}
