package controllers

import (
	"errors"
	"strconv"

	"go.uber.org/zap"

	"backend/internal/dto"
	"backend/internal/logger"
	"backend/internal/mappers"
	"backend/internal/middlewares"
	"backend/internal/services"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type TransportController struct {
	service services.TransportService
}

func NewTransportController(service services.TransportService) *TransportController {
	return &TransportController{
		service: service,
	}
}

func (ctrl *TransportController) CreateTransport(c *gin.Context) {
	currentUser, exists := middlewares.GetCurrentUser(c)
	if !exists || currentUser == nil {
		logger.Log.Warn("Không tìm thấy thông tin người dùng đăng nhập trong context")
		utils.UnauthorizedResponse(c, services.ErrMissingToken)
		return
	}

	var req dto.CreateTransportRequest
	if err := c.ShouldBindJSON(&req); err != nil {

		logger.Log.Warn(
			"Dữ liệu gửi lên không hợp lệ",
			zap.Error(err),
		)

		utils.InvalidDataResponse(c, err)
		return
	}

	transportModel := mappers.ToCreateTransport(&req)
	transportModel.IDNguoiTao = currentUser.ID

	result, err := ctrl.service.CreateTransport(c.Request.Context(), transportModel)
	if err != nil {
		if errors.Is(err, services.ErrEmpty) ||
			errors.Is(err, services.ErrNegative) ||
			errors.Is(err, services.ErrInvalidQuantity) ||
			errors.Is(err, services.ErrInvalidTask) ||
			errors.Is(err, services.ErrNegativeStock) {

			logger.Log.Warn(
				"Dữ liệu tạo phiếu vận chuyển không hợp lệ",
				zap.Error(err),
			)

			utils.InvalidDataResponse(c, err)
			return
		}

		if errors.Is(err, services.ErrAccountDisabled) {
			logger.Log.Warn(
				"Tài khoản người tạo đã bị vô hiệu hóa",
				zap.Int64("id_nguoi_tao", transportModel.IDNguoiTao),
				zap.Error(err),
			)
			utils.ForbiddenResponse(c, err)
			return
		}

		if errors.Is(err, services.ErrProductDisabled) {
			logger.Log.Warn(
				"Sản phẩm đã bị vô hiệu hóa",
				zap.Int64("id_san_pham", transportModel.IDSanPham),
				zap.Error(err),
			)
			utils.InvalidDataResponse(c, err)
			return
		}

		if errors.Is(err, services.ErrUserNotFound) {
			logger.Log.Warn(
				"Không tìm thấy người dùng với ID đã cung cấp",
				zap.Int64("id_nguoi_tao", transportModel.IDNguoiTao),
				zap.Error(err),
			)
			utils.NotFoundDataResponse(c, err)
			return
		}

		if errors.Is(err, services.ErrProductNotFound) {
			logger.Log.Warn(
				"Không tìm thấy sản phẩm với ID đã cung cấp",
				zap.Int64("id_san_pham", transportModel.IDSanPham),
				zap.Error(err),
			)
			utils.NotFoundDataResponse(c, err)
			return
		}

		logger.Log.Error(
			"Lỗi hệ thống khi tạo phiếu vận chuyển",
			zap.Error(err),
		)

		utils.NoConnectDataResponse(c, err)
		return
	}

	logger.Log.Info(
		"Tạo phiếu vận chuyển thành công",
		zap.Any("transport_id", result.ID),
	)

	utils.SuccessResponse(c, "Tạo phiếu vận chuyển", mappers.ToTransportResponse(result))
}

func (ctrl *TransportController) GetAllTransports(c *gin.Context) {
	transports, err := ctrl.service.GetAllTransports(c.Request.Context())
	if err != nil {

		logger.Log.Error(
			"Lỗi hệ thống khi lấy danh sách phiếu vận chuyển",
			zap.Error(err),
		)

		utils.NoConnectDataResponse(c, err)
		return
	}

	logger.Log.Info(
		"Lấy danh sách phiếu vận chuyển thành công",
		zap.Any("count", len(transports)),
	)

	responses := make([]dto.TransportResponse, 0, len(transports))
	for _, t := range transports {
		responses = append(responses, *mappers.ToTransportResponse(&t))
	}
	utils.SuccessResponse(c, "Lấy phiếu vận chuyển", responses)
}

func (ctrl *TransportController) GetTransportsByUserID(c *gin.Context) {
	currentUser, exists := middlewares.GetCurrentUser(c)
	if !exists || currentUser == nil {
		logger.Log.Warn("Không tìm thấy thông tin người dùng đăng nhập trong context")
		utils.UnauthorizedResponse(c, services.ErrMissingToken)
		return
	}

	idParam := c.Param("id")
	if idParam == "" {
		idParam = c.Param("userId")
	}

	userID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		logger.Log.Warn(
			"ID người dùng không hợp lệ",
			zap.Error(err),
		)

		utils.InvalidDataResponse(c, err)
		return
	}

	transports, err := ctrl.service.GetTransportsByUserID(c.Request.Context(), currentUser, userID)
	if err != nil {
		if errors.Is(err, services.ErrForbidden) {
			logger.Log.Warn(
				"Người dùng không có quyền xem phiếu vận chuyển của người khác",
				zap.Int64("current_user_id", currentUser.ID),
				zap.Int64("target_user_id", userID),
			)

			utils.ForbiddenResponse(c, err)
			return
		}

		if errors.Is(err, services.ErrUserNotFound) {
			logger.Log.Warn(
				"Không tìm thấy người dùng với ID đã cung cấp",
				zap.Int64("id_nguoi_tao", userID),
				zap.Error(err),
			)

			utils.NotFoundDataResponse(c, err)
			return
		}

		logger.Log.Error(
			"Lỗi hệ thống khi lấy danh sách phiếu vận chuyển theo người dùng",
			zap.Error(err),
		)

		utils.NoConnectDataResponse(c, err)
		return
	}

	logger.Log.Info(
		"Lấy danh sách phiếu vận chuyển theo người dùng thành công",
		zap.Int64("id_nguoi_tao", userID),
		zap.Any("count", len(transports)),
	)

	responses := make([]dto.TransportResponse, 0, len(transports))
	for _, t := range transports {
		responses = append(responses, *mappers.ToTransportResponse(&t))
	}
	utils.SuccessResponse(c, "Lấy phiếu vận chuyển", responses)
}
