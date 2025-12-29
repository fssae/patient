package web

import (
	"classroom-analysis/internal/service"
	"time"

	"gitee.com/huahua20414/pkgx/ginx"
	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	svc service.FileServiceInterface
}

func NewFileHandler(svc service.FileServiceInterface) *FileHandler {
	return &FileHandler{
		svc: svc,
	}
}

func (f *FileHandler) RegisterRoutes(server gin.IRouter) {
	server.POST("/upload/image", ginx.Wrap(f.UploadImage))
	server.POST("/upload/video", ginx.Wrap(f.UploadVideo))
}

// UploadImage 上传图片
// @Summary      上传图片
// @Description  上传图片文件
// @Tags         文件上传
// @Accept       multipart/form-data
// @Produce      json
// @Param        file     formData  file  true  "图片文件"
// @Success      200      {object}  map[string]interface{}  "上传成功"
// @Failure      400      {object}  map[string]interface{}  "上传失败"
// @Router       /upload/image [post]
func (f *FileHandler) UploadImage(c *gin.Context) (ginx.Response, error) {
	file, err := c.FormFile("file")
	if err != nil {
		return ginx.ErrorMess("上传失败", nil), err
	}
	url, fileType, err := f.svc.Upload(c, file)
	if err != nil {
		return ginx.ErrorMess("上传图片失败", nil), err
	}

	// 保存图片信息
	imageId, err := f.svc.UploadFileMessage(c, url, fileType)
	return ginx.SuccessMess("上传成功", gin.H{
		"imageId":    imageId,
		"url":        url,
		"fileName":   file.Filename,
		"size":       file.Size,
		"uploadTime": time.Now(),
	}), nil
}

// UploadVideo 上传视频
// @Summary      上传视频
// @Description  上传视频文件
// @Tags         文件上传
// @Accept       multipart/form-data
// @Produce      json
// @Param        file     formData  file  true  "视频文件"
// @Success      200      {object}  map[string]interface{}  "上传成功"
// @Failure      400      {object}  map[string]interface{}  "上传失败"
// @Router       /upload/video [post]
func (f *FileHandler) UploadVideo(c *gin.Context) (ginx.Response, error) {
	file, err := c.FormFile("file")
	if err != nil {
		return ginx.ErrorMess("上传失败", nil), err
	}
	url, fileType, err := f.svc.Upload(c, file)
	if err != nil {
		return ginx.ErrorMess("上传视频失败", nil), err
	}

	// 保存视频信息
	videoId, err := f.svc.UploadFileMessage(c, url, fileType)
	return ginx.SuccessMess("上传成功", gin.H{
		"videoId":    videoId,
		"url":        url,
		"fileName":   file.Filename,
		"size":       file.Size,
		"uploadTime": time.Now(),
	}), nil
}
