package biz

import (
	"context"
	"fmt"
	"mic-trainning-lesson-part2/proto/pb"
	"testing"
)

/* 测试时，记得将测试的工作目录与主程序main的工作目录设置相同，否则会导致nacos配置信息所在的yaml文件读取失败（默认从当前工作目录找起）
例如：此处调用viper_config.go的init函数时，会在当前工作目录下寻找名为filename的yaml文件，即寻找失败出错
*/

func TestProductServer_CreateBrand(t *testing.T) {
	brands := []string{
		"a", "b", "c", "d",
	}
	for _, item := range brands {
		res, err := client.CreateBrand(context.Background(), &pb.BrandItemReq{
			Id:   0,
			Name: item,
			Logo: "XXX/XXX",
		})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res.Id)
	}

}

func TestProductServer_BrandList(t *testing.T) {
	res, err := client.BrandList(context.Background(), &pb.BrandPagingReq{
		PageNo:   2,
		PageSize: 3,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.Total)
	fmt.Println(res.ItemList)
}

func TestProductServer_DeleteBrand(t *testing.T) {
	res, err := client.DeleteBrand(context.Background(), &pb.BrandItemReq{
		Id: 11,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)

}

func TestProductServer_UpdateBrand(t *testing.T) {
	res, err := client.UpdateBrand(context.Background(), &pb.BrandItemReq{
		Id:   8,
		Name: "CCC",
		Logo: "XX//xX",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}
