package repo

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"medrussia_news_bot/internal/pkg/postgres"

	"github.com/jackc/pgx/v5"
)

type Repo struct {
	client postgres.Client
	logger *slog.Logger
}

// UserDialog представляет запись из таблицы users_dialog
type UserDialog struct {
	ID                 int64         `sql:"id"`
	UserID             int64         `sql:"user_id"`
	LastAdminMessageID sql.NullInt64 `sql:"last_admin_message_id"`
	LastUserMessageID  sql.NullInt64 `sql:"last_user_message_id"`
	Available          bool          `sql:"available"`
}

// GetUser получает запись пользователя по ID
func (r *Repo) GetUser(ctx context.Context, userID int64) (*UserDialog, error) {
	sql := `select id, user_id, last_admin_message_id, last_user_message_id, available 
				from users_dialog where user_id = $1`

	var user UserDialog
	err := r.client.QueryRow(ctx, sql, userID).Scan(
		&user.ID,
		&user.UserID,
		&user.LastAdminMessageID,
		&user.LastUserMessageID,
		&user.Available,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// CreateUser создает новую запись пользователя
func (r *Repo) CreateUser(ctx context.Context, userID int64) error {
	sql := `insert into users_dialog (user_id, available)
				values ($1, false)  
				on conflict (user_id) do nothing
	`

	_, err := r.client.Exec(ctx, sql, userID)

	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}

// UpdateLastAdminMessage обновляет last_admin_message_id
func (r *Repo) UpdateLastAdminMessage(ctx context.Context, userID, messageID int64) error {
	sql := `update users_dialog set last_admin_message_id = $1 where user_id = $2`
	_, err := r.client.Exec(ctx, sql, messageID, userID)
	return err
}

// UpdateLastUserMessage обновляет last_user_message_id
func (r *Repo) UpdateLastUserMessage(ctx context.Context, userID, messageID int64) error {
	sql := `update users_dialog set last_user_message_id = $1 where user_id = $2`
	_, err := r.client.Exec(ctx, sql, messageID, userID)
	return err
}

// UpdateAvailable обновляет available
func (r *Repo) UpdateAvailable(ctx context.Context, userID int64, available bool) error {
	sql := `update users_dialog set available = $1 where user_id = $2`
	_, err := r.client.Exec(ctx, sql, available, userID)
	return err
}

func NewRepo(client postgres.Client, logger *slog.Logger) *Repo {
	return &Repo{
		client: client,
		logger: logger,
	}
}
