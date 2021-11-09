package io

import (
	"core/src/lib/errs"
	"errors"
	"fmt"
	"io"
	"os"
)

func SaveFileIfNotExist(path string, bin []byte) (err error)  {
	if _, err = os.Stat(path); err == nil {
		// 已經存在就不用存了
		return
	} else if errors.Is(err, os.ErrNotExist) {
		err = SaveAsNewFile(path, bin)
		if err != nil {
			err = errs.WithLine(err)
			return
		}
		return
	} else {
		err = errs.WithLine(err)
		return
	}
}


// 將檔案存到硬碟
func SaveAsNewFile(filePath string, bin []byte) (err error) {
	fo, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error when os.Create: %w", err)
	}
	defer func(fo *os.File) {
		e := fo.Close()
		if e != nil {
			fmt.Printf("error when close file: %s", e.Error())
		}
	}(fo)

	_, err = fo.Write(bin)
	if err != nil {
		return fmt.Errorf("error Write file: %w", err)
	}
	return nil
}

func CopyFile(src, dst string) (err error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		err = errs.WithLine(err)
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		err = errs.WithLine(fmt.Errorf("%s is not a regular file", src))
		return err
	}

	source, err := os.Open(src)
	if err != nil {
		err = errs.WithLine(err)
		return
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		err = errs.WithLine(err)

	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return
}
