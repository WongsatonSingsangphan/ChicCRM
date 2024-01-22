package handlers

import (
	"chicCRM/modules/finalCode/handlers/errors"
	"chicCRM/modules/finalCode/models"
	"chicCRM/modules/finalCode/services"
	"chicCRM/modules/finalCode/services/mail"
	"database/sql"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type HandlerPort interface {
	CheckFeatureOrganizeHandlers(c *gin.Context)
	CheckFeatureMemberHandlers(c *gin.Context)
	GetPolicyDataHandlers(c *gin.Context)
	GetSecuredocActivityByOrganizeMemberUUIDHandlers(c *gin.Context)
	GetSecuredocActivityByTeamleadIDHandlers(c *gin.Context)
	GetPolicyAuthorizationByTeamleadHandlers(c *gin.Context)
	GetRequestStatusByReqIDHandlers(c *gin.Context)
	AddActivity(c *gin.Context)
	GETUser_File(c *gin.Context)
	File_encrypt(c *gin.Context)
	Send_mail(c *gin.Context)
}

type handlerAdapter struct {
	s  services.ServicePort
	db *sql.DB // add db for h.db
}

//	func NewHanerhandlerAdapter(s services.ServicePort) HandlerPort {
//		return &handlerAdapter{s: s}
//	}
func NewHanerhandlerAdapter(s services.ServicePort, db *sql.DB) HandlerPort {
	return &handlerAdapter{s: s, db: db}
}

// Check Tenant Feature Organize
func (h *handlerAdapter) CheckFeatureOrganizeHandlers(c *gin.Context) {
	// เปลี่ยนที่นี่เพื่อดึงค่า ID จาก URL
	id := c.Param("id")
	println(id)
	data, err := h.s.CheckFeatureOrganizeServices(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		return
	}
	// ส่งข้อมูลกลับไป
	c.IndentedJSON(http.StatusOK, gin.H{"data": data})
}

// Check Tenant Feature Member
func (h *handlerAdapter) CheckFeatureMemberHandlers(c *gin.Context) {
	id := c.Param("id")

	tenantFeatureMember, memberAuthorization, err := h.s.CheckFeatureMemberServices(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		return
	}

	// ส่งข้อมูลกลับไป
	c.IndentedJSON(http.StatusOK, gin.H{
		"tenantFeatureMember": tenantFeatureMember,
		"memberAuthorization": memberAuthorization,
	})
}

// Get Policy Data
func (h *handlerAdapter) GetPolicyDataHandlers(c *gin.Context) {
	id := c.Param("id")
	println(id)

	data, err := h.s.GetPolicyDataServices(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		return
	}
	// ส่งข้อมูลกลับไป
	c.IndentedJSON(http.StatusOK, gin.H{"data": data})
}

func (h *handlerAdapter) GetSecuredocActivityByOrganizeMemberUUIDHandlers(c *gin.Context) {
	organizeMemberUUID := c.Param("uuid")
	dataResponse, err := h.s.GetSecuredocActivityByOrganizeMemberUUIDServices(organizeMemberUUID)
	if err != nil {
		switch err.Error() {
		case "memberID does not match":
			c.JSON(http.StatusNotFound, gin.H{"status": "Error", "message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "logSecuredocActivityMember": dataResponse})
}

func (h *handlerAdapter) GetSecuredocActivityByTeamleadIDHandlers(c *gin.Context) {
	teamleadID := c.Param("uuid")
	dataResponse, err := h.s.GetSecuredocActivityByTeamleadIDServices(teamleadID)
	if err != nil {
		switch err.Error() {
		case "users do not have permission to access this information; only managers and founders are authorized to do so":
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Error", "message": err.Error()})
		case "user not found":
			c.JSON(http.StatusNotFound, gin.H{"status": "Error", "message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "logSecuredogActivityTeamlead": dataResponse})
}

func (h *handlerAdapter) GetPolicyAuthorizationByTeamleadHandlers(c *gin.Context) {
	teamleadID := c.Param("uuid")
	dataResponse, err := h.s.GetPolicyAuthorizationByTeamleadServices(teamleadID)
	if err != nil {
		switch err.Error() {
		case "organizeID does not match":
			c.JSON(http.StatusNotFound, gin.H{"status": "Error", "message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "policyByTeamlead": dataResponse})
}

func (h *handlerAdapter) GetRequestStatusByReqIDHandlers(c *gin.Context) {
	requestID := c.Param("requestID")
	// fmt.Println(requestID)
	dataResponse, err := h.s.GetRequestStatusByReqIDServices(requestID)
	if err != nil {
		switch err.Error() {
		case "requestID not found":
			c.JSON(http.StatusNotFound, gin.H{"status": "Error", "message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "StatusByRequestID": dataResponse})
}

func (h *handlerAdapter) AddActivity(c *gin.Context) {
	// ดึงข้อมูลทั้งหมดจากฟอร์ม multipart
	if err := c.Request.ParseMultipartForm(0); err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseStatus{
			Status:  "Error",
			Message: "Failed to parse multipart form",
		})
		return
	}

	// ตรวจสอบและดึงข้อมูล scdact_binary และ scdact_command จากฟอร์ม
	files := c.Request.MultipartForm.File["scdact_binary"]
	commands := c.Request.MultipartForm.Value["scdact_command"]

	if len(files) != len(commands) {
		c.JSON(http.StatusBadRequest, models.ResponseStatus{
			Status:  "Error",
			Message: "Mismatched number of files and commands",
		})
		return
	}

	// วนลูปผ่านไฟล์ที่ถูกส่งมา
	for i, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ResponseStatus{
				Status:  "Error",
				Message: "Failed to open file",
			})
			return
		}
		defer file.Close()

		reqBinary, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ResponseStatus{
				Status:  "Error",
				Message: "Failed to read file content",
			})
			return
		}

		// สร้าง RequestActivity ใหม่ในแต่ละลูป
		req := models.RequestActivity{
			ScdactBinary:       reqBinary,
			ScdactCommand:      commands[i],
			ScdactFilename:     c.Request.MultipartForm.Value["scdact_filename"][i],
			ScdactFiletype:     c.Request.MultipartForm.Value["scdact_filetype"][i],
			ScdactFilehash:     c.Request.MultipartForm.Value["scdact_filehash"][i],
			ScdactFilesize:     c.Request.MultipartForm.Value["scdact_filesize"][i],
			ScdactFilecreated:  c.Request.MultipartForm.Value["scdact_filecreated"][i],
			ScdactFilemodified: c.Request.MultipartForm.Value["scdact_filemodified"][i],

			ScdactAction: c.PostForm("scdact_action"),

			OrgmbatOrgmbid: c.PostForm("uuid_member"),

			ScdactStatus:         c.PostForm("scdact_status"),
			ScdactReqID:          c.PostForm("scdact_reqid"),
			ScdactFilelocation:   c.Request.MultipartForm.Value["scdact_filelocation"][i],
			ScdactName:           c.PostForm("scdact_name"),
			ScdactType:           c.PostForm("scdact_type"),
			ScdactReciepient:     c.PostForm("scdact_reciepient"),
			ScdactStartTime:      c.PostForm("scdact_starttime"),
			ScdactEndTime:        c.PostForm("scdact_endtime"),
			ScdactPeriodDay:      c.PostForm("scdact_periodday"),
			ScdactPeriodHour:     c.PostForm("scdact_periodhour"),
			ScdactNumberOpen:     c.PostForm("scdact_numberopen"),
			ScdactNoLimit:        c.PostForm("scdact_nolimit") == "true",
			ScdactCvtOriginal:    c.PostForm("scdact_cvtoriginal") == "true",
			ScdactEdit:           c.PostForm("scdact_edit") == "true",
			ScdactPrint:          c.PostForm("scdact_print") == "true",
			ScdactCopy:           c.PostForm("scdact_copy") == "true",
			ScdactScrWatermark:   c.PostForm("scdact_scrwatermark") == "true",
			ScdactWatermark:      c.PostForm("scdact_watermark") == "true",
			ScdactCvtHtml:        c.PostForm("scdact_cvthtml") == "true",
			ScdactCvtFcl:         c.PostForm("scdact_cvtfcl") == "true",
			ScdactMarcro:         c.PostForm("scdact_marcro") == "true",
			ScdactEnableCv:       c.PostForm("scdact_enableconvertoriginal") == "true",
			ScdactMsgText:        c.PostForm("scdact_msgtext"),
			ScdactSubject:        c.PostForm("scdact_subject"),
			ScdactSender:         c.PostForm("scdact_sender"),
			ScdactCreateLocation: c.PostForm("scdact_createlocation"),
			ScdactUpdateLocation: c.PostForm("scdact_updatelocation"),
			ScdactActionTime:     c.PostForm("scdact_actiontime"),
		}

		// Call the AddActivity method with the populated RequestActivity instance
		if err := h.s.AddActivity(req); err != nil {
			c.JSON(http.StatusInternalServerError, models.ResponseStatus{
				Status:  "Error",
				Message: err.Error(),
			})
			return
		}
	}

	// ส่ง response สำเร็จ
	c.JSON(http.StatusOK, models.ResponseStatus{
		Status:  "OK",
		Message: "Activity Successfully",
	})
}

func (h *handlerAdapter) GETUser_File(c *gin.Context) {
	// ดึงค่า UUID จาก URL
	scdactID := c.Param("scdact_id")

	// เรียกใช้ GETUser_File ใน services
	scdactBinary, err := h.s.GETUser_File(scdactID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseStatus{
			Status:  "Error",
			Message: err.Error(),
		})
		return
	}

	// ตรวจสอบ Content-Type ของไฟล์
	contentType := http.DetectContentType(scdactBinary[0])

	// ตั้งค่า Content-Type ใน header
	c.Header("Content-Type", contentType)

	// ส่ง binary data กลับไปที่ client
	c.Data(http.StatusOK, "application/octet-stream", scdactBinary[0])
}

func (h *handlerAdapter) File_encrypt(c *gin.Context) {
	// ดึงข้อมูล order_id จาก request
	orderID := c.PostForm("order_id")

	err := h.s.CFolder(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseStatus{
			Status:  "Error",
			Message: err.Error(),
		})
		return
	}

	// ดึงข้อมูลไฟล์จาก Multipart Form
	files, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseStatus{
			Status:  "Error",
			Message: err.Error(),
		})
		return
	}

	// สำหรับทุกไฟล์ที่อัปโหลด
	for _, file := range files.File["file"] {
		// บันทึกไฟล์ลงในระบบไฟล์ของคุณ
		err := h.s.FilesUpload(c, file, orderID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ResponseStatus{
				Status:  "Error",
				Message: err.Error(),
			})
			return
		}
	}
	// ทำงานกับคำสั่งทั้งหมด
	cmdStrList := c.PostFormArray("cmd")

	// ทำงานกับทุกคำสั่ง
	for _, cmdStr := range cmdStrList {
		_, exitCode, err := h.s.ExecuteCommand(cmdStr)
		if err != nil {
			_, message := errors.ErrorCommand(err, exitCode)
			c.JSON(http.StatusInternalServerError, models.ResponseStatusFC{
				Status: "Error",
				Message: map[string]string{
					"Finalcode_result": message,
				},
			})

			h.s.DeleteFolder(orderID)
			return
		}
	}

	c.JSON(http.StatusOK, models.ResponseStatusFC{
		Status: "OK",
		Message: map[string]string{
			"Finalcode_result": "The operation completed successfully",
		},
	})
}

func (h *handlerAdapter) Send_mail(c *gin.Context) {
	orderID := c.PostForm("order_id")

	// สร้างโฟลเดอร์ตาม order_id
	// err := h.s.CFolder(orderID)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, models.ResponseStatus{
	// 		Status:  "Error",
	// 		Message: err.Error(),
	// 	})
	// 	return
	// }

	ScdactStatus := c.PostForm("action")
	ScdactReqID := c.PostForm("order_id")
	teamleadID := c.PostForm("teamleadID")

	var Level int

	// ตรวจสอบและอัปเดตสถานะเฉพาะเมื่อ ScdactStatus เป็น "Reject" หรือ "Approve"
	if ScdactStatus == "Rejected" || ScdactStatus == "Approved" {
		var err error
		Level, err = h.s.AddStatus(ScdactStatus, ScdactReqID, teamleadID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ResponseStatus{
				Status:  "Error",
				Message: "Cannot update status",
			})
			return
		}
	}

	var status1, status2, status3 string
	approvedCount := 0

	query := `
SELECT scdwflap_status_1, scdwflap_status_2, scdwflap_status_3
FROM securedoc_workflow_approvers
WHERE scdwflap_reqid = $1`
	err := h.db.QueryRow(query, orderID).Scan(&status1, &status2, &status3)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseStatus{
			Status:  "Error",
			Message: "Failed to check approver statuses: " + err.Error(),
		})
		return
	}
	if status1 == "Approved" {
		approvedCount++
	}
	if status2 == "Approved" {
		approvedCount++
	}
	if status3 == "Approved" {
		approvedCount++
	}
	// fmt.Println(Level)

	// ตรวจสอบว่าได้รับ 'Approved' อย่างน้อยสองครั้งหรือไม่
	// if approvedCount < 2 {
	if Level == 3 || approvedCount >= 2 {

		emailInput := c.PostForm("email")

		// สร้างโฟลเดอร์ตามอีเมลที่รับมา
		err = h.s.CFolderEmail(orderID, emailInput)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ResponseStatus{
				Status:  "Error",
				Message: err.Error(),
			})
			return
		}

		emailList := strings.Split(emailInput, ",")

		var AttachFile = make(map[string][]string)
		// เรียกใช้ AttachFile ฟังก์ชัน
		AttachFile, err = h.s.AttachFile(orderID, services.Datapath, emailList)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ResponseStatus{
				Status:  "Error",
				Message: err.Error(),
			})

			h.s.DeleteFolder(orderID)
			return
		}

		// ตรวจสอบสถานะ ScdactStatus ก่อนที่จะส่งอีเมล
		if ScdactStatus != "Reject" {
			sender := c.PostForm("sender")
			subject := c.PostForm("subject")

			mail.SendEmail(sender, subject, emailList, AttachFile)

			if err != nil {
				c.JSON(http.StatusInternalServerError, models.ResponseStatus{
					Status:  "Error",
					Message: err.Error(),
				})

				h.s.DeleteFolder(orderID)
				return
			}
		}

		responseMessage := "Successfully"
		if ScdactStatus == "Reject" {
			responseMessage = "Successfully Rejected"
		}

		c.JSON(http.StatusOK, models.ResponseStatusFC{
			Status: "OK",
			Message: map[string]string{
				"Send_mail": responseMessage,
			},
		})

		h.s.DeleteFolder(orderID)
	} else {
		// ถ้าไม่, ส่ง response กลับไปว่ายังไม่ได้รับการอนุมัติทั้งหมด
		c.JSON(http.StatusOK, models.ResponseStatus{
			Status:  "OK",
			Message: "All approvers have not approved the request.",
		})
		return
	}
}
