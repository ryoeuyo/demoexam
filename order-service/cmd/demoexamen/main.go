package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ryoeuyo/demoexamen/internal/domain/order"
	"github.com/ryoeuyo/demoexamen/internal/storage"
)

var id int64 = 1

type ID struct {
	Id int64 `uri:"id" binding:"required"`
}

const (
	NewOrder   order.Status = "новая заявка"
	Processing order.Status = "в процессе ремонта"
	Completed  order.Status = "завершена"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	}
}

func main() {
	storage, err := storage.NewStorage()
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Use(CORSMiddleware())
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/orders", func(ctx *gin.Context) {
		orders, err := storage.FetchAll()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "error",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"orders": orders,
		})
	})

	r.POST("/order", func(ctx *gin.Context) {
		var order order.Order
		if err := ctx.ShouldBindJSON(&order); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "error",
			})

			return
		}

		order.CreatedAt = time.Now()
		order.UpdatedAt = time.Now()
		order.Status = NewOrder
		order.ID = id
		id++

		id, err := storage.Create(&order)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "error",
			})

			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "created",
			"id":      id,
		})
	})

	r.PATCH("/order/:id", func(ctx *gin.Context) {
		var id ID
		if err := ctx.ShouldBindUri(&id); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "error",
			})

			return
		}

		var order order.Order
		if err := ctx.ShouldBindJSON(&order); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "error",
			})

			return
		}

		if err := storage.Update(id.Id, order); err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status": http.StatusNotFound,
				"error":  err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "updated",
			"id":      id,
		})
	})

	r.Run()
}
