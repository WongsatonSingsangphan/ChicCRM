package handlers

import (
	"chicCRM/pkg/utilts/uploadBinary/models"
	"chicCRM/pkg/utilts/uploadBinary/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerPort interface {
	UploadBinaryChicCRMHandlers(c *gin.Context)
	GetBinaryChicCRMHandlders(c *gin.Context)
	EditLogoProfileChiCCRMHandlers(c *gin.Context)
	EditPersonalProfileChicCRMHandlers(c *gin.Context)
	GetPersonalProfileChicCRMHandlers(c *gin.Context)
}

type handlerAdapter struct {
	s services.ServicePort
}

func NewHanerhandlerAdapter(s services.ServicePort) HandlerPort {
	return &handlerAdapter{s: s}
}

func (h *handlerAdapter) UploadBinaryChicCRMHandlers(c *gin.Context) {
	var requestOrganizeBinary models.RequestOrganizeBinary
	if err := c.ShouldBind(&requestOrganizeBinary); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.s.UploadBinaryChicCRMServices(requestOrganizeBinary, form.File["file"])
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "OK", "message": "Registered Successfully. You can edit the company profile later", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Registered Successfully"})
}

func (h *handlerAdapter) GetBinaryChicCRMHandlders(c *gin.Context) {
	id := c.Param("organizeID")
	fileBytes, fileType, err := h.s.GetBinaryChicCRMServices(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(fileBytes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No file found"})
		return
	}
	c.Header("Content-Type", fileType)
	c.Header("Content-Disposition", "inline")
	c.Data(http.StatusOK, fileType, fileBytes)
}

func (h *handlerAdapter) EditLogoProfileChiCCRMHandlers(c *gin.Context) {
	var editLogoRequest models.EditLogoRequest
	if err := c.ShouldBind(&editLogoRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.s.EditLogoProfileChiCCRMServices(editLogoRequest, form.File["file"])
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Edit Company Profile Succesfully."})
}

func (h *handlerAdapter) EditPersonalProfileChicCRMHandlers(c *gin.Context) {
	var editPersonalProfileRequest models.EditPersonalProfileRequest
	if err := c.ShouldBind(&editPersonalProfileRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.s.EditPersonalProfileChicCRMServices(editPersonalProfileRequest, form.File["file"])
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Update personal Profile Successfully."})
}

func (h *handlerAdapter) GetPersonalProfileChicCRMHandlers(c *gin.Context) {
	id := c.Param("personalID")
	fileBytes, fileType, err := h.s.GetPersonalProfileChicCRMServices(id)
	if err != nil {
		switch err.Error() {
		case "personalID does not match":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err.Error()})
		}
		return
	}
	if len(fileBytes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No file found"})
		return
	}
	c.Header("Content-Type", fileType)
	c.Header("Content-Disposition", "inline")
	c.Data(http.StatusOK, fileType, fileBytes)
}
