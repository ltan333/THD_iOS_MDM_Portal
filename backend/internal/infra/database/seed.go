package database

import (
	"context"
	"fmt"

	"github.com/thienel/tlog"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/thienel/go-backend-template/internal/domain/entity"
)

// SeedUser tạo tài khoản admin mặc định nếu chưa có người dùng nào
func SeedUser() error {
	ctx := context.Background()
	count, err := client.User.Query().Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}

	if count > 0 {
		tlog.Info("Users already exist, skipping seed", zap.Int("count", count))
		return nil
	}

	tlog.Info("Seeding default admin user...")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := client.User.Create().
		SetUsername("admin").
		SetEmail("admin@thd.vn").
		SetPassword(string(hashedPassword)).
		SetStatus(entity.UserStatusActive).
		Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to create default user: %w", err)
	}

	tlog.Info("Default admin user seeded successfully",
		zap.Uint("id", user.ID),
		zap.String("username", user.Username),
	)

	return nil
}
