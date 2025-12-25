package consts

import (
	"github.com/gogf/gf/v2/errors/gcode"
)

var (
	ErrorCodeUserNotFound         = gcode.New(1001, "User not found", nil)
	ErrorCodeIncorrectPassword    = gcode.New(1002, "Incorrect password", nil)
	ErrorCodeUnauthorized         = gcode.New(1003, "Unauthorized", nil)
	ErrorCodeRefreshTokenRequired = gcode.New(1004, "Refresh token required", nil)
	ErrorCodeRefreshTokenInvalid  = gcode.New(1005, "Refresh token invalid", nil)
)
