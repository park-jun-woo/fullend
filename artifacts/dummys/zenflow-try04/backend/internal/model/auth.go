package model

// CurrentUser is the authenticated user extracted by JWT middleware.
type CurrentUser struct {
	Email string
	Role string
	UserID int64
}

// Authorizer checks permissions.
type Authorizer interface {
	Check(user *CurrentUser, action, resource string, input interface{}) error
}
