package services

import (
	"archive/zip"
	"chicCRM/modules/finalCode/models"
	"chicCRM/modules/finalCode/repositories"
	"chicCRM/modules/finalCode/services/copy"
	"chicCRM/pkg/utilts/decrypt"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"errors"

	"github.com/gin-gonic/gin"
)

// "github.com/kballard/go-shellquote"
type ServicePort interface {
	CheckFeatureOrganizeServices(id string) (models.TenantFeature, error)
	CheckFeatureMemberServices(id string) (models.TenantFeatureMember, models.OrganizeMemberAuthorization, error)
	GetPolicyDataServices(id string) (models.PolicyData, error)
	GetSecuredocActivityByOrganizeMemberUUIDServices(organizeMemberUUID string) ([]models.DataResponse, error)
	GetSecuredocActivityByTeamleadIDServices(teamleadID string) ([]models.TeamleadResponse, error)
	GetPolicyAuthorizationByTeamleadServices(teamleadID string) (models.PolicyResponse, error)
	GetRequestStatusByReqIDServices(requestID string) (models.StatusResponseByRequestID, error)
	AddActivity(req models.RequestActivity) error
	GETUser_File(scdact_id string) ([][]byte, error)
	FilesUpload(c *gin.Context, file *multipart.FileHeader, orderID string) error
	ExecuteCommand(cmdStr string) (string, int, error)
	CFolder(orderID string) error
	CFolderEmail(orderID, emailInput string) error
	DeleteFolder(orderID string) error
	AttachFile(orderID, baseDir string, emailList []string) (map[string][]string, error)
	AddStatus(ScdactStatus string, ScdactReqID string, teamleadID string) (int, error)
}

type serviceAdapter struct {
	r repositories.RepositoryPort
}

func NewServiceAdapter(r repositories.RepositoryPort) ServicePort {
	return &serviceAdapter{r: r}
}

// Check Tenant Feature Organize
func (s *serviceAdapter) CheckFeatureOrganizeServices(id string) (models.TenantFeature, error) {
	tenantFeature, err := s.r.CheckFeatureOrganizeRepositories(id)
	if err != nil {
		err = errors.New("invalid parameter id")
		return models.TenantFeature{}, err
	}

	return tenantFeature, nil
}

// Check Tenant Feature Member
func (s *serviceAdapter) CheckFeatureMemberServices(id string) (models.TenantFeatureMember, models.OrganizeMemberAuthorization, error) {
	tenantFeatureMember, memberAuthorization, err := s.r.CheckFeatureMemberRepositories(id)
	if err != nil {
		return models.TenantFeatureMember{}, models.OrganizeMemberAuthorization{}, errors.New("invalid parameter id")
	}

	return tenantFeatureMember, memberAuthorization, nil
}

// Get Policy Data
func (s *serviceAdapter) GetPolicyDataServices(id string) (models.PolicyData, error) {
	policyData, err := s.r.GetPolicyDataRepositories(id)
	if err != nil {
		err = errors.New("invalid parameter id")
		return models.PolicyData{}, err
	}

	return policyData, nil
}

