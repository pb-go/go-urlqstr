package urlqstr

import (
	"fmt"
	"reflect"
	"net/url"
)

func (uqstr *UQueryString) Marshal(DataIntf interface{}) (string,error){
	uqstr.init()
	uqstr.DataIntf = DataIntf
	uqstr.extractDataFromStruct(DataIntf)
}

// extractDataFromStruct is an internal function which used to recursively extract
// data from any structure and parse the tag value. Since HTTP QueryString is only a single level,
// we use coroutine to speed up.
func (uqstr *UQueryString) extractDataFromStruct(DataIntf interface{}){
	// make channel as signal to avoid use sync.WaitGroup
	// https://golang.org/ref/spec#Receive_operator
	sig1 := make(chan interface{})
	sig2 := make(chan interface{})
	ifaceTyps := reflect.TypeOf(DataIntf)
	ifaceVals := reflect.ValueOf(DataIntf)

	// dereference Ptr to avoid have an address introduced
	if ifaceTyps.Kind() == reflect.Ptr {
		ifaceTyps = ifaceTyps.Elem()
		ifaceVals = ifaceVals.Elem()
	}

	// get name of all structs fields and its corresponding tag names
	// gkd gkd
	go func() {
		defer close(sig1)
		for i := 0; i < ifaceTyps.NumField(); i++ {
			curfd := ifaceTyps.Field(i)
			uqstr.f2Tbn.mu.Lock()
			dataA, exists := curfd.Tag.Lookup(TAGKEY)
			if exists {
				var afterProcDataA = dataA
				// invalid symbols and delimiters should not be directly put into the querystring
				// url encode first
				if !isValidTag(dataA) {
					afterProcDataA = url.QueryEscape(dataA)
				}
				uqstr.f2Tbn.val[curfd.Name] = afterProcDataA
			}
		}
	}()

	// recursive resolve the struct value and corresponding fields
	// gkd gkd
	go func() {
		defer close(sig2)
		typeOfDts := ifaceVals.Type()
		for i:= 0; i < ifaceVals.NumField(); i++ {
			curfd := ifaceVals.Field(i)
			// dereference if any ptr exists
			if curfd.Kind() == reflect.Ptr {
				curfd = curfd.Elem()
			}
			switch curfd.Kind(){
			case reflect.Struct:
				uqstr.extractDataFromStruct(curfd.Interface())
			default:
				fdname := typeOfDts.Field(i).Name
				// transform interface to string, DO NOT USE `.String()` directly!
				// which will result in "<int value>" not "2"
				fdval := fmt.Sprintf("%v", curfd.Interface())
				uqstr.tn2Tv.mu.Lock()
				uqstr.tn2Tv.val[fdname] = fdval
				uqstr.tn2Tv.mu.Unlock()
			}
		}
	}()

	// wait until all coroutine finished to return
	<-sig1
	<-sig2
}

