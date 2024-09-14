package biz

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"mic-trainning-lesson-part2/internal"
	"mic-trainning-lesson-part2/proto/pb"
	"os"
	"testing"
)

// 为当前类测试写一个init函数，连接srv服务器
var client pb.ProductServiceClient

func init() {
	//	设置测试用的工作路径
	err := os.Chdir("D:/goProject/mic-trainning-lesson-part2")
	if err != nil {
		zap.S().Error("测试工作目录设置失败")
	}
	addr := fmt.Sprintf("%s:%d", "192.168.150.11", internal.AppConf.ProductSrvConfig.Port)
	fmt.Println("测试的地址为：" + addr)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return
	}
	client = pb.NewProductServiceClient(conn)
}

func TestProductServer_BatchGetProduct(t *testing.T) {
	res, err := client.BatchGetProduct(context.Background(), &pb.BatchProductIdReq{Ids: []int32{1, 2, 5}})
	if err != nil {
		zap.S().Error(err)
	}
	fmt.Println("total: ", res.Total)
	fmt.Println("ItemList: ", res.ItemList)
}

func TestProductServer_ProductList(t *testing.T) {
	list, err := client.ProductList(context.Background(), &pb.ProductConditionReq{
		PageNo:   1,
		PageSize: 3,
	})
	if err != nil {
		zap.S().Error(err)
	}
	fmt.Println(list.Total)
	fmt.Println(list.ItemList)
}

func TestProductServer_GetProductDetail(t *testing.T) {
	detail, err := client.GetProductDetail(context.Background(), &pb.ProductItemReq{Id: 1})
	if err != nil {
		zap.S().Error(err)
	}
	fmt.Println(detail)
}

func TestProductServer_CreateProduct(t *testing.T) {
	for i := 0; i < 2; i++ {
		res, err := client.CreateProduct(context.Background(), &pb.CreateProductItem{
			Name:        fmt.Sprintf("Laptop A%d", i),
			Sn:          "123456",
			Stocks:      5000,
			Price:       399.00,
			RealPrice:   199.00,
			ShortDesc:   "",
			ProductDesc: "",
			Images:      nil,
			DescImages:  nil,
			CoverImage:  "http://image.png",
			IsNew:       true,
			IsPop:       true,
			Selling:     true,
			BrandId:     2,
			FavNum:      4856,
			SoldNum:     7928,
			CategoryId:  3,
			IsShipFree:  true,
		})
		if err != nil {
			zap.S().Error(err.Error())
		}
		fmt.Println(res)
	}
	fmt.Sprintln("什么鬼")
}

func TestProductServer_DeleteProduct(t *testing.T) {
	res, err := client.DeleteProduct(context.Background(), &pb.ProductDelItem{Id: 4})
	if err != nil {
		zap.S().Error(err)
	}
	fmt.Println(res)
}

func TestProductServer_UpdateProduct(t *testing.T) {
	res, err := client.UpdateProduct(context.Background(), &pb.CreateProductItem{
		Id:         3,
		Name:       "Laptop Model Z2",
		BrandId:    2,
		CategoryId: 3,
		FavNum:     800,
		SoldNum:    1000,
	})
	if err != nil {
		zap.S().Error(err)
	}
	fmt.Println(res)
}
