package main

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

func forwardRouter(cacheList []CacheListModel) http.Handler {
	engine := gin.New()
	engine.Use(gin.Recovery())
	// 似乎没有办法指定特定host？应该可以用分端口的方式
	// https://chenyitian.gitbooks.io/gin-web-framework/content/docs/37.html
	engine.Any("/*pathInfo", func(context *gin.Context) {

		// 处理域名
		// 检查是否在列表中
		cacheModel := checkModel(context.Request.Host, cacheList)

		// TODO 检查路径是否在黑名单
		// 检查是否缓存
		if cacheModel.NickName != ""  {
			cacheDir := "cache/" + cacheModel.NickName + "/" +
				cacheName(cacheModel.SeparateMobile, context.Request)
			if cacheUsable, cacheType := checkAndGetType(cacheDir); cacheUsable {
				// TODO 有缓存直接返回

				fileData, err := ioutil.ReadFile(cacheDir + "/file.data")
				if err != nil {
					log.Println("获取缓存数据时出现了问题")
					context.JSON(http.StatusBadGateway, gin.H{
						"message": "获取缓存数据" + cacheDir + "时出现了问题",
					})
					return
				}
				context.Writer.Header().Add("Content-Type", cacheType)
				context.Writer.Write(fileData)
				return
			} else {
				// 无缓存请求并存储
				// 参考自：https://blog.csdn.net/mengxinghuiku/article/details/65448600
				// 有就懒得造轮子 =w=
				transport :=  http.DefaultTransport
				outReq := new(http.Request)
				*outReq = *context.Request // 浅层复制

				if clientIP, _, err := net.SplitHostPort(context.Request.RemoteAddr); err == nil {
					if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
						clientIP = strings.Join(prior, ", ") + ", " + clientIP
					}
					outReq.Header.Set("X-Forwarded-For", clientIP)
				}


				outReq.URL.Host = cacheModel.Source
				outReq.URL.Scheme = cacheModel.Scheme
				outReq.Host = cacheModel.Source
				response, err := transport.RoundTrip(outReq)
				// 原因不明的EOF
				// 解决方法：outReq.Host也要设置
				if err != nil {
					log.Println(err)
					context.JSON(http.StatusBadGateway, gin.H {
						"error": err,
					})
				}
				// 设置Header
				for key, value := range response.Header {
					for _, v := range value {
						context.Writer.Header().Add(key, v)
					}
				}
				// 写到网页里
				context.Writer.WriteHeader(response.StatusCode)
				responseBody, err := ioutil.ReadAll(response.Body)
				context.Writer.Write(responseBody)
				//io.Copy(c.Writer, res.Body)
				defer response.Body.Close()



				_, err = saveCache(cacheModel, context.Request, response, responseBody)
				if err != nil {
					log.Println("Error while trying to save cache: ", err)
				}

			}
		} else {
			context.JSON(404, gin.H{
				"error": "Cannot Find Host",
			})
		}
	})
	return engine
}

func managementRouter() http.Handler {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.GET("/api/:method", func(context *gin.Context) {
		method := context.Param("method")
		context.JSON(http.StatusOK, gin.H{
			"message": "Hello World! You visit " + method,
		})
	})
	return engine
}