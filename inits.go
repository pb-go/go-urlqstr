package urlqstr

import "sync"

const TAGKEY = "qstr"

type Marshaller interface {
	Marshal(DataIntf interface{}) (string,error)
}

type UnMarshaller interface {
	UnMarshal() (interface{}, error)
}

type ExtraFunc interface {
	Sort()
	InsertKey(key string, value string)
	AppendKey(key string, value string)
	DeleteKey(key string)
}

type internalFunc interface {
	generateQStr()
	extractDataFromStruct(DataIntf interface{})
	writeStructFromQStr()
}

type Field2TagByName struct{
	mu sync.Mutex
	val map[string]string
}

type TagNm2TagVal struct {
	mu sync.Mutex
	val map[string]string
}

type UQueryString struct {
	f2Tbn *Field2TagByName
	tn2Tv *TagNm2TagVal
	DataIntf interface{}
	DataStr string
}

