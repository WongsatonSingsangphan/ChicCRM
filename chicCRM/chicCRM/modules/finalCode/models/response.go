package models

type DataResponse struct {
	Scdact_id           string `json:"scdact_id"`
	Scdact_action       string `json:"scdact_action"`
	Scdact_actiontime   string `json:"scdact_actiontime"`
	Scdact_status       string `json:"scdact_status"`
	Scdact_reqid        string `json:"scdact_reqid"`
	Scdact_command      string `json:"scdact_command"`
	Scdact_filename     string `json:"scdact_filename"`
	Scdact_filetype     string `json:"scdact_filetype"`
	Scdact_filehash     string `json:"scdact_filehash"`
	Scdact_filesize     string `json:"scdact_filesize"`
	Scdact_filecreated  string `json:"scdact_filecreated"`
	Scdact_filemodified string `json:"scdact_filemodified"`
	Scdact_filelocation string `json:"scdact_filelocation"`
	Scdact_name         string `json:"scdact_name"`
	Scdact_type         string `json:"scdact_type"`
	Scdact_reciepient   string `json:"scdact_reciepient"`
	Scdact_starttime    string `json:"scdact_starttime"`
	Scdact_endtime      string `json:"scdact_endtime"`
	Scdact_periodday    int    `json:"scdact_periodday"`
	Scdact_periodhour   int    `json:"scdact_periodhour"`
	Scdact_numberopen   int    `json:"scdact_numberopen"`
	Scdact_nolimit      bool   `json:"scdact_nolimit"`
	Scdact_cvtoriginal  bool   `json:"scdact_cvtoriginal"`
	Scdact_edit         bool   `json:"scdact_edit"`
	Scdact_print        bool   `json:"scdact_print"`
	Scdact_copy         bool   `json:"scdact_copy"`
	Scdact_scrwatermark bool   `json:"scdact_scrwatermark"`
	Scdact_watermark    bool   `json:"scdact_watermark"`
	Scdact_cvthtml      bool   `json:"scdact_cvthtml"`
	Scdact_cvtfcl       bool   `json:"scdact_cvtfcl"`
	Scdact_marcro       bool   `json:"scdact_marcro"`
	Scdact_msgtext      string `json:"scdact_msgtext"`
	Scdact_subject      string `json:"scdact_subject"`
	// Scdact_teamlead_id           string `json:"scdact_teamlead_id"`
	Scdact_timestamp             string `json:"scdact_timestamp"`
	Scdact_sender                string `json:"scdact_sender"`
	Scdact_enableconvertoriginal bool   `json:"scdact_enableconvertoriginal"`
}

type TeamleadResponse struct {
	Scdact_id                    string `json:"scdact_id"`
	Scdact_action                string `json:"scdact_action"`
	Scdact_actiontime            string `json:"scdact_actiontime"`
	Scdact_status                string `json:"scdact_status"`
	Scdact_reqid                 string `json:"scdact_reqid"`
	Scdact_command               string `json:"scdact_command"`
	Scdact_filename              string `json:"scdact_filename"`
	Scdact_filetype              string `json:"scdact_filetype"`
	Scdact_filehash              string `json:"scdact_filehash"`
	Scdact_filesize              string `json:"scdact_filesize"`
	Scdact_filecreated           string `json:"scdact_filecreated"`
	Scdact_filemodified          string `json:"scdact_filemodified"`
	Scdact_filelocation          string `json:"scdact_filelocation"`
	Scdact_name                  string `json:"scdact_name"`
	Scdact_type                  string `json:"scdact_type"`
	Scdact_reciepient            string `json:"scdact_reciepient"`
	Scdact_starttime             string `json:"scdact_starttime"`
	Scdact_endtime               string `json:"scdact_endtime"`
	Scdact_periodday             int    `json:"scdact_periodday"`
	Scdact_periodhour            int    `json:"scdact_periodhour"`
	Scdact_numberopen            int    `json:"scdact_numberopen"`
	Scdact_nolimit               bool   `json:"scdact_nolimit"`
	Scdact_cvtoriginal           bool   `json:"scdact_cvtoriginal"`
	Scdact_edit                  bool   `json:"scdact_edit"`
	Scdact_print                 bool   `json:"scdact_print"`
	Scdact_copy                  bool   `json:"scdact_copy"`
	Scdact_scrwatermark          bool   `json:"scdact_scrwatermark"`
	Scdact_watermark             bool   `json:"scdact_watermark"`
	Scdact_cvthtml               bool   `json:"scdact_cvthtml"`
	Scdact_cvtfcl                bool   `json:"scdact_cvtfcl"`
	Scdact_marcro                bool   `json:"scdact_marcro"`
	Scdact_msgtext               string `json:"scdact_msgtext"`
	Scdact_subject               string `json:"scdact_subject"`
	Scdact_timestamp             string `json:"scdact_timestamp"`
	Scdact_actor                 string `json:"scdact_actor"`
	Scdact_sender                string `json:"scdact_sender"`
	Scdact_enableconvertoriginal bool   `json:"scdact_enableconvertoriginal"`
}

type PolicyResponse struct {
	TeamleadID string `json:"orgmbat_teamlead_id"`
	FeatureID  string `json:"orgmbat_tntft_id"`
	Feature    string `json:"orgmbat_feature"`
	Right      string `json:"orgmbat_right"`
	Level      string `json:"level"`
}

// type StatusResponseByRequestID struct { // <<<< before fix because return as a array frontend more happy
//
//		RequestID string  `json:"requestID"`
//		Approver1 *string `json:"approver1"`
//		Status1   *string `json:"status1"`
//		Approver2 *string `json:"approver2"`
//		Status2   *string `json:"status2"`
//		Approver3 *string `json:"approver3"`
//		Status3   *string `json:"status3"`
//		Approver4 *string `json:"approver4"`
//		Status4   *string `json:"status4"`
//	}

// type ApproverStatus struct { // in case we want "" but we want null
//
//		Approver string `json:"approver"`
//		Status   string `json:"status"`
//	}

type ApproverStatus struct {
	Approver *string `json:"approver"`
	Status   *string `json:"status"`
	Level    *int    `json:"level"`
}

type StatusResponseByRequestID struct {
	RequestID string           `json:"requestID"`
	Approvals []ApproverStatus `json:"approvals"`
}
