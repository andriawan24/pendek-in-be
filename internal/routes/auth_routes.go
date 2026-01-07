package routes

import (
	"database/sql"
	"errors"

	"github.com/andriawan24/link-short/internal/database"
	"github.com/andriawan24/link-short/internal/models/requests"
	"github.com/andriawan24/link-short/internal/models/responses"
	"github.com/andriawan24/link-short/internal/services"
	"github.com/andriawan24/link-short/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authRoutes struct {
	userService services.UserService
}

func NewAuthRoutes(userService services.UserService) authRoutes {
	return authRoutes{
		userService: userService,
	}
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user with email and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body requests.LoginParam true "Login credentials"
// @Success      200  {object}  responses.BaseResponse{data=responses.LoginResponse}
// @Failure      401  {object}  responses.ErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /auth/login [post]
func (r *authRoutes) Login(ctx *gin.Context) {
	var param requests.LoginParam

	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	user, err := r.userService.FindUserByEmail(param.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.RespondUnauthorized(ctx, "invalid email or password")
			return
		}

		utils.HandleErrorResponse(ctx, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(param.Password))
	if err != nil {
		utils.RespondUnauthorized(ctx, "invalid email or password, mismatch")
		return
	}

	accessToken, accessClaim, err := utils.GenerateJwtToken(user)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	refreshToken, refreshClaim, err := utils.GenerateRefreshToken(user)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	response := responses.LoginResponse{
		Token:                 accessToken,
		TokenExpiredAt:        accessClaim.ExpiresAt.Time,
		RefreshToken:          refreshToken,
		RefreshTokenExpiredAt: refreshClaim.ExpiresAt.Time,
		User: responses.UserResponse{
			ID:         user.ID,
			Name:       user.Name,
			Email:      user.Email,
			IsActive:   user.IsActive,
			IsVerified: user.IsVerified,
		},
	}

	utils.RespondOK(ctx, "successfully login", response)
}

// Refresh godoc
// @Summary      Refresh token
// @Description  Get new access token using refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body requests.RefreshParam true "Refresh token"
// @Success      200  {object}  responses.BaseResponse{data=responses.LoginResponse}
// @Failure      401  {object}  responses.ErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /auth/refresh [post]
func (r *authRoutes) Refresh(ctx *gin.Context) {
	var param requests.RefreshParam

	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	refreshClaims, err := utils.ParseRefreshToken(param.RefreshToken)
	if err != nil {
		utils.RespondUnauthorized(ctx, "invalid refresh token")
		return
	}

	user, err := r.userService.GetUserByID(refreshClaims.UserId)
	if err != nil {
		// Don't leak whether a user exists.
		utils.RespondUnauthorized(ctx, "invalid refresh token")
		return
	}

	accessToken, accessClaim, err := utils.GenerateAccessToken(user)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	newRefreshToken, newRefreshClaim, err := utils.GenerateRefreshToken(user)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	response := responses.LoginResponse{
		Token:                 accessToken,
		TokenExpiredAt:        accessClaim.ExpiresAt.Time,
		RefreshToken:          newRefreshToken,
		RefreshTokenExpiredAt: newRefreshClaim.ExpiresAt.Time,
	}

	utils.RespondOK(ctx, "successfully refresh token", response)
}

// Register godoc
// @Summary      Register new user
// @Description  Create a new user account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body requests.RegisterParam true "Registration details"
// @Success      200  {object}  responses.BaseResponse{data=responses.LoginResponse}
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /auth/register [post]
func (r *authRoutes) Register(ctx *gin.Context) {
	var param requests.RegisterParam

	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(param.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	registerParam := database.InsertUserParams{
		Name:  param.Name,
		Email: param.Email,
		PasswordHash: sql.NullString{
			String: string(hashedPassword),
			Valid:  true,
		},
	}

	user, err := r.userService.InsertUser(registerParam)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	accessToken, accessClaim, err := utils.GenerateJwtToken(user)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	refreshToken, refreshClaim, err := utils.GenerateRefreshToken(user)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	response := responses.LoginResponse{
		Token:                 accessToken,
		TokenExpiredAt:        accessClaim.ExpiresAt.Time,
		RefreshToken:          refreshToken,
		RefreshTokenExpiredAt: refreshClaim.ExpiresAt.Time,
		User: responses.UserResponse{
			ID:         user.ID,
			Name:       user.Name,
			Email:      user.Email,
			IsActive:   user.IsActive,
			IsVerified: user.IsVerified,
		},
	}

	utils.RespondOK(ctx, "successfully register", response)
}

// Profile godoc
// @Summary      Get user profile
// @Description  Get current authenticated user profile
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  responses.BaseResponse{data=responses.UserResponse}
// @Failure      401  {object}  responses.ErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /auth/me [get]
func (r *authRoutes) Profile(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)

	user, err := r.userService.GetUserByID(userId)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	response := responses.UserResponse{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
	}

	utils.RespondOK(ctx, "successfully get profile", response)
}

// UpdateProfile godoc
// @Summary      Update user profile
// @Description  Update current authenticated user profile
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body requests.UpdateProfileParam true "Profile update details"
// @Success      200  {object}  responses.BaseResponse{data=responses.UserResponse}
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      401  {object}  responses.ErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /auth/update-profile [put]
func (r *authRoutes) UpdateProfile(ctx *gin.Context) {
	var param requests.UpdateProfileParam

	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	userId := ctx.MustGet("user_id").(uuid.UUID)
	user, err := r.userService.GetUserByID(userId)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	if param.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(param.Password), bcrypt.DefaultCost)
		if err != nil {
			utils.HandleErrorResponse(ctx, err)
			return
		}

		user.PasswordHash.String = string(hashedPassword)
	}

	if param.Name != "" {
		user.Name = param.Name
	}

	if param.Email != "" && param.Email != user.Email {
		user.Email = param.Email
		user.IsVerified = false
	}

	updateUserParam := database.UpdateUserParams{
		ID:           userId,
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		IsVerified:   user.IsVerified,
	}

	updatedUser, err := r.userService.UpdateUser(updateUserParam)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	response := responses.UserResponse{
		ID:         updatedUser.ID,
		Name:       updatedUser.Name,
		Email:      updatedUser.Email,
		IsActive:   updatedUser.IsActive,
		IsVerified: updatedUser.IsVerified,
	}

	utils.RespondOK(ctx, "successfully update profile", response)
}
