package main

import (
	"encoding/json"
	"fmt"
	"github.com/rs/cors"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	mecab "github.com/shogo82148/go-mecab"
	"net/http"
	"strings"
	"os/exec"
)

type Morpheme struct {
	Surface            string //表層形
	PartOfSpeech       string //品詞
	Subclassification1 string //品詞細分類1
	Subclassification2 string //品詞細分類2
	Subclassification3 string //品詞細分類3
	InflectedForms     string //活用形
	UtilizationType    string //活用型
	OriginalForm       string //原型
	Syllabary          string //読み仮名
}

var (
	arg string

)

func apiRequest(c web.C, w http.ResponseWriter, r *http.Request) {
	var results []Morpheme

	sentence := c.URLParams["sentence"]
	tagger, err := mecab.New(map[string]string{"dicdir": arg})
	defer tagger.Destroy()

	var m Morpheme

	if err != nil {
		fmt.Println(err.Error())

	}

	tagger.Parse("")
	node, err := tagger.ParseToNode(sentence)

	if err != nil {
		fmt.Println(err.Error())

	}

	for ; !node.IsZero(); node = node.Next() {

		surface := node.Surface()
		// fmt.Println(surface)
		if surface == "" {
			continue
		}
		cols := strings.Split(node.Feature(), ",")

		m = Morpheme{
			Surface:            surface,
			PartOfSpeech:       tryGet(cols, 0),
			Subclassification1: tryGet(cols, 1),
			Subclassification2: tryGet(cols, 2),
			Subclassification3: tryGet(cols, 3),
			InflectedForms:     tryGet(cols, 4),
			UtilizationType:    tryGet(cols, 5),
			OriginalForm:       tryGet(cols, 6),
			Syllabary:          tryGet(cols, 7),
		}
		results = append(results, m)


	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(results)

}

func tryGet(cols []string, index int) string {

	if len(cols)-1 <= index {
		return ""
	}
	return cols[index]

}
func convNewline(str, nlcode string) string {
	    return strings.NewReplacer(
		            "\r\n", nlcode,
			            "\r", nlcode,
				            "\n", nlcode,
					        ).Replace(str)
					}
func main() {

	out, _ := exec.Command("mecab-config", "--dicdir").Output()
	arg = convNewline(string(out), "") + "/mecab-ipadic-neologd"
	fmt.Println(arg)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})
	goji.Use(c.Handler)
	goji.Get("/api/:sentence", apiRequest)
	goji.Serve()
}
