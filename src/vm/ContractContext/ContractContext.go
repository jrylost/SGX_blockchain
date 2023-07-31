package ContractContext

import (
	"errors"
	"time"
)

type Context struct {
	stub Stub
}

type Stub struct {
	Strdb          map[string]string
	clientIdentity *ClientIdentity
	txHashWith0x   string
}

type ClientIdentity struct {
	MSPID string
}

type QueryContent struct {
	Selector string `json:"selector"`
}

func Initial(strdb map[string]string, txHashWith0x string) *Context {
	return &Context{stub: Stub{
		Strdb:          strdb,
		clientIdentity: &ClientIdentity{MSPID: "MSPID"},
		txHashWith0x:   txHashWith0x,
	}}
}

func (s *Stub) ClientIdentity() *ClientIdentity {
	return s.clientIdentity
}

func (i *ClientIdentity) GetMSPID() string {
	return i.MSPID
}

func (s *Stub) PutStringState(key, value string) {
	s.Strdb[key] = value
}

func (s *Stub) GetStringState(key string) (string, error) {
	val, ok := s.Strdb[key]
	if !ok {
		return "", errors.New("invalid key")
	}
	return val, nil
}

func (s *Stub) GetTxId() string {
	return s.txHashWith0x
}

func (s *Stub) GetQueryResult(query string) []string {
	//var queryContent QueryContent
	//err := json.Unmarshal([]byte(query), &queryContent)
	//if err != nil {
	//	return []string{""}
	//}
	//selector := queryContent.Selector
	//keys := gjson.Get(selector, "@keys")
	//for i, key := range keys.Array() {
	//	keyString := key.String()
	//	query := gjson.Get(selector, keyString)
	//}
	return []string{"0x11111"}
}

func (s *Stub) GetTxTimestamp() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05Z")
}

func (c *Context) Getstub() *Stub {
	return &c.stub
}
