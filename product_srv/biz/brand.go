package biz

import (
	"context"
	"errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"mic-trainning-lesson-part2/custom_error"
	"mic-trainning-lesson-part2/internal"
	"mic-trainning-lesson-part2/model"
	"mic-trainning-lesson-part2/proto/pb"
)

func (p ProductServer) BrandList(ctx context.Context, req *pb.BrandPagingReq) (*pb.BrandRes, error) {
	//	存放数据库获取的品牌数据列表
	var brandList []model.Brand
	//	存放经转化过的响应请求的品牌数据列表
	var brands []*pb.BrandItemRes
	//	存放响应请求的整合数据（数据列表 + 总共的数量）
	var brandRes pb.BrandRes

	//	非分页的写法
	//r := internal.DB.Find(&brandList)
	//fmt.Println(r.RowsAffected)
	//
	//for _, item := range brandList {
	//	brands = append(brands, ConvertBrandModel2Pb(item))
	//}
	//brandRes.ItemList = brands
	//brandRes.Total = int32(r.RowsAffected)

	//	分页
	//	存放满足条件的记录总数
	//var count int64
	//	需要跳过第几页（分页偏移量）
	//skip := (req.PageNo - 1) * req.PageSize
	//r := internal.DB.Model(&model.Brand{}).Count(&count).Offset(int(skip)).Limit(int(req.PageSize)).Find(brandList)
	//if r.RowsAffected < 1 {
	//	//TODO 可以进一步判断，根据业务需求
	//}
	//for _, item := range brandList {
	//	brands = append(brands, ConvertBrandModel2Pb(item))
	//}
	//brandRes.ItemList = brands
	//brandRes.Total = int32(count)

	//	分页2：通过自定义分页模板实现
	var count int64
	//paging := internal.MyPaging(int(req.PageNo), int(req.PageSize))
	//r := paging(internal.DB).Find(&brandList)
	r := internal.DB.Scopes(internal.MyPaging(int(req.PageNo), int(req.PageSize))).Count(&count).Find(&brandList)
	if r.RowsAffected < 1 {
		//TODO 可以进一步判断，根据业务需求
	}
	for _, item := range brandList {
		brands = append(brands, ConvertBrandModel2Pb(item))
	}
	brandRes.ItemList = brands
	brandRes.Total = int32(count)
	return &brandRes, nil
}

func (p ProductServer) CreateBrand(ctx context.Context, req *pb.BrandItemReq) (*pb.BrandItemRes, error) {
	var brand model.Brand
	r := internal.DB.Find("name=? and logo=?", req.Name, req.Logo)
	if r.RowsAffected > 0 {
		return nil, errors.New(custom_error.BrandAlreadyExists)
	}
	brand.Name = req.Name
	brand.Logo = req.Logo
	internal.DB.Save(&brand)
	return ConvertBrandModel2Pb(brand), nil
}

func (p ProductServer) DeleteBrand(ctx context.Context, req *pb.BrandItemReq) (*emptypb.Empty, error) {
	r := internal.DB.Delete(&model.Brand{}, req.Id)
	if r.Error != nil {
		return nil, errors.New(custom_error.DelBrandFail)
	}
	//	声明一个空信息
	return &emptypb.Empty{}, nil
}

func (p ProductServer) UpdateBrand(ctx context.Context, req *pb.BrandItemReq) (*emptypb.Empty, error) {
	var brand model.Brand
	r := internal.DB.Find(&brand, req.Id)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.BrandNotExists)
	}
	if req.Name != "" {
		brand.Name = req.Name
	}
	if req.Logo != "" {
		brand.Logo = req.Logo
	}
	internal.DB.Save(&brand)
	return &emptypb.Empty{}, nil
}

// 转换数据类型 该对象将用于请求响应
func ConvertBrandModel2Pb(item model.Brand) *pb.BrandItemRes {
	brand := &pb.BrandItemRes{
		Name: item.Name,
		Logo: item.Logo,
	}
	if item.ID > 0 {
		brand.Id = item.ID
	}
	return brand
}
