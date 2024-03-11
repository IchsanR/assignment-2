package controllers

import (
	"assigntment2/core"
	"assigntment2/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	repo *repository.Repository
}

func NewOrderController(repo *repository.Repository) *OrderController {
	return &OrderController{
		repo: repo,
	}
}

func (oc *OrderController) CreateOrder(c *gin.Context) {
	var order core.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if order.CustomerName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer name is required"})
		return
	}

	if len(order.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No items provided"})
		return
	}

	if order.OrderedAt.IsZero() {
		order.OrderedAt = time.Now()
	}

	orderID, err := oc.repo.InsertOrder(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert order", "details": err.Error()})
		return
	}

	// Construct the response JSON with the order and its ID
	response := gin.H{
		"id":           orderID,
		"orderedAt":    order.OrderedAt,
		"customerName": order.CustomerName,
		"items":        order.Items,
	}

	c.JSON(http.StatusCreated, response)
}
