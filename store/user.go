package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/furon-kuina/go_todo_app/entity"
	"github.com/go-sql-driver/mysql"
)

func (r *Repository) RegisterUser(ctx context.Context, db Execer, u *entity.User) error {
	u.Created = r.Clocker.Now()
	u.Modified = r.Clocker.Now()
	sql := `INSERT INTO user (
	name, password, role, created, modified
	) VALUES (?, ?, ?, ?, ?)`
	result, err := db.ExecContext(ctx, sql, u.Name, u.Password, u.Role, u.Created, u.Modified)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeMySQLDuplicateEntry {
			return fmt.Errorf("duplicated user not allowed: %w", ErrAlreadyEntry)
		}
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = entity.UserID(id)
	return nil
}
