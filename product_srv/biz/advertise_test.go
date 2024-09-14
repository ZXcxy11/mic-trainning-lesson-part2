package biz

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"mic-trainning-lesson-part2/proto/pb"
	"testing"
)

func TestProductServer_AdvertiseList(t *testing.T) {
	list, err := client.AdvertiseList(context.Background(), &emptypb.Empty{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("total: ", list.Total)
	fmt.Println(list.ItemList)
}
func TestProductServer_CreateAdvertise(t *testing.T) {
	req := pb.AdvertiseReq{
		Index: 12345,
		Image: "XXX.com",
		Url:   "XXX/XXX",
	}
	req2 := pb.AdvertiseReq{
		Index: 123,
		Image: "XXX.com2",
		Url:   "XXX/XXX2",
	}
	res1, err := client.CreateAdvertise(context.Background(), &req)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(res1.Id)
	res2, err := client.CreateAdvertise(context.Background(), &req2)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(res2.Id)
}
func TestProductServer_DeleteAdvertise(t *testing.T) {
	req := pb.AdvertiseReq{Id: 11}
	_, err := client.DeleteAdvertise(context.Background(), &req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func TestProductServer_UpdateAdvertise(t *testing.T) {
	_, err := client.UpdateAdvertise(context.Background(), &pb.AdvertiseReq{
		Id:    3,
		Index: 3,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}
