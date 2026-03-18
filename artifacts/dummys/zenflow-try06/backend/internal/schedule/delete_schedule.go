package schedule

import (
	"fmt"

	"github.com/park-jun-woo/fullend/pkg/session"
)

// @func deleteSchedule
// @description 워크플로우의 스케줄을 session에서 삭제한다

type DeleteScheduleRequest struct {
	WorkflowID int64
}

type DeleteScheduleResponse struct{}

func DeleteSchedule(req DeleteScheduleRequest) (DeleteScheduleResponse, error) {
	key := fmt.Sprintf("schedule:%d", req.WorkflowID)
	_, err := session.Delete(session.DeleteRequest{Key: key})
	return DeleteScheduleResponse{}, err
}
