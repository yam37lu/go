package core

import (
	"fmt"
	"strings"
	"template/utils"

	"github.com/spf13/viper"
)

const (
	Prefix = "ENC("
	Suffix = ")"
)

var keys = []string{"license.module", "gorm.password", "gorm.username", "gorm-gis.password", "gorm-gis.username", "system.regrex-key", "redis.password"}

func Encrypt(v *viper.Viper) {
	var src string
	for _, key := range keys {
		val := v.Get(key)
		if val == nil {
			fmt.Println(key)
			continue
		}
		src = val.(string)
		if strings.HasPrefix(src, Prefix) && strings.HasSuffix(src, Suffix) {
			continue
		}
		dst, _ := utils.ConfigEncrypt(src)
		v.Set(key, Prefix+dst+Suffix)
	}
}

func Decrypt(v *viper.Viper) {
	var src string
	for _, key := range keys {
		val := v.Get(key)
		if val == nil {
			fmt.Println(key)
			continue
		}
		src = val.(string)
		if strings.HasPrefix(src, Prefix) && strings.HasSuffix(src, Suffix) {
			src = strings.Replace(src, Prefix, "", 1)
			src = strings.Replace(src, Suffix, "", 1)
			dst, _ := utils.ConfigDecrypt(src)
			v.Set(key, string(dst))
		}
	}
}
