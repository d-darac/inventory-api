package services

import (
	"context"
	"time"

	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

type Price struct {
	ID        uuid.UUID         `json:"id"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Currency  database.Currency `json:"currency"`
	Amount    int32             `json:"amount"`
	Item      api.Expandable    `json:"item"`
}

type PricesService struct {
	Db *database.Queries
}

func NewPricesService(db *database.Queries) *PricesService {
	return &PricesService{
		Db: db,
	}
}

func (s *PricesService) Create(cpp database.CreatePriceParams) (*Price, error) {
	pr, err := s.Db.CreatePrice(context.Background(), cpp)
	if err != nil {
		return nil, err
	}
	return &Price{
		ID:        pr.ID,
		CreatedAt: pr.CreatedAt,
		UpdatedAt: pr.UpdatedAt,
		Currency:  pr.Currency,
		Amount:    pr.Amount,
		Item:      api.Expandable{ID: uuid.NullUUID{UUID: pr.Item, Valid: true}},
	}, nil
}

func (s *PricesService) Delete(dpp database.DeletePriceParams) error {
	return s.Db.DeletePrice(context.Background(), dpp)
}

func (s *PricesService) Get(gpp database.GetPriceParams) (*Price, error) {
	pr, err := s.Db.GetPrice(context.Background(), gpp)
	if err != nil {
		return nil, err
	}
	return &Price{
		ID:        pr.ID,
		CreatedAt: pr.CreatedAt,
		UpdatedAt: pr.UpdatedAt,
		Currency:  pr.Currency,
		Amount:    pr.Amount,
		Item:      api.Expandable{ID: uuid.NullUUID{UUID: pr.Item, Valid: true}},
	}, nil
}

func (s *PricesService) List(lpp database.ListPricesParams) (prices []*Price, hasMore bool, err error) {
	prs, err := s.Db.ListPrices(context.Background(), lpp)
	if err != nil {
		return
	}
	if lpp.Limit.Valid {
		hasMore = len(prs) > int(lpp.Limit.Int32)
	} else {
		hasMore = len(prs) > 10
	}
	if hasMore {
		if lpp.EndingBefore.Valid {
			prs = prs[1:]
		} else {
			prs = prs[:len(prs)-1]
		}
	}
	for _, p := range prs {
		prices = append(prices, &Price{
			ID:        p.ID,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
			Currency:  p.Currency,
			Amount:    p.Amount,
			Item:      api.Expandable{ID: uuid.NullUUID{UUID: p.Item, Valid: true}},
		})
	}
	return
}

func (s *PricesService) Update(upp database.UpdatePriceParams) (*Price, error) {
	pr, err := s.Db.UpdatePrice(context.Background(), upp)
	if err != nil {
		return nil, err
	}
	return &Price{
		ID:        pr.ID,
		CreatedAt: pr.CreatedAt,
		UpdatedAt: pr.UpdatedAt,
		Currency:  pr.Currency,
		Amount:    pr.Amount,
		Item:      api.Expandable{ID: uuid.NullUUID{UUID: pr.Item, Valid: true}},
	}, nil
}
