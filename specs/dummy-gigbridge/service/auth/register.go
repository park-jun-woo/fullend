package service

import "github.com/geul-org/fullend/pkg/auth"

// @call string hashedPassword = auth.HashPassword({Password: request.Password})
// @post User user = User.Create(request.Email, hashedPassword, request.Role, request.Name)
// @response {
//   user: user
// }
func Register() {}
