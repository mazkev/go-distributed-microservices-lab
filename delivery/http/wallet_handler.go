package http

import (
	"net/http"

	"gotest/domain"
	"gotest/usecase"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	walletUsecase usecase.WalletUsecase
}

func NewWalletHandler(walletUsecase usecase.WalletUsecase) *WalletHandler {
	return &WalletHandler{walletUsecase: walletUsecase}
}

func (h *WalletHandler) GetMyWallet(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uint)

	wallet, err := h.walletUsecase.GetWallet(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": wallet})
}

func (h *WalletHandler) TopUp(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uint)

	var req domain.TopUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	wallet, err := h.walletUsecase.TopUp(userID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Top up berhasil",
		"data":    wallet,
	})
}

func (h *WalletHandler) Transfer(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	senderID := userIDVal.(uint)

	var req domain.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	err := h.walletUsecase.Transfer(senderID, req.ToUserID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Transfer saldo berhasil dieksekusi secara atomik",
	})
}
