package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-gota/gota/dataframe"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

// URL samsung
var URL string = "https://finance.naver.com/item/main.nhn?code=005930" // 삼성전자
// URL1 investorDealTrendDay
var URL1 string = "https://finance.naver.com/sise/investorDealTrendDay.nhn?bizdate=" // 동향
var URL_KRX string = "http://kind.krx.co.kr/corpgeneral/corpList.do?method=download&searchType=13"
var (
	now_price  string
	prev_price string
	cs         []dataframe.DataFrame
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
	cs = dataframe.ReadHTML(strings.NewReader(htmlUTF8))

}
func main() {
	GetKrxData()
	fmt.Print(cs)
}
