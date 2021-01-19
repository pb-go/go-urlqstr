package urlqstr

import "sync"

func (uqstr *UQueryString) generateQStr() {

}

func (uqstr *UQueryString) init(){
	// get all struct fields name and its tag names
	uqstr.f2Tbn = &Field2TagByName{
		mu:  sync.Mutex{},
		val: make(map[string]string),
	}
	uqstr.tn2Tv = &TagNm2TagVal{
		mu:  sync.Mutex{},
		val: make(map[string]string),
	}
}
