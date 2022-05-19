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

var ftagPath string
var ftagNum int
var fisRandom int
var foptc string
var foptn string
var foptt string

func main() {

	//flag定義
	var (
		tagPath = flag.String("p", "./searchtag.txt", "計測に使用するタグ名が記載されているファイルのPathを指定")
		tagNum = flag.Int("s", 100, "計測に使用するタグ数を指定．-1の場合はファイルに記載されているすべてタグを使用")
		isRandom = flag.Int("r", 1, "1の場合のみ計測に使用するタグ名をランダムな順番で選択")
		optc = flag.Int("c", 5, "テストで同時に発行するリクエストの数を数値で指定")
		optn = flag.Int("n", 10, "テストで発行するリクエストの回数を数値で指定")
		optt = flag.Int("t", 2, "1リクエストのタイムアウト時間を秒単位で指定")
	)

	flag.Parse()
	
	ftagPath = *tagPath
	ftagNum = *tagNum
	fisRandom = *isRandom
	foptc = strconv.Itoa(*optc)
	foptn = strconv.Itoa(*optn)
	foptt = strconv.Itoa(*optt)

	//c > n は禁止
	if *optc > *optn {
		log.Println("<Debug> -c must be smaller than -n")
		return;
	}

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
	tmpl := template.Must(template.ParseFiles("./web/html/preindex.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Println("<Debug> can't open ./web/html/preindex.htm : ", err)
	}
}

// ajax戻り値のJSON用構造体
type measureParam struct {
	Time string
	Msg  string
}

//フォームからの入力を処理 index.jsから受け取る
func measureHandler(w http.ResponseWriter, r *http.Request) {

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
	ret.Msg, ret.Time = ab.Ab(guid.String(), url, ftagPath, ftagNum, fisRandom, foptc, foptn, foptt)


	// 構造体をJSON文字列化する
	jsonBytes, _ := json.Marshal(ret)
	// index.jsに返す
	fmt.Fprint(w, string(jsonBytes))
}



