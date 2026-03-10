package service

import "github.com/gigbridge/api/auth"

// @get User user = User.FindByEmail(request.Email)
// @empty user "User not found"
// @call auth.VerifyPassword(user.PasswordHash, request.Password)
// @call string accessToken = auth.IssueToken(user.ID, user.Email, user.Role)
// @response {
//   accessToken: accessToken
// }
func Login() {}
