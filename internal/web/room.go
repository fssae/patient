package web

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/service"

	"gitee.com/fssae/ginx"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	svc *service.RoomService
}

func NewRoomHandler(svc *service.RoomService) *RoomHandler {
	return &RoomHandler{svc: svc}
}

func (h *RoomHandler) RegisterRoutes(server gin.IRouter) {
	group := server.Group("/api/rooms")
	group.GET("", ginx.Wrap(h.GetList))
	group.GET("/:id", ginx.Wrap(h.GetById))
	group.POST("/create", ginx.WrapBody[domain.Room](h.Create))
	group.PUT("/:id", ginx.WrapBody[domain.Room](h.Update))
	group.DELETE("/:id", ginx.Wrap(h.Delete))
}

func (h *RoomHandler) Create(c *gin.Context, req domain.Room) (ginx.Result, error) {
	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("创建成功"), nil
}

func (h *RoomHandler) GetById(c *gin.Context) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	room, err := h.svc.GetById(c.Request.Context(), id)
	if err != nil {
		return ginx.Result{}, err
	}
	if room == nil {
		return ginx.Fail(404, "房间不存在"), nil
	}
	return ginx.Ok("获取成功", room), nil
}

func (h *RoomHandler) GetList(c *gin.Context) (ginx.Result, error) {
	status := c.Query("status")
	floor := 0
	if room := c.Query("room"); room != "" {
		for _, ch := range room {
			if ch >= '0' && ch <= '9' {
				floor = floor*10 + int(ch-'0')
			}
		}
	}
	page := ginx.GetPage(c)

	list, total, err := h.svc.GetList(c.Request.Context(), status, floor, page.Skip, page.Limit)
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.OkList("获取成功", list, total), nil
}

func (h *RoomHandler) Update(c *gin.Context, req domain.Room) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	if err := h.svc.Update(c.Request.Context(), id, &req); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("更新成功"), nil
}

func (h *RoomHandler) Delete(c *gin.Context) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("删除成功"), nil
}
