package repositories

import (
	"chicCRM/modules/finalCode/models"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/lib/pq"
)

type RepositoryPort interface {
	CheckFeatureOrganizeRepositories(id string) (models.TenantFeature, error)
	CheckFeatureMemberRepositories(id string) (models.TenantFeatureMember, models.OrganizeMemberAuthorization, error)
	GetPolicyDataRepositories(id string) (models.PolicyData, error)
	GetSecuredocActivityByOrganizeMemberUUIDRepositories(organizeMemberUUID string) ([]models.DataResponse, error)
	GetSecuredocActivityByTeamleadIDRepositories(teamleadID string) ([]models.TeamleadResponse, error)
	GetPolicyAuthorizationByTeamleadRepositories(teamleadID string) (models.PolicyResponse, error)
	GetRequestStatusByReqIDRepositories(requestID string) (models.StatusResponseByRequestID, error)
	AddActivity(req models.RequestActivity) error
	GETUser_File(scdact_id string) ([][]byte, error)
	// AddStatus(ScdactStatus string, ScdactReqID string, teamleadID string) error
	AddStatus(ScdactStatus string, ScdactReqID string, teamleadID string) (int, error)
}

type repositoryAdapter struct {
	db *sql.DB
}

func NewRepositoryAdapter(db *sql.DB) RepositoryPort {
	return &repositoryAdapter{db: db}
}

// Check Tenant Feature Organize
func (r *repositoryAdapter) CheckFeatureOrganizeRepositories(id string) (models.TenantFeature, error) {
	query := `
		SELECT 
			tf.tntft_id, tf.tntft_tnt_id, tf.tntft_name, tf.tntft_type, tf.tntft_start, 
			tf.tntft_exp, tf.tntft_status, tf.tntft_license, tf.tntft_unit, tf.tntft_limit, 
			tf.tntft_editor, tf.tntft_creator, tf.tntft_createlocation, tf.tntft_updatelocation, 
			tf.tntft_createtime, tf.tntft_updatetime, oma.orgmbat_feature, oma.orgmbat_right
		FROM tenant_feature tf
		JOIN organize_member_authorization oma ON tf.tntft_id = oma.orgmbat_tntft_id
		WHERE tf.tntft_tnt_id = $1 AND oma.orgmbat_orgid = $1
	`
	row := r.db.QueryRow(query, id)

	// สร้างตัวแปรที่ใช้เก็บข้อมูลที่คืนมา
	var tenantFeature models.TenantFeature
	// ใช้ Scan เพื่อ map ข้อมูลจากแถวที่ query ได้เข้ากับโครงสร้างข้อมูลโมเดล TenantFeature
	err := row.Scan(
		&tenantFeature.TntftID,
		&tenantFeature.TntftTntID,
		&tenantFeature.TntftName,
		&tenantFeature.TntftType,
		&tenantFeature.TntftStart,
		&tenantFeature.TntftExp,
		&tenantFeature.TntftStatus,
		&tenantFeature.TntftLicense,
		&tenantFeature.TntftUnit,
		&tenantFeature.TntftLimit,
		&tenantFeature.TntftEditor,
		&tenantFeature.TntftCreator,
		&tenantFeature.TntftCreateLoc,
		&tenantFeature.TntftUpdateLoc,
		&tenantFeature.TntftCreateTime,
		&tenantFeature.TntftUpdateTime,
		&tenantFeature.OrgmbatFeature,
		&tenantFeature.OrgmbatRight,
	)
	if err != nil {
		log.Println(err)
		return models.TenantFeature{}, err
	}

	return tenantFeature, nil
}

