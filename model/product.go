package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int32 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

//	商品分类

type Category struct {
	BaseModel
	Name             string `gorm:"type:varchar(32);not null"`
	ParentCategoryID int32
	ParentCategory   *Category
	SubCategory      []*Category `gorm:"foreignKey:ParentCategoryID;reference:ID"`
	Level            int32
}

//	品牌

type Brand struct {
	BaseModel
	Name string `gorm:"type:varchar(32);not null"`
	Logo string `gorm:"type:varchar(255);not null;default:''"`
}

//	广告

type Advertise struct {
	BaseModel
	Index int32  `gorm:"type:int;not null;default:1"`
	Image string `gorm:"type:varchar(255);not null"`
	Url   string `gorm:"type:varchar(255);not null"`
	Sort  int32  `gorm:"type:int;not null;default:1"`
}

//	产品

type Product struct {
	BaseModel
	CategoryID int32    `gorm:"type:int;not null"`
	Category   Category `gorm:"foreignKey:CategoryID"`

	BrandID int32 `gorm:"type:int;not null"`
	Brand   Brand `gorm:"foreignKey:BrandID"`

	Selling    bool `gorm:"default:false"`
	IsShipFree bool `gorm:"default:false"`
	IsPop      bool `gorm:"default:false"`
	IsNew      bool `gorm:"default:false"`

	Name       string  `gorm:"type:varchar(64);not null"`
	SN         string  `gorm:"type:varchar(64);not null"`
	FavNum     int32   `gorm:"type:int;default:0"`
	SoldNum    int32   `gorm:"type:int;default:0"`
	Price      float32 `gorm:"not null"`
	RealPrice  float32 `gorm:"not null"`
	ShortDesc  string  `gorm:"type:varchar(255);not null"`
	Images     MyList  `gorm:"type:varchar(1024);not null;"`
	DescImages MyList  `gorm:"type:varchar(1024);not null"`
	CoverImage string  `gorm:"type:varchar(255);not null"`
}
type ProductCategoryBrand struct {
	BaseModel
	CategoryID int32
	BrandID    int32
	Category   Category
	Brand      Brand
}

type MyList []string

// 将数组类型数据，转化为JSON数据存储在数据库中
func (myList MyList) Value() (driver.Value, error) {
	return json.Marshal(myList)
}

// 将存储在数据库中的JSON数据，转化回数组（所以使用模拟数据做测试时，模拟数据需要时JSON类型，varchar类型也能存放json类型数据）
func (myList *MyList) Scan(v interface{}) error {
	return json.Unmarshal(v.([]byte), myList)
}
