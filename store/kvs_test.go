package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/furon-kuina/go_todo_app/entity"
	"github.com/furon-kuina/go_todo_app/testutil"
)

func TestKVS_Save(t *testing.T) {
	t.Parallel()

	cli := testutil.OpenRedisForTest(t)
	sut := &KVS{Cli: cli}
	key := "TestKVS_Store"
	uid := entity.UserID(42)
	ctx := context.Background()
	t.Cleanup(func() {
		cli.Del(ctx, key)
	})
	if err := sut.Save(ctx, key, uid); err != nil {
		t.Errorf("want no error, but got %v", err)
	}
}

func TestKVS_Load(t *testing.T) {
	t.Parallel()

	cli := testutil.OpenRedisForTest(t)
	sut := &KVS{Cli: cli}
	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		key := "TestKVS_Load_OK"
		uid := entity.UserID(42)
		ctx := context.Background()
		cli.Set(ctx, key, int64(uid), 30*time.Minute)
		t.Cleanup(func() {
			cli.Del(ctx, key)
		})
		got, err := sut.Load(ctx, key)
		if err != nil {
			t.Fatalf("want nil error, got %v", err)
		}
		if got != uid {
			t.Errorf("want %d, got %d", uid, got)
		}
	})

	t.Run("notFound", func(t *testing.T) {
		t.Parallel()

		key := "TestKVS_Save_notFound"
		ctx := context.Background()
		got, err := sut.Load(ctx, key)
		if err == nil || !errors.Is(err, ErrNotFound) {
			t.Errorf("want %v, got %v(value = %d)", ErrNotFound, err, got)
		}
	})
}
