package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
)

type jobData struct {
	location string
	company string
	time string
	pay string
	regDate string
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

// getAlbaList는 한 직업에 관한 모든 매장리스트들을 반환한다. albaList를 ch 채널로 전송함.
func getAlbaList(url string, ch chan<- []jobData)  {
	albaList := []jobData{}
	res, err := http.Get(url)
	respCheck(res)
	errCheck(err)
	html, err := goquery.NewDocumentFromReader(res.Body) //브랜드 사이트
	errCheck(err)

	html.Find("tbody").Find("tr").Each(func(i int, s *goquery.Selection) {
		className, _ := s.Attr("class")
	if className != "summaryView" {
			var arr [5]string
			arr[0], _ = iconv.ConvertString(s.Find(".local").Text(), "euc-kr", "utf-8")
			arr[1], _ = iconv.ConvertString(s.Find(".company").Text(), "euc-kr", "utf-8")
			arr[2], _ = iconv.ConvertString(s.Find(".data").Text(), "euc-kr", "utf-8")
			arr[3], _ = iconv.ConvertString(s.Find(".pay").Text(), "euc-kr", "utf-8")
			arr[4], _ = iconv.ConvertString(s.Find(".regDate").Text(), "euc-kr", "utf-8")
			li := jobData{}
			li.location=arr[0]
			li.company=arr[1]
			li.time=arr[2]
			li.pay=arr[3]
			li.regDate=arr[4]
			if len(li.location) != 0 {albaList = append(albaList, li)}
		}
	})

	ch <- albaList
}

func setWriting(job jobData, wch chan<- []string){
	jobslice := []string{job.location, job.company, job.time, job.pay, job.regDate}
	wch <- jobslice
}

//메인에서는 서버를 열어 home에 원하는 직업을 입력해 그 결과를 csv파일로 얻을 수 있도록 해준다.
func main(){
	totalJobs := 0
	ch := make(chan []jobData)
	var jobs []jobData
	resp, err := http.Get("http://www.alba.co.kr")
	// var arr []jobData
	respCheck(resp)
	errCheck(err)
	defer resp.Body.Close()
	//html load
	doc, err := goquery.NewDocumentFromReader(resp.Body) // 사이트 html 얻기
	errCheck(err)
	//슈퍼브랜드 채용정보 리스트를 반복문으로 돌아 각 브랜드 점포리스트들을 배열로 얻어 csv로 저장한다.
	doc.Find("div#MainSuperBrand").Find("ul.goodsBox").Find("li").Each(func(i int,s *goquery.Selection){
		brandUrl, ok := s.Find("a.goodsBox-info").Attr("href") //슈퍼브랜드 채용주소 접근 성공
		if ok {
			go getAlbaList(brandUrl, ch) //알바리스트 배열 반환
			totalJobs+=1
			// for i := 0; i < len(arr); i++ {
			// 	fmt.Println(arr[i])
			// }
			// fmt.Println("다음 직업으로 넘어갑니다.")
		}
	})

	for i := 0; i < totalJobs; i++ {
		jobs = append(jobs, <-ch...)
	}
	
	wch := make(chan []string)
	file, err := os.Create("Jobs.csv")
	errCheck(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string {"location", "company", "Data", "Pay", ""}

	wErr := w.Write(headers)
	errCheck(wErr)

	for i := 0; i < len(jobs); i++ {
		go setWriting(jobs[i], wch)
	}

	for i := 0; i < len(jobs); i++ {
		wErr = w.Write(<-wch)
		errCheck(wErr)
	}
}