package types

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrMissingName     = errors.New("'Name' must not be empty")
	ErrMissingCity     = errors.New("'City' must not be empty")
	ErrMissingTimezone = errors.New("'Timezone' must not be empty")
	ErrorInvalidIndex  = errors.New("invalid index")
	ErrEmptyList       = errors.New("colleagues list is empty")
)

type Colleague struct {
	Name     string `json:"name"`
	City     string `json:"city"`
	Timezone string `json:"timezone"`
}

func (c *Colleague) Validate() error {
	name := strings.TrimSpace(c.Name)
	city := strings.TrimSpace(c.City)
	tz := strings.TrimSpace(c.Timezone)

	if name == "" {
		return ErrMissingName
	}

	if city == "" {
		return ErrMissingCity
	}

	if tz == "" {
		return ErrMissingTimezone
	}

	if _, err := time.LoadLocation(tz); err != nil {
		return err
	}

	return nil
}

type ColleagueList []Colleague

func (cl *ColleagueList) Add(name, city, tz string) error {
	name = strings.TrimSpace(name)
	city = strings.TrimSpace(city)
	tz = strings.TrimSpace(tz)
	newCol := Colleague{
		Name:     name,
		City:     city,
		Timezone: tz,
	}

	if err := newCol.Validate(); err != nil {
		return err
	}

	*cl = append(*cl, newCol)
	return nil
}

func (cl *ColleagueList) Delete(idx int) (Colleague, error) {
	ls := *cl

	if len(ls) == 0 {
		return Colleague{}, ErrEmptyList
	}

	if idx <= 0 || idx > len(ls) {
		return Colleague{}, fmt.Errorf("%w: %d (must be a number between 1 and %d)", ErrorInvalidIndex, idx, len(ls))
	}
	deleted := ls[idx-1]

	copy(ls[idx-1:], ls[idx:])
	ls[len(ls)-1] = Colleague{}
	*cl = ls[:len(ls)-1]
	return deleted, nil
}

func NewColleagues() *ColleagueList {
	return &ColleagueList{}
}
