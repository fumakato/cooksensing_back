package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// SendPOSTは指定されたURLに対してPOSTリクエストを送信し、feature_dataを返します
func SendPOST(downloaddata string) (string, error) {
	url := "http://127.0.0.1:5001/feature_extraction"

	// POSTリクエストのボディとして送信するデータをJSON形式にエンコード
	postData := map[string]string{"url": downloaddata}
	jsonData, err := json.Marshal(postData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// POSTリクエストを作成
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// クライアントを作成してリクエストを送信
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// ステータスコードが200以外の場合はエラーとする
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// レスポンスボディを読み取る
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// feature_dataとして返す
	return string(body), nil
}

func main() {
	// SendPOST関数の使用例
	downloaddata := "https://minio.kajilab.dev/cucumber-slices/suzaki2-acc.csv?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=hRjq2yhc1WqPrfEV%2F20240729%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20240729T102332Z&X-Amz-Expires=900&X-Amz-SignedHeaders=host&X-Amz-Signature=c8855857b5582dd3915a8ce4fcead006a2d9252213769c78eea57812991741ef"
	// downloaddata := "aa"
	featureData, err := SendPOST(downloaddata)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Received feature data: %s\n", featureData)
}
