package report

import "fmt"

// @func generateReport
// @description Generates a report key and content for workflow execution

type GenerateReportRequest struct {
	WorkflowID int64
	Status     string
}

type GenerateReportResponse struct {
	ReportKey string
}

func GenerateReport(req GenerateReportRequest) (GenerateReportResponse, error) {
	if req.WorkflowID <= 0 {
		return GenerateReportResponse{}, fmt.Errorf("invalid workflow ID: %d", req.WorkflowID)
	}
	key := fmt.Sprintf("reports/wf-%d-%s.txt", req.WorkflowID, req.Status)
	return GenerateReportResponse{ReportKey: key}, nil
}
