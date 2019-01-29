package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hunterhug/fafacms/core/config"
	. "github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
	"github.com/hunterhug/parrot/util"
	"path/filepath"
	"time"
)

var FileAllow = map[string][]string{
	"image": {
		"jpg", "jpeg", "png", "bmp", "gif"},
	"flash": {
		"swf", "flv"},
	"media": {
		"swf", "flv", "mp3", "wav", "wma", "wmv", "mid", "avi", "mpg", "asf", "rm", "rmvb"},
	"file": {
		"doc", "docx", "xls", "xlsx", "ppt", "htm", "html", "txt", "zip", "rar", "gz", "bz2", "pdf"},
	"other": {
		"jpg", "jpeg", "png", "bmp", "gif", "swf", "flv", "mp3",
		"wav", "wma", "wmv", "mid", "avi", "mpg", "asf", "rm", "rmvb",
		"doc", "docx", "xls", "xlsx", "ppt", "htm", "html", "txt", "zip", "rar", "gz", "bz2"}}

var FileBytes = 1 << 25 // (1<<25)/1000.0/1000.0 33.54 不能超出33M

type UploadResponse struct {
	FileName string `json:"file_name"`
	Size     int64  `json:"size"`
	Url      string `json:"url"`
}

/*
	file: 文件form名称
	type: 上传类型，分别为image、flash、media、file、other
	describe: 备注
*/
func Upload(c *gin.Context) {
	uid := c.GetInt("uid")
	resp := new(Resp)
	data := UploadResponse{}
	defer func() {
		JSONL(c, 200, nil, resp)
	}()

	fileType := c.DefaultPostForm("type", "other")
	describe := c.DefaultPostForm("describe", "")

	h, err := c.FormFile("file")
	if err != nil {
		Log.Errorf("upload err:%s", err.Error())
		resp.Error = &ErrorResp{
			ErrorID:  UploadFileError,
			ErrorMsg: ErrorMap[UploadFileError] + " " + err.Error(),
		}
		return
	}

	fileAllowArray, ok := FileAllow[fileType]
	if !ok {
		Log.Errorf("upload err: type not permit")
		resp.Error = &ErrorResp{
			ErrorID:  UploadFileError,
			ErrorMsg: ErrorMap[UploadFileError] + " type not permit",
		}
		return
	}

	fileSuffix := util.GetFileSuffix(h.Filename)

	if !util.InArray(fileAllowArray, fileSuffix) {
		Log.Errorf("upload err: file suffix: %s not permit", fileSuffix)
		resp.Error = &ErrorResp{
			ErrorID:  UploadFileError,
			ErrorMsg: ErrorMap[UploadFileError] + fmt.Sprintf(" file suffix: %s not permit", fileSuffix),
		}
		return
	}

	if h.Size > int64(FileBytes) {
		Log.Errorf("upload err: file size too big: %d", h.Size)
		resp.Error = &ErrorResp{
			ErrorID:  UploadFileError,
			ErrorMsg: ErrorMap[UploadFileError] + fmt.Sprintf(" file size too big: %d", h.Size),
		}
		return
	}

	f, err := h.Open()
	if err != nil {
		Log.Errorf("upload err:%s", err.Error())
		resp.Error = &ErrorResp{
			ErrorID:  UploadFileError,
			ErrorMsg: ErrorMap[UploadFileError] + fmt.Sprintf(":%s", err.Error()),
		}
		return
	}
	defer f.Close()
	fileMd5 := util.Md5FS(f)
	if fileMd5 == "" {
		Log.Errorf("upload err: file md5 down")
		resp.Error = &ErrorResp{
			ErrorID:  UploadFileError,
			ErrorMsg: ErrorMap[UploadFileError] + fmt.Sprintf(":file md5 down"),
		}
		return
	}

	fileDir := filepath.Join(config.FafaConfig.DefaultConfig.StoragePath, fileType, util.IS(uid))
	util.MakeDir(fileDir)

	fileName := fileMd5 + "." + fileSuffix
	fileAbName := filepath.Join(fileDir, fileName)
	if !util.HasFile(fileAbName) {
		err = util.CopyFS(f, fileAbName)
		if err != nil {
			Log.Errorf("upload err:%s", err.Error())
			resp.Error = &ErrorResp{
				ErrorID:  UploadFileError,
				ErrorMsg: ErrorMap[UploadFileError] + fmt.Sprintf(":%s", err.Error()),
			}
			return
		}

		p := new(model.Picture)
		p.Md = fileMd5
		p.Type = fileType
		p.FileName = fileName
		p.ReallyFileName = h.Filename
		p.CreateTime = time.Now().Unix()
		p.Status = 1
		p.Describe = describe
		p.Url = fmt.Sprintf("/storage/%s/%d/%s", fileType, uid, fileName)

		_, err = config.FafaRdb.InsertOne(p)
		if err != nil {
			Log.Errorf("upload err:%s", err.Error())
			resp.Error = &ErrorResp{
				ErrorID:  UploadFileError,
				ErrorMsg: ErrorMap[UploadFileError] + fmt.Sprintf(":%s", err.Error()),
			}
			return
		}
	}

	resp.Flag = true
	data.FileName = h.Filename
	data.Size = h.Size
	data.Url = fmt.Sprintf("/storage/%s/%d/%s", fileType, uid, fileName)
	resp.Data = data
	return
}
