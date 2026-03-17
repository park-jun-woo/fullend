package template

import (
	"github.com/example/zenflow/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

//fullend:gen ssot=service/template/list_templates.ssac contract=f585df2
func (h *Handler) ListTemplates(c *gin.Context) {
	opts := model.ParseQueryOpts(c, model.QueryOptsConfig{
		Pagination: &model.PaginationConfig{Style: "cursor", DefaultLimit: 20, MaxLimit: 100},
		Filter:     &model.FilterConfig{Allowed: []string{"category"}},
	})

	tmplCursor, err := h.TemplateModel.List(opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Template 조회 실패"})
		return
	}

	c.JSON(http.StatusOK, tmplCursor)

}
