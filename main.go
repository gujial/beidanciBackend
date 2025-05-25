package main

import (
	"math/rand"
	"net/http"
	"time"

	"beidanciBackend/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库连接
	dsn := "root:041109@tcp(127.0.0.1:3306)/beidanci?charset=utf8mb4&parseTime=True&loc=Local"
	if err := models.InitDB(dsn); err != nil {
		panic("failed to connect database: " + err.Error())
	}

	rand.Seed(time.Now().UnixNano())

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/word", func(c *gin.Context) {
		table := c.DefaultQuery("level", "CET4") // 你可以让它随机选择 CET4 或 CET6

		// 获取正确的单词和释义
		word, err := models.GetRandomWord(table)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询单词失败"})
			return
		}

		// 获取错误选项
		distractors, err := models.GetRandomDistractors(table, word.Translate, 3)
		if err != nil || len(distractors) < 3 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询释义失败"})
			return
		}

		// 构造选项并打乱
		options := append(distractors, word.Translate)
		rand.Shuffle(len(options), func(i, j int) {
			options[i], options[j] = options[j], options[i]
		})

		// 找到正确释义的下标
		var correctIndex int
		for i, opt := range options {
			if opt == word.Translate {
				correctIndex = i
				break
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"word":         word.Word,
			"options":      options,
			"correctIndex": correctIndex,
		})
	})

	r.Run(":8080")
}
