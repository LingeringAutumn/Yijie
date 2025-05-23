// Code generated by Kitex v0.12.3. DO NOT EDIT.

package videoservice

import (
	"context"
	"errors"

	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"

	video "github.com/LingeringAutumn/Yijie/kitex_gen/video"
)

var errInvalidMessageType = errors.New("invalid message type for service method handler")

var serviceMethods = map[string]kitex.MethodInfo{
	"SubmitVideo": kitex.NewMethodInfo(
		submitVideoHandler,
		newVideoServiceSubmitVideoArgs,
		newVideoServiceSubmitVideoResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"GetVideo": kitex.NewMethodInfo(
		getVideoHandler,
		newVideoServiceGetVideoArgs,
		newVideoServiceGetVideoResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"SearchVideo": kitex.NewMethodInfo(
		searchVideoHandler,
		newVideoServiceSearchVideoArgs,
		newVideoServiceSearchVideoResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"TrendVideo": kitex.NewMethodInfo(
		trendVideoHandler,
		newVideoServiceTrendVideoArgs,
		newVideoServiceTrendVideoResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"UpdateVideoHot": kitex.NewMethodInfo(
		updateVideoHotHandler,
		newVideoServiceUpdateVideoHotArgs,
		newVideoServiceUpdateVideoHotResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
}

var (
	videoServiceServiceInfo                = NewServiceInfo()
	videoServiceServiceInfoForClient       = NewServiceInfoForClient()
	videoServiceServiceInfoForStreamClient = NewServiceInfoForStreamClient()
)

// for server
func serviceInfo() *kitex.ServiceInfo {
	return videoServiceServiceInfo
}

// for stream client
func serviceInfoForStreamClient() *kitex.ServiceInfo {
	return videoServiceServiceInfoForStreamClient
}

// for client
func serviceInfoForClient() *kitex.ServiceInfo {
	return videoServiceServiceInfoForClient
}

// NewServiceInfo creates a new ServiceInfo containing all methods
func NewServiceInfo() *kitex.ServiceInfo {
	return newServiceInfo(false, true, true)
}

// NewServiceInfo creates a new ServiceInfo containing non-streaming methods
func NewServiceInfoForClient() *kitex.ServiceInfo {
	return newServiceInfo(false, false, true)
}
func NewServiceInfoForStreamClient() *kitex.ServiceInfo {
	return newServiceInfo(true, true, false)
}

func newServiceInfo(hasStreaming bool, keepStreamingMethods bool, keepNonStreamingMethods bool) *kitex.ServiceInfo {
	serviceName := "VideoService"
	handlerType := (*video.VideoService)(nil)
	methods := map[string]kitex.MethodInfo{}
	for name, m := range serviceMethods {
		if m.IsStreaming() && !keepStreamingMethods {
			continue
		}
		if !m.IsStreaming() && !keepNonStreamingMethods {
			continue
		}
		methods[name] = m
	}
	extra := map[string]interface{}{
		"PackageName": "video",
	}
	if hasStreaming {
		extra["streaming"] = hasStreaming
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.12.3",
		Extra:           extra,
	}
	return svcInfo
}

func submitVideoHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*video.VideoServiceSubmitVideoArgs)
	realResult := result.(*video.VideoServiceSubmitVideoResult)
	success, err := handler.(video.VideoService).SubmitVideo(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newVideoServiceSubmitVideoArgs() interface{} {
	return video.NewVideoServiceSubmitVideoArgs()
}

func newVideoServiceSubmitVideoResult() interface{} {
	return video.NewVideoServiceSubmitVideoResult()
}

func getVideoHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*video.VideoServiceGetVideoArgs)
	realResult := result.(*video.VideoServiceGetVideoResult)
	success, err := handler.(video.VideoService).GetVideo(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newVideoServiceGetVideoArgs() interface{} {
	return video.NewVideoServiceGetVideoArgs()
}

func newVideoServiceGetVideoResult() interface{} {
	return video.NewVideoServiceGetVideoResult()
}

func searchVideoHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*video.VideoServiceSearchVideoArgs)
	realResult := result.(*video.VideoServiceSearchVideoResult)
	success, err := handler.(video.VideoService).SearchVideo(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newVideoServiceSearchVideoArgs() interface{} {
	return video.NewVideoServiceSearchVideoArgs()
}

func newVideoServiceSearchVideoResult() interface{} {
	return video.NewVideoServiceSearchVideoResult()
}

func trendVideoHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*video.VideoServiceTrendVideoArgs)
	realResult := result.(*video.VideoServiceTrendVideoResult)
	success, err := handler.(video.VideoService).TrendVideo(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newVideoServiceTrendVideoArgs() interface{} {
	return video.NewVideoServiceTrendVideoArgs()
}

func newVideoServiceTrendVideoResult() interface{} {
	return video.NewVideoServiceTrendVideoResult()
}

func updateVideoHotHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*video.VideoServiceUpdateVideoHotArgs)
	realResult := result.(*video.VideoServiceUpdateVideoHotResult)
	success, err := handler.(video.VideoService).UpdateVideoHot(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newVideoServiceUpdateVideoHotArgs() interface{} {
	return video.NewVideoServiceUpdateVideoHotArgs()
}

func newVideoServiceUpdateVideoHotResult() interface{} {
	return video.NewVideoServiceUpdateVideoHotResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) SubmitVideo(ctx context.Context, req *video.VideoSubmissionRequest) (r *video.VideoSubmissionResponse, err error) {
	var _args video.VideoServiceSubmitVideoArgs
	_args.Req = req
	var _result video.VideoServiceSubmitVideoResult
	if err = p.c.Call(ctx, "SubmitVideo", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetVideo(ctx context.Context, req *video.VideoDetailRequest) (r *video.VideoDetailResponse, err error) {
	var _args video.VideoServiceGetVideoArgs
	_args.Req = req
	var _result video.VideoServiceGetVideoResult
	if err = p.c.Call(ctx, "GetVideo", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) SearchVideo(ctx context.Context, req *video.VideoSearchRequest) (r *video.VideoSearchResponse, err error) {
	var _args video.VideoServiceSearchVideoArgs
	_args.Req = req
	var _result video.VideoServiceSearchVideoResult
	if err = p.c.Call(ctx, "SearchVideo", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) TrendVideo(ctx context.Context, req *video.VideoTrendingRequest) (r *video.VideoTrendingResponse, err error) {
	var _args video.VideoServiceTrendVideoArgs
	_args.Req = req
	var _result video.VideoServiceTrendVideoResult
	if err = p.c.Call(ctx, "TrendVideo", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) UpdateVideoHot(ctx context.Context, req *video.VideoHotUpdateRequest) (r *video.VideoHotUpdateResponse, err error) {
	var _args video.VideoServiceUpdateVideoHotArgs
	_args.Req = req
	var _result video.VideoServiceUpdateVideoHotResult
	if err = p.c.Call(ctx, "UpdateVideoHot", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
