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
```pageURL := fmt.Sprintf("%s&page=%d", url, pageNum)
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
		defer res2.Body.Close()```


 
