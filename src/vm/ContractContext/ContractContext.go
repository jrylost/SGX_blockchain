package ContractContext

import "time"

type Context struct {
	stub Stub
}

type Stub struct {
	strdb          map[string]string
	clientIdentity *ClientIdentity
}

type ClientIdentity struct {
	MSPID string
}

func Initial() *Context {
	return &Context{}
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
	s.strdb[key] = value
}

func (s *Stub) GetStringState(key string) (string, error) {

	return s.strdb[key], nil
}

func (s *Stub) GetTxId() string {
	return "0xaaaaaaaaaaaaaaaaaaaaaaaaaa"
}

func (s *Stub) GetQueryResult(query string) []string {
	return []string{"0x11111"}
}

func (s *Stub) GetTxTimestamp() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05Z")
}

func (c *Context) Getstub() *Stub {
	db := make(map[string]string)
	return &Stub{strdb: db, clientIdentity: &ClientIdentity{MSPID: "0x111"}}
}
