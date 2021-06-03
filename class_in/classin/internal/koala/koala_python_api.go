package koala

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/alecthomas/log4go"
	"github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"zhiyuan/koalamate_statistics_server/configs"
)


//var jar, _ = cookiejar.New(nil)
var KoalaMateCookie []*http.Cookie

//func doResponse(resp *http.Response) (*simplejson.Json, error) {
//	log4go.Debug(resp.Status)
//	if resp.StatusCode != 200 {
//		return nil, errors.New(resp.Status)
//	}
//
//	//defer resp.Body.Close()
//
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		log4go.Error(err.Error())
//		return nil, errors.New("Read response body error")
//	}
//	log4go.Debug(string(body))
//
//	jdata, err := simplejson.NewJson(body)
//	if err != nil {
//		log4go.Error(err.Error())
//		return nil, errors.New("koalamate 返回报文错误")
//	}
//
//	code, _ := jdata.Get("code").Int()
//	if code != 0 {
//		desc, _ := jdata.Get("err_msg").String()
//		log4go.Error(desc)
//		return nil, errors.New(desc)
//	}
//	return jdata, nil
//}

func ReverseProxy(c *gin.Context) {
	baseUrl := configs.Gconf.KoalaApisHost
	fmt.Println(configs.Gconf.KoalaApisHost)
	//baseUrl := "http:127.0.0.1:5001"
	remote, err := url.Parse(baseUrl)
	if err != nil {
		log4go.Error(err)
		c.Writer.Write([]byte(err.Error()))
		return
	}
	fmt.Println(remote)
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(c.Writer, c.Request)

}

func ReverseProxy2(w http.ResponseWriter, r *http.Request) {
	baseUrl := configs.Gconf.KoalaApisHost
	remote, err := url.Parse(baseUrl)
	if err != nil {
		log4go.Error(err)
		w.Write([]byte(err.Error()))
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)

}

func ModifyResponse(w http.ResponseWriter, rec *httptest.ResponseRecorder, jdata *simplejson.Json, level int) {
	code, _ := jdata.Get("code").Int()
	if code != 0 {
		desc, _ := jdata.Get("err_msg").String()
		log4go.Error(desc)
		w.Write([]byte(rec.Body.Bytes()))
		return
	}

	// after this finishes, we have the response recorded
	// and can modify it before copying it to the original RW

	// we copy the original headers first
	for k, v := range rec.Header() {
		w.Header()[k] = v
	}
	// and set an additional one
	//w.Header().Set("X-We-Modified-This", "Yup")
	// only then the status code, as this call writes out the headers
	w.WriteHeader(200)

	// The body hasn't been written (to the real RW) yet,
	// so we can prepend some data.
	data := simplejson.New()
	data.Set("level", level)
	jdata.Set("data", data)
	byte_data, _ := jdata.MarshalJSON()
	log4go.Debug(string(byte_data))

	// But the Content-Length might have been set already,
	// we should modify it by adding the length
	// of our own data.
	// Ignoring the error is fine here:
	// if Content-Length is empty or otherwise invalid,
	// Atoi() will return zero,
	// which is just what we'd want in that case.
	clen := len(byte_data)
	w.Header().Set("Content-Length", strconv.Itoa(clen))

	// finally, write out our data
	w.Write(byte_data)
	// then write out the original body
	//w.Write(rec.Body.Bytes())

}

//python接口的权限控制
func AuthRequiredPython() gin.HandlerFunc {
	return func(c *gin.Context) {
		if cookie, err := c.Request.Cookie("username"); err == nil {
			value := cookie.Value
			fmt.Println(value)
			if value == "admin" || strings.Contains(value, "@") {
				c.Next()
				return
			} else {
				c.JSON(http.StatusOK, gin.H{
					"code": -100,
					"err_msg": "无权限访问!",
				})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": -100,
				"err_msg": "请先登录!",
			})
			c.Abort()
			return
		}
	}
}

//本地接口的权限控制
func AuthRequiredlocal() gin.HandlerFunc {
	return func(c *gin.Context) {
		if cookie, err := c.Request.Cookie("username"); err == nil {
			value := cookie.Value
			fmt.Println(value)
			if strings.Contains(value, "@") {
				c.Next()
				return
			} else {
				c.JSON(http.StatusOK, gin.H{
					"code": -100,
					"err_msg": "无权限访问!",
				})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": -100,
				"err_msg": "请先登录!",
			})
			c.Abort()
			return
		}
	}
}

/**
登陆伴侣
*/
func LoginZybox(level int, c *gin.Context) {
	r := c.Request
	w := c.Writer
	rec := httptest.NewRecorder()
	ReverseProxy2(rec, r)

	jdata, err := simplejson.NewJson(rec.Body.Bytes())
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -100,
			"err_msg": "登陆异常!",
		})
		return
	}

	ModifyResponse(w, rec, jdata, level)

}

//查看当前cookie中的信息
func dbgPrintCurCookies() map[string]interface{} {
	var cookieNum int = len(KoalaMateCookie)
	log4go.Info("cookieNum=%d", cookieNum)
	var info = make(map[string]interface{})
	for i := 0; i < cookieNum; i++ {
		var curCk *http.Cookie = KoalaMateCookie[i]
		//gLogger.Info("curCk.Raw=%s", curCk.Raw)
		log4go.Info("------ Cookie [%d]------", i)
		log4go.Info("Name\t=%s", curCk.Name)
		log4go.Info("Value\t=%s", curCk.Value)
		log4go.Info("Path\t=%s", curCk.Path)
		log4go.Info("Domain\t=%s", curCk.Domain)
		log4go.Info("Expires\t=%s", curCk.Expires)
		log4go.Info("RawExpires=%s", curCk.RawExpires)
		log4go.Info("MaxAge\t=%d", curCk.MaxAge)
		log4go.Info("Secure\t=%t", curCk.Secure)
		log4go.Info("HttpOnly=%t", curCk.HttpOnly)
		log4go.Info("Raw\t=%s", curCk.Raw)
		log4go.Info("Unparsed=%s", curCk.Unparsed)
		info[curCk.Name] = curCk.Value
	}
	return info
}

func LogRequest(c *gin.Context) {
	log4go.Info(c.Request.RequestURI)
	log4go.Info(c.Request.Method)
	log4go.Info(c.Request.Header)
	log4go.Info(c.Request.Header["Tag"])
	if c.Request.Form == nil {
		c.Request.ParseMultipartForm(32 << 20)
	}
	for k, v := range c.Request.Form {
		log4go.Info(k)
		log4go.Info(v)
	}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

