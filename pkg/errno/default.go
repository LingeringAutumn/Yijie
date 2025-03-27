// This file is be designed to define any common error so that we can use it in any service simply.

package errno

// 一些常用的的错误
var (
	Success = NewErrNo(SuccessCode, "ok")

	ParamVerifyError  = NewErrNo(ParamVerifyErrorCode, "parameter validation failed")
	ParamMissingError = NewErrNo(ParamMissingErrorCode, "missing parameter")

	AuthInvalid             = NewErrNo(AuthInvalidCode, "authentication failure")
	AuthAccessExpired       = NewErrNo(AuthAccessExpiredCode, "token expiration")
	AuthNoToken             = NewErrNo(AuthNoTokenCode, "lack of token")
	AuthNoOperatePermission = NewErrNo(AuthNoOperatePermissionCode, "No permission to operate")

	InternalServiceError = NewErrNo(InternalServiceErrorCode, "internal server error")
)

var (
	FileUploadError = NewErrNo(UploadFileFailed, "uploaded file is not exist")
	// FileSaveError   = NewErrNo(FileSaveFailed, "uploaded file is not exist")
)
