package controller

import (
	"myapp/database"
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
)

func GetBestAverage(c *gin.Context) {
	// var req ReqUserID

	averagePace, accelerationStdDev, err := database.AveragePaceAndAccelerationStdDev()
	if err != nil {
		log.Printf("Error calculating averages: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to calculate averages",
		})
		return

	}
	// 結果をJSONで返す
	c.JSON(http.StatusOK, gin.H{
		"average_pace":                    averagePace,
		"acceleration_standard_deviation": accelerationStdDev,
	})
}

func GetBestAll(c *gin.Context) {
	// BestDataを取得
	bestData, err := database.FindAllBestData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data"})
		return
	}

	// 成功した場合、データをJSONで返す
	c.JSON(http.StatusOK, gin.H{"data": bestData})
}
