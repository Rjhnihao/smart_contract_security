package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8081"}
	r.Use(cors.New(config))

	r.POST("/api/upload", func(c *gin.Context) {

		file, fileErr := c.FormFile("file")

		if fileErr == nil {
			f, err := file.Open()
			defer f.Close()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取文件"})
				return
			}
			data, err := ioutil.ReadAll(f)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "读取文件内容失败"})
				return
			}

			ConvertToXMLTofile(data)

		} else {
			var text string
			contentType := c.GetHeader("Content-Type")

			switch {
			case strings.Contains(contentType, "multipart/form-data"):
				text = c.PostForm("text")
			case strings.Contains(contentType, "application/x-www-form-urlencoded"),
				strings.Contains(contentType, "application/json"):
				text = c.PostForm("text")
				if text == "" {
					var jsonData map[string]string
					if err := c.BindJSON(&jsonData); err == nil {
						text = jsonData["text"]
					}
				}
			}

			if text == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "无有效文本内容"})
				return
			}

			text_data := []byte(text)
			ConvertToXMLTofile(text_data)

		}

		result, err := ProcessXml()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(result)==0{
			c.JSON(http.StatusOK, gin.H{"result": "没有检测到漏洞!"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"result": result})
	})

	r.Run()
}
