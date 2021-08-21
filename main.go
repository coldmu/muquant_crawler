package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

// URL samsung
var URL_KRX string = "http://kind.krx.co.kr/corpgeneral/corpList.do?method=download&searchType=13"
var (
	stockCodeName dataframe.DataFrame
)

func GetKrxData() {
	// krx 홈페이지에 접속하여 request
	resp, err := http.Get(URL_KRX)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// gota BUG 로 추정, ReadHTML에서 <th> 를 제목으로 인식못함. <th> -> <td> 로 변경해야함.
	re1, _ := regexp.Compile("(<[\\/]?)th")
	html := re1.ReplaceAllString(string(data), "${1}td")

	// 한글인코딩이 깨져서, 한글을 UTF8로 변환? 헷갈리는 부분.
	// 보통 NewEncoder()로 UTF8 -> euc-kr 로 변환해서 출력해야하는데...
	// 이상하게 반대로 euc-kr -> UTF8 로 변환해서 출력해야 잘 됨.
	htmlUTF8, _, err := transform.String(korean.EUCKR.NewDecoder(), html)
	if err != nil {
		panic(err)
	}

	// parser 를 이용하여 가져온 html 태그를 gota 의 ReadHTML 를 이용하여 dataframe 형식으로 전달
	// 자동으로 hasHeader는 true 로 지정되어 있음.
	// ReadHTML 를 사용하지 않고, goquery 를 사용해도 되나, 만들어진 ReadHTML를 사용하고 싶어서 삽질 후 성공!!
	cs := dataframe.ReadHTML(strings.NewReader(htmlUTF8), dataframe.DetectTypes(false), dataframe.DefaultType(series.String))[0]

	// 필요한 정보만 Select
	// code번호를 6자리로 바꾸고, 헤더를 영어로 변경
	stockCodeName = cs.Rename("company", "회사명").
		Rename("code", "종목코드").
		Select([]string{"code", "company"})

	/*
		!!! 기록용 !!!
		series 를 하나 떼어 낸 다음에, 0을 붙일 수 있다.
		Map을 이용하여 함수를 호출하여 하나씩 변경.
		하지만, series의 type이 변경되지 않기 때문에 역시나 다시 00을 붙이더라도, input 할 때 다시 0이 사라진다.
		일단은 해결하지 못했지만, 값을 읽어들일 때 string형식으로 삽입하면 해결 됨.

		sel2 := sel1.Col("code")
		zeroFill := func(e series.Element) series.Element {
			result := e.Copy()
			zeroResult := fmt.Sprintf("%06s", result.String())
			result.Set(zeroResult)
			return series.Element(result)
		}
		received := sel2.Map(zeroFill)
	*/

}
func main() {
	GetKrxData()
	fmt.Print(stockCodeName)
}
