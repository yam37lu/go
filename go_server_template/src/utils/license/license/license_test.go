package license

import (
	"fmt"
	"testing"
)

func TestLicenseService_Check(t *testing.T) {
	service := LicenseService{}
	service.EpgislicenseAPIUri = "http://172.20.20.232:10086/epgislicense"
	service.PublicKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDA4K2+VZ1mhAR4iBXUhXau+LjZ
oK6Le8TT/Q+3NANmb9LIznKgUN9TddPOx3tt+IU//Lrj8lA3Orl8UJj8i+vxEltP
7R2z7nLkMHUVmWceFe3bigwI1rXeBdZB5RGM14elQuOejJzAw+6w3TGcBFI2/l9f
JmHrmgHi2p1ZUiZOjQIDAQAB
-----END PUBLIC KEY-----`
	var moduleData ModuleData
	moduleData.ModuleId = "001"
	moduleData.ModuleName = "core_component"
	moduleData.ModuleVersion = "V1"
	service.Params = moduleData

	service.LoginErrorHandler = func() {
		fmt.Println("LoginErrorHandler...")
	}
}
