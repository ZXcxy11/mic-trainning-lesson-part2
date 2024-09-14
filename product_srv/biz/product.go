package biz

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"mic-trainning-lesson-part2/custom_error"
	"mic-trainning-lesson-part2/internal"
	"mic-trainning-lesson-part2/model"
	"mic-trainning-lesson-part2/proto/pb"
)

type ProductServer struct {
}

func (p ProductServer) ProductList(ctx context.Context, req *pb.ProductConditionReq) (*pb.ProductsRes, error) {
	//	逻辑判断
	fmt.Println("req：", req)
	var products []model.Product
	var Items []*pb.ProductItemRes
	var res pb.ProductsRes
	//	由于查询条件较为复杂，因此使用Model()指定操作数据模型更好
	m := internal.DB.Model(&products)
	//	用户选择性挑选查询条件（需要动态添加查询语句）
	if req.IsPop {
		//	注：这些构建查询条件的方法，都会返回新的gorm.db对象所以需要重新赋值
		m = m.Where("is_pop = ?", req.IsPop)
	}
	if req.IsNew {
		m = m.Where("is_new = ?", req.IsNew)
	}
	if req.MaxPrice > 0 {
		m = m.Where("max_price = ?", req.MaxPrice)
	}
	if req.MinPrice > 0 {
		m = m.Where("min_price > ?", req.MinPrice)
	}
	if req.BrandId > 0 {
		m = m.Where("brand_id = ?", req.BrandId)
	}
	if req.KeyWord != "" {
		m = m.Where("key_word like ?", "%"+req.KeyWord+"%")
	}
	if req.CategoryId > 0 {
		var category model.Category
		r := internal.DB.First(&category, req.CategoryId)
		if r.RowsAffected == 1 {
			return nil, errors.New(custom_error.CategoryNotExists)
		}
		var q string
		//	三级分类(根据不同级别分类查询，查询每个子分类下的全部商品)
		if category.Level == 1 {
			q = fmt.Sprintf("select id from category WHERE parent_category_id in (select id from category WHERE parent_category_id = %d)", req.CategoryId)
		} else if category.Level == 2 {
			q = fmt.Sprintf("select id from category WHERE parent_category_id=%d", req.CategoryId)
		} else if category.Level == 3 {
			q = fmt.Sprintf("select id from category WHERE id=%d", req.CategoryId)
		}
		m = m.Where(fmt.Sprintf("category_id in %s", q))
	}
	var count int64
	m.Count(&count)
	m.Joins("Category").Joins("Brand").Scopes(internal.MyPaging(int(req.PageNo), int(req.PageSize))).Find(&products)
	for _, item := range products {
		res := ConvertProductModel2Pb(item)
		Items = append(Items, res)
	}
	res.ItemList = Items
	res.Total = int32(count)
	fmt.Println("查询总数：", count)
	return &res, nil
}

func (p ProductServer) BatchGetProduct(ctx context.Context, req *pb.BatchProductIdReq) (*pb.ProductsRes, error) {
	var res pb.ProductsRes
	var products []model.Product

	r := internal.DB.Find(&products, req.Ids)
	res.Total = int32(r.RowsAffected)
	for _, item := range products {
		res.ItemList = append(res.ItemList, ConvertProductModel2Pb(item))
	}
	return &res, nil
}

func (p ProductServer) GetProductDetail(ctx context.Context, req *pb.ProductItemReq) (*pb.ProductItemRes, error) {
	if req.Id == 0 {
		return nil, errors.New(custom_error.ParamError)
	}
	var res *pb.ProductItemRes
	var product model.Product
	r := internal.DB.Find(&product, req.Id)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.ProductNotExists)
	}
	res = ConvertProductModel2Pb(product)
	return res, nil
}

func (p ProductServer) CreateProduct(ctx context.Context, req *pb.CreateProductItem) (*pb.ProductItemRes, error) {
	var category model.Category
	var brand model.Brand
	var res *pb.ProductItemRes
	//tx := internal.DB.Begin()
	//TODO 业务逻辑判断（实际上会更复杂）
	fmt.Println("req:", req)
	r := internal.DB.First(&category, req.CategoryId)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.BrandNotExists)
	}
	r = internal.DB.First(&brand, req.BrandId)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.CategoryNotExists)
	}
	item := ConvertReq2Model(req, category, brand)
	r = internal.DB.Save(&item)
	//defer func() {
	//	if e := tx.RowsAffected; e == 0 {
	//		tx.Rollback()
	//		fmt.Println("事务没通过")
	//	}
	//}()
	//if r.RowsAffected != 1 {
	//	zap.S().Error("插入失败")
	//	return nil, errors.New(custom_error.ProductCreateFail)
	//}
	//tx.Commit()
	res = ConvertProductModel2Pb(item)
	return res, nil
}

