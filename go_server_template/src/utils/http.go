package utils

import (
	"github.com/ddliu/go-httpclient"
	"net/http"
	"time"
)

// GetHttpReuseClient 获取客户端，http连接复用
func GetHttpReuseClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     10 * time.Minute,
		},
		Timeout: 30 * time.Second,
	}

}

func HttpRequest(url string, param map[string]string, headers map[string]string) ([]byte, error) {
	var (
		resp *httpclient.Response
		err  error
	)
	if headers == nil {
		if param == nil {
			resp, err = httpclient.NewHttpClient().Get(url)
		} else {
			resp, err = httpclient.NewHttpClient().Get(url, param)
		}
	} else {
		if param == nil {
			resp, err = httpclient.NewHttpClient().WithHeaders(headers).Get(url)
		} else {
			resp, err = httpclient.NewHttpClient().WithHeaders(headers).Get(url, param)
		}
	}
	if err != nil {
		return nil, UpstreamServiceRequestError.Wrapf(err, "call url[%s] error", url)
	}
	if resp.StatusCode != 200 {
		return nil, UpstreamServiceResponseError.Newf("call url[%s] status[%d] isn't 200", url, resp.StatusCode)
	}
	buf, err := resp.ReadAll()
	if err != nil {
		return nil, UpstreamServiceResponseError.Wrapf(err, "call url[%s] read body error", url)
	}
	return buf, err
}
