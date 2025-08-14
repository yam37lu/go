package license

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"template/utils/license/codec"
)

type LicenseService struct {
	EpgislicenseAPIUri string
	PublicKey          string
	Params             ModuleData
	LoginErrorHandler  func()
	HealthHandler      func(health bool)
}

type ModuleData struct {
	ModuleId      string `json:"moduleId"`      //模块ID
	ModuleName    string `json:"moduleName"`    //模块名称
	ModuleVersion string `json:"moduleVersion"` //模块版本
}

type EncryptData struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Result  LicenseResult `json:"result"`
}

type ActiveResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type LicenseResult struct {
	Id          string      `json:"id"`
	Timestamp   int64       `json:"timestamp"`
	Interval    int64       `json:"interval"`
	LicenseInfo LicenseInfo `json:"licenseInfo"`
}

type LicenseInfo struct {
	ExpiryDate    string `json:"expiryDate"`
	EffectiveDate string `json:"effectiveDate"`
	Industry      string `json:"industry"`
	Company       string `json:"company"`
	Project       string `json:"project"`
}

type LicenseServiceAPI interface {
	Login() error
	Logout() error
}

const (
	//licenseService = "http://172.20.20.232:10086/epgislicense"

	API_SERVERTIME = "/servertime"
	API_LOGIN      = "/login"
	API_ACTIVE     = "/active"
	API_LOGOUT     = "/logout"

	//	pUBLICK_KEY = `-----BEGIN PUBLIC KEY-----
	//
	// MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDA4K2+VZ1mhAR4iBXUhXau+LjZ
	// oK6Le8TT/Q+3NANmb9LIznKgUN9TddPOx3tt+IU//Lrj8lA3Orl8UJj8i+vxEltP
	// 7R2z7nLkMHUVmWceFe3bigwI1rXeBdZB5RGM14elQuOejJzAw+6w3TGcBFI2/l9f
	// JmHrmgHi2p1ZUiZOjQIDAQAB
	// -----END PUBLIC KEY-----`
)

var _id string = ""

func (licenseService *LicenseService) Login() error {
	serverTime, err := getServerTime(licenseService.EpgislicenseAPIUri + API_SERVERTIME)
	if err != nil {
		return errors.New(fmt.Sprintf("获取服务器时间戳异常：%v\n", err))
	}
	now := time.Now().UnixNano() / int64(time.Millisecond)
	diffTime := now - serverTime

	licenseData, err := doLogin(diffTime, licenseService)
	if err != nil {
		return errors.New(fmt.Sprintf("获取license服务返回值属性异常:%v\n", err))
	}
	if licenseData.Code != 0 {
		fmt.Printf("license服务验证不通过，错误码：%v,错误信息:%v\n", licenseData.Code, licenseData.Message)
		licenseService.LoginErrorHandler()
		//fmt.Println("login error hook called")
		return errors.New(fmt.Sprintf("license服务验证不通过，错误码：%v,错误信息:%v\n", licenseData.Code, licenseData.Message))
	} else {
		// 启动心跳检测
		//fmt.Printf("启动心跳检测,id:%s,interval:%d\n", licenseData.Result.Id, licenseData.Result.Interval)
		_id = licenseData.Result.Id
		go healthCheck(licenseData.Result.Id, diffTime, licenseData.Result.Interval, licenseService)
	}
	return nil
}