func (s *serviceAdapter) GetSecuredocActivityByOrganizeMemberUUIDServices(organizeMemberUUID string) ([]models.DataResponse, error) {
	logSecuredocActivity, err := s.r.GetSecuredocActivityByOrganizeMemberUUIDRepositories(organizeMemberUUID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return logSecuredocActivity, nil
}

func (s *serviceAdapter) GetSecuredocActivityByTeamleadIDServices(teamleadID string) ([]models.TeamleadResponse, error) {
	logSecuredogActivity, err := s.r.GetSecuredocActivityByTeamleadIDRepositories(teamleadID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return logSecuredogActivity, nil
}

func (s *serviceAdapter) GetPolicyAuthorizationByTeamleadServices(teamleadID string) (models.PolicyResponse, error) {
	policyByTeamlead, err := s.r.GetPolicyAuthorizationByTeamleadRepositories(teamleadID)
	if err != nil {
		log.Println(err)
		return models.PolicyResponse{}, err
	}
	return policyByTeamlead, nil
}

// func (s *serviceAdapter) GetRequestStatusByReqIDServices(requestID string) (models.StatusResponseByRequestID, error) { << return array but the value still tokenize we reverse to original
//
//		// fmt.Println(requestID)
//		requestStatusbyRequestID, err := s.r.GetRequestStatusByReqIDRepositories(requestID)
//		if err != nil {
//			log.Println(err)
//			return models.StatusResponseByRequestID{}, nil
//		}
//		return requestStatusbyRequestID, nil
//	}
func (s *serviceAdapter) GetRequestStatusByReqIDServices(requestID string) (models.StatusResponseByRequestID, error) {
	// Retrieve response from repository
	response, err := s.r.GetRequestStatusByReqIDRepositories(requestID)
	if err != nil {
		log.Println(err)
		return models.StatusResponseByRequestID{}, err
	}

	// Create a new slice to store decrypted information
	decryptedApprovals := make([]models.ApproverStatus, 0, len(response.Approvals))

	// Loop through every approver and decrypt or detokenize the username
	for _, approval := range response.Approvals {
		// Initialize as nil which will be JSON null if no changes occur
		var decryptedApprover *string
		var decryptedStatus *string
		var decryptedLevel *int

		if approval.Approver != nil && *approval.Approver != "" {
			decryptedEmail, err := decrypt.DetokenizationEmailForMasking(*approval.Approver)
			if err != nil {
				log.Println(err)
				continue
			}
			decodedEmailBytes, err := base64.StdEncoding.DecodeString(decryptedEmail)
			if err != nil {
				log.Println(err)
				continue
			}
			decodedEmailStr := string(decodedEmailBytes)
			decryptedApprover = &decodedEmailStr
			decryptedStatus = approval.Status
			decryptedLevel = approval.Level
		}

		// Append the decrypted approver status to the slice
		decryptedApprovals = append(decryptedApprovals, models.ApproverStatus{
			Approver: decryptedApprover,
			Status:   decryptedStatus,
			Level:    decryptedLevel,
		})
	}

	// Update the slice in the response with decrypted and decoded approvals
	response.Approvals = decryptedApprovals
	return response, nil
}

func (s *serviceAdapter) AddActivity(req models.RequestActivity) error {
	if err := s.r.AddActivity(req); err != nil {
		return err
	}
	return nil
}

func (s *serviceAdapter) GETUser_File(scdact_id string) ([][]byte, error) {
	// เรียกใช้ GETUser_File ใน Repository
	scdactBinary, err := s.r.GETUser_File(scdact_id)
	if err != nil {
		return nil, err
	}

	return scdactBinary, nil
}

var (
	BaseDirectory     = "/home/chiccrm/Desktop/chicCRM/fianlCodeConfig"
	Datapath          = filepath.Join(BaseDirectory, "API3.40R02.0001/64bit/data")
	Cmdfinal_codepath = filepath.Join(BaseDirectory, "API3.40R02.0001/64bit/bin")
)

func (s *serviceAdapter) ExecuteCommand(cmdStr string) (string, int, error) {

	// cmdArgs, err := shellquote.Split(cmdStr)
	// if err != nil {
	//  return "", 0, err
	// }

	fmt.Println("cmd", cmdStr)

	// cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...) // in case linux
	cmd := exec.Command("bash", "-c", cmdStr) // in case linux

	// กำหนดไดเร็กทอรีที่ให้ command ทำงาน
	cmd.Dir = Cmdfinal_codepath

	// รับผลลัพธ์จากการรันคำสั่ง (รับ Output และ Error แยก)
	output, err := cmd.CombinedOutput()

	// ตรวจสอบว่ามี error หรือไม่
	exitCode := 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		// ถ้ามี error, ดึง exit code จาก ExitError
		exitCode = exitErr.ExitCode()
	}

	// ส่งผลลัพธ์กลับ
	return string(output), exitCode, err
}
func (s *serviceAdapter) CFolder(orderID string) error {
	mainFolderPath := filepath.Join(Datapath, orderID)
	if err := os.MkdirAll(mainFolderPath, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (s *serviceAdapter) FilesUpload(c *gin.Context, file *multipart.FileHeader, orderID string) error {
	filePath := filepath.Join(Datapath, orderID, file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return err
	}
	return nil
}
func (s *serviceAdapter) DeleteFolder(orderID string) error {
	folderPath := filepath.Join(Datapath, orderID)
	err := os.RemoveAll(folderPath)
	if err != nil {
		// การลบโฟลเดอร์ไม่สำเร็จ ทำการ handle ตามที่คุณต้องการ
		return err
	}
	return nil
}

func (s *serviceAdapter) CFolderEmail(orderID, emailInput string) error {
	// ใช้ orderID เป็น main folder
	mainFolderPath := filepath.Join(Datapath, orderID)
	if _, err := os.Stat(mainFolderPath); os.IsNotExist(err) {
		// ถ้า Folder ไม่มีอยู่ ให้ส่ง error
		return errors.New("order_id folder does not exist")
	}

	// แยก emailInput เป็นรายการของอีเมล
	emailList := strings.Split(emailInput, ",")

	for _, email := range emailList {
		email = strings.TrimSpace(email)

		// สร้าง subfolder ด้วยชื่ออีเมลภายใน main folder
		emailFolderPath := filepath.Join(mainFolderPath, email)
		if err := os.MkdirAll(emailFolderPath, os.ModePerm); err != nil {
			return err
		}

		// คัดลอกไฟล์ HTML ไปยัง subfolder
		if err := copyHTMLFiles(Datapath, orderID, emailFolderPath, email); err != nil {
			return err
		}
	}

	return nil
}

func copyHTMLFiles(Datapath, orderID, destDir, email string) error {
	// ค้นหาไฟล์ HTML และ FCL ในไดเรกทอรีต้นทาง
	htmlFiles, err := filepath.Glob(filepath.Join(Datapath, orderID, "*.html"))
	if err != nil {
		return err
	}

	fclFiles, err := filepath.Glob(filepath.Join(Datapath, orderID, "*.fcl"))
	if err != nil {
		return err
	}

	allFiles := append(htmlFiles, fclFiles...)

	for _, file := range allFiles {
		// ดึงชื่อไฟล์โดยไม่รวมข้อมูลอีเมล
		fileName := filepath.Base(file)
		startIndex := strings.Index(fileName, "(")
		endIndex := strings.LastIndex(fileName, ")")
		if startIndex != -1 && endIndex != -1 && endIndex > startIndex {
			trimmedFileName := fileName[:startIndex] + fileName[endIndex+1:]

			// ตรวจสอบว่าชื่อไฟล์มีข้อมูลอีเมลหรือไม่
			if strings.Contains(fileName, email) {
				destPath := filepath.Join(destDir, trimmedFileName+".zip")

				// สร้าง zip ไฟล์แยกต่างหากสำหรับแต่ละไฟล์
				zipFile, err := os.Create(destPath)
				if err != nil {
					return err
				}
				defer zipFile.Close()

				zipWriter := zip.NewWriter(zipFile)
				defer zipWriter.Close()

				// เพิ่มไฟล์ลงใน zip
				zipEntry, err := zipWriter.Create(trimmedFileName)
				if err != nil {
					return err
				}

				fileContent, err := ioutil.ReadFile(file)
				if err != nil {
					return err
				}

				_, err = zipEntry.Write(fileContent)
				if err != nil {
					return err
				}

				// สแกนไฟล์ที่คัดลอกเพื่อค้นหาเนื้อหาที่เกี่ยวข้องกับอีเมล
				if err := copy.ScanFileForEmail(destPath, email); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (s *serviceAdapter) AttachFile(orderID, baseDir string, emailList []string) (map[string][]string, error) {
	AttachFile := make(map[string][]string)

	for _, email := range emailList {
		// ใช้ orderID เป็นโฟลเดอร์หลัก
		mainFolderPath := filepath.Join(baseDir, orderID)

		// เพิ่ม email เป็นโฟลเดอร์ย่อย
		folderPath := filepath.Join(mainFolderPath, email)

		dirEntries, err := os.ReadDir(folderPath)
		if err != nil {
			return nil, err
		}

		var files []string
		for _, dirEntry := range dirEntries {
			files = append(files, filepath.Join(folderPath, dirEntry.Name()))
		}

		AttachFile[email] = files
	}
	return AttachFile, nil
}

func (s *serviceAdapter) AddStatus(ScdactStatus string, ScdactReqID string, teamleadID string) (int, error) {
	// ดำเนินการเพิ่มข้อมูลลงในฐานข้อมูลโดยใช้ Repository
	level, err := s.r.AddStatus(ScdactStatus, ScdactReqID, teamleadID)
	if err != nil {
		return 0, err
	}
	return level, nil
}
