package service

import (
	"fmt"
	"strings"

	"github.com/matteo-gildone/teamtime/internals/storage"
	"github.com/matteo-gildone/teamtime/internals/types"
)

type ColleagueService struct {
	manager *storage.Manager
}

func NewColleagueService(m *storage.Manager) *ColleagueService {
	return &ColleagueService{
		manager: m,
	}
}

func (s *ColleagueService) AddColleague(name, city, tz string) (types.Colleague, error) {
	cl, err := s.manager.Load()
	if err != nil {
		return types.Colleague{}, fmt.Errorf("failed to load colleagues: %w", err)
	}

	colleague, err := types.NewColleague(name, city, tz)
	if err != nil {
		return types.Colleague{}, fmt.Errorf("invalid colleague data: %w", err)
	}

	if err := cl.Add(colleague); err != nil {
		return types.Colleague{}, fmt.Errorf("failed to add colleague to list: %w", err)
	}

	if err := s.manager.Save(cl); err != nil {
		return types.Colleague{}, fmt.Errorf("colleague added to list but failed to save: %w", err)
	}

	return colleague, nil
}

func (s *ColleagueService) RemoveColleague(idx int) (types.Colleague, error) {
	cl, err := s.manager.Load()
	if err != nil {
		return types.Colleague{}, fmt.Errorf("failed to load colleagues: %w", err)
	}

	removed, err := cl.Remove(idx)
	if err != nil {
		return types.Colleague{}, fmt.Errorf("failed to remove colleagues: %w", err)
	}

	if err := s.manager.Save(cl); err != nil {
		return types.Colleague{}, fmt.Errorf("colleague removed to list but failed to save: %w", err)
	}

	return removed, nil
}

func (s *ColleagueService) AllColleagues() ([]types.Colleague, error) {
	cl, err := s.manager.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load colleagues: %w", err)
	}

	return *cl, nil
}

func (s *ColleagueService) FindColleague(name string) ([]types.Colleague, error) {
	cl, err := s.manager.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load colleagues: %w", err)
	}

	var results types.ColleagueList
	for _, c := range *cl {
		if strings.EqualFold(c.Name, name) {
			results = append(results, c)
		}
	}

	return results, nil
}
