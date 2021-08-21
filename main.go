package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// URL samsung
var URL_KRX string = "http://kind.krx.co.kr/corpgeneral/corpList.do?method=download&searchType=13"
var (
	companyInfo dataframe.DataFrame
)

func GetCompanyInfo() {
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
	re1, _ := regexp.Compile(`(<[/]?)th`)
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
	companyInfo = cs.Rename("company", "회사명").
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

func CreateCompanyInfoTable() {
	db, err := sql.Open("mysql", "root:zhfeman@tcp(192.168.29.209:3306)/MUQUANT")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Connect and check the server version
	//var version string
	sql := `CREATE TABLE IF NOT EXISTS company_info (
            	code VARCHAR(20),
                company VARCHAR(40),
                last_update DATE,
                PRIMARY KEY(CODE) )`
	db.Exec(sql)

	sql = `CREATE TABLE IF NOT EXISTS daily_price (
        		code VARCHAR(20),
            	date DATE,
            	open BIGINT(20),
            	high BIGINT(20),
            	low BIGINT(20),
            	close BIGINT(20),
            	diff BIGINT(20),
            	volume BIGINT(20),
            	PRIMARY KEY(code, date) )`

	db.Exec(sql)

	//fmt.Println("Connected to:", version)
}

func UpdateCompanyInfo() {
	db, err := sql.Open("mysql", "root:zhfeman@tcp(192.168.29.209:3306)/MUQUANT")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	today := time.Now().Format("2006-01-02")

	for i := 0; i < companyInfo.Nrow(); i++ {
		// 선택하는 방법이 pandas에 비해서 매끄럽지 못함.
		// 2가지 방법이 있을 것 같은데.. 조금 더 범용성을 주고자 Column을 선택할 수 있게 함.
		// 추후, code, company 외의 더 자료가 필요할 경우를 대비.
		// code := companyInfo.Elem(i, 0).String()
		// code1 := companyInfo.Elem(i, 1).String()
		codeMap := companyInfo.Subset(i).Maps()[0]
		code := codeMap["code"]
		company := codeMap["company"]

		sql := fmt.Sprintf("REPLACE INTO company_info ( code, company, last_update ) VALUES ('%s', '%s', '%s')", code, company, today)
		db.Exec(sql)
	}
}

func main() {
	GetCompanyInfo()
	CreateCompanyInfoTable()
	UpdateCompanyInfo()

	fmt.Println("test")
}
