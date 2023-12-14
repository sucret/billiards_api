package qiniu

import (
	"billiards/pkg/config"
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"mime/multipart"
)

type Qiniu struct {
	accessKey string
	secretKey string
}

func UploadFile(file *multipart.FileHeader, key string) (err error) {
	src, _ := file.Open()
	defer func(src multipart.File) {
		_ = src.Close()
	}(src)

	conf := config.GetConfig().Qiniu

	putPolicy := storage.PutPolicy{
		Scope: conf.Bucket,
	}

	mac := qbox.NewMac(conf.AccessKey, conf.SecretKey)

	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Region = &storage.ZoneHuadong
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}

	err = formUploader.Put(context.Background(), &ret, upToken, key, src, file.Size, &putExtra)

	if err != nil {
		fmt.Println(err)
		return
	}
	return
}