// Check Tenant Feature Member
func (r *repositoryAdapter) CheckFeatureMemberRepositories(id string) (models.TenantFeatureMember, models.OrganizeMemberAuthorization, error) {
	memberID := id
	var tenantFeatureMember models.TenantFeatureMember
	var memberAuthorization models.OrganizeMemberAuthorization

	// Query to retrieve orgmbat_tntft_id from organize_member_authorization where orgmbat_orgmbid matches memberID
	queryOrgAuth := `
        SELECT orgmbat_tntft_id, orgmbat_feature, orgmbat_right
        FROM organize_member_authorization
        WHERE orgmbat_orgmbid = $1
    `
	row := r.db.QueryRow(queryOrgAuth, memberID)
	err := row.Scan(&memberAuthorization.TntftID, &memberAuthorization.Feature, &memberAuthorization.Right)
	if err != nil {
		log.Println(err)
		return models.TenantFeatureMember{}, models.OrganizeMemberAuthorization{}, err
	}

	// Query to retrieve tenant_feature data based on the orgmbat_tntft_id obtained above
	queryTenantFeatureMember := `
        SELECT tntft_id, tntft_tnt_id, tntft_name, tntft_type, tntft_start,
            tntft_exp, tntft_status, tntft_license, tntft_unit, tntft_limit,
            tntft_editor, tntft_creator, tntft_createlocation, tntft_updatelocation,
            tntft_createtime, tntft_updatetime
        FROM tenant_feature
        WHERE tntft_id = $1
    `
	row = r.db.QueryRow(queryTenantFeatureMember, memberAuthorization.TntftID)
	err = row.Scan(
		&tenantFeatureMember.TntftID, &tenantFeatureMember.TntftTntID, &tenantFeatureMember.TntftName, &tenantFeatureMember.TntftType,
		&tenantFeatureMember.TntftStart, &tenantFeatureMember.TntftExp, &tenantFeatureMember.TntftStatus, &tenantFeatureMember.TntftLicense,
		&tenantFeatureMember.TntftUnit, &tenantFeatureMember.TntftLimit, &tenantFeatureMember.TntftEditor, &tenantFeatureMember.TntftCreator,
		&tenantFeatureMember.TntftCreateLoc, &tenantFeatureMember.TntftUpdateLoc, &tenantFeatureMember.TntftCreateTime,
		&tenantFeatureMember.TntftUpdateTime,
	)

	if err != nil {
		log.Println(err)
		return models.TenantFeatureMember{}, models.OrganizeMemberAuthorization{}, err
	}

	return tenantFeatureMember, memberAuthorization, nil
}

func (r *repositoryAdapter) GetPolicyDataRepositories(id string) (models.PolicyData, error) {

	// Step 1: Get orgmbat_id from organize_member_authorization
	orgmbatIDQuery := `
        SELECT orgmbat_id
        FROM organize_member_authorization
        WHERE orgmbat_orgmbid = $1;
    `

	var orgmbatID string
	err := r.db.QueryRow(orgmbatIDQuery, id).Scan(&orgmbatID)
	if err != nil {
		log.Println(err)
		return models.PolicyData{}, err
	}

	// Step 2: Get scdpolusr_scdpol_id from securedoc_policy_user
	var scdpolusrID string
	scdpolusrIDQuery := `
        SELECT scdpolusr_scdpol_id
        FROM securedoc_policy_user
        WHERE scdpolusr_scdusr_id = $1;
    `
	err = r.db.QueryRow(scdpolusrIDQuery, orgmbatID).Scan(&scdpolusrID)
	if err != nil {
		log.Println(err)
		return models.PolicyData{}, err
	}
	// fmt.Println(scdpolusrID)

	// Step 3: Get policy data from securedoc_policy
	query := `
        SELECT * FROM securedoc_policy
        WHERE scdpol_id = $1
    `
	row := r.db.QueryRow(query, scdpolusrID)

	var policyData models.PolicyData
	err = row.Scan(
		&policyData.ScdpolID,
		&policyData.ScdusrOrgID,
		&policyData.ScdpolFiletype,
		&policyData.ScdpolName,
		&policyData.ScdpolType,
		&policyData.ScdpolRecipient,
		&policyData.ScdpolStartTime,
		&policyData.ScdpolEndTime,
		&policyData.ScdpolPeriodDay,
		&policyData.ScdpolPeriodHour,
		&policyData.ScdpolNumberOpen,
		&policyData.ScdpolNoLimit,
		&policyData.ScdpolCvtOriginal,
		&policyData.ScdpolEdit,
		&policyData.ScdpolPrint,
		&policyData.ScdpolCopy,
		&policyData.ScdpolScrWatermark,
		&policyData.ScdpolWatermark,
		&policyData.ScdpolCvtHTML,
		&policyData.ScdpolCvtFCL,
		&policyData.ScdpolMarcro,
		&policyData.ScdpolMsgText,
		&policyData.ScdpolCreateLocation,
		&policyData.ScdpolUpdateLocation,
		&policyData.ScdpolCreateTime,
		&policyData.ScdpolUpdateTime,
	)
	if err != nil {
		log.Println(err)
		return models.PolicyData{}, err
	}
	// fmt.Println(policyData)
	// fmt.Println(scdpolusrID)
	return policyData, nil
}

