package main

import (
	"encoding/json"
	"fmt"
	"log"
	"flag"
	"strconv"
	"net/http"
	"text/template"
	"benchmarkserver/internal/ab"
	"github.com/rs/xid"
)

func main() {

	// webフォルダにアクセスできるようにする
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./web/css/"))))
	http.Handle("/script/", http.StripPrefix("/script/", http.FileServer(http.Dir("./web/script/"))))
	http.Handle("/gif/", http.StripPrefix("/gif/", http.FileServer(http.Dir("./web/gif/"))))

	//ルーティング設定 "/"というアクセスがきたら rootHandlerを呼び出す
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/measure", measureHandler)

	log.Println("Listening...")
	// 3000ポートでサーバーを立ち上げる
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Println("<Debug> http.LinstenAndServe(:3000) : ", err)
	}
}

//main画面
func rootHandler(w http.ResponseWriter, r *http.Request) {
	//index.htmlを表示させる
	tmpl := template.Must(template.ParseFiles("./web/html/index.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Println("<Debug> can't open ./web/html/index.htm : ", err)
	}
}

// ajax戻り値のJSON用構造体
type measureParam struct {
	Time string
	Msg  string
}

//フォームからの入力を処理 index.jsから受け取る
func measureHandler(w http.ResponseWriter, r *http.Request) {

	//flag定義
	var (
		tagPath = flag.String("p", "./data/searchtag.txt", "検索タグのファイルパス")
		tagNum = flag.Int("s", -1, "検索数:-1の場合は全検索")
		isRandom = flag.Int("r", 1, "1の場合はランダムに検索する")
		optc = flag.Int("c", 1, "abコマンドの-c")
		optn = flag.Int("n", 1, "abコマンドの-n")
		optt = flag.Int("t", 1, "abコマンドの-t")
	)

	flag.Parse()


	//index.jsに返すJSONデータ変数
	var ret measureParam
	//POSTデータのフォームを解析
	err := r.ParseForm()
	if err != nil {
		log.Println("<Debug> r.ParseForm : ", err)
	}

	url := r.Form["url"][0]

	//idを設定(logを対応づけるため)
	guid := xid.New()
	log.Println("<Info> request URL: " + url + ", id: " + guid.String())

	//abコマンドで負荷をかける．計測時間を返す
	ret.Msg, ret.Time = ab.Ab(guid.String(), url, *tagPath, *tagNum, *isRandom, strconv.Itoa(*optc), strconv.Itoa(*optn), strconv.Itoa(*optt))


	// 構造体をJSON文字列化する
	jsonBytes, _ := json.Marshal(ret)
	// index.jsに返す
	fmt.Fprint(w, string(jsonBytes))
}



