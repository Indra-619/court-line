package testutil

import (
	"context"
	"testing"
	"time"

	"github.com/Indra-619/court-line/backend/internal/domain/entity"
)

func TestFakeRefreshTokenCreateFindDelete(t *testing.T) {
	ctx := context.Background()
	repo := NewFakeRefreshTokenRepository()

	record := &entity.RefreshToken{
		UserID:    "user-1",
		TokenHash: "hash-a",
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}
	if err := repo.Create(ctx, record); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if record.ID == "" {
		t.Error("expected Create to assign an ID")
	}

	got, err := repo.FindByHash(ctx, "hash-a")
	if err != nil {
		t.Fatalf("FindByHash failed: %v", err)
	}
	if got.UserID != "user-1" {
		t.Errorf("expected UserID user-1, got %s", got.UserID)
	}

	if err := repo.DeleteByHash(ctx, "hash-a"); err != nil {
		t.Fatalf("DeleteByHash failed: %v", err)
	}
	if _, err := repo.FindByHash(ctx, "hash-a"); err == nil {
		t.Error("expected FindByHash to fail after delete")
	}
}

func TestFakeRefreshTokenFindByHashFiltersExpiredAndRevoked(t *testing.T) {
	ctx := context.Background()
	repo := NewFakeRefreshTokenRepository()

	expired := &entity.RefreshToken{
		UserID:    "user-1",
		TokenHash: "hash-expired",
		ExpiresAt: time.Now().Add(-time.Hour),
		CreatedAt: time.Now().Add(-31 * 24 * time.Hour),
	}
	revoked := &entity.RefreshToken{
		UserID:    "user-1",
		TokenHash: "hash-revoked",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
		Revoked:   true,
	}
	if err := repo.Create(ctx, expired); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if err := repo.Create(ctx, revoked); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if _, err := repo.FindByHash(ctx, "hash-expired"); err == nil {
		t.Error("expected expired token to be rejected by FindByHash")
	}
	if _, err := repo.FindByHash(ctx, "hash-revoked"); err == nil {
		t.Error("expected revoked token to be rejected by FindByHash")
	}
	if _, err := repo.FindByHash(ctx, "hash-unknown"); err == nil {
		t.Error("expected unknown token to be rejected by FindByHash")
	}
}

func TestFakeRefreshTokenFindAndDeleteByHashIsSingleUse(t *testing.T) {
	ctx := context.Background()
	repo := NewFakeRefreshTokenRepository()

	if err := repo.Create(ctx, &entity.RefreshToken{
		UserID:    "user-1",
		TokenHash: "hash-live",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if err := repo.Create(ctx, &entity.RefreshToken{
		UserID:    "user-1",
		TokenHash: "hash-expired",
		ExpiresAt: time.Now().Add(-time.Hour),
		CreatedAt: time.Now().Add(-31 * 24 * time.Hour),
	}); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if err := repo.Create(ctx, &entity.RefreshToken{
		UserID:    "user-1",
		TokenHash: "hash-revoked",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
		Revoked:   true,
	}); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	record, err := repo.FindAndDeleteByHash(ctx, "hash-live")
	if err != nil {
		t.Fatalf("FindAndDeleteByHash failed: %v", err)
	}
	if record.UserID != "user-1" {
		t.Errorf("expected UserID user-1, got %s", record.UserID)
	}

	// The same hash cannot be claimed twice.
	if _, err := repo.FindAndDeleteByHash(ctx, "hash-live"); err == nil {
		t.Error("expected second FindAndDeleteByHash to fail")
	}
	if _, err := repo.FindAndDeleteByHash(ctx, "hash-expired"); err == nil {
		t.Error("expected expired token to be rejected")
	}
	if _, err := repo.FindAndDeleteByHash(ctx, "hash-revoked"); err == nil {
		t.Error("expected revoked token to be rejected")
	}
	if _, err := repo.FindAndDeleteByHash(ctx, "hash-unknown"); err == nil {
		t.Error("expected unknown token to be rejected")
	}
}

func TestFakeRefreshTokenDeleteByUserID(t *testing.T) {
	ctx := context.Background()
	repo := NewFakeRefreshTokenRepository()

	for _, hash := range []string{"h1", "h2"} {
		if err := repo.Create(ctx, &entity.RefreshToken{
			UserID:    "user-1",
			TokenHash: hash,
			ExpiresAt: time.Now().Add(24 * time.Hour),
			CreatedAt: time.Now(),
		}); err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}
	if err := repo.Create(ctx, &entity.RefreshToken{
		UserID:    "user-2",
		TokenHash: "h3",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if err := repo.DeleteByUserID(ctx, "user-1"); err != nil {
		t.Fatalf("DeleteByUserID failed: %v", err)
	}
	if _, err := repo.FindByHash(ctx, "h1"); err == nil {
		t.Error("expected h1 to be gone")
	}
	if _, err := repo.FindByHash(ctx, "h2"); err == nil {
		t.Error("expected h2 to be gone")
	}
	if _, err := repo.FindByHash(ctx, "h3"); err != nil {
		t.Error("expected h3 (other user) to survive")
	}
}

func TestFakeRevokedTokenCreateAndExists(t *testing.T) {
	ctx := context.Background()
	repo := NewFakeRevokedTokenRepository()

	if exists, err := repo.Exists(ctx, "jti-1"); err != nil || exists {
		t.Fatalf("expected unknown jti to not exist (exists=%v err=%v)", exists, err)
	}

	if err := repo.Create(ctx, &entity.RevokedToken{
		JTI:       "jti-1",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	exists, err := repo.Exists(ctx, "jti-1")
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if !exists {
		t.Error("expected revoked jti to exist in blacklist")
	}

	// Expired blacklist entries are no longer relevant (the JWT itself
	// is expired) and must not be reported as revoked.
	if err := repo.Create(ctx, &entity.RevokedToken{
		JTI:       "jti-old",
		ExpiresAt: time.Now().Add(-time.Hour),
	}); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if exists, err := repo.Exists(ctx, "jti-old"); err != nil || exists {
		t.Errorf("expected expired blacklist entry to be dropped (exists=%v err=%v)", exists, err)
	}
}
