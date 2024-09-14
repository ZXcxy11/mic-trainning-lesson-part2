package res

// 后台处理的数据包含较为隐私的数据不宜传到前端，因此自定义结构体来返回
type Account4Res struct {
	Mobile   string `json:"mobile"`
	NickName string `json:"nick_name"`
	Gender   string `json:"gender"`
}
