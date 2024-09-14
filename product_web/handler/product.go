package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"mic-trainning-lesson-part2/custom_error"
	"mic-trainning-lesson-part2/internal"
	"mic-trainning-lesson-part2/product_web/req"
	"mic-trainning-lesson-part2/proto/pb"
	"net/http"
	"strconv"
)

//	修改完代码，重启看不到效果，可以考虑使用 go clean -cache 重新编译项目避免缓存问题

var productClient pb.ProductServiceClient

func init() {
	addr := fmt.Sprintf("%s:%d", internal.AppConf.ProductSrvConfig.Host,
		internal.AppConf.ProductSrvConfig.Port)
	fmt.Printf("==========srv的地址：" + addr)
	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
		//	负载均衡
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal(err)
	}
	productClient = pb.NewProductServiceClient(conn)

}

// ProductListHandler 获取商品列表
func ProductListHandler(c *gin.Context) {
	//	绑定商品的条件
	var condition pb.ProductConditionReq
	//	ShouldBindJSON函数将JSON数据绑定到结构体
	c.ShouldBindJSON(&condition) //list?pageNo=1&pageSize=2
	minPriceStr := c.DefaultQuery("minPrice", "0")
	//strconv.Atoi将字符串转换为整数
	minPrice, err := strconv.Atoi(minPriceStr)
	if err != nil {
		zap.S().Error("minPrice error")
		//	防御性编程，不让外界知道详情
		c.JSON(http.StatusOK, gin.H{"msg": custom_error.ParamError})
		return
	}
	maxPriceStr := c.DefaultQuery("maxPrice", "0")
	maxPrice, err := strconv.Atoi(maxPriceStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": custom_error.ParamError})
		return
	}
	condition.MinPrice = int32(minPrice)
	condition.MaxPrice = int32(maxPrice)

	//	种类
	categoryIdStr := c.DefaultQuery("categoryId", "0")
	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		zap.S().Error("categoryId error")
		c.JSON(http.StatusOK, gin.H{"msg": custom_error.ParamError})
	}
	condition.CategoryId = int32(categoryId)

	//	品牌
	brandIdStr := c.DefaultQuery("brandId", "0")
	brandId, err := strconv.Atoi(brandIdStr)
	if err != nil {
		zap.S().Error("brandId error")
		c.JSON(http.StatusOK, gin.H{"msg": custom_error.ParamError})
	}
	condition.BrandId = int32(brandId)
	//	是否为新品
	isNew := c.DefaultQuery("is_new", "0")
	if isNew == "1" {
		condition.IsNew = true
	}
	//	是否为热销
	isHot := c.DefaultQuery("is_pop", "0")
	if isHot == "1" {
		condition.IsPop = true
	}
	//	页面序号
	pageNoStr := c.DefaultQuery("pageNo", "0")
	pageNo, err := strconv.Atoi(pageNoStr)
	if err != nil {
		zap.S().Error("pageNo error")
		c.JSON(http.StatusOK, gin.H{"msg": custom_error.ParamError})
	}
	condition.PageNo = int32(pageNo)
	//	页面大小
	pageSizeStr := c.DefaultQuery("pageSize", "0")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		zap.S().Error("pageSize error")
		c.JSON(http.StatusOK, gin.H{"msg": custom_error.ParamError})
	}
	condition.PageSize = int32(pageSize)
	//	关键字
	keyword := c.DefaultQuery("keyword", "")
	condition.KeyWord = keyword

	//查询
	r, err := productClient.ProductList(context.Background(), &condition)
	if err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": "产品列表查询失败",
			//	给出默认值
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":   "",
		"total": r.Total,
		"data":  r.ItemList,
	})
}

// DetailHandler 获取商品详情
func DetailHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": "参数错误",
		})
		return
	}
	res, err := productClient.GetProductDetail(context.Background(), &pb.ProductItemReq{Id: int32(id)})
	if err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": "获取详情失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "",
		"data": res,
	})
	//	可以感受到，先将srv（服务）写好以后再写web（接口）可以更加清晰得让我们进行流量的对接
}

// AddHandler 添加商品
func AddHandler(c *gin.Context) {
	var productReq req.ProductReq
	//	将JSON数据，转换为结构体数据
	err := c.ShouldBindJSON(&productReq)
	if err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": "参数解析错误",
		})
		return
	}
	fmt.Println("productReq：", productReq)
	r := ConvertProductReq2Pb(productReq)
	res, err := productClient.CreateProduct(context.Background(), r)
	if err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": "添加产品失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "",
		"data": res,
		//	添加产品注意库存问题，不能相信前端
		//	产品的图片一般的解决方法是上传到一个地方（自己的图片服务器），将url作为其值，另一种是对象存储，例：腾讯云的COS，阿里云的OSS等
	})
}

// DelHandler 删除商品
func DelHandler(c *gin.Context) {
	//	从请求路径中获取参数（JSON类型）
	idStr := c.Param("id")
	//	将JSON类型转换为整数类型
	id, err := strconv.Atoi(idStr)
	if err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": "参数错误",
		})
		return
	}
	_, err = productClient.DeleteProduct(context.Background(), &pb.ProductDelItem{Id: int32(id)})
	if err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": "删除产品失败",
		})
		return
	}
}

// UpdateHandler 更新产品
func UpdateHandler(c *gin.Context) {
	var productReq req.ProductReq
	err := c.ShouldBindJSON(&productReq)
	if err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": "参数解析错误",
		})
		return
	}
	r := ConvertProductReq2Pb(productReq)

	_, err = productClient.UpdateProduct(context.Background(), r)
	if err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": "更新产品失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "",
	})

}

// ConvertProductReq2Pb 数据类型转换
func ConvertProductReq2Pb(productReq req.ProductReq) *pb.CreateProductItem {
	item := pb.CreateProductItem{
		Name:        productReq.Name,
		Sn:          productReq.SN,
		Stocks:      productReq.Stocks,
		Price:       productReq.Price,
		RealPrice:   productReq.RealPrice,
		ShortDesc:   productReq.ShortDesc,
		ProductDesc: productReq.Desc,
		Images:      productReq.Images,
		DescImages:  productReq.DescImages,
		CoverImage:  productReq.CoverImage,
		IsNew:       productReq.IsNew,
		IsPop:       productReq.IsPop,
		Selling:     productReq.Selling,
		BrandId:     productReq.BrandId,
		FavNum:      productReq.FavNum,
		SoldNum:     productReq.SoldNum,
		CategoryId:  productReq.CategoryId,
	}
	if productReq.Id > 0 {
		item.Id = productReq.Id
	}
	return &item
}
