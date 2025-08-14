package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	regrexKey = "^[\u4e00-\u9fa5a-zA-Z0-9]+$"
)

type Rules map[string][]string

type RulesMap map[string]Rules

var CustomizeMap = make(map[string]Rules)

func RegisterRule(key string, rule Rules) (err error) {
	if CustomizeMap[key] != nil {
		return errors.New(key + "已注册,无法重复注册")
	} else {
		CustomizeMap[key] = rule
		return nil
	}
}

func String() string {
	return "string"
}

func NotSpecialString() string {
	return "notSpecialString"
}

func NotSpecialStringOrEmpty() string {
	return "notSpecialStringOrEmpty"
}

func NotEmpty() string {
	return "notEmpty"
}

func NotEmptyDeep() string {
	return "notEmptyDeep"
}

func Lt(mark string) string {
	return "lt=" + mark
}

func Le(mark string) string {
	return "le=" + mark
}

func Eq(mark string) string {
	return "eq=" + mark
}

func Ne(mark string) string {
	return "ne=" + mark
}

func Ge(mark string) string {
	return "ge=" + mark
}

func Gt(mark string) string {
	return "gt=" + mark
}

func In(mark string) string {
	return "in=" + mark
}

func verifyMap(obj reflect.Value, roleMap Rules, compareMap map[string]bool) (err error) {
	if obj.Kind() != reflect.Map {
		return errors.New("无效对象")
	}
	for _, key := range obj.MapKeys() {
		val := reflect.ValueOf(obj.MapIndex(key).Interface())
		k := strings.Title(key.String())
		if len(roleMap[k]) > 0 {
			for _, v := range roleMap[k] {
				switch {
				case v == "notEmptyDeep":
					if isBlank(val, true) {
						return errors.New(k + "值不能为空")
					}
				case v == "notEmpty":
					if isBlank(val, false) {
						return errors.New(k + "值不能为空")
					}
				case v == "string":
					if !isString(val) {
						return errors.New(k + "值不是string或为空")
					}
				case v == "notSpecialString":
					{
						if !isString(val) {
							return errors.New(k + "值不是string或为空")
						}
					}
				case v == "notSpecialStringOrEmpty":
					{
						if !isStringOrEmpty(val) {
							return errors.New(k + "值不是string或为空")
						}
						if val.Len() > 0 {
						}
					}
				case compareMap[strings.Split(v, "=")[0]]:
					if !compareVerify(val, v) {
						return errors.New(k + "长度或值不在合法范围," + v)
					}
				case strings.Split(v, "=")[0] == "in":
					if !isIn(val, strings.Split(v, "=")[1]) {
						return errors.New(k + "值不在合理范围内")
					}
				}
			}
		}
	}
	return nil
}

// @Description 校验请求参数结构体必须有字段接收
func CheckoutQuery(ctx *gin.Context, ptr any, tag string) error {
	// tag: "form"/"json"
	// ptr: 结构体或结构体指针
	// 校验： 只支持结构体或结构体指针
	ptrVal := reflect.ValueOf(ptr)
	if ptrVal.Kind() == reflect.Ptr {
		ptr_ := ptrVal.Elem().Interface()
		ptrVal = reflect.ValueOf(ptr_)
	}
	if ptrVal.Kind() != reflect.Struct {
		return ServiceInternalError.New("checkout param in struct error: please give a struct or a struct ptr")
	}
	// 获取tag对应的字段
	structtag := make(map[string]string, 0)

	tValue := ptrVal.Type()

	for i := 0; i < ptrVal.NumField(); i++ {
		sf := tValue.Field(i)
		if sf.Tag.Get(tag) == "-" { // just ignoring this field
			continue
		}
		tagValue := sf.Tag.Get(tag)
		if tagValue == "" { // default value is FieldName
			tagValue = sf.Name
		}
		if tagValue == "" { // when field is "emptyField" variable
			return OtherUnknownError.New("get tag error")
		}
		structtag[tagValue] = "1"
	}
	if tag == "form" {
		if ctx.Request.Method != http.MethodGet {
			return ServiceInternalError.New("form check only support GET")
		}
		parammap := ctx.Request.URL.Query()
		for k := range parammap {
			if _, ok := structtag[k]; !ok {
				return InvaildParams.Newf("invaild param[%s] is find", k)
			}
		}
		return nil
	}
	if tag == "json" {
		if ctx.Request.Method != http.MethodPost {
			return ServiceInternalError.New("json check only support POST")
		}
		parammap := make(map[string]interface{}, 0)
		body, _ := ioutil.ReadAll(ctx.Request.Body)
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body)) //重新写入
		defer ctx.Request.Body.Close()
		if err := json.Unmarshal(body, &parammap); err != nil {
			return ServiceInternalError.New("unmarshal to json error")
		}
		for k := range parammap {
			if _, ok := structtag[k]; !ok {
				return InvaildParams.Newf("invaild param[%s] is find", k)
			}
		}
		return nil
	}
	return ServiceInternalError.New("unsupported query type")
}