func (r *repositoryAdapter) GetSecuredocActivityByOrganizeMemberUUIDRepositories(organizeMemberUUID string) ([]models.DataResponse, error) {
	// var dataResponse models.DataResponse ** only one activity
	// var organizeMemberAuthorizationUUID string
	// err := r.db.QueryRow("SELECT orgmbat_id FROM organize_member_authorization WHERE orgmbat_orgmbid = $1", organizeMemberUUID).Scan(&organizeMemberAuthorizationUUID)
	// if err != nil {
	// 	log.Printf("Failed to retriving orgmbat_id from organize_member_authorization:%v\n", err)
	// 	return models.DataResponse{}, err
	// }
	// err = r.db.QueryRow(`SELECT scdact_id, scdact_action, scdact_actiontime, scdact_status, scdact_reqid, scdact_command, scdact_filename, scdact_filetype, scdact_filehash,
	// scdact_filesize, scdact_filecreated, scdact_filemodified, scdact_filelocation, scdact_name, scdact_type, scdact_reciepient, scdact_starttime, scdact_endtime, scdact_periodday,
	// scdact_periodhour, scdact_numberopen, scdact_nolimit, scdact_cvtoriginal, scdact_edit, scdact_print, scdact_copy, scdact_scrwatermark, scdact_watermark, scdact_cvthtml, scdact_cvtfcl,
	// scdact_marcro, scdact_msgtext, scdact_subject, scdact_teamlead_id FROM securedoc_activity WHERE scdact_actor = $1`, organizeMemberAuthorizationUUID).Scan(&dataResponse.Scdact_id, &dataResponse.Scdact_action,
	// 	&dataResponse.Scdact_actiontime, &dataResponse.Scdact_status, &dataResponse.Scdact_reqid, &dataResponse.Scdact_command, &dataResponse.Scdact_filename, &dataResponse.Scdact_filetype,
	// 	&dataResponse.Scdact_filehash, &dataResponse.Scdact_filesize, &dataResponse.Scdact_filecreated, &dataResponse.Scdact_filemodified, &dataResponse.Scdact_filelocation, &dataResponse.Scdact_name,
	// 	&dataResponse.Scdact_type, &dataResponse.Scdact_reciepient, &dataResponse.Scdact_starttime, &dataResponse.Scdact_endtime, &dataResponse.Scdact_periodday, &dataResponse.Scdact_periodhour,
	// 	&dataResponse.Scdact_numberopen, &dataResponse.Scdact_nolimit, &dataResponse.Scdact_cvtoriginal, &dataResponse.Scdact_edit, &dataResponse.Scdact_print, &dataResponse.Scdact_copy,
	// 	&dataResponse.Scdact_scrwatermark, &dataResponse.Scdact_watermark, &dataResponse.Scdact_cvthtml, &dataResponse.Scdact_cvtfcl, &dataResponse.Scdact_marcro, &dataResponse.Scdact_msgtext, &dataResponse.Scdact_subject,
	// 	&dataResponse.Scdact_teamlead_id)
	// if err != nil {
	// 	log.Printf("failed to query logSecuredocActivity from securedoc_activity:%v\n", err)
	// 	return models.DataResponse{}, err
	// }

	// return dataResponse, nil
	var dataResponses []models.DataResponse
	var organizeMemberAuthorizationUUID string
	err := r.db.QueryRow("SELECT orgmbat_id FROM organize_member_authorization WHERE orgmbat_orgmbid = $1", organizeMemberUUID).Scan(&organizeMemberAuthorizationUUID)
	if err != nil {
		log.Printf("Failed to retrive orgmbat_id from organize_member_authorization: %v\n", err)
		return nil, errors.New("memberID does not match")
	}

	rows, err := r.db.Query(`SELECT scdact_id, scdact_action, scdact_actiontime, scdact_status, scdact_reqid, scdact_command, scdact_filename, scdact_filetype, scdact_filehash,
        scdact_filesize, scdact_filecreated, scdact_filemodified, scdact_filelocation, scdact_name, scdact_type, scdact_reciepient, scdact_starttime, scdact_endtime, scdact_periodday,
        scdact_periodhour, scdact_numberopen, scdact_nolimit, scdact_cvtoriginal, scdact_edit, scdact_print, scdact_copy, scdact_scrwatermark, scdact_watermark, scdact_cvthtml, scdact_cvtfcl,
        scdact_marcro, scdact_msgtext, scdact_subject, scdact_timestamp, scdact_sender, scdact_enableconvertoriginal FROM securedoc_activity WHERE scdact_actor = $1`, organizeMemberAuthorizationUUID)
	// scdact_marcro, scdact_msgtext, scdact_subject, scdact_teamlead_id, scdact_timestamp, scdact_sender, scdact_enableconvertoriginal FROM securedoc_activity WHERE scdact_actor = $1`, organizeMemberAuthorizationUUID)

	if err != nil {
		log.Printf("Failed to query logSecuredocActivity from securedoc_activity: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dataResponse models.DataResponse
		err := rows.Scan(&dataResponse.Scdact_id, &dataResponse.Scdact_action,
			&dataResponse.Scdact_actiontime, &dataResponse.Scdact_status, &dataResponse.Scdact_reqid, &dataResponse.Scdact_command, &dataResponse.Scdact_filename, &dataResponse.Scdact_filetype,
			&dataResponse.Scdact_filehash, &dataResponse.Scdact_filesize, &dataResponse.Scdact_filecreated, &dataResponse.Scdact_filemodified, &dataResponse.Scdact_filelocation, &dataResponse.Scdact_name,
			&dataResponse.Scdact_type, &dataResponse.Scdact_reciepient, &dataResponse.Scdact_starttime, &dataResponse.Scdact_endtime, &dataResponse.Scdact_periodday, &dataResponse.Scdact_periodhour,
			&dataResponse.Scdact_numberopen, &dataResponse.Scdact_nolimit, &dataResponse.Scdact_cvtoriginal, &dataResponse.Scdact_edit, &dataResponse.Scdact_print, &dataResponse.Scdact_copy,
			&dataResponse.Scdact_scrwatermark, &dataResponse.Scdact_watermark, &dataResponse.Scdact_cvthtml, &dataResponse.Scdact_cvtfcl, &dataResponse.Scdact_marcro, &dataResponse.Scdact_msgtext, &dataResponse.Scdact_subject,
			&dataResponse.Scdact_timestamp, &dataResponse.Scdact_sender, &dataResponse.Scdact_enableconvertoriginal)
		// &dataResponse.Scdact_teamlead_id, &dataResponse.Scdact_timestamp, &dataResponse.Scdact_sender, &dataResponse.Scdact_enableconvertoriginal)
		if err != nil {
			log.Printf("Failed to scan row: %v\n", err)
			return nil, err
		}
		dataResponses = append(dataResponses, dataResponse)
	}

	return dataResponses, nil
}

func (r *repositoryAdapter) GetSecuredocActivityByTeamleadIDRepositories(teamleadID string) ([]models.TeamleadResponse, error) {
	var teamleadResponses []models.TeamleadResponse
	var orgmbatLevel int
	err := r.db.QueryRow("SELECT orgmbat_level FROM organize_member_authorization WHERE orgmbat_orgmbid = $1", teamleadID).Scan(&orgmbatLevel)
	if err != nil {
		log.Println(err)
		return []models.TeamleadResponse{}, errors.New("user not found")
	}
	if orgmbatLevel == 2 || orgmbatLevel == 3 {
		employeeLevel := 1
		rows, err := r.db.Query(`SELECT scdact_id, scdact_action, scdact_actiontime, scdact_status, scdact_reqid, scdact_command, scdact_filename, scdact_filetype, scdact_filehash,
        scdact_filesize, scdact_filecreated, scdact_filemodified, scdact_filelocation, scdact_name, scdact_type, scdact_reciepient, scdact_starttime, scdact_endtime, scdact_periodday,
        scdact_periodhour, scdact_numberopen, scdact_nolimit, scdact_cvtoriginal, scdact_edit, scdact_print, scdact_copy, scdact_scrwatermark, scdact_watermark, scdact_cvthtml, scdact_cvtfcl,
        scdact_marcro, scdact_msgtext, scdact_subject, scdact_actor, scdact_timestamp, scdact_sender, scdact_enableconvertoriginal FROM securedoc_activity WHERE scdact_level = $1`, employeeLevel)
		if err != nil {
			log.Printf("Failed to query logSecuredocActivity from securedoc_activity: %v\n", err)
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var teamleadResponse models.TeamleadResponse
			err := rows.Scan(&teamleadResponse.Scdact_id, &teamleadResponse.Scdact_action,
				&teamleadResponse.Scdact_actiontime, &teamleadResponse.Scdact_status, &teamleadResponse.Scdact_reqid, &teamleadResponse.Scdact_command, &teamleadResponse.Scdact_filename, &teamleadResponse.Scdact_filetype,
				&teamleadResponse.Scdact_filehash, &teamleadResponse.Scdact_filesize, &teamleadResponse.Scdact_filecreated, &teamleadResponse.Scdact_filemodified, &teamleadResponse.Scdact_filelocation, &teamleadResponse.Scdact_name,
				&teamleadResponse.Scdact_type, &teamleadResponse.Scdact_reciepient, &teamleadResponse.Scdact_starttime, &teamleadResponse.Scdact_endtime, &teamleadResponse.Scdact_periodday, &teamleadResponse.Scdact_periodhour,
				&teamleadResponse.Scdact_numberopen, &teamleadResponse.Scdact_nolimit, &teamleadResponse.Scdact_cvtoriginal, &teamleadResponse.Scdact_edit, &teamleadResponse.Scdact_print, &teamleadResponse.Scdact_copy,
				&teamleadResponse.Scdact_scrwatermark, &teamleadResponse.Scdact_watermark, &teamleadResponse.Scdact_cvthtml, &teamleadResponse.Scdact_cvtfcl, &teamleadResponse.Scdact_marcro, &teamleadResponse.Scdact_msgtext, &teamleadResponse.Scdact_subject,
				&teamleadResponse.Scdact_actor, &teamleadResponse.Scdact_timestamp, &teamleadResponse.Scdact_sender, &teamleadResponse.Scdact_enableconvertoriginal)
			if err != nil {
				log.Printf("Failed to scan row: %v\n", err)
				return nil, err
			}
			teamleadResponses = append(teamleadResponses, teamleadResponse)
		}
		return teamleadResponses, nil
	} else {
		log.Println(err)
		return nil, errors.New("users do not have permission to access this information; only managers and founders are authorized to do so")
	}
}

// func (r *repositoryAdapter) GetRequestStatusByReqIDRepositories(requestID string) (models.StatusResponseByRequestID, error) {
// 	var statusResponseByRequestID models.StatusResponseByRequestID
// 	err := r.db.QueryRow("SELECT scdwflap_reqid, scdwflap_1, scdwflap_status_1, scdwflap_2, scdwflap_status_2, scdwflap_3, scdwflap_status_3, scdwflap_4, scdwflap_status_4 FROM securedoc_workflow_approvers WHERE scdwflap_reqid = $1", requestID).Scan(
// 		&statusResponseByRequestID.RequestID, &statusResponseByRequestID.Approver1, &statusResponseByRequestID.Status1, &statusResponseByRequestID.Approver2, &statusResponseByRequestID.Status2,
// 		&statusResponseByRequestID.Approver3, &statusResponseByRequestID.Status3, &statusResponseByRequestID.Approver4, &statusResponseByRequestID.Status4)
// 	if err != nil {
// 		log.Printf("Failed retriving scdwflap_reqid from securedoc_workflow_approvers table. Error:%v", err)
// 		return models.StatusResponseByRequestID{}, errors.New("requestID not found")
// 	}
// 	// fmt.Println(statusResponseByRequestID)

//		return statusResponseByRequestID, nil
//	}
// func (r *repositoryAdapter) GetRequestStatusByReqIDRepositories(requestID string) (models.StatusResponseByRequestID, error) { // <<< we going to return array to help frontend
// 	var statusResponseByRequestID models.StatusResponseByRequestID
// 	err := r.db.QueryRow(`
// 	SELECT
//     a.scdwflap_reqid,
//     b1.orgtl_trac_username AS approver1_email,
//     a.scdwflap_status_1,
//     b2.orgtl_trac_username AS approver2_email,
//     a.scdwflap_status_2,
//     b3.orgtl_trac_username AS approver3_email,
//     a.scdwflap_status_3,
//     b4.orgtl_trac_username AS approver4_email,
//     a.scdwflap_status_4
// FROM
//     securedoc_workflow_approvers a
// LEFT JOIN organize_teamlead_trac b1 ON a.scdwflap_1 = b1.orgtl_trac_id
// LEFT JOIN organize_teamlead_trac b2 ON a.scdwflap_2 = b2.orgtl_trac_id
// LEFT JOIN organize_teamlead_trac b3 ON a.scdwflap_3 = b3.orgtl_trac_id
// LEFT JOIN organize_teamlead_trac b4 ON a.scdwflap_4 = b4.orgtl_trac_id
// WHERE
//     a.scdwflap_reqid = $1;`, requestID).Scan(
// 		&statusResponseByRequestID.RequestID, &statusResponseByRequestID.Approver1, &statusResponseByRequestID.Status1, &statusResponseByRequestID.Approver2, &statusResponseByRequestID.Status2,
// 		&statusResponseByRequestID.Approver3, &statusResponseByRequestID.Status3, &statusResponseByRequestID.Approver4, &statusResponseByRequestID.Status4)
// 	if err != nil {
// 		log.Printf("Failed retrieving status by request ID from securedoc_workflow_approvers table. Error: %v", err)
// 		return models.StatusResponseByRequestID{}, errors.New("requestID not found")
// 	}

//		return statusResponseByRequestID, nil
//	}
func (r *repositoryAdapter) GetRequestStatusByReqIDRepositories(requestID string) (models.StatusResponseByRequestID, error) {
	var response models.StatusResponseByRequestID
	response.RequestID = requestID

	var approver1, approver2, approver3, approver4 *string
	var status1, status2, status3, status4 *string
	var level1, level2, level3, level4 *int
	err := r.db.QueryRow(`
    SELECT 
        a.scdwflap_reqid, 
        b1.orgtl_trac_username, 
        a.scdwflap_status_1, 
		b1.orgtl_level, 
        b2.orgtl_trac_username, 
        a.scdwflap_status_2, 
		b2.orgtl_level, 
        b3.orgtl_trac_username, 
        a.scdwflap_status_3, 
		b3.orgtl_level, 
        b4.orgtl_trac_username, 
        a.scdwflap_status_4,
		b4.orgtl_level
    FROM 
        securedoc_workflow_approvers a
    LEFT JOIN organize_teamlead_trac b1 ON a.scdwflap_1 = b1.orgtl_trac_id
    LEFT JOIN organize_teamlead_trac b2 ON a.scdwflap_2 = b2.orgtl_trac_id
    LEFT JOIN organize_teamlead_trac b3 ON a.scdwflap_3 = b3.orgtl_trac_id
    LEFT JOIN organize_teamlead_trac b4 ON a.scdwflap_4 = b4.orgtl_trac_id
    WHERE 
        a.scdwflap_reqid = $1`, requestID).Scan(
		&response.RequestID, &approver1, &status1, &level1, &approver2, &status2, &level2, &approver3, &status3, &level3, &approver4, &status4, &level4)
	// &response.RequestID, &approver1, &status1, &approver2, &status2, &approver3, &status3, &approver4, &status4)

	if err != nil {
		log.Printf("Failed retrieving status by request ID from securedoc_workflow_approvers table. Error: %v", err)
		return models.StatusResponseByRequestID{}, errors.New("requestID not found")
	}

	// Define a slice to hold the approver and status pair
	var approvals []models.ApproverStatus

	// Append each approver and status to the slice, checking for nil
	// approvals = append(approvals, models.ApproverStatus{Approver: utility.DereferenceString(approver1), Status: utility.DereferenceString(status1)}) // fix want to retrn null instead of ""
	// approvals = append(approvals, models.ApproverStatus{Approver: utility.DereferenceString(approver2), Status: utility.DereferenceString(status2)})
	// approvals = append(approvals, models.ApproverStatus{Approver: utility.DereferenceString(approver3), Status: utility.DereferenceString(status3)})
	// approvals = append(approvals, models.ApproverStatus{Approver: utility.DereferenceString(approver4), Status: utility.DereferenceString(status4)})

	// approvals = append(approvals, models.ApproverStatus{Approver: (approver1), Status: (status1)}) // ** work but no level we want to return level
	// approvals = append(approvals, models.ApproverStatus{Approver: (approver2), Status: (status2)})
	// approvals = append(approvals, models.ApproverStatus{Approver: (approver3), Status: (status3)})
	// approvals = append(approvals, models.ApproverStatus{Approver: (approver4), Status: (status4)})

	approvals = append(approvals, models.ApproverStatus{Approver: (approver1), Status: (status1), Level: (level1)})
	approvals = append(approvals, models.ApproverStatus{Approver: (approver2), Status: (status2), Level: (level2)})
	approvals = append(approvals, models.ApproverStatus{Approver: (approver3), Status: (status3), Level: (level3)})
	approvals = append(approvals, models.ApproverStatus{Approver: (approver4), Status: (status4), Level: (level4)})

	// Assign the slice to the response
	response.Approvals = approvals
	return response, nil
}

func (r *repositoryAdapter) GetPolicyAuthorizationByTeamleadRepositories(teamleadID string) (models.PolicyResponse, error) {
	var policyResponse models.PolicyResponse
	var organizeTeamleadTracID string
	err := r.db.QueryRow("SELECT orgtl_trac_id FROM organize_teamlead_trac WHERE orgtl_trac_id = $1", teamleadID).Scan(&organizeTeamleadTracID)
	if err != nil {
		log.Printf("Failed to retriving orgtl_trac_id from organize_teamlead_trac table. Error:%v", err)
		return models.PolicyResponse{}, errors.New("organizeID does not match")
	}
	err = r.db.QueryRow("SELECT orgmbat_tntft_id, orgmbat_feature, orgmbat_right, orgmbat_level FROM organize_member_authorization WHERE orgmbat_orgmbid = $1", organizeTeamleadTracID).Scan(&policyResponse.FeatureID, &policyResponse.Feature, &policyResponse.Right, &policyResponse.Level)
	if err != nil {
		log.Printf("Failed to query logSecuredocActivity from securedoc_activity: %v\n", err)
		return models.PolicyResponse{}, err
	}
	policyResponse.TeamleadID = organizeTeamleadTracID
	return policyResponse, nil
}

func (r *repositoryAdapter) AddActivity(req models.RequestActivity) error {
	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("unable to start transaction: %v", err)
	}

	// var uuid_Member  string
	var orgmbatID, orgmbat_orgmbid, orgmbat_orgid string
	// var orgmbatTeamLeadID string
	err = tx.QueryRow(`
		SELECT orgmbat_id, orgmbat_orgmbid, orgmbat_orgid
		FROM organize_member_authorization 
		WHERE orgmbat_orgmbid = $1`, req.OrgmbatOrgmbid).Scan(&orgmbatID, &orgmbat_orgmbid, &orgmbat_orgid)

	if err != nil {
		// Rollback the transaction in case of an error
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("unable to rollback transaction: %v, original error: %v", rollbackErr, err)
		}
		return fmt.Errorf("unable to fetch orgmbat_id: %v", err)
	}
	rows, err := tx.Query(`
    SELECT orgmbat_orgmbid 
    FROM organize_member_authorization 
    WHERE orgmbat_level > 1;
`)
	if err != nil {
		log.Printf("Failed to query teamlead level above than level 1. Error:%v", err)
		return err
	}
	defer rows.Close()

	var orgmbat_orgmbIDs []string
	for rows.Next() {
		var orgmbat_orgmbID string
		if err := rows.Scan(&orgmbat_orgmbID); err != nil {
			log.Printf("Failed to scan orgmbat_id. Error: %v", err)
			return err
		}
		orgmbat_orgmbIDs = append(orgmbat_orgmbIDs, orgmbat_orgmbID)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Failed during rows iteration. Error: %v", err)
		return err
	}
	var securedocWorkflowUUID string
	err = tx.QueryRow("INSERT INTO securedoc_workflow (scdwfl_actor, scdwfl_org) VALUES ($1, $2) RETURNING scdwfl_id", orgmbatID, orgmbat_orgid).Scan(&securedocWorkflowUUID)
	if err != nil {
		log.Printf("Failed to insert data and retriving scdwfl_id from table securedoc_workflow. Error: %v", err)
		return err
	}
	var securedocWorkflowApproversUUID string
	pending := "Pending"
	err = tx.QueryRow("INSERT INTO securedoc_workflow_approvers (scdwflap_sequence, scdwflap_1, scdwflap_2, scdwflap_3, scdwflap_reqid, scdwflap_status_1, scdwflap_status_2, scdwflap_status_3) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING scdwflap_id", securedocWorkflowUUID, orgmbat_orgmbIDs[0], orgmbat_orgmbIDs[1], orgmbat_orgmbIDs[2], req.ScdactReqID, pending, pending, pending).Scan(&securedocWorkflowApproversUUID)
	if err != nil {
		log.Printf("Failed to insert data and retriving scdwflap_id from table securedoc_workflow_approvers. Error:%v", err)
		return err
	}
	// fmt.Println("Type of securedocWorkflowApproversUUID:", reflect.TypeOf(securedocWorkflowApproversUUID))
	// fmt.Println(securedocWorkflowApproversUUID)
	// fmt.Println(securedocWorkflowUUID)
	_, err = tx.Exec("UPDATE securedoc_workflow SET scdwfl_approver = $1 WHERE scdwfl_id = $2", securedocWorkflowApproversUUID, securedocWorkflowUUID)
	if err != nil {
		log.Printf("Failed to update scdwfl_approver in securedoc_workflow. Error: %v", err)
		return err
	}
	// fmt.Println("Type of securedocWorkflowApproversUUID:", reflect.TypeOf(securedocWorkflowApproversUUID))
	// fmt.Println(securedocWorkflowApproversUUID)
	_, err = tx.Exec(`
	INSERT INTO securedoc_activity (
	    scdact_binary,
		scdact_actor,
		scdact_action,
		scdact_status,
	    scdact_reqid,
	    scdact_command,
	    scdact_filename,
	    scdact_filetype,
	    scdact_filehash,
	    scdact_filesize,
	    scdact_filecreated,
	    scdact_filemodified,
	    scdact_filelocation,
	    scdact_name,
	    scdact_type,
	    scdact_reciepient,
	    scdact_starttime,
	    scdact_endtime,
	    scdact_periodday,
	    scdact_periodhour,
	    scdact_numberopen,
	    scdact_nolimit,
	    scdact_cvthtml,
	    scdact_cvtfcl,
	    scdact_marcro,
	    scdact_cvtoriginal,
	    scdact_edit,
	    scdact_print,
	    scdact_copy,
	    scdact_scrwatermark,
	    scdact_watermark,
	    scdact_msgtext,
	    scdact_subject,
		scdact_sender,
	    scdact_createlocation,
	    scdact_updatelocation,
		scdact_actiontime,
		scdact_enableconvertoriginal,
		scdact_approverid
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39)`,
		pq.Array([][]byte{req.ScdactBinary}),
		orgmbatID,
		req.ScdactAction,
		req.ScdactStatus,
		req.ScdactReqID,
		req.ScdactCommand,
		req.ScdactFilename,
		req.ScdactFiletype,
		req.ScdactFilehash,
		req.ScdactFilesize,
		req.ScdactFilecreated,
		req.ScdactFilemodified,
		req.ScdactFilelocation,
		req.ScdactName,
		req.ScdactType,
		req.ScdactReciepient,
		req.ScdactStartTime,
		req.ScdactEndTime,
		req.ScdactPeriodDay,
		req.ScdactPeriodHour,
		req.ScdactNumberOpen,
		req.ScdactNoLimit,
		req.ScdactCvtHtml,
		req.ScdactCvtFcl,
		req.ScdactMarcro,
		req.ScdactCvtOriginal,
		req.ScdactEdit,
		req.ScdactPrint,
		req.ScdactCopy,
		req.ScdactScrWatermark,
		req.ScdactWatermark,
		req.ScdactMsgText,
		req.ScdactSubject,
		req.ScdactSender,
		req.ScdactCreateLocation,
		req.ScdactUpdateLocation,
		req.ScdactActionTime,
		req.ScdactEnableCv,
		securedocWorkflowApproversUUID)
	// fmt.Println(securedocWorkflowApproversUUID)
	// fmt.Println("Type of securedocWorkflowApproversUUID:", reflect.TypeOf(securedocWorkflowApproversUUID))

	if err != nil {
		// Rollback the transaction in case of an error
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("unable to rollback transaction: %v, original error: %v", rollbackErr, err)
		}
		return fmt.Errorf("unable to execute query: %v", err)
	}

	// Commit the transaction if the query is successful
	if commitErr := tx.Commit(); commitErr != nil {
		return fmt.Errorf("unable to commit transaction: %v", commitErr)
	}

	return nil
}