func doLogin(diffTime int64, licenseService *LicenseService) (EncryptData, error) {
	data := EncryptData{}
	params, err := buildLoginParams(diffTime, licenseService.PublicKey, licenseService.Params)
	if err != nil {
		return data, errors.New(fmt.Sprintf("拼凑license服务请求参数异常：%v\n", err))
	}

	body, err := loginLicense(licenseService.EpgislicenseAPIUri+API_LOGIN, params)
	if err != nil {
		return data, errors.New(fmt.Sprintf("登录license服务异常：%v\n", err))
	}

	m := make(map[string]string)
	err = json.Unmarshal(body, &m)
	if err != nil {
		return data, errors.New(fmt.Sprintf("解析license服务返回值异常:%v\n", err))
	}

	loginResult, _ := m["encryptData"]
	decodeLoginResult, err := codec.RSA.String(loginResult, codec.MODE_PUBKEY_DECRYPT)
	if err != nil {
		return data, errors.New(fmt.Sprintf("rsa解密license服务返回值异常:%v\n", err))
	}
	//fmt.Printf("encryptData decode result:%v\n", decodeLoginResult)

	var licenseData EncryptData
	err = json.Unmarshal([]byte(decodeLoginResult), &licenseData)
	if err != nil {
		return data, errors.New(fmt.Sprintf("解析license服务返回值属性异常:%v\n", err))
	}
	return licenseData, nil
}

func healthCheck(id string, diffTime int64, interval int64, licenseService *LicenseService) {
	ticker := time.NewTicker(time.Millisecond * time.Duration(interval))
	count := 0
	go func() {
		for _ = range ticker.C {
			activeResponse, err := activeTask(id, diffTime, licenseService)
			if err != nil {
				fmt.Printf("health check error : count=%v, %v\n", count, err)
				count = count + 1
				if count >= 3 {
					count = 0
					healthHandler := licenseService.HealthHandler

					licenseData, err := doLogin(diffTime, licenseService)
					if err != nil {
						activeResponse.Message = fmt.Sprintf("%v", err)
						//*licenseService.ActiveResult <- activeResponse
						fmt.Printf("healthCheck relogin error:%v\n", err)
						healthHandler(false)
						continue
					}
					if licenseData.Code != 0 {
						activeResponse.Code = licenseData.Code
						activeResponse.Message = licenseData.Message
						//*licenseService.ActiveResult <- activeResponse
						fmt.Printf("healthCheck relogin error:%+v\n", licenseData)
						healthHandler(false)
					} else {
						ticker.Stop()
						//fmt.Println("healthCheck login success, ticker stop")
						healthHandler(true)
						// 启动心跳检测
						go healthCheck(licenseData.Result.Id, diffTime, licenseData.Result.Interval, licenseService)
						return
					}
				}
			}
			//fmt.Printf("id:%v tick end %v\n", id, time.Now())
		}
	}()
}

func activeTask(id string, diffTime int64, licenseService *LicenseService) (ActiveResponse, error) {
	activeResponse := ActiveResponse{-1, ""}

	params, err := buildActiveParams(id, diffTime, licenseService.PublicKey)
	if err != nil {
		activeResponse.Message = fmt.Sprintf("拼凑license心跳服务请求参数异常：%v\n", err)
		return activeResponse, errors.New(fmt.Sprintf("拼凑license心跳服务请求参数异常：%v\n", err))
	}
	body, err := licenseActive(licenseService.EpgislicenseAPIUri+API_ACTIVE, params)
	if err != nil {
		activeResponse.Message = fmt.Sprintf("访问license心跳服务异常：%v\n", err)
		return activeResponse, errors.New(fmt.Sprintf("访问license心跳服务异常：%v\n", err))
	}

	m := make(map[string]string)
	err = json.Unmarshal(body, &m)
	if err != nil {
		activeResponse.Message = fmt.Sprintf("解析license心跳服务返回值异常:%v\n", err)
		return activeResponse, errors.New(fmt.Sprintf("解析license心跳服务返回值异常:%v\n", err))
	}

	activeResult, _ := m["encryptData"]
	result, err := codec.RSA.String(activeResult, codec.MODE_PUBKEY_DECRYPT)
	if err != nil {
		activeResponse.Message = fmt.Sprintf("rsa解密license心跳服务返回值异常:%v\n", err)
		return activeResponse, errors.New(fmt.Sprintf("rsa解密license心跳服务返回值异常:%v\n", err))
	}
	//fmt.Printf("encryptData decode result:%v\n", result)

	err = json.Unmarshal([]byte(result), &activeResponse)
	if err != nil {
		activeResponse.Message = fmt.Sprintf("解析license心跳服务返回值属性异常:%v\n", err)
		return activeResponse, errors.New(fmt.Sprintf("解析license心跳服务返回值属性异常:%v\n", err))
	}
	if activeResponse.Code != 0 {
		//fmt.Printf("license心跳服务验证不通过，错误码：%v,错误信息:%v\n", activeResponse.Code, activeResponse.Message)
		activeResponse.Message = fmt.Sprintf("license心跳服务验证不通过，错误码：%v,错误信息:%v\n", activeResponse.Code, activeResponse.Message)
		return activeResponse, errors.New(fmt.Sprintf("license心跳服务验证不通过，错误码：%v,错误信息:%v\n", activeResponse.Code, activeResponse.Message))
	}
	//fmt.Printf("active task data:%+v\n", activeResponse)
	return activeResponse, nil
}

