package repository_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/williamkoller/payment-system/internal/payment/domain"
	"github.com/williamkoller/payment-system/internal/payment/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm DB: %v", err)
	}
	return gormDB, mock
}

func TestPaymentRepository_Create_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := repository.NewPaymentRepository(gormDB)

	p := &domain.Payment{
		ID:             "id‑123",
		StripeID:       "stripe_1",
		Amount:         1000,
		Currency:       "USD",
		Status:         "PENDING",
		Email:          "user@example.com",
		PaymentMethod:  "card",
		IdempotencyKey: "idem‑123",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "payments"`).
		WithArgs(
			p.ID, p.StripeID, p.Amount, p.Currency, p.Status,
			p.Email, p.PaymentMethod, p.IdempotencyKey,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	created, err := repo.Save(p)
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, p.ID, created.ID)
	assert.Equal(t, p.Email, created.Email)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %v", err)
	}
}

func TestPaymentRepository_FindByID_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := repository.NewPaymentRepository(gormDB)

	id := "id‑123"
	rows := sqlmock.NewRows([]string{
		"id", "stripe_id", "amount", "currency", "status", "email", "payment_method", "idempotency_key",
	}).AddRow(
		id, "stripe_1", 1000, "USD", "PENDING", "user@example.com", "card", "idem‑123",
	)

	mock.ExpectQuery(`SELECT \* FROM "payments"`).
		WithArgs(id, sqlmock.AnyArg()).
		WillReturnRows(rows)

	found, err := repo.FindByID(id)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, id, found.ID)
	assert.Equal(t, "user@example.com", found.Email)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %v", err)
	}
}

func TestPaymentRepository_FindAll_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := repository.NewPaymentRepository(gormDB)

	id := "id‑123"
	email := "user@example.com"

	rows := sqlmock.NewRows([]string{
		"id", "stripe_id", "amount", "currency", "status", "email", "payment_method", "idempotency_key",
	}).AddRow(
		id, "stripe_1", 1000, "USD", "PENDING", email, "card", "idem‑123",
	)

	mock.ExpectQuery(`SELECT \* FROM "payments"`).
		WillReturnRows(rows)

	found, err := repo.FindAll()
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Len(t, found, 1)

	assert.Equal(t, id, found[0].ID)
	assert.Equal(t, email, found[0].Email)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %v", err)
	}
}

func TestPaymentRepository_Update_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := repository.NewPaymentRepository(gormDB)

	p := &domain.Payment{
		ID:             "id‑123",
		StripeID:       "stripe_1",
		Amount:         1000,
		Currency:       "USD",
		Status:         "COMPLETED",
		Email:          "user@example.com",
		PaymentMethod:  "card",
		IdempotencyKey: "idem‑123",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "payments"`).
		WithArgs(
			p.StripeID,
			p.Amount,
			p.Currency,
			p.Status,
			p.Email,
			p.PaymentMethod,
			p.IdempotencyKey,
			sqlmock.AnyArg(),
			p.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(p)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestPaymentRepository_Remove_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := repository.NewPaymentRepository(gormDB)

	id := "id‑123"

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "payments" WHERE id = \$1`).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Remove(id)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestPaymentRepository_Create_Error(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := repository.NewPaymentRepository(gormDB)

	p := &domain.Payment{
		ID: "id‑123",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "payments"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("db insert error"))
	mock.ExpectRollback()

	created, err := repo.Save(p)
	assert.Error(t, err)
	assert.Nil(t, created)
}

func TestPaymentRepository_FindByID_Error(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := repository.NewPaymentRepository(gormDB)

	id := "id‑123"

	mock.ExpectQuery(`SELECT \* FROM "payments"`).
		WithArgs(id, sqlmock.AnyArg()).
		WillReturnError(errors.New("db find error"))

	found, err := repo.FindByID(id)
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestPaymentRepository_FindAll_Error(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := repository.NewPaymentRepository(gormDB)

	mock.ExpectQuery(`SELECT \* FROM "payments"`).
		WillReturnError(errors.New("db find all error"))

	found, err := repo.FindAll()
	assert.Error(t, err)
	assert.Nil(t, found)
}
