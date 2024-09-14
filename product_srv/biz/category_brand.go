package biz

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"mic-trainning-lesson-part2/custom_error"
	"mic-trainning-lesson-part2/internal"
	"mic-trainning-lesson-part2/model"
	"mic-trainning-lesson-part2/proto/pb"
)

// CategoryBrandList 返回分类和品牌的列表，支持分页
func (p ProductServer) CategoryBrandList(ctx context.Context, req *pb.PagingReq) (*pb.CategoryBrandListRes, error) {
	// 变量声明
	var count int64                        // 用于存储总记录数
	var items []model.ProductCategoryBrand // 用于存储查询到的产品分类品牌数据
	var resList []*pb.CategoryBrandRes     // 用于存储转换后的响应数据
	var res pb.CategoryBrandListRes        // 最终返回的响应结构

	// 查询产品分类品牌的总数
	internal.DB.Model(&model.ProductCategoryBrand{}).Count(&count)
	// 预加载相关的分类和品牌数据，应用分页查询
	internal.DB.Preload("Category").Preload("Brand").Scopes(internal.MyPaging(int(req.PageNo), int(req.PageSize))).Find(&items)

	fmt.Println("items, ", items)

	// 遍历查询到的产品分类品牌数据，将其转换为响应格式
	for _, item := range items {
		pcb := CovertProductCategoryBrand2Pb(item) // 调用转换函数，将 model 转换为 protobuf 结构
		resList = append(resList, pcb)             // 将转换后的结果添加到响应列表中
	}

	// 设置最终响应的记录总数和具体项列表
	res.Total = int32(count) // 设置总记录数到响应中
	res.ItemList = resList   // 设置项列表到响应中
	return &res, nil         // 返回响应结构和 nil（表示没有错误）
}

// GetCategoryBrandList 获取指定分类的品牌列表
func (p ProductServer) GetCategoryBrandList(ctx context.Context, req *pb.CategoryItemReq) (*pb.BrandRes, error) {
	// 声明变量
	var itemList []model.ProductCategoryBrand // 用于存储查询到的产品分类品牌数据
	var itemListRes []*pb.BrandItemRes        // 用于存储响应格式的品牌数据
	var category model.Category               // 用于存储查询到的分类数据
	var res pb.BrandRes                       // 最终返回的响应结构

	// 查询指定 ID 的分类数据
	r := internal.DB.First(&category, req.Id)
	if r.RowsAffected == 0 { // 如果没有找到分类
		return nil, errors.New(custom_error.ProductCategoryBrandNotFound) // 返回错误
	}

	// 查询属于指定父分类的产品品牌
	r = internal.DB.Preload("Brand").Where(&model.ProductCategoryBrand{CategoryID: req.ParentCategoryId}).Find(&itemList)
	if r.RowsAffected > 0 { // 如果找到了产品分类品牌
		res.Total = int32(r.RowsAffected) // 设置响应中的总记录数
	}

	// 遍历查询到的产品分类品牌，将其转换为响应格式
	for _, item := range itemList {
		itemListRes = append(itemListRes, &pb.BrandItemRes{
			Id:   item.Brand.ID,   // 品牌 ID
			Name: item.Brand.Name, // 品牌名称
			Logo: item.Brand.Logo, // 品牌 logo
		})
	}

	// 将转换后的品牌数据列表添加到响应中
	res.ItemList = itemListRes
	return &res, nil // 返回响应结构和 nil（表示没有错误）
}

func (p ProductServer) CreateCategoryBrand(ctx context.Context, req *pb.CategoryBrandReq) (*pb.CategoryBrandRes, error) {
	var res pb.CategoryBrandRes
	var item model.ProductCategoryBrand
	var category model.Category
	var brand model.Brand
	//	分类判断
	r := internal.DB.First(&category, req.CategoryId)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.CategoryNotExists)
	}
	//	品牌判断
	r = internal.DB.First(&brand, req.BrandId)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.BrandNotExists)
	}
	//	是否已经存在关系（待处理）
	item.Category.ID = req.CategoryId
	item.Brand.ID = req.BrandId
	r = internal.DB.Save(&item)
	if r.Error != nil {
		fmt.Println(r.Error)
	}
	res.Id = item.ID
	return &res, nil
}

func (p ProductServer) DeleteCategoryBrand(ctx context.Context, req *pb.CategoryBrandReq) (*emptypb.Empty, error) {
	r := internal.DB.Delete(&model.ProductCategoryBrand{}, req.Id)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.DelProductCategoryBrandFail)
	}
	return &emptypb.Empty{}, nil
}

func (p ProductServer) UpdateCategoryBrand(ctx context.Context, req *pb.CategoryBrandReq) (*emptypb.Empty, error) {
	var category model.Category
	var brand model.Brand
	//	将前端传回的数据进行查询，判断参数是否合理
	r := internal.DB.First(&category, req.CategoryId)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.CategoryNotExists)
	}
	r = internal.DB.First(&brand, req.BrandId)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.BrandNotExists)
	}
	//	将前端的参数放入结构体中
	productCategoryBrand := model.ProductCategoryBrand{}
	//	先从数据库获取记录映射到结构体中，在修改结构体信息，最后将结构体映射到数据库
	internal.DB.First(&productCategoryBrand, req.Id)
	productCategoryBrand.Category.ID = req.CategoryId
	productCategoryBrand.Brand.ID = req.BrandId
	//	调用数据库
	internal.DB.Save(&productCategoryBrand)
	//	emptypb.Empty 可用于必须有返回值，而不需要传递任何实际数据的情况（空消息结构体）
	return &emptypb.Empty{}, nil
}

// CovertProductCategoryBrand2Pb 数据转换(利于复用)
func CovertProductCategoryBrand2Pb(pcb model.ProductCategoryBrand) *pb.CategoryBrandRes {
	cb := pb.CategoryBrandRes{
		Id: pcb.ID,
		Brand: &pb.BrandItemRes{
			Id:   pcb.Brand.ID,
			Name: pcb.Brand.Name,
			Logo: pcb.Brand.Logo,
		},
		Category: &pb.CategoryItemRes{
			Id:               pcb.Category.ID,
			Name:             pcb.Category.Name,
			ParentCategoryId: pcb.Category.ParentCategoryID,
			Level:            pcb.Category.Level,
		},
	}
	return &cb
}
