package biz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"mic-trainning-lesson-part2/custom_error"
	"mic-trainning-lesson-part2/internal"
	"mic-trainning-lesson-part2/model"
	"mic-trainning-lesson-part2/proto/pb"
)

// 获取全部分类（包括全部层）

func (p ProductServer) GetAllCategoryList(ctx context.Context, empty *emptypb.Empty) (*pb.CategoriesRes, error) {
	var categoryList []model.Category
	internal.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categoryList)
	fmt.Println("categoryList: ", categoryList)
	var res pb.CategoriesRes
	//	若需要获取每一层的分类商品数据，需要逐层遍历，此处直接返回json数据观察效果
	//var items []*pb.CategoryItemRes
	//for _, c := range categoryList {
	//	items = append(items, ConvertCategoryModel2Pb(c))
	//}
	marshal, err := json.Marshal(categoryList)
	if err != nil {
		return nil, errors.New(custom_error.CategoryMarshalFail)
	}
	//res.InfoResList = items
	res.CategoryJsonFormat = string(marshal)
	return &res, nil
}

// 获取当前分类的所有子类(下一层)

func (p ProductServer) GetSubCategory(ctx context.Context, req *pb.CategoriesReq) (*pb.SubCategoriesRes, error) {
	var category model.Category
	var subItemList []*pb.CategoryItemRes
	var res pb.SubCategoriesRes
	r := internal.DB.First(&category, req.Id)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.CategoryNotExists)
	}
	pre := "SubCategory"
	if category.Level == 1 {
		pre = "SubCategory.SubCategory"
	}
	var subCategoryList []model.Category
	internal.DB.Where(&model.Category{ParentCategoryID: req.Id}).
		Preload(pre).
		Find(&subCategoryList)
	for _, c := range subCategoryList {
		subItemList = append(subItemList, ConvertCategoryModel2Pb(c))
	}
	b, err := json.Marshal(subItemList)
	if err != nil {
		return nil, errors.New(custom_error.CategoryMarshalFail)
	}
	res.SubCategoryList = subItemList
	res.CategoryJsonFormat = string(b)
	return &res, nil
}

//	新增分类

func (p ProductServer) CreteCategory(ctx context.Context, req *pb.CategoryItemReq) (*pb.CategoryItemRes, error) {
	category := model.Category{}
	category.Name = req.Name
	category.Level = req.Level
	category.ParentCategoryID = req.ParentCategoryId
	//	若存在父分类，则加上父分类ID
	if category.Level > 1 {
		category.ParentCategoryID = req.ParentCategoryId
	}
	internal.DB.Save(&category)
	res := ConvertCategoryModel2Pb(category)
	return res, nil
}

// 删除分类

func (p ProductServer) DeleteCategory(ctx context.Context, req *pb.CategoryDelReq) (*emptypb.Empty, error) {
	var level int
	internal.DB.Select("level").Where("id = ?", req.Id).First(&level)
	res := internal.DB.Delete(&model.Category{}, req.Id)
	if res.RowsAffected < 1 {
		return nil, errors.New(custom_error.CategoryNotExists)
	}
	if res.Error != nil {
		log.Fatal(res.Error)
		return nil, res.Error
	}
	//	逻辑判断
	//	如果删除的是一级分类，其它子分类一并删除
	if level == 1 {
		r := internal.DB.Where("ParentCategoryID = ?", req.Id).Delete(&model.Category{})
		if r.Error != nil {
			log.Fatal(r.Error)
			return nil, r.Error
		}
		fmt.Println(r.RowsAffected)
	}
	return &emptypb.Empty{}, nil
}

// 更新分类

func (p ProductServer) UpdateCategory(ctx context.Context, req *pb.CategoryItemReq) (*emptypb.Empty, error) {
	var cg model.Category
	r := internal.DB.Find(&cg, req.Id)
	//	删除前检查是否存在
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.CategoryNotExists)
	}
	if req.Name != "" {
		cg.Name = req.Name
	}
	if req.Level > 0 {
		cg.Level = req.Level
	}
	if req.ParentCategoryId > 0 {
		cg.ParentCategoryID = req.ParentCategoryId
	}
	fmt.Println("cg.ID, ", cg.ID)
	r = internal.DB.Updates(&cg)
	if r.Error != nil {
		return nil, r.Error
	}
	return &emptypb.Empty{}, nil
}

// 数据类型转换

func ConvertCategoryModel2Pb(c model.Category) *pb.CategoryItemRes {
	item := &pb.CategoryItemRes{
		Id:    c.ID,
		Name:  c.Name,
		Level: c.Level,
	}
	if c.Level > 1 {
		item.ParentCategoryId = c.ParentCategoryID
	}
	return item
}
