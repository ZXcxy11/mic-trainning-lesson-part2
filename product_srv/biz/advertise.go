package biz

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"mic-trainning-lesson-part2/custom_error"
	"mic-trainning-lesson-part2/internal"
	"mic-trainning-lesson-part2/model"
	"mic-trainning-lesson-part2/proto/pb"
)

func (p ProductServer) AdvertiseList(ctx context.Context, empty *emptypb.Empty) (*pb.AdvertisesRes, error) {
	var adList []model.Advertise
	var adItemList []*pb.AdvertiseItemRes
	var advertisesRes pb.AdvertisesRes
	r := internal.DB.Find(&adList)
	if r.Error != nil {
		log.Fatal(r.Error)
	}
	for _, v := range adList {
		adItemList = append(adItemList, ConvertAdModel2Pb(v))
	}
	advertisesRes.ItemList = adItemList
	advertisesRes.Total = int32(r.RowsAffected)
	return &advertisesRes, nil
}

func (p ProductServer) CreateAdvertise(ctx context.Context, req *pb.AdvertiseReq) (*pb.AdvertiseItemRes, error) {
	ad := model.Advertise{
		Index: req.Index,
		Image: req.Image,
		Url:   req.Url,
	}
	internal.DB.Save(&ad)
	return ConvertAdModel2Pb(ad), nil
}

func (p ProductServer) DeleteAdvertise(ctx context.Context, req *pb.AdvertiseReq) (*emptypb.Empty, error) {
	r := internal.DB.Delete(&model.Advertise{}, req.Id)
	fmt.Println("r: ", r)
	//	TODO 需要业务判断失败
	if r.RowsAffected < 1 {

	}
	return &emptypb.Empty{}, nil
}

func (p ProductServer) UpdateAdvertise(ctx context.Context, req *pb.AdvertiseReq) (*emptypb.Empty, error) {
	var ad model.Advertise
	r := internal.DB.Find(&ad, req.Id)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.AdvertiseNotExists)
	}
	if req.Image != "" {
		ad.Image = req.Image
	}
	if req.Url != "" {
		ad.Url = req.Url
	}
	if req.Index > 0 {
		ad.Index = req.Index
	}
	internal.DB.Save(&ad)
	return &emptypb.Empty{}, nil
}

func ConvertAdModel2Pb(ad model.Advertise) *pb.AdvertiseItemRes {
	adRes := &pb.AdvertiseItemRes{
		Index: ad.Index,
		Image: ad.Image,
		Url:   ad.Url,
	}
	if ad.ID > 0 {
		adRes.Id = ad.ID
	}
	return adRes
}
