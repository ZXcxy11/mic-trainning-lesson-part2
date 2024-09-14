package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"mic-trainning-lesson-part2/proto/pb"
	"testing"
)

func TestProductServer_CategoryBrandList(t *testing.T) {
	req := pb.PagingReq{
		PageSize: 3,
		PageNo:   2,
	}
	res, err := client.CategoryBrandList(context.Background(), &req)
	if err != nil {
		zap.S().Error(err.Error())
	}
	fmt.Println("查询拿出来的总数：", res.Total)
	//使用JSON格式显示更明了
	prettyJSON, _ := json.Marshal(res.ItemList)
	fmt.Println("查询出来的数据：", string(prettyJSON))
}
func TestProductServer_GetCategoryBrandList(t *testing.T) {
	req := pb.CategoryItemReq{
		Id:               3,
		ParentCategoryId: 2,
	}
	res, err := client.GetCategoryBrandList(context.Background(), &req)
	if err != nil {
		zap.S().Error(err.Error())
	}
	prettyJSON, _ := json.Marshal(res.ItemList)
	fmt.Println("查询出来的数据：", string(prettyJSON))
	fmt.Println("总数为：", res.Total)
}
func TestProductServer_CreateCategoryBrand(t *testing.T) {
	res, err := client.CreateCategoryBrand(context.Background(), &pb.CategoryBrandReq{
		CategoryId: 2,
		BrandId:    4,
	})
	if err != nil {
		//	Fatal方法，一般用于单元测试中，为结构体t中的方法 报告错误并立即终止测试。
		t.Fatal(err)
	}
	fmt.Println(res.Id)
}
func TestProductServer_DeleteCategoryBrand(t *testing.T) {
	req := pb.CategoryBrandReq{
		Id: 8,
	}
	res, err := client.DeleteCategoryBrand(context.Background(), &req)
	if err != nil {
		zap.S().Error(err.Error())
	}
	fmt.Println("返回信息：", res)
}
func TestProductServer_UpdateCategoryBrand(t *testing.T) {
	req := pb.CategoryBrandReq{
		Id:         7,
		BrandId:    5,
		CategoryId: 2,
	}
	res, err := client.UpdateCategoryBrand(context.Background(), &req)
	if err != nil {
		zap.S().Error(err.Error())
	}
	fmt.Println(res)
}