func (r *repositoryAdapter) GETUser_File(scdact_id string) ([][]byte, error) {
	// SQL query ที่ใช้ parameterized query
	var scdactBinary [][]byte
	err := r.db.QueryRow("SELECT scdact_binary FROM securedoc_activity WHERE scdact_id = $1",
		scdact_id).Scan(pq.Array(&scdactBinary))
	if err != nil {
		return nil, err
	}

	return scdactBinary, nil
}

// func (r *repositoryAdapter) AddStatus(ScdactStatus string, ScdactReqID string, teamleadID string) error {

// 	var level int
// 	err := r.db.QueryRow(`
//         SELECT orgmbat_level FROM organize_member_authorization WHERE orgmbat_orgmbid = $1`, teamleadID).Scan(&level)
// 	if err != nil {
// 		log.Printf("Failed to query level from organize_member_authorization table. Error: %v", err)
// 		return err
// 	}
// 	if level != 3 {
// 		_, err = r.db.Exec(`
//             UPDATE securedoc_workflow_approvers
//             SET
//                 scdwflap_status_1 = CASE WHEN scdwflap_1 = $1 THEN $2 ELSE scdwflap_status_1 END,
//                 scdwflap_status_2 = CASE WHEN scdwflap_2 = $1 THEN $2 ELSE scdwflap_status_2 END,
//                 scdwflap_status_3 = CASE WHEN scdwflap_3 = $1 THEN $2 ELSE scdwflap_status_3 END
//             WHERE scdwflap_reqid = $3`, teamleadID, ScdactStatus, ScdactReqID)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