func Verify(st interface{}, roleMap Rules) (err error) {
	compareMap := map[string]bool{
		"lt": true,
		"le": true,
		"eq": true,
		"ne": true,
		"ge": true,
		"gt": true,
	}

	typ := reflect.TypeOf(st)
	val := reflect.ValueOf(st) // 获取reflect.Type类型

	kd := val.Kind() // 获取到st对应的类别
	if kd == reflect.Map {
		return verifyMap(val, roleMap, compareMap)
	}
	if kd != reflect.Struct {
		return errors.New("expect struct")
	}
	num := val.NumField()
	// 遍历结构体的所有字段
	for i := 0; i < num; i++ {
		tagVal := typ.Field(i)
		fieldVal := val.Field(i)
		if tagVal.Anonymous {
			fieldInterface := fieldVal.Interface()
			fieldVal = reflect.ValueOf(fieldInterface)
			if fieldVal.Kind() == reflect.Pointer {
				fieldVal = fieldVal.Elem()
			}
			if er := Verify(fieldVal.Interface(), roleMap); er != nil {
				return er
			}
			continue
		}
		if len(roleMap[tagVal.Name]) > 0 {
			for _, v := range roleMap[tagVal.Name] {
				switch {
				case v == "notEmptyDeep":
					if isBlank(fieldVal, true) {
						return errors.New(tagVal.Tag.Get("json") + "值不能为空")
					}
				case v == "notEmpty":
					if isBlank(fieldVal, false) {
						return errors.New(tagVal.Tag.Get("json") + "值不能为空")
					}
				case v == "notSpecialString":
					{
						if !isString(fieldVal) {
							return errors.New(tagVal.Tag.Get("json") + "值不是string或为空")
						}
					}
				case v == "notSpecialStringOrEmpty":
					{
						if !isStringOrEmpty(fieldVal) {
							return errors.New(tagVal.Tag.Get("json") + "值不是string或为空")
						}
						if fieldVal.Len() > 0 {
						}
					}
				case compareMap[strings.Split(v, "=")[0]]:
					if !compareVerify(fieldVal, v) {
						return errors.New(tagVal.Tag.Get("json") + "长度或值不在合法范围," + v)
					}
				case strings.Split(v, "=")[0] == "in":
					if !isIn(fieldVal, strings.Split(v, "=")[1]) {
						return errors.New(tagVal.Tag.Get("json") + "值不在合理范围内")
					}
				}
			}
		}
	}
	return nil
}

func compare2(value reflect.Value, VerifyStr string) bool {
	switch value.Kind() {
	case reflect.String:
		return value.String() == VerifyStr
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		VInt, VErr := strconv.Atoi(VerifyStr)
		if VErr != nil {
			return false
		}
		return value.Uint() == uint64(VInt)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		VInt, VErr := strconv.ParseInt(VerifyStr, 10, 64)
		if VErr != nil {
			return false
		}
		return value.Int() == VInt
	case reflect.Float32, reflect.Float64:
		VFloat, VErr := strconv.ParseFloat(VerifyStr, 64)
		if VErr != nil {
			return false
		}
		return value.Float() == VFloat
	default:
		return false
	}
}

