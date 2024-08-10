package auth

import (
	"context"
	"sync"

	"github.com/furon-kuina/go_todo_app/entity"
)

type StoreMock struct {
	LoadFunc func(ctx context.Context, key string) (entity.UserID, error)
	SaveFunc func(ctx context.Context, key string, userID entity.UserID) error
	calls    struct {
		Load []struct {
			Ctx context.Context
			Key string
		}
		Save []struct {
			Ctx    context.Context
			Key    string
			UserID entity.UserID
		}
	}
	lockLoad sync.RWMutex
	lockSave sync.RWMutex
}

func (mock *StoreMock) Save(ctx context.Context, key string, userID entity.UserID) error {
	if mock.SaveFunc == nil {
		panic("StoreMock.SaveFunc: method is nil but Store.Save was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		Key    string
		UserID entity.UserID
	}{
		Ctx:    ctx,
		Key:    key,
		UserID: userID,
	}
	mock.lockSave.Lock()
	mock.calls.Save = append(mock.calls.Save, callInfo)
	mock.lockSave.Unlock()
	return mock.SaveFunc(ctx, key, userID)
}

func (mock *StoreMock) Load(ctx context.Context, key string) (entity.UserID, error) {
	if mock.LoadFunc == nil {
		panic("StoreMock.LoadFunc: method is nil but Store.Load was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Key string
	}{
		Ctx: ctx,
		Key: key,
	}
	mock.lockLoad.Lock()
	mock.calls.Load = append(mock.calls.Load, callInfo)
	mock.lockLoad.Unlock()
	return mock.LoadFunc(ctx, key)
}
