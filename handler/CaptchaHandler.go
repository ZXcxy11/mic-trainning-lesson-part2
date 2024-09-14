package handler

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"mic-trainning-lesson-part2/internal"
	"net/http"
	"os"
	"time"
)

// gin+redis生成并保存验证码
func CaptchaHandler(c *gin.Context) {
	//	根据每个电话号码得到一个验证码
	mobile, ok := c.GetQuery("mobile")
	fmt.Println(mobile)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}

	filename := "data.png"
	f, err := os.Create(filename)
	if err != nil {
		zap.S().Error("GenCaptcha()中Create()执行失败")
		return
	}
	defer f.Close()
	var w io.WriterTo
	//	获取默认长度的随机字节数
	d := captcha.RandomDigits(captcha.DefaultLen)
	//	声明绘制验证码图片的对象
	w = captcha.NewImage("", d, captcha.StdWidth, captcha.StdHeight)
	//	执行绘图
	_, err = w.WriteTo(f)
	if err != nil {
		zap.S().Error("GenCaptcha()中WriteTo()执行失败")
		return
	}
	fmt.Println(d)
	//	拼接验证码
	captcha := ""
	for _, item := range d {
		captcha += fmt.Sprintf("%d", item)
	}
	fmt.Println(captcha)

	//	调用redis存储验证码相关信息
	internal.RedisClient.Set(context.Background(), mobile, captcha, 120*time.Second)

	b64, err := GetBase64(filename)
	if err != nil {
		zap.S().Error("GenCaptcha()中GetBase64()执行失败")
		return
	}
	fmt.Println(b64)
	c.JSON(http.StatusOK, gin.H{
		"captcha": b64,
	})
}

// 将文件编码为为64位字符存储
func GetBase64(fileName string) (string, error) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	//	os.Stat()用于获取属性的属性
	fileInfo, err := os.Stat(fileName)
	//	通过返回的fileInfo对象获取各个信息
	fmt.Println("文件字节大小：", fileInfo.Size())
	fmt.Println("文件名称：", fileInfo.Name())
	//	获取一个字节数组（通过文件大小合理设置数组大小），用于存储编码后的文件内容
	base64FileSize := 102400
	//int64(math.Ceil(float64(fileInfo.Size()*4.0) / 3.0))
	fmt.Println("编码后文件大小：", base64FileSize)
	b := make([]byte, base64FileSize)
	//	将文件内容转换Base64字符串
	base64.StdEncoding.Encode(b, file)
	s := string(b)
	fmt.Println(s)
	return s, nil
}
