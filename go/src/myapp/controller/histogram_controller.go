package controller

import (
	"myapp/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHistogramAll(c *gin.Context) {
	// BestDataを取得
	histogramData, err := database.FindAllHistogram()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data"})
		return
	}

	// 成功した場合、データをJSONで返す
	c.JSON(http.StatusOK, gin.H{"data": histogramData})
}
