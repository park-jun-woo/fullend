package workflow

import "fmt"

// @func nextVersion
// @description Increments the workflow version number

type NextVersionRequest struct {
	CurrentVersion int64
}

type NextVersionResponse struct {
	NextVersion int64
}

func NextVersion(req NextVersionRequest) (NextVersionResponse, error) {
	if req.CurrentVersion < 0 {
		return NextVersionResponse{}, fmt.Errorf("invalid version number: %d", req.CurrentVersion)
	}
	next := req.CurrentVersion + 1
	return NextVersionResponse{NextVersion: next}, nil
}
