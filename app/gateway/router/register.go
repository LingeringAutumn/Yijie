// Code generated by hertz generator. DO NOT EDIT.

package router

import (
	api_user "Yijie/app/gateway/router/api/user"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// GeneratedRegister registers routers generated by IDL.
func GeneratedRegister(r *server.Hertz) {
	//INSERT_POINT: DO NOT DELETE THIS LINE!
	api_user.Register(r)
}
