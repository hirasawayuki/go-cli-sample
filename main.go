package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
)

// API ドキュメント
// http://zipcloud.ibsnet.co.jp/doc/api
const APIURL = "https://zipcloud.ibsnet.co.jp/api/search"

// コマンドラインオプションを定義
var zipcode string

func init() {
	// 〇〇Var 関数の引数は以下の通り
	// 第一引数はコマンドラインオプションで定義した変数のポインタ
	// 第二引数はフラグ名
	// 第三引数はデフォルト値
	// 第四引数はフラグの説明
	flag.StringVar(&zipcode, "z", "", "郵便番号")
}

func main() {
	flag.Parse()

	// ※1 ここから先は皆さんの作りたいものに合わせて実装してください
	getAddress(zipcode)
}

type SearchResponse struct {
	Status  int
	Message string
	Results []Result
}

type Result struct {
	Zipcode  string
	Prefcode string
	Address1 string
	Address2 string
	Address3 string
	kana1    string
	kana2    string
	kana3    string
}

// curl "https://zipcloud.ibsnet.co.jp/api/search" -d "zipcode={zipcode}" と同じ
func getAddress(code string) {
	// http.Request 構造体を生成
	req, err := http.NewRequest(http.MethodGet, APIURL, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// クエリパラメータを追加
	// https://zipcloud.ibsnet.co.jp/api/search?zipcode={zipcode}
	q := req.URL.Query()
	q.Add("zipcode", zipcode)
	req.URL.RawQuery = q.Encode()

	// HTTP クライアントを生成
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	s := SearchResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(s.Results) == 0 {
		fmt.Println("No address matching zip code found.")
		os.Exit(1)
	}

	for _, r := range s.Results {
		fmt.Printf("%s%s%s\n", r.Address1, r.Address2, r.Address3)
	}
}
