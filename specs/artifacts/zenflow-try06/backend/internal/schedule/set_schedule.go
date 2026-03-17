package schedule

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/fullend/pkg/session"
)

// @func setSchedule
// @error 400
// @description 워크플로우에 cron 스케줄을 등록하고 session에 저장한다

type SetScheduleRequest struct {
	WorkflowID int64
	Cron       string
}

type SetScheduleResponse struct {
	Cron    string
	NextRun string
}

func SetSchedule(req SetScheduleRequest) (SetScheduleResponse, error) {
	parts := strings.Fields(req.Cron)
	if len(parts) != 5 {
		return SetScheduleResponse{}, fmt.Errorf("invalid cron expression: expected 5 fields, got %d", len(parts))
	}
	key := fmt.Sprintf("schedule:%d", req.WorkflowID)
	_, err := session.Set(session.SetRequest{Key: key, Value: req.Cron, TTL: 0})
	if err != nil {
		return SetScheduleResponse{}, err
	}
	return SetScheduleResponse{Cron: req.Cron, NextRun: "2026-03-19T00:00:00Z"}, nil
}
