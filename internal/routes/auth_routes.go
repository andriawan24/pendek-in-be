package routes

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	userService  services.UserService
	oauthService services.OAuthService
}

func NewAuthRoutes(userService services.UserService, oauthService services.OAuthService) authRoutes {
	return authRoutes{
		userService:  userService,
		oauthService: oauthService,
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

	user, err := r.userService.FindUserByEmail(ctx.Request.Context(), param.Email)
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
			ID:              user.ID,
			Name:            user.Name,
			Email:           user.Email,
			IsActive:        user.IsActive,
			IsVerified:      user.IsVerified,
			ProfileImageUrl: user.ProfileImageUrl.String,
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

	user, err := r.userService.GetUserByID(ctx.Request.Context(), refreshClaims.UserId)
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

	user, err := r.userService.InsertUser(ctx.Request.Context(), registerParam)
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
			ID:              user.ID,
			Name:            user.Name,
			Email:           user.Email,
			IsActive:        user.IsActive,
			IsVerified:      user.IsVerified,
			ProfileImageUrl: user.ProfileImageUrl.String,
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

	user, err := r.userService.GetUserByID(ctx.Request.Context(), userId)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	response := responses.UserResponse{
		ID:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		IsActive:        user.IsActive,
		IsVerified:      user.IsVerified,
		ProfileImageUrl: user.ProfileImageUrl.String,
	}

	utils.RespondOK(ctx, "successfully get profile", response)
}

// UpdateProfile godoc
// @Summary      Update user profile
// @Description  Update current authenticated user profile with optional profile image upload
// @Tags         Auth
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        name formData string false "User name"
// @Param        email formData string false "User email"
// @Param        password formData string false "User password"
// @Param        profile_image formData file false "Profile image file (jpg, jpeg, png, gif)"
// @Success      200  {object}  responses.BaseResponse{data=responses.UserResponse}
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      401  {object}  responses.ErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /auth/update-profile [put]
func (r *authRoutes) UpdateProfile(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	user, err := r.userService.GetUserByID(ctx.Request.Context(), userId)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	name := ctx.PostForm("name")
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")

	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			utils.HandleErrorResponse(ctx, err)
			return
		}
		user.PasswordHash.String = string(hashedPassword)
	}

	if name != "" {
		user.Name = name
	}

	if email != "" && email != user.Email {
		user.Email = email
		user.IsVerified = false
	}

	file, err := ctx.FormFile("profile_image")
	if err == nil && file != nil {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true}
		if !allowedExts[ext] {
			utils.RespondBadRequest(ctx, "invalid file type. Allowed: jpg, jpeg, png, gif")
			return
		}

		const maxFileSize = 5 << 20 // 5MB
		if file.Size > maxFileSize {
			utils.RespondBadRequest(ctx, "file size exceeds 5MB limit")
			return
		}

		uploadDir := "uploads/profiles"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			utils.HandleErrorResponse(ctx, err)
			return
		}

		filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		filePath := filepath.Join(uploadDir, filename)

		if err := ctx.SaveUploadedFile(file, filePath); err != nil {
			utils.HandleErrorResponse(ctx, err)
			return
		}

		user.ProfileImageUrl.String = "/" + filePath
		user.ProfileImageUrl.Valid = true
	}

	updateUserParam := database.UpdateUserParams{
		ID:              userId,
		Name:            user.Name,
		Email:           user.Email,
		PasswordHash:    user.PasswordHash,
		IsVerified:      user.IsVerified,
		ProfileImageUrl: user.ProfileImageUrl,
	}

	updatedUser, err := r.userService.UpdateUser(ctx.Request.Context(), updateUserParam)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	response := responses.UserResponse{
		ID:              updatedUser.ID,
		Name:            updatedUser.Name,
		Email:           updatedUser.Email,
		IsActive:        updatedUser.IsActive,
		IsVerified:      updatedUser.IsVerified,
		ProfileImageUrl: updatedUser.ProfileImageUrl.String,
	}

	utils.RespondOK(ctx, "successfully update profile", response)
}

// GoogleAuth godoc
// @Summary      Google OAuth authentication
// @Description  Authenticate or register user via Google OAuth. Use without code param to get redirect URL, with code param to complete authentication.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        code  query     string  false  "OAuth authorization code from Google callback"
// @Success      200   {object}  responses.BaseResponse{data=responses.LoginResponse}
// @Failure      400   {object}  responses.ErrorResponse
// @Failure      500   {object}  responses.ErrorResponse
// @Router       /auth/google [get]
func (r *authRoutes) GoogleAuth(ctx *gin.Context) {
	code := ctx.Query("code")

	if code == "" {
		state := uuid.New().String()
		authURL := r.oauthService.GetGoogleAuthURL(state)
		response := responses.LoginResponse{
			AuthURL: authURL,
			State:   state,
		}
		utils.RespondOK(ctx, "redirect to Google for authentication", response)
		return
	}

	googleUser, err := r.oauthService.GetGoogleUserInfo(ctx.Request.Context(), code)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	user, err := r.userService.FindUserByGoogleID(ctx.Request.Context(), googleUser.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			insertParam := database.InsertUserWithGoogleParams{
				Name:  googleUser.Name,
				Email: googleUser.Email,
				GoogleID: sql.NullString{
					String: googleUser.ID,
					Valid:  true,
				},
				ProfileImageUrl: sql.NullString{
					String: googleUser.Picture,
					Valid:  googleUser.Picture != "",
				},
			}

			user, err = r.userService.InsertUserWithGoogle(ctx.Request.Context(), insertParam)
			if err != nil {
				utils.HandleErrorResponse(ctx, err)
				return
			}
		} else {
			utils.HandleErrorResponse(ctx, err)
			return
		}
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
			ID:              user.ID,
			Name:            user.Name,
			Email:           user.Email,
			IsActive:        user.IsActive,
			IsVerified:      user.IsVerified,
			ProfileImageUrl: user.ProfileImageUrl.String,
		},
	}

	utils.RespondOK(ctx, "successfully authenticated with Google", response)
}
