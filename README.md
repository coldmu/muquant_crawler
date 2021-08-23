# muquant_crawler(돌연변이 퀀트) - My first GO project
GO 를 이용하여 취미로 시작하는 퀀트 프로젝트
 - Python pandas, BeautifulSoup를 이용하여 프로젝트 진행 중, 크롤링이 너무너무 오래 걸려 go lang으로 코딩해보기로 결심.


### Step1) krx 종목코드 크롤링
 - request 이용하여 크롤링
 - 인코딩해결
 - gota 의 ReadHTML 사용
 - 필요한 정보만 파싱(종목코드, 회사명?) >> 2021-08-21 21h
 - test로 MariaDB에 넣어보기 >> 2021-08-21 22h 50m
 - **mongoDB 에 넣어보기 (나중에)**

### Step2) 네이버 크롤링
 - 종목코드를 이용하여 전체 크롤링 >> 2021-08-23 01h
 - 크롤링한 정보를 DB에 삽입
 - 병렬처리 >> 2021-08-24 12h 40m (진행중)
 - 파이썬을 이용하여 크롤링하였으나, 너무 오래걸려서 학습겸 GO로 변경 잘될것으로 기대..

<br>

## Python -> Go 대체
 - beautifulsoup : github.com/PuerkitoBio/goquery
 - pandas : github.com/go-gota/gota/dataframe, github.com/go-gota/gota/series

<br>

## Memo
- 디버깅시 unused 에러를 피하기 위해선, ` _ = df ` 와 같이 임시로 사용

- 이중 slice 만들기 / 데이터 삽입은 append 이용
```
crawlStockInfo := make([][]string, 7)
```
- request 헤더 변경
```
pageURL := fmt.Sprintf("%s&page=%d", url, pageNum)
req2, err := http.NewRequest("GET", pageURL, nil)
if err != nil {
	log.Fatal(err)
}
req2.Header.Add("User-agent", "Mozilla/5.0")

client2 := &http.Client{}
res2, _ := client2.Do(req2)
if res.StatusCode != 200 {
	log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
}
defer res2.Body.Close()
```
- request 로 소스불러오기
```
resp, err := http.Get(URL_KRX)
if err != nil {
	panic(err)
}
defer resp.Body.Close()

data, err := ioutil.ReadAll(resp.Body)
if err != nil {
	panic(err)
}
```
- jQuery 에서 원하는 정보 Find 후, 문자열 숫자로 변경
```
tempHref, _ := doc.Find(".pgRR a").Attr("href")
tempPgrr := strings.Split(tempHref, "=")
pgrr, _ := strconv.Atoi(tempPgrr[len(tempPgrr)-1])
```

- jQuery 에서 Find, Each 사용
```
// tr을 찾은 후, tr의 자식노드(td)가 7인 경우만 로드를 시작.
// 자식노드(td)의 값을 하나씩 불러와서 화이트스페이스 제거 후, 값만 추출하여 slice에 삽입
doc2.Find("tr").Each(func(idx int, sel *goquery.Selection) {
	if sel.Find("td").Length() == 7 {
		sel.Find("td").Each(func(idx2 int, sel2 *goquery.Selection) {
			if idx2 == 2 {
				tt := strings.ReplaceAll(sel2.Text(), "\n", "")
				tt = strings.ReplaceAll(tt, "\t", "")
				crawlStockInfo[idx2] = append(crawlStockInfo[idx2], tt)
			} else {
				crawlStockInfo[idx2] = append(crawlStockInfo[idx2], sel2.Text())
			}
		})
	}
})
```
- golang slice 사용 (make 로 미리 할당하나, append로 하나씩 넣어주나 시간복잡도 비슷한 것으로 추정)<br>
`crawlStockInfo := make([][]string, 7)`

- 파일출력
```
file, err := os.Create("output.txt") // output.txt 파일 열기
if err != nil {
	panic(err)
}
defer file.Close() // main 함수가 끝나기 직전에 파일을 닫음
fmt.Fprint(file, "출력할 내용")
```

