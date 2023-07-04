//go:build !release

package config

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
)

var LOG_REQUEST_AND_RESPONSE = true

func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rec := httptest.NewRecorder()

		// 打印请求的方法，URL，头部和主体
		fmt.Println("Request:")
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(dump))

		// 调用处理函数
		next(rec, r)
		w.WriteHeader(rec.Code)
		if rec.Body != nil {
			w.Write(rec.Body.Bytes())
		}
		// 打印响应的状态码，头部和主体
		fmt.Println("Response:")
		dump, err = httputil.DumpResponse(rec.Result(), true)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//fmt.Println(string(dump))
	}
}
