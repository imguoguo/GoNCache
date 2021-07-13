package main

import (
	"crypto/md5"
	"encoding/hex"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)


// cache 文件格式
type CacheInfo struct {
	OriginUri string `yaml:"origin_uri"`
	ExpiredAt int64 `yaml:"expired_at"`
	Type string `yaml:"type"`
}

func saveCache(model CacheListModel, request *http.Request, response *http.Response, responseBody []byte) (bool, error) {

	// 处理缓存
	//  判断是否在黑名单内
	// TODO 是否为黑名单文件类型
	// TODO 写入缓存

	contentType := response.Header.Get("Content-Type")
	if isStringContain(model.TypeBlackList, contentType) {
		// 黑名单文件类型不缓存
		return false, nil
	}

	if response.StatusCode != http.StatusOK {
		// 非正常状态码不缓存
		return false, nil
	}

	cacheDir := "cache/" + model.NickName + "/" +
		cacheName(model.SeparateMobile, request)
	_, err := ExistOrMkdir(cacheDir)
	if err != nil {
		return false, err
	}
	_, err = storeIOFile(cacheDir + "/file.data", responseBody)
	if err != nil {
		return false, err
	}

	// 写入cache信息
	cacheInfo := &CacheInfo{
		OriginUri: request.RequestURI,
		ExpiredAt: time.Now().Unix() + int64(model.ExpiredTime),
		Type: response.Header.Get("Content-Type"),
	}

	data, err := yaml.Marshal(cacheInfo)
	if err != nil {
		return false, err
	}

	err = ioutil.WriteFile(cacheDir + "/info.yaml", data, 0777)
	return true, nil
}


func checkAndGetType(cacheDir string) (bool, string) {
	// TODO 检查缓存是否存在
	// TODO 检查缓存是否过期

	if exist, _ := pathExists(cacheDir); exist != true {
		return false, ""
	}

	yamlFile, err := ioutil.ReadFile(cacheDir + "/info.yaml")
	if err != nil {
		return false, ""
	}
	var cache CacheInfo
	err = yaml.Unmarshal(yamlFile, &cache)
	if err != nil {
		return false, ""
	}
	if cache.ExpiredAt < time.Now().Unix() {
		return false, ""
	}
	return true, cache.Type
}

func cacheName(separateMobile bool, r *http.Request) string {
	// md5(路径+参数+是否手机UA)
	originInfo := r.RequestURI

	// 当手机电脑分开缓存的情况
	if separateMobile {
		userAgent := r.Header.Get("User-Agent")

		if strings.Contains(userAgent, "Mobile") {
			originInfo += "1"
		} else {
			originInfo += "0"
		}
	}

	// 输出 md5 值
	hash := md5.Sum([]byte(originInfo))
	return hex.EncodeToString(hash[:])

}

func checkModel(requestHost string, cacheList []CacheListModel) CacheListModel{
	var cacheModel CacheListModel
	for _, cache := range cacheList {
		for _, host := range cache.Host{
			if host == requestHost {
				cacheModel = cache
				return cacheModel
			}
		}
	}
	return cacheModel
}


