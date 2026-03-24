package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RegionController 区域管理控制器
// 注意: 当前无对应数据库表, 提供基础 API 结构供后续扩展
type RegionController struct {
}

// RegionProvinceList 省份列表
func (c *RegionController) RegionProvinceList(ctx *gin.Context) {
	// 中国省份基础数据
	provinces := []map[string]interface{}{
		{"code": "110000", "name": "北京市"},
		{"code": "120000", "name": "天津市"},
		{"code": "130000", "name": "河北省"},
		{"code": "140000", "name": "山西省"},
		{"code": "150000", "name": "内蒙古自治区"},
		{"code": "210000", "name": "辽宁省"},
		{"code": "220000", "name": "吉林省"},
		{"code": "230000", "name": "黑龙江省"},
		{"code": "310000", "name": "上海市"},
		{"code": "320000", "name": "江苏省"},
		{"code": "330000", "name": "浙江省"},
		{"code": "340000", "name": "安徽省"},
		{"code": "350000", "name": "福建省"},
		{"code": "360000", "name": "江西省"},
		{"code": "370000", "name": "山东省"},
		{"code": "410000", "name": "河南省"},
		{"code": "420000", "name": "湖北省"},
		{"code": "430000", "name": "湖南省"},
		{"code": "440000", "name": "广东省"},
		{"code": "450000", "name": "广西壮族自治区"},
		{"code": "460000", "name": "海南省"},
		{"code": "500000", "name": "重庆市"},
		{"code": "510000", "name": "四川省"},
		{"code": "520000", "name": "贵州省"},
		{"code": "530000", "name": "云南省"},
		{"code": "540000", "name": "西藏自治区"},
		{"code": "610000", "name": "陕西省"},
		{"code": "620000", "name": "甘肃省"},
		{"code": "630000", "name": "青海省"},
		{"code": "640000", "name": "宁夏回族自治区"},
		{"code": "650000", "name": "新疆维吾尔自治区"},
		{"code": "710000", "name": "台湾省"},
		{"code": "810000", "name": "香港特别行政区"},
		{"code": "820000", "name": "澳门特别行政区"},
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": provinces})
}

// RegionCityList 城市列表
func (c *RegionController) RegionCityList(ctx *gin.Context) {
	provinceCode := ctx.Query("province_code")
	if provinceCode == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "省份代码不能为空"})
		return
	}
	// 返回空, 实际应根据 province_code 查询
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": []interface{}{}})
}

// RegionDistrictList 区县列表
func (c *RegionController) RegionDistrictList(ctx *gin.Context) {
	cityCode := ctx.Query("city_code")
	if cityCode == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "城市代码不能为空"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": []interface{}{}})
}

// RegionSearch 区域搜索
func (c *RegionController) RegionSearch(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	if keyword == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "关键词不能为空"})
		return
	}
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": []interface{}{}, "total": 0, "page": page, "page_size": pageSize,
	}})
}
