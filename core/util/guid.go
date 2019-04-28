package util

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"strings"
)

// GetGUID 生成GUID
// todo
func GetGUID() (valueGUID string) {
	objID, _ := uuid.NewV4()
	objidStr := objID.String()
	objidStr = strings.Replace(objidStr, "-", "", -1)
	valueGUID = objidStr
	return valueGUID
}

// todo MD5变成数字来索引
// sha256  256 bit防止碰撞
func Md5(raw []byte) (string, error) {
	h := md5.New()
	num, err := h.Write(raw)
	if err != nil {
		return "", err
	}
	if num == 0 {
		return "", errors.New("num 0")
	}
	data := h.Sum([]byte("hunterhug"))
	return fmt.Sprintf("%x", data), nil
}
