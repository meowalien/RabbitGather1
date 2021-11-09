package files

import (
	"core/src/conf"
	"core/src/lib/errs"
	"core/src/lib/hash"
	"core/src/lib/io"
	"encoding/base64"
	"fmt"
	"net/url"
	"path/filepath"
)

type File struct {
	Bin           []byte
	ExtensionName string
}

//func encodeToFileName (b)

func uploadFileName(uploadFile File) string {
	fn :=  fmt.Sprintf("%s.%s", base64.URLEncoding.EncodeToString(hash.NewSHA256(uploadFile.Bin)), uploadFile.ExtensionName)
	fmt.Println("uploadFileName: ",fn)
	return fn
}

func SaveUploadFileIfNotExist(uploadFile File) (path string, err error) {
	// 存在非公開資料夾
	path = filepath.Join(UploadSaveFilePath, uploadFileName(uploadFile))
	err = io.SaveFileIfNotExist(path,uploadFile.Bin)
	if err != nil {
		err = errs.WithLine(err)
	}
	return
}

// save and get the file url
func TakeUploadFileURL(uploadFile File) (theFileURL string, err error) {
	// 直接存在公開資料夾
	path := filepath.Join(ServeStaticFilePath, uploadFileName(uploadFile))
	err = io.SaveFileIfNotExist(path,uploadFile.Bin)
	if err != nil {
		err = errs.WithLine(err)
		return
	}

	theFileURL, err = TakeFileURLOnDisk(path)
	if err != nil {
		err = errs.WithLine(err)
		return
	}
	return
}

var ServeStaticFilePath string
var ServeFileURL *url.URL
var UploadSaveFilePath string

func init() {
	var err error
	UploadSaveFilePath, err = filepath.Abs(conf.GlobalConfig.Files.UploadSaveFilePath)
	if err != nil {
		panic("error when parse url on UploadSaveFilePath: " + err.Error())
	}
}
func init() {
	var err error
	ServeFileURL, err = url.Parse(conf.GlobalConfig.Files.ServeFileURL)
	if err != nil {
		panic("error when parse url on ServeFileURL: " + err.Error())
	}
}
func init() {
	var err error
	ServeStaticFilePath, err = filepath.Abs(conf.GlobalConfig.Files.ServeStaticFilePath)
	if err != nil {
		panic("error when get abs on ServeStaticFilePath: " + err.Error())
	}
}

func TakeFileURLOnDisk(filePathOnDisk string) (theFileURL string, err error) {
	theDir := filepath.Dir(filePathOnDisk)
	absDir, err := filepath.Abs(theDir)
	if err != nil {
		err = errs.WithLine(err)
		return
	}
	match, err := filepath.Match(filepath.Join(ServeStaticFilePath, "*"), absDir)
	if err != nil {
		err = errs.WithLine(err)
		return
	}

	var u *url.URL
	fileName := filepath.Base(filePathOnDisk)
	u, err = ServeFileURL.Parse(fileName)
	if err != nil {
		err = errs.WithLine(err)
		return
	}
	theFileURL = u.String()
	if !match {

		err = io.CopyFile(filePathOnDisk, filepath.Join(ServeStaticFilePath, fileName))
		if err != nil {
			err = errs.WithLine(err)
			return
		}
		return
	}

	return

}