func isIn(value reflect.Value, str string) bool {
	vals := strings.Split(str, ",")
	for _, v := range vals {
		if compare2(value, v) {
			return true
		}
	}
	return false
}

func compareVerify(value reflect.Value, VerifyStr string) bool {
	switch value.Kind() {
	case reflect.String:
		return compare(value.String(), VerifyStr)
	case reflect.Slice, reflect.Array:
		return compare(value.Len(), VerifyStr)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return compare(value.Uint(), VerifyStr)
	case reflect.Float32, reflect.Float64:
		return compare(value.Float(), VerifyStr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return compare(value.Int(), VerifyStr)
	default:
		return false
	}
}

func isString(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() != 0
	default:
		return false
	}
}

func isStringOrEmpty(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return true
	default:
		return false
	}
}

func isBlank(value reflect.Value, deep bool) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Slice, reflect.Array:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	if deep {
		return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
	} else {
		return false
	}
}

func compare(value interface{}, VerifyStr string) bool {
	VerifyStrArr := strings.Split(VerifyStr, "=")
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		VInt, VErr := strconv.ParseInt(VerifyStrArr[1], 10, 64)
		if VErr != nil {
			return false
		}
		switch {
		case VerifyStrArr[0] == "lt":
			return val.Int() < VInt
		case VerifyStrArr[0] == "le":
			return val.Int() <= VInt
		case VerifyStrArr[0] == "eq":
			return val.Int() == VInt
		case VerifyStrArr[0] == "ne":
			return val.Int() != VInt
		case VerifyStrArr[0] == "ge":
			return val.Int() >= VInt
		case VerifyStrArr[0] == "gt":
			return val.Int() > VInt
		default:
			return false
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		VInt, VErr := strconv.Atoi(VerifyStrArr[1])
		if VErr != nil {
			return false
		}
		switch {
		case VerifyStrArr[0] == "lt":
			return val.Uint() < uint64(VInt)
		case VerifyStrArr[0] == "le":
			return val.Uint() <= uint64(VInt)
		case VerifyStrArr[0] == "eq":
			return val.Uint() == uint64(VInt)
		case VerifyStrArr[0] == "ne":
			return val.Uint() != uint64(VInt)
		case VerifyStrArr[0] == "ge":
			return val.Uint() >= uint64(VInt)
		case VerifyStrArr[0] == "gt":
			return val.Uint() > uint64(VInt)
		default:
			return false
		}
	case reflect.Float32, reflect.Float64:
		VFloat, VErr := strconv.ParseFloat(VerifyStrArr[1], 64)
		if VErr != nil {
			return false
		}
		switch {
		case VerifyStrArr[0] == "lt":
			return val.Float() < VFloat
		case VerifyStrArr[0] == "le":
			return val.Float() <= VFloat
		case VerifyStrArr[0] == "eq":
			return val.Float() == VFloat
		case VerifyStrArr[0] == "ne":
			return val.Float() != VFloat
		case VerifyStrArr[0] == "ge":
			return val.Float() >= VFloat
		case VerifyStrArr[0] == "gt":
			return val.Float() > VFloat
		default:
			return false
		}
	case reflect.String:
		switch {
		case VerifyStrArr[0] == "lt":
			return val.String() < VerifyStrArr[1]
		case VerifyStrArr[0] == "le":
			return val.String() <= VerifyStrArr[1]
		case VerifyStrArr[0] == "eq":
			return val.String() == VerifyStrArr[1]
		case VerifyStrArr[0] == "ne":
			return val.String() != VerifyStrArr[1]
		case VerifyStrArr[0] == "ge":
			return val.String() >= VerifyStrArr[1]
		case VerifyStrArr[0] == "gt":
			return val.String() > VerifyStrArr[1]
		default:
			return false
		}
	default:
		return false
	}
}
