package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
)
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

func getAlbaList(url string) {
	res, err := http.Get(url)
	respCheck(res)
	errCheck(err)
	html, err := goquery.NewDocumentFromReader(res.Body) //브랜드 사이트
	errCheck(err)

	html.Find("tbody").Find("tr").Each(func(i int, s *goquery.Selection) {
		className, _ := s.Attr("class")
		if className != "summaryView" {
			s.Find("td").Each(func(j int, sel* goquery.Selection){
				out, _ := iconv.ConvertString(sel.Text(), "euc-kr", "utf-8")
				fmt.Println(out)
			})
		}
	})

}


//메인에서는 서버를 열어 home에 원하는 직업을 입력해 그 결과를 csv파일로 얻을 수 있도록 해준다.
func main(){
	resp, err := http.Get("http://www.alba.co.kr")
	
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
			getAlbaList(brandUrl) //알바리스트 배열 반환
		}
	})
		
}