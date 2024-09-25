// handlers/handlers.go
package handlers

import (
	"net/http"
	"test-go/models"

	"github.com/gin-gonic/gin"
)

// CreateDynamicCredentialHandler handles POST /dyncreds
func CreateDynamicCredentialHandler(c *gin.Context) {
	var req models.CreateDynamicCredentialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cred, err := services.CreateDynamicCredential(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create dynamic credential"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Dynamic credential created successfully",
		"dyncred": cred,
	})
}

// GetDynamicCredentialHandler handles GET /dyncreds/:dyncredId
func GetDynamicCredentialHandler(c *gin.Context) {
	id := c.Param("dyncredId")
	cred, err := services.GetDynamicCredential(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"dyncred": cred,
	})
}

// UpdateDynamicCredentialHandler handles PUT /dyncreds/:dyncredId
func UpdateDynamicCredentialHandler(c *gin.Context) {
	id := c.Param("dyncredId")
	var req models.UpdateDynamicCredentialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cred, err := services.UpdateDynamicCredential(id, req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Dynamic credential updated successfully",
		"dyncred": cred,
	})
}

// DeleteDynamicCredentialHandler handles DELETE /dyncreds/:dyncredId
func DeleteDynamicCredentialHandler(c *gin.Context) {
	id := c.Param("dyncredId")
	err := services.DeleteDynamicCredential(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Dynamic credential deleted successfully",
		"dyncredId": id,
	})
}

// PatchDynamicCredentialHandler handles PATCH /dyncreds/:dyncredId
func PatchDynamicCredentialHandler(c *gin.Context) {
	id := c.Param("dyncredId")
	var req models.UpdateTTLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update TTL in the credential
	cred, err := services.UpdateDynamicCredentialTTL(id, req.TTL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Update TTL across all Terraform workspaces
	err = services.UpdateTTLForAllWorkspaces(id, req.TTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update TTL in workspaces"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "TTL updated successfully for all workspaces",
		"dyncredId": id,
		"ttl":       cred.TTL,
	})
}
