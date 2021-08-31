package handler

import (
	"bytes"
	"encoding/json"
	"github.com/ELOATS/studyMem/account/model"
	"github.com/ELOATS/studyMem/account/model/apperrors"
	"github.com/ELOATS/studyMem/account/model/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Signup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Email and Password Required", func(t *testing.T) {
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Signup", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*model.User")).Return(nil)

		rr := httptest.NewRecorder()

		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// create a request body
		reqBody, err := json.Marshal(gin.H{
			"email": "",
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockUserService.AssertNotCalled(t, "Signup")
	})

	t.Run("Password is too short", func(t *testing.T) {
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Signup", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*model.User")).Return(nil)

		rr := httptest.NewRecorder()

		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		reqBody, err := json.Marshal(gin.H{
			"email":    "2514307815@qq.com",
			"password": "123",
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockUserService.AssertNotCalled(t, "Signup")
	})

	t.Run("Email has existed", func(t *testing.T) {
		u := &model.User{
			Email:    "2514307815@qq.com",
			Password: "root123",
		}

		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Signup", mock.AnythingOfType("*gin.Context"), u).Return(apperrors.NewConflict("User already exists.", u.Email))

		rr := httptest.NewRecorder()

		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		reqBody, err := json.Marshal(gin.H{
			"email":    u.Email,
			"password": u.Password,
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)
		assert.Equal(t, http.StatusConflict, rr.Code)
		mockUserService.AssertExpectations(t)
	})

	t.Run("Successful token creation", func(t *testing.T) {
		u := &model.User{
			Email:    "2514307815@qq.com",
			Password: "root123",
		}

		mockTokenResp := &model.TokenPair{
			IDToken:      "idToken",
			RefreshToken: "refreshToken",
		}

		mockUserService, mockTokenService := new(mocks.MockUserService), new(mocks.MockTokenService)
		mockUserService.On("Signup", mock.AnythingOfType("*gin.Context"), u).Return(nil)
		mockTokenService.On("NewPairFromUser", mock.AnythingOfType("*gin.Context"), u, "").Return(mockTokenResp, nil)

		// a response record for getting written http response
		rr := httptest.NewRecorder()

		router := gin.Default()

		NewHandler(&Config{
			R:            router,
			UserService:  mockUserService,
			TokenService: mockTokenService,
		})

		reqBody, err := json.Marshal(gin.H{
			"email":    u.Email,
			"password": u.Password,
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"tokens": mockTokenResp,
		})
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockUserService.AssertExpectations(t)
		mockTokenService.AssertExpectations(t)
	})

	t.Run("Failed token creation", func(t *testing.T) {
		u := &model.User{
			Email:    "2514307815@qq.com",
			Password: "root123",
		}

		mockErrorResponse := apperrors.NewInternal()

		mockUserService, mockTokenService := new(mocks.MockUserService), new(mocks.MockTokenService)
		mockUserService.On("Signup", mock.AnythingOfType("*gin.Context"), u).Return(nil)
		mockTokenService.On("NewPairFromUser", mock.AnythingOfType("*gin.Context"), u, "").Return(nil, mockErrorResponse)

		rr := httptest.NewRecorder()

		router := gin.Default()

		NewHandler(&Config{
			R:            router,
			UserService:  mockUserService,
			TokenService: mockTokenService,
		})

		reqBody, err := json.Marshal(gin.H{
			"email":    u.Email,
			"password": u.Password,
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockErrorResponse,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockErrorResponse.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockUserService.AssertExpectations(t)
		mockTokenService.AssertExpectations(t)
	})
}