func buildActiveParams(id string, diffTime int64, publicKey string) (map[string][]string, error) {
	params := url.Values{}
	encryptData := fmt.Sprintf("{\"id\":\"%s\"}", id)
	//fmt.Printf("jsonData buildActiveParams:%v\n", string(encryptData))

	pubErr := codec.RSA.Init(publicKey)
	if pubErr != nil {
		//fmt.Printf("初始化rsa加密错误：%v\n", pubErr)
		return nil, errors.New(fmt.Sprintf("初始化rsa加密错误：%v\n", pubErr))
	}
	data, err := codec.RSA.String(encryptData, codec.MODE_PUBKEY_ENCRYPT)
	if err != nil {
		//fmt.Printf("rsa加密请求数据错误：%v\n", err)
		return nil, errors.New(fmt.Sprintf("rsa加密请求数据错误：%v\n", err))
	}
	params.Add("encryptData", data) //rsa加密后已base64编码

	hash := md5.New()
	hash.Write([]byte(encryptData))
	cipherText := hash.Sum(nil)
	hexText := make([]byte, 32)
	hex.Encode(hexText, cipherText)
	params.Add("signData", string(hexText))

	now := time.Now().UnixNano() / int64(time.Millisecond)
	timestamp := now - diffTime
	//fmt.Println("timestamp:", timestamp)
	timeData, err := codec.RSA.String(string([]byte(strconv.FormatInt(timestamp, 10))), codec.MODE_PUBKEY_ENCRYPT)
	if err != nil {
		//fmt.Printf("rsa加密时间戳错误：%v\n", err)
		return nil, errors.New(fmt.Sprintf("rsa加密时间戳错误：%v\n", err))
	}
	params.Add("time", string(timeData))
	//fmt.Printf("buildActiveParams:%v\n", params)
	return params, nil
}

func licenseActive(activeAPIUrl string, params map[string][]string) ([]byte, error) {
	licenseActiveResp, err := http.PostForm(activeAPIUrl, params)
	if err != nil {
		//fmt.Printf("调用license心跳服务异常：%v\n", err)
		return nil, errors.New(fmt.Sprintf("调用license心跳服务异常：%v\n", err))
	}
	defer func() ([]byte, error) {
		err := licenseActiveResp.Body.Close()
		if err == nil {
		}
		//logs.Error(err)
		return nil, err
	}()
	body, err := ioutil.ReadAll(licenseActiveResp.Body)
	if err != nil {
		//fmt.Printf("读取license心跳服务返回值异常：%v\n", err)
		return nil, errors.New(fmt.Sprintf("读取license心跳服务返回值异常：%v\n", err))
	}
	//fmt.Printf("license active response:%v\n", string(body))
	return body, nil
}

