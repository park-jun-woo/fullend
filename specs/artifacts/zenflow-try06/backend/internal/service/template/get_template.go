package template

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//fullend:gen ssot=service/template/get_template.ssac contract=b62e85c
func (h *Handler) GetTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid path parameter"})
		return
	}

	tmpl, err := h.TemplateModel.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Template 조회 실패"})
		return
	}

	if tmpl == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	org, err := h.OrganizationModel.FindByID(tmpl.OrgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Organization 조회 실패"})
		return
	}

	if org == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"org_name": org.Name,
		"template": tmpl,
	})

}
