package sys_message

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"backend/api/sys_message/v1"
)

func (c *ControllerV1) SetRead(ctx context.Context, req *v1.SetReadReq) (res *v1.SetReadRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
