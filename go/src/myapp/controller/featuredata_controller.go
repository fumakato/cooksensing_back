package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"myapp/database"
	"myapp/model"
	"net/http"
	"strconv"

	"time"

	"log"

	"github.com/gin-gonic/gin"
)

// リクエストボディの構造体を定義
type CreateFeatureRequest struct {
	Fileurl string `json:"fileurl" binding:"required"` // nameフィールドが必須
	Uid     string `json:"uid" binding:"required"`     // valueフィールドが必須
	Date    string `json:"date" binding:"required"`    // dateフィールドが必須
}

// JSONレスポンスの構造体を定義
type Response struct {
	AveAcc  float64 `json:"aveAcc"`
	AvePace float64 `json:"avePace"`
	Stdev   float64 `json:"stdev"`
}

func CreateFeatureDatas(c *gin.Context) {
	var requestBody CreateFeatureRequest

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	// 🔸 Pythonのレスポンスを取得（ステータス付き）
	status, featureDataRaw, err := SendPOST(requestBody.Fileurl)
	if err != nil {
		// Pythonからのエラー内容（JSON）をそのまま返す
		var pyErr map[string]string
		if unmarshalErr := json.Unmarshal([]byte(err.Error()), &pyErr); unmarshalErr == nil {
			c.JSON(status, pyErr)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract feature data", "code": "FEATURE_EXTRACTION_ERROR"})
		}
		return
	}

	// 🔸 正常時のパース
	var response Response
	if err := json.Unmarshal(featureDataRaw, &response); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse feature data: " + err.Error(), "code": "INVALID_FEATURE_DATA_FORMAT"})
		return
	}

	i, err := strconv.Atoi(requestBody.Uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid uid: must be a number", "code": "INVALID_UID"})
		return
	}
	uid := uint(i)

	parsedDate, err := time.Parse("2006-01-02T15:04", requestBody.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid date format: must be YYYY-MM-DDTHH:mm",
			"code":  "INVALID_DATE_FORMAT",
		})
		return
	}
	parsedDate = parsedDate.In(time.UTC) // 保存前にUTC化
	newDate := parsedDate

	feature := model.FeatureData{
		UserID:                        uid,
		ActionID:                      1,
		Date:                          newDate,
		AveragePace:                   float32(response.AvePace),
		AccelerationStandardDeviation: float32(response.Stdev),
	}

	if err := database.AddFeatureData(feature); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save feature data: " + err.Error(), "code": "DB_SAVE_ERROR"})
		return
	}

	// database.UpdateBestDataFromFeatureData()
	// database.GenerateAndStoreHistogramData()

	err = database.UpdateBestDataFromFeatureData()
	if err != nil {
		log.Println("UpdateBestData error:", err)
	}

	err = database.GenerateAndStoreHistogramData()
	if err != nil {
		log.Println("GenerateHistogram error:", err)
	}
	// ⬇ ユーザーID指定で呼び出し
	// err = database.AssignBestClassByUserID(uid)
	err = database.AssignBestClassToAll()
	if err != nil {
		log.Println("AssignBestClass error:", err)
	}

	c.JSON(http.StatusOK, gin.H{"status": feature})
}

