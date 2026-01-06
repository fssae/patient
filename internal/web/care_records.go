package web

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/service"
	"classroom-analysis/internal/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CareRecordHandler struct {
	svc *service.CareRecordService
}

func NewCareRecordHandler(svc *service.CareRecordService) *CareRecordHandler {
	return &CareRecordHandler{
		svc: svc,
	}
}

func (h *CareRecordHandler) RegisterRoutes(server gin.IRouter) {
	group := server.Group("/api/care-records")
	group.GET("", h.GetList)
	group.GET("/:id", h.GetById)
	group.POST("", h.Create)
	group.POST("/:id/records", h.AddRecord) // 追加单条护理记录
	group.PUT("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
}

// Create 创建护理记录
// @Summary      创建护理记录
// @Description  为客户创建新的护理记录文档
// @Tags         护理记录管理
// @Accept       json
// @Produce      json
// @Param        request  body      domain.CareRecords  true  "护理记录信息"
// @Success      200      {object}  map[string]interface{}  "创建成功"
// @Router       /care-records [post]
func (h *CareRecordHandler) Create(c *gin.Context) {
	var req domain.CareRecords
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "创建成功",
		"success": true,
	})
}

// GetById 获取详情
// @Summary      获取详情
// @Description  根据ID获取护理记录详情
// @Tags         护理记录管理
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "记录ID"
// @Success      200  {object}  map[string]interface{}  "获取成功"
// @Router       /care-records/{id} [get]
func (h *CareRecordHandler) GetById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "无效的ID",
		})
		return
	}

	record, err := h.svc.GetById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}
	if record == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "记录不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "获取成功",
		"success": true,
		"data":    record,
	})
}

// GetList 获取列表
// @Summary      获取列表
// @Description  分页获取护理记录列表
// @Tags         护理记录管理
// @Accept       json
// @Produce      json
// @Param        customer_name  query     string  false  "客户姓名"
// @Param        skip           query     int     false  "跳过数量"  default(0)
// @Param        limit          query     int     false  "每页数量"  default(20)
// @Success      200            {object}  map[string]interface{}  "获取成功"
// @Router       /care-records [get]
func (h *CareRecordHandler) GetList(c *gin.Context) {
	customerName := c.Query("customer_name")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	customerID := c.Query("customer_id")
	skipStr := c.DefaultQuery("skip", "0")
	limitStr := c.DefaultQuery("limit", "20")

	skip, _ := strconv.ParseInt(skipStr, 10, 64)
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	list, total, err := h.svc.GetList(c.Request.Context(), customerName, skip, limit, startDate, endDate, customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "获取成功",
		"success": true,
		"data":    list,
		"total":   total,
	})
}

// Update 更新记录
// @Summary      更新记录
// @Description  根据ID更新护理记录信息
// @Tags         护理记录管理
// @Accept       json
// @Produce      json
// @Param        id       path      string              true  "记录ID"
// @Param        request  body      domain.CareRecords  true  "更新信息"
// @Success      200      {object}  map[string]interface{}  "更新成功"
// @Router       /care-records/{id} [put]
func (h *CareRecordHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "无效的ID",
		})
		return
	}

	// 使用 util.Validate 进行部分更新字段绑定（参考 CustomerHandler）
	updates, err := util.Validate(domain.CareRecords{}, c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	err = h.svc.Update(c.Request.Context(), id, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "更新成功",
		"success": true,
	})
}

// Delete 删除记录
// @Summary      删除记录
// @Description  根据ID删除护理记录
// @Tags         护理记录管理
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "记录ID"
// @Success      200  {object}  map[string]interface{}  "删除成功"
// @Router       /care-records/{id} [delete]
func (h *CareRecordHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "无效的ID",
		})
		return
	}

	err = h.svc.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "删除成功",
		"success": true,
	})
}

// AddRecord 追加护理记录子项
// @Summary      追加护理记录
// @Description  向指定护理记录文档中追加一条详细记录
// @Tags         护理记录管理
// @Accept       json
// @Produce      json
// @Param        id       path      string              true  "记录ID"
// @Param        request  body      domain.RecordItems  true  "详细记录信息"
// @Success      200      {object}  map[string]interface{}  "追加成功"
// @Router       /care-records/{id}/records [post]
func (h *CareRecordHandler) AddRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "无效的ID",
		})
		return
	}

	var req domain.RecordItems
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	err = h.svc.AddRecord(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "追加成功",
		"success": true,
	})
}
