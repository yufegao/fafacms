package util

import (
	"strings"

	"github.com/satori/go.uuid"
)

// GetGUID 生成GUID
func GetGUID() (valueGUID string) {
	objID, _ := uuid.NewV4()
	objidStr := objID.String()
	objidStr = strings.Replace(objidStr, "-", "", -1)
	valueGUID = objidStr
	return valueGUID
}
