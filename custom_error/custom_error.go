package custom_error

//	集中声明错误类型

const (
	BrandAlreadyExists           = "品牌已存在"
	DelBrandFail                 = "品牌删除失败"
	BrandNotExists               = "品牌不存在"
	AdvertiseNotExists           = "广告不存在"
	AdvertiseAlreadyExists       = "广告已存在"
	CategoryNotExists            = "分类不存在"
	CategoryAlreadyExists        = "分类已存在"
	CategoryMarshalFail          = "分类序列化错误"
	ProductAlreadyExists         = "产品已存在"
	DelProductFail               = "产品删除失败"
	ProductNotExists             = "产品不存在"
	DelCategoryBrandFail         = "分类品牌删除失败"
	DelProductCategoryBrandFail  = "产品分类品牌删除失败"
	ProductCategoryBrandNotFound = "产品分类品牌找不到"
	ParamError                   = "参数错误"
	ProductCreateFail            = "新增产品失败"
)
