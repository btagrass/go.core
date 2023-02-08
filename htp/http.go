package htp

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/btagrass/go.core/app"
	"github.com/btagrass/go.core/utl"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

var (
	Ip      string        // Ip地址
	Port    uint          // 端口
	Timeout time.Duration // 超时
)

// 初始化
func init() {
	Ip = viper.GetString("http.ip")
	if Ip == "" {
		Ip, _ = utl.GetIp()
	}
	Port = viper.GetUint("http.port")
	var err error
	Timeout, err = time.ParseDuration(viper.GetString("http.timeout"))
	if err != nil {
		Timeout = 7 * time.Second
	}
}

// 获取
func Get(url string, r ...any) (string, error) {
	req := resty.New().
		SetTimeout(Timeout).
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		}).
		OnAfterResponse(respond).
		R().
		ForceContentType("application/json")
	if len(r) > 0 {
		req.SetResult(r[0])
	}
	resp, err := req.Get(GetUrl(url))
	logrus.Debugf("method: %s, url: %s -> %s", req.Method, req.URL, resp)

	return resp.String(), err
}

// 获取网址
func GetUrl(url string) string {
	if strings.HasPrefix(url, "http") {
		return url
	}
	if strings.HasPrefix(url, "/") {
		return fmt.Sprintf("http://%s:%d%s", Ip, Port, url)
	}

	return fmt.Sprintf("http://%s:%d/%s", Ip, Port, url)
}

// 提交
func Post(url string, data any, r ...any) (string, error) {
	req := resty.New().
		SetTimeout(Timeout).
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		}).
		OnAfterResponse(respond).
		R().
		SetHeader("Accept", "*/*").
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		ForceContentType("application/json")
	if len(r) > 0 {
		req.SetResult(r[0])
	}
	resp, err := req.Post(GetUrl(url))
	logrus.Debugf("method: %s, url: %s, data: %s -> %s", req.Method, req.URL, data, resp)

	return resp.String(), err
}

// 提交文件
func PostFile(url string, data map[string]string, r ...any) error {
	req := resty.New().
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		}).
		OnAfterResponse(respond).
		R().
		SetFiles(data).
		ForceContentType("application/json")
	if len(r) > 0 {
		req.SetResult(r[0])
	}
	resp, err := req.Post(GetUrl(url))
	logrus.Debugf("method: %s, url: %s, data: %s -> %s", req.Method, req.URL, data, resp)

	return err
}

// 提交表单
func PostForm(url string, data map[string]string, r ...any) (string, error) {
	req := resty.New().
		SetTimeout(Timeout).
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		}).
		OnAfterResponse(respond).
		R().
		SetHeader("Accept", "*/*").
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(data).
		ForceContentType("application/json")
	if len(r) > 0 {
		req.SetResult(r[0])
	}
	resp, err := req.Post(GetUrl(url))
	logrus.Debugf("method: %s, url: %s, data: %s -> %s", req.Method, req.URL, data, resp)

	return resp.String(), err
}

// 保存文件
func SaveFile(url string, file ...string) error {
	filePath := filepath.Base(url)
	if len(file) > 0 {
		filePath = file[0]
	}
	req := resty.New().
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		}).
		OnAfterResponse(respond).
		R().
		SetOutput(filePath)
	resp, err := req.Get(GetUrl(url))
	logrus.Debugf("method: %s, url: %s -> %s", req.Method, req.URL, resp)

	return err
}

// 响应
func respond(c *resty.Client, resp *resty.Response) error {
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf(resp.Status())
	} else {
		r := cast.ToStringMap(resp.String())
		code, ok := r["code"]
		if !ok {
			code, ok = r["error_code"]
		}
		if ok {
			code = cast.ToInt(code)
			if code != 0 && code != http.StatusOK {
				msg, ok := r["msg"]
				if !ok {
					msg = r["desp"]
				}
				return fmt.Errorf(cast.ToString(msg))
			}
		}
	}

	return nil
}