func (p ProductServer) DeleteProduct(ctx context.Context, item *pb.ProductDelItem) (*emptypb.Empty, error) {
	//	注：由于是软删除，所以不要忘记添加结构体实例的指针否则无法修改源记录中的删除时间字段（大部分设计修改的操作一定要添加结构体实例，否则会操作失败）
	r := internal.DB.Delete(&model.Product{}, item.Id)
	if r.RowsAffected < 1 {
		fmt.Println("r.Error: ", r.Error.Error())
		return nil, errors.New(custom_error.DelProductFail)
	}
	zap.S().Info("产品删除成功")
	return &emptypb.Empty{}, nil
}

func (p ProductServer) UpdateProduct(ctx context.Context, req *pb.CreateProductItem) (*emptypb.Empty, error) {
	//	业务逻辑判断(判断参数是否合理)
	//if req.Id == 0 || req.BrandId == 0 || req.CategoryId == 0 {
	//	return nil, errors.New(custom_error.ParamError)
	//}

	var product model.Product
	var brand model.Brand
	var category model.Category

	//	使用手动事务，注：需要使用返回对象 tx 来进行事务操作
	//tx := internal.DB.Begin()

	//	若出现错误则回滚数据
	//defer func() {
	//	//	recover()：若出现panic，则会捕获；此种写法可以根据整个程序是否正常运行来决定是否回滚（包括数据库问题）
	//	if p := recover(); p != nil {
	//		tx.Rollback()
	//	}
	//	//	另一种常见判断：若是数据库操作（添加失败等）出现错误时，可以使用 tx.Error() 来进行判断。
	//	//	这个的作用也可用于，当修改时，因为部分参数正确而部分错误，导致的部分字段修改情况不一致的情况（确保全部一起修改）
	//	//if e := tx.Error(); e != nil {
	//	//	tx.Rollback()
	//	//}
	//}()

	//r := tx.Find(&product, req.Id)
	//
	//if r.RowsAffected < 1 {
	//	return nil, errors.New(custom_error.ProductNotExists)
	//}
	//r = tx.First(&brand, req.BrandId)
	//if r.RowsAffected < 1 {
	//	return nil, errors.New(custom_error.BrandNotExists)
	//}
	//r = tx.First(&category, req.CategoryId)
	//if r.RowsAffected < 1 {
	//	return nil, errors.New(custom_error.CategoryNotExists)
	//}
	//
	//product = ConvertReq2Model(req, category, brand)
	//fmt.Println("product： ", product)
	//r = tx.Updates(&product)
	//
	//if r.RowsAffected == 0 {
	//	return nil, errors.New("修改失败")
	//}
	//	提交事务
	//tx.Commit()
	fmt.Println("req：", req)
	r := internal.DB.Find(&product, req.Id)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.ProductNotExists)
	}
	r = internal.DB.First(&brand, req.BrandId)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.BrandNotExists)
	}
	r = internal.DB.First(&category, req.CategoryId)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.CategoryNotExists)
	}
	product = ConvertReq2Model(req, category, brand)
	fmt.Println("product： ", product)
	//product.ID = req.Id
	r = internal.DB.Updates(&product)
	if r.RowsAffected < 1 {
		return nil, errors.New("删除失败")
	}
	return &emptypb.Empty{}, nil
}

// ConvertReq2Model ConvertProductModel2Pb 转换数据类型
func ConvertReq2Model(req *pb.CreateProductItem, category model.Category, brand model.Brand) model.Product {
	return model.Product{
		CategoryID: req.CategoryId,
		Category:   category,
		BrandID:    req.BrandId,
		Brand:      brand,
		Selling:    req.Selling,
		IsShipFree: req.IsShipFree,
		IsPop:      req.IsPop,
		IsNew:      req.IsNew,
		Name:       req.Name,
		SN:         req.Sn,
		FavNum:     req.FavNum,
		SoldNum:    req.SoldNum,
		Price:      req.Price,
		RealPrice:  req.RealPrice,
		ShortDesc:  req.ShortDesc,
		Images:     req.Images,
		DescImages: req.DescImages,
		CoverImage: req.CoverImage,
	}
}

func ConvertProductModel2Pb(pro model.Product) *pb.ProductItemRes {
	return &pb.ProductItemRes{
		CategoryId: pro.CategoryID,
		Name:       pro.Name,
		Sn:         pro.SN,
		SoldNum:    pro.SoldNum,
		FavNum:     pro.FavNum,
		Price:      pro.Price,
		RealPrice:  pro.RealPrice,
		ShortDesc:  pro.ShortDesc,
		Images:     pro.Images,
		DescImages: pro.DescImages,
		CoverImage: pro.CoverImage,
		IsNew:      pro.IsNew,
		IsPop:      pro.IsPop,
		Selling:    pro.Selling,
		Category: &pb.CategoryItemRes{
			Id:   pro.Category.ID,
			Name: pro.Category.Name,
		},
		Brand: &pb.BrandItemRes{
			Id:   pro.Brand.ID,
			Name: pro.Brand.Name,
			Logo: pro.Brand.Logo,
		},
	}
}
