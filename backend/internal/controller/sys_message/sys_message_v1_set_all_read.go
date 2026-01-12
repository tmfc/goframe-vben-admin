package sys_message

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"backend/api/sys_message/v1"
)

func (c *ControllerV1) SetAllRead(ctx context.Context, req *v1.SetAllReadReq) (res *v1.SetAllReadRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
