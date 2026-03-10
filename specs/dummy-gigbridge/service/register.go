package service

import "github.com/gigbridge/api/auth"

// @call string hashedPassword = auth.HashPassword(request.Password)
// @post User user = User.Create(request.Email, hashedPassword, request.Role, request.Name)
// @response {
//   user: user
// }
func Register() {}
