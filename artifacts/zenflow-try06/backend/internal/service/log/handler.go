package log

import "github.com/example/zenflow/internal/model"

// Handler handles requests for the log domain.
type Handler struct {
	ExecutionLogModel model.ExecutionLogModel
	WorkflowModel model.WorkflowModel
}
