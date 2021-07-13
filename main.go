package main

import (
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
)

var (
	g errgroup.Group
)

func main() {
	// 获取配置文件
	var conf Config
	err := conf.getConf("config.yaml")
	if err != nil {
		log.Println("Error while Getting Configuration: ", err)
		return
	}
	// 检查目录，没有就创建
	// 考虑采用 freecache 充当缓存？
	//for _, cache := range conf.CacheList {
	//	pathName := "cache/" + cache.NickName
	//	_, err := ExistOrMkdir(pathName)
	//	if err != nil {
	//		log.Println("Error while mkdir: ", err)
	//		return
	//	}
	//}

	forwardServer := &http.Server {
		Addr:    conf.ListenAddress,
		Handler: forwardRouter(conf.CacheList),
	}

	managementServer := &http.Server {
		Addr: conf.ManagementAddress,
		Handler: managementRouter(),
	}

	g.Go(func() error {
		return forwardServer.ListenAndServe()
	})


	g.Go(func() error {
		return managementServer.ListenAndServe()
	})

	// 传入配置信息
	//err = router(conf)
	if err := g.Wait(); err != nil {
		log.Println("Error while Running Router: ", err)
		return
	}
}