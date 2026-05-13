package main

import (
	"github.com/lexiaowenn/md2wechat-new/internal/config"
	"github.com/lexiaowenn/md2wechat-new/internal/image"
	"github.com/lexiaowenn/md2wechat-new/internal/wechat"
)

func newRuntimeImageProcessor() *image.Processor {
	return newRuntimeImageProcessorWithConfig(cfg)
}

func newRuntimeImageProcessorWithConfig(runtimeCfg *config.Config) *image.Processor {
	if runtimeCfg == nil {
		runtimeCfg = cfg
	}
	svc := wechat.NewService(runtimeCfg, log)
	return image.NewProcessor(
		runtimeCfg,
		log,
		image.WithDownloadFunc(wechat.DownloadFile),
		image.WithUploadFunc(func(filePath string) (*image.UploadResult, error) {
			result, err := svc.UploadMaterialWithRetry(filePath, 3)
			if err != nil {
				return nil, err
			}
			return &image.UploadResult{
				MediaID:   result.MediaID,
				WechatURL: result.WechatURL,
				Width:     result.Width,
				Height:    result.Height,
			}, nil
		}),
	)
}
