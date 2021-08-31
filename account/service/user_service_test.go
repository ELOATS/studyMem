package service

import (
	"context"
	"fmt"
	"github.com/ELOATS/studyMem/account/model"
	"github.com/ELOATS/studyMem/account/model/apperrors"
	"github.com/ELOATS/studyMem/account/model/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestUserService_Get(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUserResp := &model.User{
			UID:   uid,
			Email: "2514307815@qq.com",
			Name:  "Cha",
		}

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})
		mockUserRepository.On("FindByID", mock.Anything, uid).Return(mockUserResp, nil)

		u, err := us.Get(context.TODO(), uid)

		assert.NoError(t, err)
		assert.Equal(t, u, mockUserResp)
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})
		mockUserRepository.On("FindByID", mock.Anything, uid).Return(nil, fmt.Errorf("Some error down the call chain.\n"))

		u, err := us.Get(context.TODO(), uid)
		assert.Nil(t, u)
		assert.Error(t, err)
		mockUserRepository.AssertExpectations(t)
	})
}

func TestUserService_Signup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUser := &model.User{
			Email:    "2514307815@qq.com",
			Password: "root123",
		}

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})

		mockUserRepository.On("Create", mock.AnythingOfType("*context.emptyCtx"), mockUser).
			Run(func(args mock.Arguments) {
				userArg := args.Get(1).(*model.User)
				userArg.UID = uid
			}).Return(nil)

		err := us.Signup(context.TODO(), mockUser)
		assert.NoError(t, err)
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockUser := &model.User{
			Email:    "2514307815@qq.com",
			Password: "root123",
		}

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})

		mockError := apperrors.NewConflict("email", mockUser.Email)

		mockUserRepository.On("Create", mock.AnythingOfType("*context.emptyCtx"), mockUser).Return(mockError)

		err := us.Signup(context.TODO(), mockUser)
		assert.EqualError(t, err, mockError.Error())
		mockUserRepository.AssertExpectations(t)
	})
}
