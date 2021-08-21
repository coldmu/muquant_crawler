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
 - 종목코드를 이용하여 전체 크롤링
 - 파이썬을 이용하여 크롤링하였으나, 너무 오래걸려서 학습겸 GO로 변경 잘될것으로 기대..

<대체>
 - beautifulsoup : github.com/anaskhan96/soup
 - pandas : github.com/go-gota/gota/dataframe, github.com/go-gota/gota/series
 