func SendPOST(downloaddata string) (int, json.RawMessage, error) {
	url := "http://python_app:5001/feature_extraction"

	postData := map[string]string{"url": downloaddata}
	jsonData, err := json.Marshal(postData)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf(`{"error": "Failed to encode request: %s", "code": "ENCODE_ERROR"}`, err.Error())
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf(`{"error": "Failed to create request: %s", "code": "REQUEST_CREATION_ERROR"}`, err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf(`{"error": "Failed to send request: %s", "code": "SEND_REQUEST_ERROR"}`, err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf(`{"error": "Failed to read response: %s", "code": "READ_RESPONSE_ERROR"}`, err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		// エラーの内容はそのまま文字列にして返す（JSON形式前提）
		return resp.StatusCode, nil, errors.New(string(body))
	}

	return resp.StatusCode, json.RawMessage(body), nil
}

// // 新しいデータを作成
// func CreateFeatureDatas(c *gin.Context) {
// 	fmt.Println("") //改行
// 	fmt.Println("CreateFeatureDatas")
// 	var requestBody CreateFeatureRequest

// 	// リクエストボディを構造体にバインドし、バリデーションを行う
// 	if err := c.ShouldBindJSON(&requestBody); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	fmt.Printf("requestBody: %s\n", requestBody)

// 	featureData, err := SendPOST(requestBody.Fileurl)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}

// 	fmt.Printf("featureData: %s\n", featureData)

// 	// json.RawMessageをResponse構造体にアンマーシャル

// 	var response Response
// 	if err := json.Unmarshal(featureData, &response); err != nil {
// 		fmt.Println("Error unmarshalling response:", err)
// 		return
// 	}

// 	// aveAcc, avePace, stdevの値を取得
// 	aveAcc := response.AveAcc
// 	avePace := response.AvePace
// 	stdev := response.Stdev

// 	// これらの値を以降のコードで使用
// 	fmt.Printf("Average Accuracy: %f\n", aveAcc)
// 	fmt.Printf("Average Pace: %f\n", avePace)
// 	fmt.Printf("Standard Deviation: %f\n", stdev)

// 	// db := database.ConnectDB()
// 	// defer db.Close()

// 	var i int
// 	i, _ = strconv.Atoi(requestBody.Uid)
// 	uid := uint(i)

// 	// 文字列から日付を解析
// 	parsedDate, err := time.Parse("2006-01-02", requestBody.Date)
// 	if err != nil {
// 		fmt.Println("Error parsing date:", err)
// 		return
// 	}

// 	// 解析した日付を使用して新しいtime.Dateを作成
// 	newDate := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 12, 0, 0, 0, time.Local)

// 	var featuredatalast = model.FeatureData{
// 		UserID: uid, ActionID: 1, Date: newDate, AveragePace: float32(avePace), AccelerationStandardDeviation: float32(stdev),
// 	}

// 	// db.Create(&featuredatalast)
// 	database.AddFeatureData(featuredatalast)

// 	// ベストデータの更新
// 	database.UpdateBestDataFromFeatureData()

// 	// ヒストグラムの更新
// 	database.GenerateAndStoreHistogramData()

// 	// 成功レスポンスを受け取ったパラメータと共に返す
// 	c.JSON(http.StatusOK, gin.H{"status": featuredatalast})
// }

// // SendPOSTは指定されたURLに対してPOSTリクエストを送信し、feature_dataを返します
// func SendPOST(downloaddata string) (json.RawMessage, error) {
// 	url := "http://127.0.0.1:5001/feature_extraction"

// 	// POSTリクエストのボディとして送信するデータをJSON形式にエンコード
// 	postData := map[string]string{"url": downloaddata}
// 	jsonData, err := json.Marshal(postData)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
// 	}

// 	// POSTリクエストを作成
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create request: %w", err)
// 	}
// 	req.Header.Set("Content-Type", "application/json")

// 	// クライアントを作成してリクエストを送信
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to send request: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	// ステータスコードが200以外の場合はエラーとする
// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
// 	}

// 	// レスポンスボディを読み取る
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read response body: %w", err)
// 	}

// 	// feature_dataとして返す
// 	return json.RawMessage(body), nil
// }

type Return struct {
	Date        []string  `json:"date"`
	AveragePace []float32 `json:"average_pace"`
}

type ReqUserID struct {
	UserID uint `json:"user_id" binding:"required"`
}

func GetFeatureDatasByUserID(c *gin.Context) {
	var req ReqUserID

	// リクエストボディをパースしてUserIDを取得
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// FeatureDataを取得
	featureData, err := database.GetFeatureDataByUserID(req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data"})
		return
	}

	// 成功した場合、データをJSONで返す
	c.JSON(http.StatusOK, gin.H{"data": featureData})
}

type ReqUserDays struct {
	UserID uint `json:"user_id" binding:"required"`
	Days   int  `json:"days" `
}

// GetFeatureDatasByUserID はリクエストを処理してデータを取得
func GetFeatureDatasByUserIDWithinDays(c *gin.Context) {
	var req ReqUserDays
	// リクエストボディをパース
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// FeatureDataを取得
	featureData, err := database.GetFeatureDataByUserIDWithinDays(req.UserID, req.Days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data"})
		return
	}

	// 成功した場合、データをJSONで返す
	c.JSON(http.StatusOK, gin.H{"data": featureData})
}
