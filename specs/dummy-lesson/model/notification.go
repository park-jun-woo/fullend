package model

// Notification은 알림 발송 컴포넌트다.
// SSaC에서 @component notification으로 참조된다.
type Notification interface {
	Send(userID int64, message string) error
}
