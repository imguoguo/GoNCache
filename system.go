package main

import (
	"errors"
	"os"
)

func pathExists(path string) (bool, error) {

	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func ExistOrMkdir(path string) (bool, error) {
	pathExist, err := pathExists(path)

	if err != nil {
		return false, err
	}
	if !pathExist {
		err = os.MkdirAll(path, 0777)
	}
		return pathExist, err

}

func storeIOFile(fileName string, body []byte) (bool, error){
	if fileName == "" {
		return false, errors.New("file name is null.")
	}
	// Windows 会有 Access Denied 的问题？
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		//log.Fatalf("fatal here")
		return false, err
	}

	// 由于 io.copy 是从缓冲区拿数据，经过一次copy以后缓冲区已经没有了，网页不能从缓冲区中取得数据
	// 可以用 ioutil.ReadAll(body) 的方法来取得其中的内容并多次赋值
	// 参考链接：https://www.zhihu.com/question/56191086， https://www.jianshu.com/p/ad2e2ad7dd07
	//_, err = io.Copy(file, body)
	_, err = file.Write(body)
	if err != nil {
		return false, err
	}

	return true, nil
}

func isStringContain(items []string, item string) bool {
	// 参考链接： https://blog.csdn.net/m0_37422289/article/details/103570799
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}