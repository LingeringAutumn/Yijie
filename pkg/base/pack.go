package base

import (
	"github.com/LingeringAutumn/Yijie/kitex_gen/model"
	"github.com/LingeringAutumn/Yijie/pkg/errno"
)

var SuccessBase = model.BaseResp{Code: errno.SuccessCode, Msg: errno.SuccessMsg}

func BuildBaseResp(err error) *model.BaseResp {
	if err == nil {
		return &model.BaseResp{
			Code: errno.SuccessCode,
			Msg:  errno.Success.ErrorMsg,
		}
	}
	Errno := errno.ConvertErr(err)
	return &model.BaseResp{
		Code: Errno.ErrorCode,
		Msg:  Errno.ErrorMsg,
	}
}

func BuildSuccessResp() *model.BaseResp {
	return BuildBaseResp(nil) // 直接调用原始函数，传入 nil 表示无错误
}
