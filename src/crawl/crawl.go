package crawl

import (
	"bytes"
	"fmt"
	_ "github.com/lib/pq"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Acao struct {
	Posicao 	int
	Papel 		string
	Empresa    	string
	Oscil_dia 	string
	Valor_merc 	float64
}

type Parser interface {
	ParsePage(*goquery.Document) Acao
}

func getRequest(url string) (*http.Response, error) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func extractLinks(doc *goquery.Document) []string {
	foundUrls := []string{}
	r1, _ := regexp.Compile("detalhes\\.php\\?papel")
	r2, _ := regexp.Compile("\\?papel.*")
	if doc != nil {
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			res, _ := s.Attr("href")
			if r1.FindString(res) != "" {
				foundUrls = append(foundUrls, r2.FindString(res))
			}

		})
		return foundUrls
	}
	return foundUrls
}

func convrtToUTF8(str string, origEncoding string) string {
	strBytes := []byte(str)
	byteReader := bytes.NewReader(strBytes)
	reader, _ := charset.NewReaderLabel(origEncoding, byteReader)
	strBytes, _ = ioutil.ReadAll(reader)
	return string(strBytes)
}

func extractData(doc *goquery.Document) Acao {
	foundData := Acao{}

	papelRegex, _ := regexp.Compile(".*Papel(.*?)\\?")
	empresaRegex, _ := regexp.Compile(".*Empresa(.*?)\\?")
	oscilDiaRegex, _ := regexp.Compile(".*Dia(.*?)\\?")
	valorMercRegex, _ := regexp.Compile(".*Valor de mercado(.*?)\\?")

	if doc != nil {
		d := doc.Find("span")
		text := d.Text()
		if text != "" {
			//Papel
			textp := papelRegex.FindStringSubmatch(text)[1]
			if textp != "" {
				foundData.Papel = textp

				//Empresa
				texte := empresaRegex.FindStringSubmatch(text)[1]
				foundData.Empresa = convrtToUTF8(texte, "latin1")

				//Oscilacao diaria
				texto := oscilDiaRegex.FindStringSubmatch(text)[1]
				foundData.Oscil_dia = texto

				//Valor de mercado
				textv := valorMercRegex.FindStringSubmatch(text)[1]
				textvf := strings.Replace(textv, ".", "",-1)

				floatv, _ := strconv.ParseFloat(textvf, 64)
				foundData.Valor_merc = floatv

				return foundData
			}
		}
		return foundData
	}
	return foundData

}

func sortResults (results []Acao) []Acao{
	resultsSorted := []Acao{}
	resultsSorted = results

	sort.SliceStable(resultsSorted, func(i, j int) bool {
		return resultsSorted[i].Valor_merc > resultsSorted[j].Valor_merc
	})

	return resultsSorted
}

func crawlMainPage(targetURL string) ([]string) {
	fmt.Println("Requesting: ", targetURL)
	resp1, _ := getRequest(targetURL)

	doc1, _ := goquery.NewDocumentFromReader(resp1.Body)
	links := extractLinks(doc1)

	foundUrls := links
	return foundUrls
}

func crawlSecondPage(baseURL string, targetURL string) (Acao) {
	foundData := Acao{}
	url := baseURL + targetURL
	fmt.Println("Requesting: ", url)
	resp2, _ := getRequest(url)

	doc2, _ := goquery.NewDocumentFromReader(resp2.Body)
	data := extractData(doc2)

	if data.Papel != "" {
		foundData = data
	}

	return  foundData
}

func parseStartURL(u string) string {
	parsed, _ := url.Parse(u)
	return fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)
}

func Crawl(startURL string) []Acao {
	results := []Acao{}
	ListURL := []string{}

	var d = 0

	seen := make(map[string]bool)

	foundLinks := crawlMainPage(startURL)

	if foundLinks != nil {
		ListURL = foundLinks
	}

	for _, link := range ListURL {
		//link := ListURL[2]
		if !seen[link] {
			seen[link] = true

			foundData := crawlSecondPage(startURL, link)
			if foundData.Papel != "" {
				results = append(results, foundData)
			}
		}

	}

	resultsSort := sortResults(results)

	results10 := []Acao{}

	for d<10 {
		resultsSort[d].Posicao = d+1
		results10 = append(results10,resultsSort[d])
		d++
	}

	return results10
}

