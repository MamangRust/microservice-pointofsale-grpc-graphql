package repository

import (
	"context"

	db "github.com/MamangRust/microservice-point-of-sale-cashier/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/cashier_errors"
)

type cashierCommandRepository struct {
	db *db.Queries
}

func NewCashierCommandRepository(db *db.Queries) CashierCommandRepository {
	return &cashierCommandRepository{
		db: db,
	}
}

func (r *cashierCommandRepository) CreateCashier(ctx context.Context, request *requests.CreateCashierRequest) (*db.Cashier, error) {
	req := db.CreateCashierParams{
		MerchantID: int32(request.MerchantID),
		UserID:     int32(request.UserID),
		Name:       request.Name,
	}

	cashier, err := r.db.CreateCashier(ctx, req)
	if err != nil {
		return nil, cashier_errors.ErrCreateCashier
	}

	return cashier, nil
}

func (r *cashierCommandRepository) UpdateCashier(ctx context.Context, request *requests.UpdateCashierRequest) (*db.Cashier, error) {
	req := db.UpdateCashierParams{
		CashierID: int32(*request.CashierID),
		Name:      request.Name,
	}

	res, err := r.db.UpdateCashier(ctx, req)
	if err != nil {
		return nil, cashier_errors.ErrUpdateCashier
	}

	return res, nil
}

func (r *cashierCommandRepository) TrashedCashier(ctx context.Context, cashier_id int) (*db.Cashier, error) {
	res, err := r.db.TrashCashier(ctx, int32(cashier_id))
	if err != nil {
		return nil, cashier_errors.ErrTrashedCashier
	}

	return res, nil
}

func (r *cashierCommandRepository) RestoreCashier(ctx context.Context, cashier_id int) (*db.Cashier, error) {
	res, err := r.db.RestoreCashier(ctx, int32(cashier_id))
	if err != nil {
		return nil, cashier_errors.ErrRestoreCashier
	}

	return res, nil
}

func (r *cashierCommandRepository) DeleteCashierPermanent(ctx context.Context, cashier_id int) (bool, error) {
	err := r.db.DeleteCashierPermanently(ctx, int32(cashier_id))
	if err != nil {
		return false, cashier_errors.ErrDeleteCashierPermanent
	}

	return true, nil
}

func (r *cashierCommandRepository) RestoreAllCashier(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllCashiers(ctx)
	if err != nil {
		return false, cashier_errors.ErrRestoreAllCashiers
	}

	return true, nil
}

func (r *cashierCommandRepository) DeleteAllCashierPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentCashiers(ctx)
	if err != nil {
		return false, cashier_errors.ErrDeleteAllCashiersPermanent
	}

	return true, nil
}
