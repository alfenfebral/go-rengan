package httpdelivery

import (
	"net/http"

	tracing "go-rengan/pkg/tracing"
	validator "go-rengan/pkg/validator"
	"go-rengan/todo/models"
	"go-rengan/todo/service"
	paginationutil "go-rengan/utils/pagination"
	responseutil "go-rengan/utils/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type HTTPHandler interface {
	RegisterRoutes(router *chi.Mux)
	GetAll(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type HTTPHandlerImpl struct {
	tracing     tracing.Tracing
	todoService service.Service
}

// New - make http handler
func New(tracing tracing.Tracing, service service.Service) HTTPHandler {
	return &HTTPHandlerImpl{
		tracing:     tracing,
		todoService: service,
	}
}

func (handler *HTTPHandlerImpl) RegisterRoutes(router *chi.Mux) {
	router.Get("/todo", handler.GetAll)
	router.Get("/todo/{id}", handler.GetByID)
	router.Post("/todo", handler.Create)
	router.Put("/todo/{id}", handler.Update)
	router.Delete("/todo/{id}", handler.Delete)
}

// GetAll - get all todo http handler
func (h *HTTPHandlerImpl) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracing.GetTracerProvider().Tracer("todoHandler").Start(r.Context(), "todoHandler.GetAll")
	defer span.End()

	qQuery := r.URL.Query().Get("q")
	pageQuery := r.URL.Query().Get("page")
	perPageQuery := r.URL.Query().Get("per_page")

	err := validator.ValidateStruct(&models.TodoListRequest{
		Keywords: &models.SearchForm{
			Keywords: qQuery,
		},
		Page:    pageQuery,
		PerPage: perPageQuery,
	})
	if err != nil {
		h.tracing.LogError(span, err)

		responseutil.ResponseErrorValidation(w, r, err)
		return
	}

	currentPage := paginationutil.CurrentPage(pageQuery)
	perPage := paginationutil.PerPage(perPageQuery)
	offset := paginationutil.Offset(currentPage, perPage)

	results, totalData, err := h.todoService.GetAll(ctx, qQuery, perPage, offset)
	if err != nil {
		h.tracing.LogError(span, err)

		responseutil.ResponseError(w, r, err)
		return
	}
	totalPages := paginationutil.TotalPage(totalData, perPage)

	responseutil.ResponseOKList(w, r, &responseutil.ResponseSuccessList{
		Data: results,
		Meta: &responseutil.Meta{
			PerPage:     perPage,
			CurrentPage: currentPage,
			TotalPage:   totalPages,
			TotalData:   totalData,
		},
	})
}

// GetByID - get todo by id http handler
func (h *HTTPHandlerImpl) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracing.GetTracerProvider().Tracer("todoHandler").Start(r.Context(), "todoHandler.GetByID")
	defer span.End()

	// Get and filter id param
	id := chi.URLParam(r, "id")

	// Get detail
	result, err := h.todoService.GetByID(ctx, id)
	if err != nil {
		h.tracing.LogError(span, err)

		if err.Error() == "not found" {
			responseutil.ResponseNotFound(w, r, "Item not found")
			return
		}

		responseutil.ResponseError(w, r, err)
		return
	}

	responseutil.ResponseOK(w, r, &responseutil.ResponseSuccess{
		Data: result,
	})

}

// Create - create todo http handler
func (h *HTTPHandlerImpl) Create(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracing.GetTracerProvider().Tracer("todoHandler").Start(r.Context(), "todoHandler.Create")
	defer span.End()

	data := &models.TodoRequest{}
	if err := render.Bind(r, data); err != nil {
		h.tracing.LogError(span, err)

		if err.Error() == "EOF" {
			responseutil.ResponseBodyError(w, r, err)
			return
		}

		responseutil.ResponseErrorValidation(w, r, err)
		return
	}

	result, err := h.todoService.Create(ctx, &models.Todo{
		Title:       data.Title,
		Description: data.Description,
	})
	if err != nil {
		h.tracing.LogError(span, err)

		responseutil.ResponseError(w, r, err)
		return
	}

	responseutil.ResponseCreated(w, r, &responseutil.ResponseSuccess{
		Data: result,
	})
}

// Update - update todo by id http handler
func (h *HTTPHandlerImpl) Update(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracing.GetTracerProvider().Tracer("todoHandler").Start(r.Context(), "todoHandler.Update")
	defer span.End()

	// Get and filter id param
	id := chi.URLParam(r, "id")

	data := &models.TodoRequest{}
	if err := render.Bind(r, data); err != nil {
		h.tracing.LogError(span, err)

		if err.Error() == "EOF" {
			responseutil.ResponseBodyError(w, r, err)
			return
		}

		responseutil.ResponseErrorValidation(w, r, err)
		return
	}

	// Edit data
	_, err := h.todoService.Update(ctx, id, &models.Todo{
		Title:       data.Title,
		Description: data.Description,
	})

	if err != nil {
		h.tracing.LogError(span, err)

		if err.Error() == "not found" {
			responseutil.ResponseNotFound(w, r, "Item not found")
			return
		}

		responseutil.ResponseError(w, r, err)
		return
	}

	responseutil.ResponseOK(w, r, &responseutil.ResponseSuccess{
		Data: responseutil.H{
			"id": id,
		},
	})
}

// Delete - delete todo by id http handler
func (h *HTTPHandlerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracing.GetTracerProvider().Tracer("todoHandler").Start(r.Context(), "todoHandler.Delete")
	defer span.End()

	// Get and filter id param
	id := chi.URLParam(r, "id")

	// Delete record
	err := h.todoService.Delete(ctx, id)
	if err != nil {
		h.tracing.LogError(span, err)

		if err.Error() == "not found" {
			responseutil.ResponseNotFound(w, r, "Item not found")
			return
		}

		responseutil.ResponseError(w, r, err)
		return
	}

	responseutil.ResponseOK(w, r, &responseutil.ResponseSuccess{
		Data: responseutil.H{
			"id": id,
		},
	})
}
