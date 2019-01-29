package controllers

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	. "github.com/hunterhug/fafacms/core/flog"
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

type Sizer interface {
	Size() int64
}

type UploadResponse struct {
	FileName string `json:"file_name"`
	Size     int    `json:"size"`
	Url      string `json:"url"`
}

/*
	file: 文件form名称
	type: 上传类型，分别为image、flash、media、file、other
*/
func Upload(c *gin.Context) {
	resp := new(Resp)

	defer func() {
		JSONL(c, 200, nil, resp)
	}()

	fileType := c.DefaultPostForm("type", "other")

	h, err := c.FormFile("file")
	if err != nil {
		Log.Errorf("upload err:%s", err.Error())
		resp.Error = &ErrorResp{
			ErrorID:  UploadFileError,
			ErrorMsg: ErrorMap[UploadFileError] + " " + err.Error(),
		}
	}

	//判断文件是否允许被添加
	//dir类型正确
	fileallowarray, ok := FileAllow[fileType]
	if !ok{

	}
	if ok {
		//得到文件后缀
		filesuffix := GetFileSuffix(h.Filename)
		//是否后缀正确
		if InArray(fileallowarray, filesuffix) {
			//获取大小
			if fileSizer, ok := f.(Sizer); ok {
				fileSize := fileSizer.Size()
				// fmt.Printf("上传%v文件的大小为: %v", fileSize, h.Filename)
				if fileSize > int64(FileBytes) {
					message = "获取上传文件错误:文件大小超出33M"
					goto END
				}
			} else {
				message = "获取上传文件错误:无法读取文件大小"
			}
			//读取二进制
			temp, err := ioutil.ReadAll(f)
			if err != nil {
				message = "读取文件错误：" + err.Error()
				goto END
			}
			filemd5 := Md5FS(bytes.NewReader(temp))
			if filemd5 == "" {
				message = "filemd5 empty"
				goto END
			}
			//创建文件夹
			subdir := Hashcode(filemd5)
			dirpath, err = MakeFileDir(filetype + "/" + subdir)
			if err != nil {
				message = "创建文件夹失败：" + err.Error()
				goto END
			} else {
				//新建文件名
				filename = filemd5 + "." + filesuffix
				// 重名没关系，因为文件相同,可以忽略
				if HasFile(dirpath + "/" + filename) {
					message = "文件重名"
					fileerror = 0
					goto END
				}
				//复制文件
				err = ioutil.WriteFile(dirpath+"/"+filename, temp, 0777)
				if err != nil {
					message = err.Error()
					goto END
				} else {
					fileerror = 0
					goto END
				}
			}
		} else {
			message = "文件后缀不被允许"
		}
	} else {
		message = "dir参数不允许"
	}

END:
	if mark == 1 {
		this.Data["json"] = &map[string]interface{}{"error": fileerror, "message": message}
	} else {
		name := dirpath + "/" + filename
		//http://lulijuan505.blog.163.com/blog/static/308369112015322102455860/
		//Base64产生的/ + =出现在url会有问题
		/*
			base64
			1、包含A-Z a-z 0-9 和加号“+”，斜杠“/” 用来作为开始的64个数字. 等号“=”用来作为后缀用途。
			2、2进制的.
			3、要比源数据多33%。
			4、常用于邮件。
			urlencode
			除了 -_. 之外的所有非字母数字字符都将被替换成百分号（%）后跟两位十六进制数，空格则编码为加号（+）
			  在神马情况下用
		*/
		token := Base64E(UrlE(name))
		//urlstring := "/public/file/getfile?token=" + token
		//fmt.Println(name)
		this.Data["json"] = &map[string]interface{}{"error": fileerror, "url": "/" + name, "token": token}
	}
	this.ServeJSON()
}
