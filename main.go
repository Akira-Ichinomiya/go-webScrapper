package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

//getMaxPageNumber은 검색한 해당 사이트의 최대번쨰 페이지 값을 int로 반환해준다.
func getMaxPageNumber(){
	
}

//getJobsPage는 페이지에 존재하는 직업을 리스트로 가져와 준다. 각 원소들은 title, company, location, info를 string값으로 가지고 있는 구조체이다.
func getJobsPage(){

}


//getAllJobs는 주어진 getMaxPageNumber를 이용해 반복문을 거쳐 각페이지에서 얻게되는 모든 직업리스트들을 한 리스트에 합한다.
func getAllJobs(){

}

func errCheck(err error) { // 에러 검사
	if err != nil {
		log.Fatal(err)
		fmt.Println("에러 발생")
	}
}

func respCheck(resp *http.Response) {
	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
		fmt.Println("에러 발생")
	}
}

//메인에서는 서버를 열어 home에 원하는 직업을 입력해 그 결과를 csv파일로 얻을 수 있도록 해준다.
func main(){
	resp, err := http.Get("https://kr.indeed.com/jobs?q=python&limit=50&radius=25&start=0")
	
	respCheck(resp)
	errCheck(err)
	defer resp.Body.Close()
	//html load
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	errCheck(err)
	lists := doc.Find(".jobsearch-SerpJobCard")
	fmt.Println(lists.Length())

	lists.Each(func(idx int, sel *goquery.Selection) {
		title := sel.Find(".title>a").Text()
		fmt.Println(title)
	})

}