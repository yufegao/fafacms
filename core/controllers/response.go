package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/hunterhug/fafacms/core/config"
	"github.com/hunterhug/fafacms/core/model"
	"github.com/hunterhug/fafacms/core/util"
	"io/ioutil"
	"runtime"
	"strings"
	"time"
	. "github.com/hunterhug/fafacms/core/flog"
)

func ParseJSON(c *gin.Context, req interface{}) *ErrorResp {
	pc, _, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	requestBody, _ := ioutil.ReadAll(c.Request.Body)

	ip := c.ClientIP()

	Log.Debugf("%s ParseJSON [%v,line:%v]:%s", ip, f.Name(), line, string(requestBody))
	if err := json.Unmarshal(requestBody, req); err != nil {
		Log.Errorf("%s ParseJSON Unmarshal err:%s", ip, err.Error())
		c.Set("notlog", true)
		return &ErrorResp{
			ErrorID:  ParseJsonError,
			ErrorMsg: ErrorMap[ParseJsonError],
		}
	}
	return nil
}

func JSONL(c *gin.Context, code int, req interface{}, obj *Resp) {
	if c.GetBool("notlog") {
		c.Render(code, render.JSON{Data: obj})
		return
	}
	record := new(model.Log)
	record.Ip = c.ClientIP()
	record.Url = c.Request.URL.Path
	record.LogTime = time.Now().Unix()
	record.Ua = c.Request.UserAgent()
	record.UserId = c.GetInt("uid")
	flag := obj.Flag
	if !flag && obj.Error != nil {
		errstr := obj.Error.Error()
		errstrr := strings.Split(errstr, "|")
		if len(errstrr) >= 2 {
			record.ErrorId = errstrr[0]
			record.ErrorMessage = strings.Join(errstrr[1:], "|")
		}
	}
	record.Flag = flag

	if req != nil {
		in, _ := json.Marshal(req)
		if len(in) > 0 {
			record.In = string(in)
		}
	}

	if obj != nil {
		out, _ := json.Marshal(obj)
		if len(out) > 0 {
			record.Out = string(out)
		}
	}
	cid := util.GetGUID()
	record.Cid = cid

	_, err := config.FafaRdb.InsertOne(record)
	if err != nil {
		Log.Errorf("insert log record:", err.Error())
	}

	obj.Cid = cid
	c.Render(code, render.JSON{Data: obj})
}

func JSON(c *gin.Context, code int, obj *Resp) {
	c.Render(code, render.JSON{Data: obj})
}
