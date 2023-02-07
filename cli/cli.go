package cli

import (
	"sync"
	"net/http"
	"log"
	"github.com/Jeffail/gabs"
)

type RequestBody struct{
	SourceLang string
	TargetLang string
	SourceText string
}

const translateUrl = "https://translate.googleapis.com/translate_a/single"

func RequestTranslate(body *RequestBody, str chan string, wg *sync.WaitGroup){
	client := &http.Client{}
	req, err := http.NewRequest("GET", translateUrl, nil)

	query := req.URL.Query()

	query.Add("client", "gtx")
	query.Add("sl", body.SourceLang)
	query.Add("tl", body.TargetLang)
	query.Add("dt", "t")
	query.Add("q", body.SourceText)

	req.URL.RawQuery = query.Encode()

	if err != nil{
		log.Fatal("1: There was a prblem:%s", err)
	}

	res, err := client.Do(req)
	if err != nil{
		log.Fatal("2: There was a problem: %s", err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusTooManyRequests{
		str <- " You ghave be rate limited, try again"
		wg.Done()
		return
	}

	parsedJson, err := gabs.ParseJSONBUFFER(res.Body)

	if err != nil {
		log.Fatalf("3: There was a problem %s", err)
	}

	nestedOne, err := parsedJson.ArrayElement(0)

	if err != nil  {
		log.Fatalf("4: There was a problem %v", err)
	}
	nestedTwo, err := nestedOne.ArrayElement(0)
	if err != nil {
		log.Fatalf("5: There was a problem %s", err)
	}
	translatedStr, err := nestedTwo.ArrayElement(0)

	if err != nil{
		log.Fatalf("6: There was a problem %s", err)
	}

	str <- translatedStr.Data().(string)
	wg.Done()
}