package util

import (
	"fmt"
	"goutil/http"
	"testing"

	http "github.com/goUtil/http"
)

func TestHttp(t *testing.T) {
	res := http.HttpJsonPost(
		"http://jwxt.shufe-zj.edu.cn/jwglxt/jxdmgl/jsjxdm_cxJxbxxByJxbid.html?gnmkdm=N254350&su=021560",
		"{jxb_id: 38626,rq: '2020-05-29'}",
		map[string]string{
			"Cookie": "JSESSIONID=3D75076946D36C470484078F1CEBC9E8; route=5f3ea6db761cde4977dfa34d61160a36",
			"Origin": "http://jwxt.shufe-zj.edu.cn",
		},
	)
	fmt.Println(string(res))
}