func getServerTime(serverTimeURL string) (int64, error) {
	serverTimeResp, err := http.PostForm(serverTimeURL, nil)
	if err != nil {
		//fmt.Printf("获取服务器时间戳异常：%v\n", err)
		return 0, errors.New(fmt.Sprintf("获取服务器时间戳异常：%v\n", err))
	}
	defer func() (int64, error) {
		err := serverTimeResp.Body.Close()
		if err == nil {
		}
		//logs.Error(err)
		return 0, err
	}()
	serverTimeBody, err := ioutil.ReadAll(serverTimeResp.Body)
	if err != nil {
		//fmt.Printf("读取服务器时间返回值异常：%v\n", err)
		return 0, errors.New(fmt.Sprintf("读取服务器时间返回值异常：%v\n", err))
	}
	serverTime, err := strconv.ParseInt(string(serverTimeBody), 10, 64)
	if err != nil {
		//fmt.Printf("转换服务器时间返回值异常：%v\n", err)
		return 0, errors.New(fmt.Sprintf("转换服务器时间返回值异常：%v\n", err))
	}
	//fmt.Println("server time result:", serverTime)
	return serverTime, nil
}

func buildLoginParams(diffTime int64, publicKey string, moduleData ModuleData) (map[string][]string, error) {
	params := url.Values{}
	moduleJsonData, err := json.Marshal(moduleData)
	if err != nil {
		//fmt.Printf("转码成json串错误：%v\n", err)
		return nil, errors.New(fmt.Sprintf("转码成json串错误：%v\n", err))
	}
	//fmt.Printf("jsonData buildLoginParams:%v\n", string(moduleJsonData))

	pubErr := codec.RSA.Init(publicKey)
	if pubErr != nil {
		//fmt.Printf("初始化rsa加密错误：%v\n", pubErr)
		return nil, errors.New(fmt.Sprintf("初始化rsa加密错误：%v\n", pubErr))
	}
	data, err := codec.RSA.String(string(moduleJsonData), codec.MODE_PUBKEY_ENCRYPT)
	if err != nil {
		//fmt.Printf("rsa加密请求数据错误：%v\n", err)
		return nil, errors.New(fmt.Sprintf("rsa加密请求数据错误：%v\n", err))
	}
	params.Add("encryptData", data) //rsa加密后已base64编码

	hash := md5.New()
	hash.Write(moduleJsonData)
	cipherText := hash.Sum(nil)
	hexText := make([]byte, 32)
	hex.Encode(hexText, cipherText)
	params.Add("signData", string(hexText))

	now := time.Now().UnixNano() / int64(time.Millisecond)
	timestamp := now - diffTime
	timeData, err := codec.RSA.String(string([]byte(strconv.FormatInt(timestamp, 10))), codec.MODE_PUBKEY_ENCRYPT)
	if err != nil {
		//fmt.Printf("rsa加密时间戳错误：%v\n", err)
		return nil, errors.New(fmt.Sprintf("rsa加密时间戳错误：%v\n", err))
	}
	params.Add("time", string(timeData))
	//fmt.Printf("buildLoginParams:%v\n", params)
	return params, nil
}

func loginLicense(loginURL string, params map[string][]string) ([]byte, error) {
	licenseCheckResp, err := http.PostForm(loginURL, params)
	if err != nil {
		//fmt.Printf("登录license服务异常：%v\n", err)
		return nil, errors.New(fmt.Sprintf("登录license服务异常：%v\n", err))
	}
	defer func() ([]byte, error) {
		err := licenseCheckResp.Body.Close()
		if err == nil {
		}
		//logs.Error(err)
		return nil, err
	}()
	body, err := ioutil.ReadAll(licenseCheckResp.Body)
	if err != nil {
		//fmt.Printf("读取license服务返回值异常：%v\n", err)
		return nil, errors.New(fmt.Sprintf("读取license服务返回值异常：%v\n", err))
	}
	//fmt.Printf("license login response:%v\n", string(body))
	return body, nil
}

