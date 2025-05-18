package usecase

import (
    ctx "context"
    "testing"
    "time"

    "github.com/bytedance/mockey"
    "github.com/smartystreets/goconvey/convey"

    "github.com/LingeringAutumn/Yijie/app/video/domain/model"
    "github.com/LingeringAutumn/Yijie/app/video/domain/service"
)

func TestVideoUseCase_GetVideo(t *testing.T) {
    defer mockey.UnPatchAll()

    videoID := int64(1001)
    createdAt := time.Now().Unix()

    testCases := []struct {
        Name               string
        MockGetRedisResult *model.VideoProfile
        MockGetRedisError  error
        MockViews          int64
        MockLikes          int64
        MockGetDBResult    *model.VideoProfile
        MockGetDBError     error
        ExpectedHotScore   float64
        ExpectedErr        error
    }{
        {
            Name: "CacheHit",
            MockGetRedisResult: &model.VideoProfile{
                VideoID:   videoID,
                CreatedAt: createdAt,
            },
            MockGetRedisError: nil,
            MockViews:         100,
            MockLikes:         20,
            MockGetDBResult: &model.VideoProfile{
                HotScore: 123.45,
            },
            MockGetDBError:   nil,
            ExpectedHotScore: 123.45,
            ExpectedErr:      nil,
        },
    }

    for _, tc := range testCases {
        mockey.PatchConvey(tc.Name, t, func() {
            mockSvc := new(service.VideoService)
            uc := &videoUseCase{svc: mockSvc}

            mockey.Mock((*service.VideoService).GetVideoRedis).Return(tc.MockGetRedisResult, tc.MockGetRedisError).Build()
            mockey.Mock((*service.VideoService).GetViews).Return(tc.MockViews, nil).Build()
            mockey.Mock((*service.VideoService).GetLikes).Return(tc.MockLikes, nil).Build()
            mockey.Mock((*service.VideoService).GetVideoDB).Return(tc.MockGetDBResult, tc.MockGetDBError).Build()

            profile, err := uc.GetVideo(ctx.Background(), videoID)
            if tc.ExpectedErr != nil {
                convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedErr.Error())
            } else {
                convey.So(err, convey.ShouldBeNil)
            }
            convey.So(profile.VideoID, convey.ShouldEqual, videoID)
            convey.So(profile.Views, convey.ShouldEqual, tc.MockViews)
            convey.So(profile.Likes, convey.ShouldEqual, tc.MockLikes)
            convey.So(profile.HotScore, convey.ShouldEqual, tc.ExpectedHotScore)
        })
    }
}
