package schedule

import (
	"fmt"

	"github.com/park-jun-woo/fullend/pkg/session"
)

// @func getSchedule
// @description 워크플로우의 현재 스케줄을 session에서 조회한다

type GetScheduleRequest struct {
	WorkflowID int64
}

type GetScheduleResponse struct {
	Cron    string
	NextRun string
}

func GetSchedule(req GetScheduleRequest) (GetScheduleResponse, error) {
	key := fmt.Sprintf("schedule:%d", req.WorkflowID)
	resp, err := session.Get(session.GetRequest{Key: key})
	if err != nil {
		return GetScheduleResponse{}, nil
	}
	return GetScheduleResponse{Cron: resp.Value, NextRun: ""}, nil
}
