package main

import (
	"SGX_blockchain/src/db"
	"SGX_blockchain/src/server"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

//var (
//	testmsg, _     = hex.DecodeString("ce0677bb30baa8cf067c88db9811f4333d131bf8bcf12fe7065d211dce971008")
//	testsig, _     = hex.DecodeString("90f27b8b488db00b00606796d2987f6a5f59ae62ea05effe84fef5b8b0e549984a691139ad57a3f0b906637673aa2f63d1f55cb1a69199d4009eea23ceaddc9301")
//	testpubkey, _  = hex.DecodeString("04e32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a0a2b2667f7e725ceea70c673093bf67663e0312623c8e091b13cf2c0f11ef652")
//	testpubkeyc, _ = hex.DecodeString("02e32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a")
//)

func logMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 打印请求的方法，URL，头部和主体
		fmt.Println("Request:")
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(dump))

		// 调用处理函数
		next(w, r)

		// 打印响应的状态码，头部和主体
		//fmt.Println("Response:")
		//dump, err = httputil.DumpResponse(*w.(http.Response), true)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//fmt.Println(string(dump))
	}
}

func main() {

	//cert, priv := utils.CreateCertificate()
	//hash := sha256.Sum256(cert)
	//report, err := utils.GetRemoteReport(hash[:])
	//tlsCfg := tls.Config{
	//	Certificates: []tls.Certificate{
	//		{
	//			Certificate: [][]byte{cert},
	//			PrivateKey:  priv,
	//		},
	//	},
	//}
	//http.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request) {
	//	type Report struct {
	//		Status string `json:"status"`
	//		Data   struct {
	//			Report string `json:"report"`
	//			Cert   string `json:"cert"`
	//		} `json:"data"`
	//	}
	//	rep := base64.StdEncoding.EncodeToString(report)
	//	c := base64.StdEncoding.EncodeToString(cert)
	//	resp := &Report{
	//		Status: "ok",
	//		Data: struct {
	//			Report string `json:"report"`
	//			Cert   string `json:"cert"`
	//		}{
	//			Report: rep,
	//			Cert:   c,
	//		},
	//	}
	//	b, _ := json.Marshal(resp)
	//	w.Write(b)
	//})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(""))
	})

	d := db.InitMemorydb()
	h := server.NewHandler(d)

	//testserver := httptest.NewServer(http.HandlerFunc(h.AccountInfoHandler))

	http.HandleFunc("/account/info", logMiddleware(h.AccountInfoHandler))
	http.HandleFunc("/files/store", logMiddleware(h.FileStoreHandler))
	http.HandleFunc("/files/retrieve", logMiddleware(h.FileRetrieveHandler))
	http.HandleFunc("/kv/store", logMiddleware(h.KVStoreHandler))
	http.HandleFunc("/kv/retrieve", logMiddleware(h.KVRetrieveHandler))
	http.HandleFunc("/block/info", logMiddleware(h.BlockInfoHandler))
	http.HandleFunc("/transaction/info", logMiddleware(h.TransactionInfoHandler))
	http.HandleFunc("/contract/deploy", logMiddleware(h.ContractDeployHandler))
	http.HandleFunc("/contract/call", logMiddleware(h.ContractCallHandler))

	//httpServer := http.Server{Addr: "127.0.0.1:8888", TLSConfig: &tlsCfg}
	httpServer := http.Server{Addr: "0.0.0.0:8888"}

	fmt.Println("listening ...")
	httpServer.ListenAndServe()
	//err = httpServer.ListenAndServeTLS("", "")
	//fmt.Println(err)
}
