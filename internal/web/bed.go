package web

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/service"

	"gitee.com/fssae/ginx"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BedHandler struct {
	svc *service.BedService
}

func NewBedHandler(svc *service.BedService) *BedHandler {
	return &BedHandler{svc: svc}
}

func (h *BedHandler) RegisterRoutes(server gin.IRouter) {
	group := server.Group("/api/beds")
	group.GET("", ginx.Wrap(h.GetList))
	group.GET("/options/beds", ginx.Wrap(h.GetBedOptions))
	group.GET("/options/rooms", ginx.Wrap(h.GetRoomOptions))
	group.GET("/:id", ginx.Wrap(h.GetById))
	group.GET("/room/:room_id", ginx.Wrap(h.GetByRoomID))
	group.POST("/create", ginx.WrapBody[domain.CreateBed](h.Create))
	group.PUT("/:id", ginx.WrapBody[domain.Bed](h.Update))
	group.PUT("/:id/release", ginx.Wrap(h.Release))
	group.DELETE("/:id", ginx.Wrap(h.Delete))
}

func (h *BedHandler) Create(c *gin.Context, req domain.CreateBed) (ginx.Result, error) {
	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.Ok("创建成功", gin.H{"room_id": req.RoomId}), nil
}

func (h *BedHandler) GetById(c *gin.Context) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	bed, err := h.svc.GetById(c.Request.Context(), id)
	if err != nil {
		return ginx.Result{}, err
	}
	if bed == nil {
		return ginx.Fail(404, "床位不存在"), nil
	}
	return ginx.Ok("获取成功", bed), nil
}

func (h *BedHandler) GetByRoomID(c *gin.Context) (ginx.Result, error) {
	roomIDStr := c.Param("room_id")
	roomID, err := primitive.ObjectIDFromHex(roomIDStr)
	if err != nil {
		return ginx.Fail(400, "无效的房间ID"), nil
	}
	beds, err := h.svc.GetByRoomID(c.Request.Context(), roomID)
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Ok("获取成功", beds), nil
}

func (h *BedHandler) GetList(c *gin.Context) (ginx.Result, error) {
	roomNumber := c.Query("room_number")
	bedNumber := c.Query("bed_number")
	status := c.Query("status")
	page := ginx.GetPage(c)

	list, total, err := h.svc.GetList(c.Request.Context(), roomNumber, bedNumber, status, page.Skip, page.Limit)
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.OkList("获取成功", list, total), nil
}

func (h *BedHandler) Update(c *gin.Context, req domain.Bed) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	if err := h.svc.Update(c.Request.Context(), id, &req); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("更新成功"), nil
}

func (h *BedHandler) Release(c *gin.Context) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	if err := h.svc.Release(c.Request.Context(), id); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("释放成功"), nil
}

func (h *BedHandler) Delete(c *gin.Context) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("删除成功"), nil
}

func (h *BedHandler) GetBedOptions(c *gin.Context) (ginx.Result, error) {
	options, err := h.svc.GetBedOptions(c.Request.Context())
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Ok("获取成功", options), nil
}

func (h *BedHandler) GetRoomOptions(c *gin.Context) (ginx.Result, error) {
	options, err := h.svc.GetRoomOptions(c.Request.Context())
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Ok("获取成功", options), nil
}
