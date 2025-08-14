package main

import (
	"flag"
	"fmt"
	"os"
	"template/core"
	"template/global"
	"template/initialize"
	"time"

	"go.uber.org/zap"
)

var (
	GoVersion,
	Branch,
	GitCommit,
	BuildTime string
)

func main() {
	showVersion := flag.Bool("v", false, "版本信息")
	flag.Parse()
	if *showVersion == true {
		fmt.Printf("Go version：%v\r\nBranch：%v\r\nGit commit：%v\r\nBuild time：%v\r\n", GoVersion, Branch, GitCommit, BuildTime)
		os.Exit(0)
	}
	begin := time.Now()
	global.SYS_VP = core.Viper() // 初始化Viper
	global.SYS_LOG = core.Zap()  // 初始化zap日志库
	// 许可校验
	//if os.Getenv("DEBUG") != "1" {
	//	fmt.Println("license check...")
	//	module := "template"
	//	if global.SYS_CONFIG.License.Module != "" {
	//		module = global.SYS_CONFIG.License.Module
	//	}
	//	license := client.NewLicense().InitLicenseService(global.SYS_CONFIG.License.Url, module, "")
	//	license.Login()
	//
	//	if !license.Active {
	//		global.SYS_LOG.Fatal("have no license")
	//	}
	//	global.SYS_LICENSE = license
	//}
	global.SYS_DB = initialize.Gorm() // gorm连接数据库
	if global.SYS_DB != nil {
		initialize.MysqlTables(global.SYS_DB) // 初始化表
		// 程序结束前关闭数据库链接
		db, _ := global.SYS_DB.DB()
		defer db.Close()
	} else {
		global.SYS_LOG.Error("Gorm init error")
	}
	// redis初始化
	//if err := initialize.Redis(); err != nil {
	//	global.SYS_LOG.Error(err.Error())
	//}
	global.SYS_LOG.Info("main", zap.Any("dur", time.Now().Sub(begin).Milliseconds()))
	address := fmt.Sprintf(":%d", global.SYS_CONFIG.System.Addr)
	core.RunWindowsServer(address, global.SYS_CONFIG.System.ContextPath)
}
