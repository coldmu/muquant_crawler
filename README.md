# muquant_crawler(돌연변이 퀀트)
GO 를 이용하여 취미로 시작하는 퀀트 프로젝트

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
 - 병렬처리
 - 파이썬을 이용하여 크롤링하였으나, 너무 오래걸려서 학습겸 GO로 변경 잘될것으로 기대..

<br>

## Python -> Go 대체
 - beautifulsoup : github.com/PuerkitoBio/goquery
 - pandas : github.com/go-gota/gota/dataframe, github.com/go-gota/gota/series

<br>
## Memo
- 디버깅시 unused 에러를 피하기 위해선, ` _ = df ` 와 같이 임시로 사용


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

- jQuery 에서 원하는 정보 Find 후, 문자열 숫자로 변경
```
tempHref, _ := doc.Find(".pgRR a").Attr("href")
tempPgrr := strings.Split(tempHref, "=")
pgrr, _ := strconv.Atoi(tempPgrr[len(tempPgrr)-1])
```

- golang slice 사용 (make 로 미리 할당하나, append로 하나씩 넣어주나 시간복잡도 비슷한 것으로 추정)
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
- gota dataframe 만들기
- mariaDB 접속 및 삽입
- 정규표현식으로 문자열 수정
- strings.ReplaceAll 사용
- request 로 소스불러오기
- gota 데이터선택(Series 선택) 
- gota 데이터변경(Series Map 사용)
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

 
