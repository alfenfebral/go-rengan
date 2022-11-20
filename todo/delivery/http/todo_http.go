package todo_http

import (
	"net/http"

	pkg_tracing "go-rengan/pkg/tracing"
	pkg_validator "go-rengan/pkg/validator"
	"go-rengan/todo/models"
	"go-rengan/todo/service"
	"go-rengan/utils"
	response "go-rengan/utils/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type TodoHTTPHandler interface {
	RegisterRoutes(router *chi.Mux)
	GetAll(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type TodoHTTPHandlerImpl struct {
	tracing     pkg_tracing.Tracing
	todoService service.TodoService
}

// NewTodoHTTPHandler - make http handler
func NewTodoHTTPHandler(tracing pkg_tracing.Tracing, service service.TodoService) TodoHTTPHandler {
	return &TodoHTTPHandlerImpl{
		tracing:     tracing,
		todoService: service,
	}
}

func (handler *TodoHTTPHandlerImpl) RegisterRoutes(router *chi.Mux) {
	router.Get("/todo", handler.GetAll)
	router.Get("/todo/{id}", handler.GetByID)
	router.Post("/todo", handler.Create)
	router.Put("/todo/{id}", handler.Update)
	router.Delete("/todo/{id}", handler.Delete)
}

// GetAll - get all todo http handler
func (handler *TodoHTTPHandlerImpl) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, span := handler.tracing.GetTracerProvider().Tracer("todoHandler").Start(r.Context(), "todoHandler.GetAll")
	defer span.End()

	qQuery := r.URL.Query().Get("q")
	pageQuery := r.URL.Query().Get("page")
	perPageQuery := r.URL.Query().Get("per_page")

	err := pkg_validator.ValidateStruct(&models.TodoListRequest{
		Keywords: &models.SearchForm{
			Keywords: qQuery,
		},
		Page:    pageQuery,
		PerPage: perPageQuery,
	})
	if err != nil {
		handler.tracing.LogError(span, err)

		response.ResponseErrorValidation(w, r, err)
		return
	}

	currentPage := utils.CurrentPage(pageQuery)
	perPage := utils.PerPage(perPageQuery)
	offset := utils.Offset(currentPage, perPage)

	results, totalData, err := handler.todoService.GetAll(ctx, qQuery, perPage, offset)
	if err != nil {
		handler.tracing.LogError(span, err)

		response.ResponseError(w, r, err)
		return
	}
	totalPages := utils.TotalPage(totalData, perPage)

	response.ResponseOKList(w, r, &response.ResponseSuccessList{
		Data: results,
		Meta: &response.Meta{
			PerPage:     perPage,
			CurrentPage: currentPage,
			TotalPage:   totalPages,
			TotalData:   totalData,
		},
	})
}

// GetByID - get todo by id http handler
func (handler *TodoHTTPHandlerImpl) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := handler.tracing.GetTracerProvider().Tracer("todoHandler").Start(r.Context(), "todoHandler.GetByID")
	defer span.End()

	// Get and filter id param
	id := chi.URLParam(r, "id")

	// Get detail
	result, err := handler.todoService.GetByID(ctx, id)
	if err != nil {
		handler.tracing.LogError(span, err)

		if err.Error() == "not found" {
			response.ResponseNotFound(w, r, "Item not found")
			return
		}

		response.ResponseError(w, r, err)
		return
	}

	response.ResponseOK(w, r, &response.ResponseSuccess{
		Data: result,
	})

}

// Create - create todo http handler
func (handler *TodoHTTPHandlerImpl) Create(w http.ResponseWriter, r *http.Request) {
	ctx, span := handler.tracing.GetTracerProvider().Tracer("todoHandler").Start(r.Context(), "todoHandler.Create")
	defer span.End()

	data := &models.TodoRequest{}
	if err := render.Bind(r, data); err != nil {
		handler.tracing.LogError(span, err)

		if err.Error() == "EOF" {
			response.ResponseBodyError(w, r, err)
			return
		}

		response.ResponseErrorValidation(w, r, err)
		return
	}

	result, err := handler.todoService.Create(ctx, &models.Todo{
		Title:       data.Title,
		Description: data.Description,
	})
	if err != nil {
		handler.tracing.LogError(span, err)

		response.ResponseError(w, r, err)
		return
	}

	response.ResponseCreated(w, r, &response.ResponseSuccess{
		Data: result,
	})
}

// Update - update todo by id http handler
func (handler *TodoHTTPHandlerImpl) Update(w http.ResponseWriter, r *http.Request) {
	ctx, span := handler.tracing.GetTracerProvider().Tracer("todoHandler").Start(r.Context(), "todoHandler.Update")
	defer span.End()

	// Get and filter id param
	id := chi.URLParam(r, "id")

	data := &models.TodoRequest{}
	if err := render.Bind(r, data); err != nil {
		handler.tracing.LogError(span, err)

		if err.Error() == "EOF" {
			response.ResponseBodyError(w, r, err)
			return
		}

		response.ResponseErrorValidation(w, r, err)
		return
	}

	// Edit data
	_, err := handler.todoService.Update(ctx, id, &models.Todo{
		Title:       data.Title,
		Description: data.Description,
	})

	if err != nil {
		handler.tracing.LogError(span, err)

		if err.Error() == "not found" {
			response.ResponseNotFound(w, r, "Item not found")
			return
		}

		response.ResponseError(w, r, err)
		return
	}

	response.ResponseOK(w, r, &response.ResponseSuccess{
		Data: response.H{
			"id": id,
		},
	})
}

// Delete - delete todo by id http handler
func (handler *TodoHTTPHandlerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, span := handler.tracing.GetTracerProvider().Tracer("todoHandler").Start(r.Context(), "todoHandler.Delete")
	defer span.End()

	// Get and filter id param
	id := chi.URLParam(r, "id")

	// Delete record
	err := handler.todoService.Delete(ctx, id)
	if err != nil {
		handler.tracing.LogError(span, err)

		if err.Error() == "not found" {
			response.ResponseNotFound(w, r, "Item not found")
			return
		}

		response.ResponseError(w, r, err)
		return
	}

	response.ResponseOK(w, r, &response.ResponseSuccess{
		Data: response.H{
			"id": id,
		},
	})
}
