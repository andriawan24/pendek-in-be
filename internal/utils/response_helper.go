package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/andriawan24/link-short/internal/models/responses"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
)

func RespondOK(ctx *gin.Context, message string, data any) {
	ctx.JSON(http.StatusOK, responses.BaseResponse{
		Message: message,
		Data:    data,
	})
}

func ResponsdJson(ctx *gin.Context, status int, message string, data any) {
	ctx.JSON(status, responses.BaseResponse{
		Message: message,
		Data:    data,
	})
}

func RespondUnauthorized(ctx *gin.Context, message string) {
	respondError(ctx, http.StatusUnauthorized, message, nil)
}

func RespondBadRequest(ctx *gin.Context, message string) {
	respondError(ctx, http.StatusBadRequest, message, nil)
}

func respondError(ctx *gin.Context, status int, message string, err any) {
	ctx.JSON(status, responses.ErrorResponse{
		Message: message,
		Error:   err,
	})
}

func HandleErrorResponse(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, context.Canceled):
		respondError(ctx, http.StatusRequestTimeout, "request canceled", nil)
	case errors.Is(err, context.DeadlineExceeded):
		respondError(ctx, http.StatusGatewayTimeout, "request timeout", nil)
	case errors.Is(err, io.EOF):
		respondError(ctx, http.StatusBadRequest, "request body is required", nil)
	case errors.As(err, new(*json.SyntaxError)):
		respondError(ctx, http.StatusBadRequest, "invalid JSON", nil)
	case errors.As(err, new(*json.UnmarshalTypeError)):
		respondError(ctx, http.StatusBadRequest, "invalid JSON field type", nil)
	case errors.As(err, new(validator.ValidationErrors)):
		ve := err.(validator.ValidationErrors)
		fieldErrors := make(map[string]string, len(ve))
		for _, fe := range ve {
			field := fe.Field()
			switch fe.Tag() {
			case "required":
				fieldErrors[field] = "is required"
			case "email":
				fieldErrors[field] = "must be a valid email"
			case "min":
				fieldErrors[field] = "must be at least " + fe.Param()
			case "max":
				fieldErrors[field] = "must be at most " + fe.Param()
			default:
				fieldErrors[field] = fe.Tag()
			}
		}
		respondError(ctx, http.StatusBadRequest, "validation error", fieldErrors)
	case errors.As(err, new(*pq.Error)):
		pqErr := err.(*pq.Error)
		switch string(pqErr.Code) {
		case "23505": // unique_violation
			msg := "duplicate resource"
			respondError(ctx, http.StatusConflict, msg, gin.H{
				"constraint": pqErr.Constraint,
			})
		case "23503": // foreign_key_violation
			respondError(ctx, http.StatusConflict, "related resource not found", gin.H{
				"constraint": pqErr.Constraint,
			})
		case "23502": // not_null_violation
			respondError(ctx, http.StatusBadRequest, "missing required field", gin.H{
				"column": pqErr.Column,
			})
		case "23514": // check_violation
			respondError(ctx, http.StatusBadRequest, "invalid field value", gin.H{
				"constraint": pqErr.Constraint,
			})
		case "22P02": // invalid_text_representation (e.g., bad uuid)
			respondError(ctx, http.StatusBadRequest, "invalid request parameter", nil)
		default:
			// Unknown PG error: treat as server error, but don't leak details in non-debug.
			internal := any(nil)
			if gin.IsDebugging() {
				internal = err
			}
			respondError(ctx, http.StatusInternalServerError, "internal server error", internal)
		}
	case errors.Is(err, sql.ErrNoRows):
		respondError(ctx, http.StatusNotFound, "resource not found", nil)
	default:
		internal := any(nil)
		if gin.IsDebugging() {
			internal = err.Error()
		}
		respondError(ctx, http.StatusInternalServerError, "internal server error", internal)
	}
}
