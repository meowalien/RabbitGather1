package util

import "io/ioutil"

var plainTextMap = map[string]string{}

// GetFileStoredPlainText will read the file content at the first time and then cache it.
func GetFileStoredPlainText(fileName string) string {
	if resStr, exist := plainTextMap[fileName]; exist {
		return resStr
	}
	btarr, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err.Error())
	}
	var resStr = string(btarr)
	plainTextMap[fileName] = resStr
	return resStr
}
