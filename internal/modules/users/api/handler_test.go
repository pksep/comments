package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	sharedmodel "github.com/pksep/location_search_server/internal/modules/shared/model" // <--- здесь BaseModel
	"github.com/pksep/location_search_server/internal/modules/users/model"              // <--- здесь модель User
	"github.com/stretchr/testify/assert"
)

// Mock репозиторий
type MockUserRepo struct{}

func (m *MockUserRepo) Create(ctx context.Context, user *model.User) (*model.User, error) {
	user.ID = "mock-id"
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return user, nil
}

func (m *MockUserRepo) List(ctx context.Context) ([]model.User, error) {
	return []model.User{
		{ID: "1", Initials: "AB", BaseModel: sharedmodel.BaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}},
	}, nil
}

func (m *MockUserRepo) GetByID(ctx context.Context, id string) (*model.User, error) {
	if id == "notfound" {
		return nil, assert.AnError
	}
	return &model.User{ID: id, Initials: "AB", BaseModel: sharedmodel.BaseModel{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}}, nil
}

func (m *MockUserRepo) Update(ctx context.Context, user *model.User) (*model.User, error) {
	user.UpdatedAt = time.Now()
	return user, nil
}

func (m *MockUserRepo) Delete(ctx context.Context, id string) error {
	if id == "error" {
		return assert.AnError
	}
	return nil
}

// --- Тесты ---

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := NewUserHandler(&MockUserRepo{})
	handler.RegisterRoutes(router.Group("/"))
	return router
}

func TestCreateUser(t *testing.T) {
	router := setupRouter()
	body := `{"initials":"AB"}`
	req, _ := http.NewRequest("POST", "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "mock-id")
}

func TestListUsers(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "AB")
}

func TestGetUser(t *testing.T) {
	router := setupRouter()

	// Существующий пользователь
	req, _ := http.NewRequest("GET", "/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "AB")

	// Не найден
	req2, _ := http.NewRequest("GET", "/users/notfound", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusNotFound, w2.Code)
}

func TestUpdateUser(t *testing.T) {
	router := setupRouter()
	body := `{"initials":"CD"}`
	req, _ := http.NewRequest("PUT", "/users/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "CD")
}

func TestDeleteUser(t *testing.T) {
	router := setupRouter()

	// Успешное удаление
	req, _ := http.NewRequest("DELETE", "/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// Ошибка удаления
	req2, _ := http.NewRequest("DELETE", "/users/error", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusInternalServerError, w2.Code)
}
