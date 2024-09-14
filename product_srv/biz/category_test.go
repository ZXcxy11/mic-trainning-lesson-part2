package biz

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"mic-trainning-lesson-part2/proto/pb"
	"testing"
)

func TestProductServer_CreteCategory(t *testing.T) {
	// 一级分类
	res, err := client.CreteCategory(context.Background(), &pb.CategoryItemReq{
		Name:             "鲜肉",
		ParentCategoryId: 1,
		Level:            1,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)

	//	二级分类
	res2, err := client.CreteCategory(context.Background(), &pb.CategoryItemReq{
		Name:             "牛肉",
		ParentCategoryId: res.Id,
		Level:            2,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res2)
	//	三级分类
	res3, err := client.CreteCategory(context.Background(), &pb.CategoryItemReq{
		Name:             "牛排",
		ParentCategoryId: res2.Id,
		Level:            3,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res3)
}

func TestProductServer_GetAllCategoryList(t *testing.T) {
	res, err := client.GetAllCategoryList(context.Background(), &emptypb.Empty{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.CategoryJsonFormat)
}

func TestProductServer_GetSubCategory(t *testing.T) {
	res, err := client.GetSubCategory(context.Background(), &pb.CategoriesReq{
		Id:    2,
		Level: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.CategoryJsonFormat)
	fmt.Println(res.SubCategoryList)
}

func TestProductServer_DeleteCategory(t *testing.T) {
	category, err := client.DeleteCategory(context.Background(), &pb.CategoryDelReq{
		Id: 33,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(category)
}

func TestProductServer_UpdateCategory(t *testing.T) {
	category, err := client.UpdateCategory(context.Background(), &pb.CategoryItemReq{
		Id:               6,
		Name:             "金枪鱼",
		ParentCategoryId: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(category)
}
