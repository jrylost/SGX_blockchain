package ContractContext

import (
	"SGX_blockchain/src/utils"
	"crypto/rand"
	"time"
)

type Context struct {
	stub Stub
}

type Stub struct {
	Strdb          map[string]string
	clientIdentity *ClientIdentity
}

type ClientIdentity struct {
	MSPID string
}

func Initial(strdb map[string]string) *Context {
	return &Context{stub: Stub{
		Strdb:          strdb,
		clientIdentity: &ClientIdentity{MSPID: "MSPID"},
	}}
}

func (s *Stub) ClientIdentity() *ClientIdentity {
	return s.clientIdentity
}

func (i *ClientIdentity) GetMSPID() string {
	return i.MSPID
}

//type Stub interface {
//    GetStringState(key string) string
//    PutStringState(key, value string)
//    GetTxId() string
//    GetTxTimestamp() string
//}

func (s *Stub) PutStringState(key, value string) {
	s.Strdb[key] = value
	//fmt.Println("??????")
}

func (s *Stub) GetStringState(key string) (string, error) {
	return s.Strdb[key], nil
}

func (s *Stub) GetTxId() string {
	b := make([]byte, 32)
	rand.Read(b)
	res := utils.EncodeBytesToHexStringWith0x(b)
	return res
}

func (s *Stub) GetQueryResult(query string) []string {
	return []string{"0x11111"}
}

func (s *Stub) GetTxTimestamp() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05Z")
}

func (c *Context) Getstub() *Stub {
	return &c.stub
}
