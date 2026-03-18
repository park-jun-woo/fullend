package schedule

import "github.com/example/zenflow/internal/model"

// Handler handles requests for the schedule domain.
type Handler struct {
	WorkflowModel model.WorkflowModel
}
