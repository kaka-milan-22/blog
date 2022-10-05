package main

import (
	"blog/app/article"
	"blog/app/comment"
	"blog/app/model"
	"blog/app/user"
	"blog/config"
	"blog/routers"
	"log"
	"net/http"
	"time"
)

var sc = config.Cfg.Server

// @title Blog API
// @version 1.0
// @description  个人博客后端API
// @termsOfService https://www.google.com

// @contact.name www.google.com
// @contact.url https://www.google.com
// @contact.email kakazhang10@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

func main() {
	defer model.DB.Close()
	// 加载多个APP的路由配置
	routers.Include(
		article.Routers,
		user.Routers,
		comment.Routers,
	)
	// 初始化路由
	// model.ConnectDatabase()

	// r := routers.Init()

	// fmt.Print(config.Cfg.Server)

	// if err := r.Run(); err != nil {
	// 	fmt.Println("startup service failed, err:%v\n", err)
	// }
	server := &http.Server{
		Addr:           sc.Addr,
		Handler:        routers.Init(),
		ReadTimeout:    time.Duration(sc.ReadTimeout * int(time.Second)), // 转换成时间数据结构
		WriteTimeout:   time.Duration(sc.WriteTimeout * int(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(server.ListenAndServe())

}
