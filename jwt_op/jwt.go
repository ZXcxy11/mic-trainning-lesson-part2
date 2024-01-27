package jwt_op

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"mic-trainning-lesson-part2/log"
	"time"
)

//	todo JWT的生成

// 声明JWT的错误变量
const (
	TokenExpired     = "Token已过期"
	TokenNotValidYet = "Token不再有效"
	TokenMalformed   = "Token非法"
	TokenInvalid     = "Token无效"
)

// 声明一个包含JWT声明的结构体对象，用于编到令牌中
type CustomClaims struct {
	jwt.StandardClaims
	ID          int32
	NickName    string
	AuthorityId int32
}

type JWT struct {
	SigningKey []byte
}

// 生成JWT对象
func NewJWT() *JWT {
	//return &JWT{SigningKey: []byte(conf.AppConf.JWTConfig.SigningKey)}
	key := []byte("my_secret_key")
	return &JWT{SigningKey: key}
}

// 生成JWT的声明
func (j *JWT) GenerateJWT(claim CustomClaims) (string, error) {
	/*
		jwt.NewWithClaims()：该函数用于创建一个JWT对象，参数需要指定签名算法和包含JWT的结构体
		jwt.SigningMethodES256：这是JWT的一个签名算法
		token.SignedString()：该函数适用于生成一个签名字符串，参数是一个密钥
	*/

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claim)
	tokenStr, err := token.SignedString(j.SigningKey)
	fmt.Println(j.SigningKey)
	if err != nil {
		log.Logger.Error("生成JWT错误：" + err.Error())
		return "", err
	}
	return tokenStr, nil
}

// 解析Token
func (j *JWT) ParseToken(tokenStr string) (*CustomClaims, error) {
	/*
		jwt.ParseWithClaims()：该函数适用于解析和验证JWT
			参数：
			tokenStr，JWT的字符串，永固解析和验证
			&CustomClaims{}，自定义生命结构，用于将JWT声明解析为自定义的结构体
			func(){}，用于提供验证JWT用的密钥
		该程序的作用：tokenStr与SigningKey进行验证，并将JWT解析为自定义声明，最后返回的token包含JWT的各个部分
	*/
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	//	判断出现什么错误
	if err != nil {
		if result, ok := err.(jwt.ValidationError); ok {
			if result.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New(TokenMalformed)
			} else if result.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New(TokenExpired)
			} else if result.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New(TokenNotValidYet)
			} else {
				return nil, errors.New(TokenInvalid)
			}
		}
	}
	//	类型断言：将token中的声明转换为自定义声明，结果ok返回是否转换成功
	if token != nil {
		//	不仅转换成功，而且还要保证token有效
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, errors.New(TokenInvalid)
	} else {
		return nil, errors.New(TokenInvalid)
	}
}

// 刷新Token(更新用户的Token有效时间)
func (j *JWT) RefreshToken(tokenStr string) (string, error) {
	/*
		jwt.TimeFunc 是jwt中函数类型的全局变量，用于提供时间相关操作
		此处是：通过匿名函数为该函数变量赋值，功能是返回Unix时间戳的起始时间
	*/
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	//	解析验证JWT
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	//	判断token是否有效
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		//	重新获取起始时间戳，刷新有效时间
		jwt.TimeFunc = time.Now
		//	设置JWT的过期时间字段(ExpiresAt)，此处表示有效期截止时间是从现在开始的一个星期
		claims.StandardClaims.ExpiresAt = time.Now().Add(7 * 24 * time.Hour).Unix()
		//	设置好后再从新生成JWT
		return j.GenerateJWT(*claims)
	}
	return "", errors.New(TokenInvalid)
}
