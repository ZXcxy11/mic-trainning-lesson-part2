package middleware

import (
	"github.com/gin-gonic/gin"
	"mic-trainning-lesson-part2/jwt_op"
	"net/http"
)

// 从http请求头获取token并进行认证的中间件函数

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//	从HTTP请求头获取token
		token := c.Request.Header.Get("token")
		//	判断token是否为空
		if token == "" || len(token) == 0 {
			//	若为空，表示无token，验证失败，返回信息给前端
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "认证失败，需要登陆",
			})
			//	Abort()函数，用于提前终止请求的处理，阻止后续中间件函数或者处理函数执行
			c.Abort()
			return
		}
		//	有token就解析
		j := jwt_op.NewJWT()
		parseToken, err := j.ParseToken(token)
		//	判断解析是否有问题
		if err != nil {
			//	token过期
			if err.Error() == jwt_op.TokenExpired {
				c.JSON(http.StatusUnauthorized, gin.H{
					"msg": jwt_op.TokenExpired,
				})
				c.Abort()
				return
			}
			//	其他错误
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "认证失败，请重新登陆",
			})
			c.Abort()
			return
		}
		//	c.Set()函数，用于存储自定义数据到上下文中去，多用于中间件或处理函数之间的数据通信
		//	若该token无问题，则将其jwt声明，放入上下文中
		c.Set("claim", parseToken)
	}
}
