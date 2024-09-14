package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 该跨域的处理方式，是通过在Gin中设置CORS头，以实现跨域资源共享（CORS）
func CrossDomain(c *gin.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		//	获取请求的方法（GET、POST、。。。）
		method := c.Request.Method

		//	设置响应头
		//	1. 允许任何域的请求
		c.Header("Access-Control-Allow-Origin", "*")
		//	2. 浏览器允许的请求头
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSrF-Token,Authorization,Token,x-token")
		//	3. 指定服务器支持的HTTP请求方法
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
		//	4. 浏览器在响应中可以访问的请求头
		c.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type")
		//	5. 允许在跨域请求中使用凭据
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			//	中止请求
			c.AbortWithStatus(http.StatusNoContent)
		}

	}
}
