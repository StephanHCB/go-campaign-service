package inmemorydb

import (
	"context"
	"fmt"
	"github.com/StephanHCB/go-campaign-service/internal/entity"
	"github.com/StephanHCB/go-campaign-service/internal/repository/database/dbrepo"
	"sync/atomic"
)

type InMemoryRepository struct {
	campaigns  map[uint]*entity.Campaign
	idSequence uint32
}

func Create() dbrepo.Repository {
	return &InMemoryRepository{}
}

func (r *InMemoryRepository) Open() {
	r.campaigns = make(map[uint]*entity.Campaign)
}

func (r *InMemoryRepository) Close() {
	r.campaigns = nil
}

func (r *InMemoryRepository) Migrate() {
	// nothing to do
}

func (r *InMemoryRepository) AddCampaign(ctx context.Context, e *entity.Campaign) (uint, error) {
	// check that no entry with this subject already exists - uniqueness constraint
	if r.countBySubjectExcludingId(e.Subject, 0) > 0 {
		return 0, fmt.Errorf("uniqueness constraint violated - an entry with this Subject already exists: %s", e.Subject)
	}

	newId := uint(atomic.AddUint32(&r.idSequence, 1))
	e.ID = newId

	// copy the entry, so later modifications won't also modify it in the simulated db
	copiedEntry := *e
	// TODO copy recipients - might use something like https://github.com/ulule/deepcopier or just implement it directly
	r.campaigns[newId] = &copiedEntry
	return newId, nil
}

func (r *InMemoryRepository) UpdateCampaign(ctx context.Context, e *entity.Campaign) error {
	if _, ok := r.campaigns[e.ID]; ok {
		if r.countBySubjectExcludingId(e.Subject, e.ID) > 0 {
			return fmt.Errorf("uniqueness constraint violated - another entry with this Subject already exists: %s", e.Subject)
		}

		// copy the entry, so later modifications won't also modify it in the simulated db
		copiedEntry := *e
		// TODO copy recipients - might use something like https://github.com/ulule/deepcopier or just implement it directly
		r.campaigns[e.ID] = &copiedEntry
		return nil
	} else {
		return fmt.Errorf("cannot update campaign %d - not present", e.ID)
	}
}

func (r *InMemoryRepository) GetCampaignById(ctx context.Context, id uint) (*entity.Campaign, error) {
	if entry, ok := r.campaigns[id]; ok {
		// copy the entry, so later modifications won't also modify it in the simulated db
		copiedEntry := *entry
		// TODO copy recipients - might use something like https://github.com/ulule/deepcopier or just implement it directly
		return &copiedEntry, nil
	} else {
		return &entity.Campaign{}, fmt.Errorf("cannot get campaign %d - not present", id)
	}
}

func (r *InMemoryRepository) CountCampaignsBySubject(ctx context.Context, subject string) (uint, error) {
	return r.countBySubjectExcludingId(subject, 0), nil
}

// helper to implement uniqueness of Subject

func (r *InMemoryRepository) countBySubjectExcludingId(subject string, ignoreId uint) uint {
	var count uint
	for _, v := range r.campaigns {
		if subject == v.Subject && ignoreId != v.ID {
			count++
		}
	}
	return count
}
