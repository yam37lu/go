package utils

import "fmt"

// @Description 字符串切片去重+去空字符串
func RemoveDuplicatAndEmpty(strSlice []string) []string {
	var duplicated = make(map[string]bool, len(strSlice))
	var ret = make([]string, 0, len(strSlice))
	for _, v := range strSlice {
		if duplicated[v] || v == "" {
			continue
		}
		duplicated[v] = true
		ret = append(ret, v)
	}
	return ret
}

func SubDistrictSql(parent string, start, end int, subs bool) string {
	sql := "select code from tb_district td where parent_code in(?)"
	if subs {
		sql = fmt.Sprintf("select code from tb_district td where parent_code in(%s)", parent)
	}
	if start == end {
		return sql
	}
	return SubDistrictSql(sql, start+1, end, true)
}
