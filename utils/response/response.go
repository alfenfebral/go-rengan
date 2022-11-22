package response

import (
	validator "go-rengan/pkg/validator"
	"net/http"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

// H is a shortcut for map[string]interface{}
type H map[string]interface{}

type SuccessList struct {
	Data interface{} `json:"data"`
	Meta *Meta       `json:"meta"`
}

type Meta struct {
	PerPage     int `json:"per_page"`
	CurrentPage int `json:"page"`
	TotalPage   int `json:"page_count"`
	TotalData   int `json:"total_count"`
}

type Success struct {
	Data interface{} `json:"data"`
}

// ErrorValidation - when error validation
func ErrorValidation(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, H{
		"success": false,
		"code":    http.StatusBadRequest,
		"message": "Validation errors in your request",
		"errors":  validator.ValidatonError(err).Errors,
	})
}

// ErrorBody - when error body eof
func ErrorBody(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, H{
		"success": false,
		"code":    http.StatusBadRequest,
		"message": "Validation errors in your request",
		"error":   "Check your body request",
	})
}

// ErrorInternal - when error internal server
func ErrorInternal(w http.ResponseWriter, r *http.Request, err error) {
	logrus.Error(err)

	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, H{
		"success": false,
		"code":    http.StatusInternalServerError,
		"message": "There is something error",
	})
}

// NotFound - when request not found
func NotFound(w http.ResponseWriter, r *http.Request, message string) {
	render.Status(r, http.StatusNotFound)
	render.JSON(w, r, H{
		"success": false,
		"code":    http.StatusNotFound,
		"message": message,
	})
}

// Created - when success created
func Created(w http.ResponseWriter, r *http.Request, data *Success) {
	render.Status(r, http.StatusCreated)

	render.JSON(w, r, H{
		"success": true,
		"code":    http.StatusCreated,
		"data":    data.Data,
	})
}

// ResponseOK - when success and return single data
func ResponseOK(w http.ResponseWriter, r *http.Request, data *Success) {
	render.Status(r, http.StatusOK)

	render.JSON(w, r, H{
		"success": true,
		"code":    http.StatusOK,
		"data":    data.Data,
	})
}

// ResponseOKList - when success and return array of data
func ResponseOKList(w http.ResponseWriter, r *http.Request, data *SuccessList) {
	render.Status(r, http.StatusOK)

	render.JSON(w, r, H{
		"success": true,
		"code":    http.StatusOK,
		"data":    data.Data,
		"meta":    data.Meta,
	})
}