func (licenseService *LicenseService) Logout() error {
	serverTime, err := getServerTime(licenseService.EpgislicenseAPIUri + API_SERVERTIME)
	if err != nil {
		return errors.New(fmt.Sprintf("获取服务器时间戳异常：%v\n", err))
	}
	now := time.Now().UnixNano() / int64(time.Millisecond)
	diffTime := now - serverTime

	doLogout(_id, diffTime, licenseService)

	return nil
}

func buildLogoutParams(id string, diffTime int64, publicKey string) (map[string][]string, error) {
	params := url.Values{}
	encryptData := fmt.Sprintf("{\"id\":\"%s\"}", id)
	//fmt.Printf("jsonData buildLogoutParams:%v\n", string(encryptData))

	pubErr := codec.RSA.Init(publicKey)
	if pubErr != nil {
		//fmt.Printf("初始化rsa加密错误：%v\n", pubErr)
		return nil, errors.New(fmt.Sprintf("初始化rsa加密错误：%v\n", pubErr))
	}
	data, err := codec.RSA.String(string(encryptData), codec.MODE_PUBKEY_ENCRYPT)
	if err != nil {
		//fmt.Printf("rsa加密请求数据错误：%v\n", err)
		return nil, errors.New(fmt.Sprintf("rsa加密请求数据错误：%v\n", err))
	}
	params.Add("encryptData", data) //rsa加密后已base64编码

	hash := md5.New()
	hash.Write([]byte(encryptData))
	cipherText := hash.Sum(nil)
	hexText := make([]byte, 32)
	hex.Encode(hexText, cipherText)
	params.Add("signData", string(hexText))

	now := time.Now().UnixNano() / int64(time.Millisecond)
	timestamp := now - diffTime
	timeData, err := codec.RSA.String(string([]byte(strconv.FormatInt(timestamp, 10))), codec.MODE_PUBKEY_ENCRYPT)
	if err != nil {
		//fmt.Printf("rsa加密时间戳错误：%v\n", err)
		return nil, errors.New(fmt.Sprintf("rsa加密时间戳错误：%v\n", err))
	}
	params.Add("time", string(timeData))
	//fmt.Printf("buildLogoutParams:%v\n", params)
	return params, nil
}

func doLogout(id string, diffTime int64, licenseService *LicenseService) error {
	if id == "" {
		return errors.New("没有登录ID，退出失败！")
	}
	params, err := buildLogoutParams(id, diffTime, licenseService.PublicKey)
	if err != nil {
		//fmt.Printf("拼凑退出license服务请求参数异常：%v\n", err)
		return errors.New(fmt.Sprintf("拼凑退出license服务请求参数异常：%v\n", err))
	}

	err = logoutLicense(licenseService.EpgislicenseAPIUri+API_LOGOUT, params)
	if err != nil {
		//fmt.Printf("退出license服务异常：%v\n", err)
		return errors.New(fmt.Sprintf("退出license服务异常：%v\n", err))
	}
	return nil
}

func logoutLicense(logoutAPIUrl string, params map[string][]string) error {
	logoutLicenseResp, err := http.PostForm(logoutAPIUrl, params)
	if err != nil {
		//fmt.Printf("退出license服务异常：%v\n", err)
		return errors.New(fmt.Sprintf("退出license服务异常：%v\n", err))
	}
	defer func() error {
		err := logoutLicenseResp.Body.Close()
		if err == nil {
		}
		//logs.Error(err)
		return err
	}()
	body, err := ioutil.ReadAll(logoutLicenseResp.Body)
	if err != nil {
		//fmt.Printf("读取退出license服务返回值异常：%v\n", err)
		return errors.New(fmt.Sprintf("读取退出license服务返回值异常：%v\n", err))
	}
	fmt.Printf("license logout response:%v\n", string(body))
	return nil
}