- 시간비교
```
startTime := time.Now()
endTime := time.Now()
fmt.Println(endTime.Sub(startTime))
```

- gota Rbind 이용하여 삽입
```
addRow := dataframe.New(
	series.New(crawlStockInfo[0], series.String, "date"),
	series.New(crawlStockInfo[1], series.String, "close"),
	series.New(crawlStockInfo[2], series.String, "diff"),
	series.New(crawlStockInfo[3], series.String, "open"),
	series.New(crawlStockInfo[4], series.String, "high"),
	series.New(crawlStockInfo[5], series.String, "low"),
	series.New(crawlStockInfo[6], series.String, "volume"),
)
df = df.RBind(addRow)
```
- gota dataframe 만들기
```
df := dataframe.New(
	series.New(crawlStockInfo[0], series.String, "date"),
	series.New(crawlStockInfo[1], series.String, "close"),
	series.New(crawlStockInfo[2], series.Int, "diff"),
	series.New(crawlStockInfo[3], series.Int, "open"),
	series.New(crawlStockInfo[4], series.Int, "high"),
	series.New(crawlStockInfo[5], series.Int, "low"),
	series.New(crawlStockInfo[6], series.Int, "volume"),
)
```
- gota 데이터변경(Series Map 사용) 
```
// fill zero-padding
sel2 := sel1.Col("code")
zeroFill := func(e series.Element) series.Element {
	result := e.Copy()
	zeroResult := fmt.Sprintf("%06s", result.String())
	result.Set(zeroResult)
	return series.Element(result)
}
received := sel2.Map(zeroFill)
```
- gota 데이터 이름 변경 및 선택
```
companyInfo = cs.Rename("company", "회사명").
	Rename("code", "종목코드").
	Select([]string{"code", "company"})
```
- gota ReadHTML (option 넣어서 기본 타입을 String으로 설정) --> ReadHTML의 경우 th인식불가, td 안에 불필요한 데이터 있을 시 인식못함.<br>
`cs := dataframe.ReadHTML(strings.NewReader(htmlUTF8), dataframe.DetectTypes(false), dataframe.DefaultType(series.String))[0]`
- mariaDB 접속 및 삽입
```
db, err := sql.Open("mysql", "root:password@tcp(ip:port)/DBNAME")
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
```
- 정규표현식으로 문자열 수정
```
re, _ := regexp.Compile(`(<[/]?)th`)
html := re.ReplaceAllString(string(data), "${1}td")
```
- strings.ReplaceAll 사용<br>
`strings.ReplaceAll("\n\nTEST\n\n", "\n", "")`


- 한글인코딩
```
// 한글인코딩이 깨져서, 한글을 UTF8로 변환? 헷갈리는 부분.
// 보통 NewEncoder()로 UTF8 -> euc-kr 로 변환해서 출력해야하는데...
// 이상하게 반대로 euc-kr -> UTF8 로 변환해서 출력해야 잘 됨.
htmlUTF8, _, err := transform.String(korean.EUCKR.NewDecoder(), html)
if err != nil {
	panic(err)
}
```

- `struct {}` 의 경우 사이즈가 0이므로, 유용하게 사용할 수 있다.

- 병렬처리
	* golang에서는 goroutine을 사용하여 병렬처리(동시성 시분할)
	* 병렬처리랑은 다르기 때문에, CPU 갯수를 늘려주어야 함<br>
	`runtime.GOMAXPROCS(runtime.NumCPU())`
	* 또한 채널을 이용하여 수신받고, 송신하여 파이프 생성
```
runtime.GOMAXPROCS(runtime.NumCPU())
guard := make(chan struct{}, runtime.NumCPU())

seriesCompanyCode := companyInfo.Col("code")

for i := 0; i < seriesCompanyCode.Len(); i++ {
	code := seriesCompanyCode.Val(i).(string)
	guard <- struct{}{}		// 버퍼 해소시까지, 대기
	go func(s string) {
		ReadNaver(s)
		<-guard				// 수신을 함으로써 버퍼 1개 해소.
	}(code)
}
```