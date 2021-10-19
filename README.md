# crawler
정기적으로 데이터를 수집하는 코드들의 집합

## 1. 구조
- 개요
    - 최상위 디렉토리 하위에 위치한 디렉토리마다 수집하는 데이터 성격이 다르다.
    - 따라서 위 디렉토리들은 서로 다른 mod.go를 갖게 된다.
- 최상위 디렉토리 하위에 위치한 디렉토리 설명
    - kamis_crawling
        - ATKamis에서 "일별 품목별 도·소매가격정보"를 수집한다.
        - ATKamis 링크
            - https://www.kamis.or.kr/customer/reference/openapi_list.do
        - 주 코드 작성자
            - 개발자 한무룡
    - media_crawling
        - 미디어 데이터를 수집한다.
        - 수집 가능 데이터
            - 네이버 뉴스
            - 네이버 블로그
            - 네이버 카페
            - 유튜브
        - 수집 중인 데이터
            - 네이버 뉴스
            - 네이버 블로그
            - 유튜브  
        - 주 코드 작성자
            - 개발자 최완재

## 2. kamis_crawling  

## 3. media_crawling  
- 빠른 시작
    - media_crawling 디렉토리 안으로 이동
        - ```
            $ cd media_crawling
            ```
        - 앞으로 나오는 모든 명령은 media_crawling 안에서 실행한다고 가정한다.

    - main 파일 설명
        - main/blog_news.go
            - 네이버 블로그 및 뉴스를 저장한다.
            - 서버에서 정제 후 운영 DB로 바로 삽입
            - 연관 테이블
                - 검색어
                    - operation/AGRIPA/AGRI_QUERY
                - 블로그
                    - operation/AGRIPA/AGRI_BLOG
                - 뉴스
                    - operation/AGRIPA/AGRI_NEWS
                    - operation/AGRIPA/AGRI_PRESS
        - main/youtube_by_channel.go
            - 지정된 요리 채널의 유튜브 영상 정보를 수집
            - 서버에서 정제 후 운영 DB로 바로 삽입
            - 연관 테이블
                - 검색어
                    - operation/AGRIPA/AGRI_QUERY
                - 유튜브
                    - operation/AGRIPA/AGRI_CHANNEL
                    - operation/AGRIPA/AGRI_YOUTUBE
        - main/youtube_by_query.go
            - 서버에서 정제 후 운영 DB로 바로 삽입
            - AGRI_QUERY에서 특정 검색어를 뽑아 유튜브 영상 정보를 수집
            - 인자로 부류코드번호의 첫번째 숫자를 받을 수 있다(1~6)
            - 해당 부류의 품목에 해당하는 검색어로 검색된 유튜브 영상 정보만 수집한다.
            - 연관 테이블
                - 검색어
                    - operation/AGRIPA/AGRI_QUERY
                - 유튜브
                    - operation/AGRIPA/AGRI_YOUTUBE
        - main/code.go
            - 각종 코드를 최신화한다.
            - 수집 DB에서 수집 및 정제 후 운영 DB로 옮김
            - 대상 코드
                - 가락시장코드
                - 도매시장코드
                - 도매시장법인코드
                - 표준등급코드
                - 표준단위코드
                - 표준산지코드
                - 표준품종코드
            - 실제로 운영에 쓰이는 코드
                - 표준품종코드
                - 표준단위코드
            - 연관 테이블
                - 가락시장코드
                    - collection/AGRIPA_COLLECTION/GARAK_CODE
                - 도매시장코드
                    - collection/AGRIPA_COLLECTION/WHOLESALE_MARKET_CODE
                - 도매시장법인코드
                    - collection/AGRIPA_COLLECTION/WHOLESALE_MARKET_CO_CODE
                - 표준등급코드
                    - collection/AGRIPA_COLLECTION/STD_GRADE_CODE
                - 표준단위코드
                    - collection/AGRIPA_COLLECTION/STD_UNIT_CODE
                - 표준산지코드
                    - collection/AGRIPA_COLLECTION/PLACE_ORIGIN_CODE
                - 표준품종코드
                    - collection/AGRIPA_COLLECTION/STD_SPECIES_CODE
                - 표준품목코드
                    - collection/AGRIPA_COLLECTION/STD_ITEM_CODE
                    - operation/AGRIPA/STD_ITEM_CODE
                - 품목 키워드
                    - collection/AGRIPA_COLLECTION/STD_ITEM_KEYWORD
                    - operation/AGRIPA/STD_ITEM_KEYWORD
                - 품목 매핑
                    - collection/AGRIPA_COLLECTION/ITEM_MAPPING
                    - operation/AGRIPA/ITEM_MAPPING
        - main/mafra_adj_auction_stats.go
            - mafra API를 사용하여 정산 경락가격 요약 정보를 수집 및 정제한다.
            - 수집 DB에서 수집 및 정제 후 운영 DB로 옮김
            - "정산 경락가격 요약 정보 원본" 테이블로 데이터 수집 후  
              "정산 경락가격 요약 정보" 테이블에 조사가격 품목코드가 존재하는 rows만 추출한 후  
              (collection/AGRIPA_COLLECTION/MAP_STD_EXAM_ITEM 테이블을 참조하여 추출(파일로부터 생성한 테이블))  
              "정산 경락 거래량" 테이블에 날짜별 거래량이 많은 순으로 TOP 3에 해당하는 품종의 거래량만 수집한다.
            - 연관 테이블
                - 정산 경락가격 요약 정보 원본(연도 별로 테이블명을 바꿔가며 수집)
                    - collection/AGRIPA_COLLECTION/MAFRA_ADJ_AUCTION_STATS_${yyyy}
                - 정산 경락가격 요약 정보
                    - collection/AGRIPA_COLLECTION/MAFRA_ADJ_AUCTION_STATS
                - 정산 경락 거래량
                    - collection/AGRIPA_COLLECTION/MAFRA_ADJ_AUCTION_QUANTITY
                    - operation/AGRIPA/MAFRA_ADJ_AUCTION_QUANTITY
        - main/mafra_examination.go
            - mafra API를 사용하여 조사가격 정보를 수집 및 정제한다.
            - 수집 DB에서 수집 및 정제 후 운영 DB로 옮김
            - "조사가격 정보 원본" 테이블로 데이터 수집 후  
              "도매시장 가격"과 "소매시장 가격" 테이블에 각각 삽입한다.
            - 연관 테이블
                - "조사가격 정보 원본"(연도 별로 테이블명을 바꿔가며 수집)
                    - collection/AGRIPA_COLLECTION/MAFRA_EXAMINATION_${yyyy}
                - 도매시장 가격
                    - collection/AGRIPA_COLLECTION/MAFRA_WHOLE_PRICE
                    - operation/AGRIPA/MAFRA_WHOLE_PRICE
                - 소매시장 가격
                    - collection/AGRIPA_COLLECTION/MAFRA_RETAIL_PRICE
                    - operation/AGRIPA/MAFRA_RETAIL_PRICE
        - main/trade.go
            - 수출입 중량(weight) 및 금액(amount) 정보를 수집 및 정제
            - 수집 DB에서 수집 및 정제 후 운영 DB로 옮김                  
            - 수출입 정보를 파일로 제공되기 때문에 자동화되지 않아 서버에서 사용하지 않는다.  
            - 연관 테이블
                - 수입
                    - collection/AGRIPA_COLLECTION/AGRI_IMPORTATION
                - 수출
                    - collection/AGRIPA_COLLECTION/AGRI_EXPORTATION
                - 수출입
                    - collection/AGRIPA_COLLECTION/TRADE
                    - operation/AGRIPA/TRADE
        - main/weather.go
            - 산지 날씨 정보 수집
            - 수집 DB에만 저장
            - 이 데이터에서 사용하고 있는 산지 및 품목 코드가 표준산지 및 표준품목코드가 아니라서 사용 보류 중이다.
            - 현재 사용하지 않고 수집만 하는 중
            - PK 잡기가 어려워 현재 PK 설정이 되어 있지 않다.
            - 그래서 과거 데이터를 분리해두고, 새로운 데이터만 AGRI_WEATHER 테이블에 저장한다.
            - 연관 테이블
                - 날씨
                    - collection/AGRIPA_COLLECTION/AGRI_WEATHER
                    - collection/AGRIPA_COLLECTION/AGRI_WEATHER_HISTORICAL
        - (미사용)main/dataGoKr_adj_auction.go
            - DataGoKr API를 사용하여 "정산 경락가격 요약 정보"를 수집 및 정제
            - 산지 코드가 존재하지 않아 사용이 유보되었다.
            - 현재는 수집하지 않고 있다.
            - 연관 테이블
                - 정산 경락가격 요약 정보 원본(연도 별로 테이블명을 바꿔가며 수집)
                    - collection/AGRIPA_COLLECTION/AGRI_ADJ_AUCTION_${yyyy}
                - 실시간 경락가격 요략 정보 원본(연도 별로 테이블명을 바꿔가며 수집)
                    - collection/AGRIPA_COLLECTION/AGRI_REAL_AUCTION_${yyyy}
        - (미사용)main/breaking_auction.go
            - 실시간 경락 거래 속보 데이터 수집 및 정제
            - 우선순위에 의해 기획에서 후순위로 밀려 미사용 중이다.
            - 연관 테이블
                - 거래 속보
                    - collection/AGRIPA_COLLECTION/BREAKING_NEWS
        - (미사용)main/remove_old_breaking_auction.go
            - 오래된 실시간 경락 거래 속보 삭제
            - 거래 속보는 휘발성이기 때문에 오래된 행을 삭제해준다.
            - 마찬가지로 우선순위에 의해 기획에서 후순위로 밀려 미사용 중이다.
            - 연관 테이블
                - 거래 속보
                    - collection/AGRIPA_COLLECTION/BREAKING_NEWS
    - github에서 pull 및 빌드 파일 생성
        - 설명
            - github에서 코드가 갱신되어 pull로 코드를 최신화 시켜도 빌드 파일이 최신화되지는 않는다.  
            - 이를 막기 위해 git_pull.sh를 실행하여 코드 및 빌드파일 최신화를 한다.
        - 명령어
            - ```
                $ bash git_pull.sh
                ```
    - main 바이너리 파일 실행
        - 설명
            - periodically_run.sh 를 실행시켜 빌드 및 해당 main 파일을 실행할 수 있다.
            - periodically_run.sh는 두 개의 파라미터를 받는다.
                - 첫 번째 파라미터
                    - 실행하고자 하는 main 디렉토리 하위의 package main인 go 파일명
                - 두 번째 파라미터(생략가능)
                    - 첫 번째 파라미터의 go 파일에 들어가는 input 파라미터
        - 명령어
            - ```
                # main/code.go 실행
                $ bash periodically_run.sh code

                # main/youtube_by_query.go 실행
                $ bash periodically_run.sh youtube_by_query 1
                ```
    - crontab에 등록
        - 설명
            - build 파일을 정기적으로 실행하게 하려면 crontab에 등록을 해야한다.
            - periodically_run.sh 통해 build 파일을 실행시키면  
              실행 시 날짜와 함께 간략한 문자열을 남겨 제대로 실행되었는지 확인할 수 있다.
        - crontab 등록 명령어
            - ```
                $ crontab -e
                ```
        - crontab 등록 예시
            - ```
                #media youtube-garakdong365
                7 */2 * * * bash /home/ubuntu/go/src/crawler/media_crawling/periodically_run.sh main_every_2_hours

                #media news, blog
                5 * * * * bash /home/ubuntu/go/src/crawler/media_crawling/periodically_run.sh main_every_hour

                #media youtube
                8 8 * * * bash /home/ubuntu/go/src/crawler/media_crawling/periodically_run.sh main_every_day
                ```
            - #media youtube-garakdong365 부분은 2시간마다 7분에 수집이 시작된다.
            - #media news, blog은 1시간마다 5분에 수집이 시작된다.
            - #media youtube는 매일 8시 8분에 수집이 시작된다.
    - 로그 확인
        - 설명
            - 프로그램이 한 번이라도 실행되었다면 main 디렉토리 하위에  
              main 파일과 동일한 이름의 .log 확장자 파일이 있을 것이다.
            - 해당 파일을 확인하면 해당 main 파일이 언제 실행되었는지 확인이 가능하다.
            - 로그는 periodically_run.sh로 실행 했을 경우만 확인할 수 있다.
            - 에러 체크를 위한 시스템 로그는 logs/{실행일}/{실행 main 파일 명}.txt로 저장된다.
            - 시스템 로그 자동 제거 기능은 없기 때문에 적당히 오래된 건 지워준다.
        - 명령어
            - ```
                ## 2020.10.01일의 main/code.go의 로그 확인
                # 실행 로그 확인
                $ cat main/code.log

                # 시스템 로그 확인
                $ cat logs/20201001/code.txt
                ```
-  Windows 에서 DB Test를 위한 Docker 개발환경 구축
    - Docker 설치 및 실행
        - 개발 pc의 운영체제 맞는 Docker를 설치 및 실행
    - Docker image 생성
        - 터미널의 현재 위치는 crawler 디렉토리이며, 이미지 이름은 crawler라고 가정
        - ```
            # 윈도우
            docker build -t crawler %cd%

            # 리눅스 및 맥os
            docker build -t crawler .
            ```
    - Docker container 생성 및 bash shell 실행
        - 터미널의 현재 위치는 crawler 디렉토리라고 가정
        - container의 8080 포트를 개발 pc의 8080 포트와 연결
        - container는 종료 시 자동 삭제
        - 현재 개발 pc 터미널 위치와 container의 /crawler 위치를 동기화
        - ```
            # 윈도우
            docker run -it --rm -p 8080:8080 -v %cd%/:/crawler crawler /bin/bash 

            # 리눅스 및 맥os
            docker run -it --rm -p 8080:8080 -v ./:/crawler crawler /bin/bash 
            ```
