package req

// 声明用户登录请求体的结构体
type LoginByPassword struct {
	//	tag标签，将结构体与json数据中的键进行绑定
	Mobile   string `json:"mobile" binding:"required"`
	Password string `json:"password" binding:"required,min=6,max=16"`
}