func (r *repositoryAdapter) AddStatus(ScdactStatus string, ScdactReqID string, teamleadID string) (int, error) {
	var level int
	err := r.db.QueryRow(`
        SELECT orgmbat_level FROM organize_member_authorization WHERE orgmbat_orgmbid = $1`, teamleadID).Scan(&level)
	if err != nil {
		log.Printf("Failed to query level from organize_member_authorization table. Error: %v", err)
		return 0, err
	}

	if level != 3 {
		// ถ้า level ไม่เท่ากับ 3, ทำการ update
		_, err = r.db.Exec(`
            UPDATE securedoc_workflow_approvers
            SET 
                scdwflap_status_1 = CASE WHEN scdwflap_1 = $1 THEN $2 ELSE scdwflap_status_1 END,
                scdwflap_status_2 = CASE WHEN scdwflap_2 = $1 THEN $2 ELSE scdwflap_status_2 END,
                scdwflap_status_3 = CASE WHEN scdwflap_3 = $1 THEN $2 ELSE scdwflap_status_3 END
            WHERE scdwflap_reqid = $3`, teamleadID, ScdactStatus, ScdactReqID)
		if err != nil {
			return 0, err
		}
	} else {
		// ถ้า level เท่ากับ 3, ทำการ update สถานะที่ตรงกับ teamleadID
		_, err = r.db.Exec(`
				UPDATE securedoc_workflow_approvers
				SET 
					scdwflap_status_1 = CASE WHEN scdwflap_1 = $1 THEN $2 ELSE scdwflap_status_1 END,
					scdwflap_status_2 = CASE WHEN scdwflap_2 = $1 THEN $2 ELSE scdwflap_status_2 END,
					scdwflap_status_3 = CASE WHEN scdwflap_3 = $1 THEN $2 ELSE scdwflap_status_3 END
				WHERE scdwflap_reqid = $3`, teamleadID, ScdactStatus, ScdactReqID)
		if err != nil {
			return 0, err
		}
	}
	// fmt.Println(level)
	return level, nil
}
