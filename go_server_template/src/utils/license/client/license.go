/*
author：Mr.C
cretetime:2018.9.6
description:初始化许可服务对象，心跳状态Active
*/
package client

import (
	"os"
	"template/utils/license/license"
)

const publicKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDA4K2+VZ1mhAR4iBXUhXau+LjZ
oK6Le8TT/Q+3NANmb9LIznKgUN9TddPOx3tt+IU//Lrj8lA3Orl8UJj8i+vxEltP
7R2z7nLkMHUVmWceFe3bigwI1rXeBdZB5RGM14elQuOejJzAw+6w3TGcBFI2/l9f
JmHrmgHi2p1ZUiZOjQIDAQAB
-----END PUBLIC KEY-----`

var Active = true

type License struct {
	Active         bool
	LicenseService *license.LicenseService
}

func NewLicense() *License {
	return &License{
		Active:         true,
		LicenseService: &license.LicenseService{},
	}
}

func (s *License) InitLicenseService(uri string, moduleName string, moduleVersion string) *License {
	s.LicenseService.EpgislicenseAPIUri = uri
	s.LicenseService.PublicKey = publicKey
	var moduleData license.ModuleData
	moduleData.ModuleName = moduleName
	moduleData.ModuleVersion = moduleVersion
	s.LicenseService.Params = moduleData

	s.LicenseService.LoginErrorHandler = func() {
		os.Stdout.WriteString("登录license服务返回错误，系统退出 ...\n")
		os.Exit(-1)
	}
	s.LicenseService.HealthHandler = func(health bool) {
		s.Active = health
		Active = s.Active
		if !s.Active {
			os.Stdout.WriteString("license心跳检测返回错误，系统停止服务...\n")
		} else {
			os.Stdout.WriteString("license心跳检测正常，系统开始服务...\n")
		}
	}
	return s
}

// 启动鉴权服务是否被许可认证
func (s *License) Login() {
	// commons.Info("启动许可认证监听 ...")
	err := s.LicenseService.Login()
	if err != nil {
		s.Active = false
		Active = s.Active
	}
}

func (s *License) Logout() {
	s.LicenseService.Logout()
}
