package dao

import (
	"fmt"
	"template/global"
	"template/model/request"
)

func ImportCount(param *request.ImportCount) (int64, error) {
	if param.Table == "" {
		param.Table = "tb_addr_std"
	}
	var total int64
	sql := fmt.Sprintf("select count(1) from %s %s ", param.Table, param.Where)
	if err := global.SYS_DB.Raw(sql).Scan(&total).Error; err != nil {
		return total, err
	}
	return total, nil
}
