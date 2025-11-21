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

type colleague struct {
	Name     string `json:"name"`
	City     string `json:"city"`
	Timezone string `json:"timezone"`
}

func (c *colleague) Validate() error {
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

	if _, err := time.LoadLocation(c.Timezone); err != nil {
		return err
	}

	return nil
}

type ColleagueList []colleague

func (cl *ColleagueList) Add(name, city, tz string) error {
	name = strings.TrimSpace(name)
	city = strings.TrimSpace(city)
	tz = strings.TrimSpace(tz)
	newCol := colleague{
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

func (cl *ColleagueList) Delete(i int) error {
	ls := *cl

	if len(ls) == 0 {
		return ErrEmptyList
	}

	if i <= 0 || i > len(ls) {
		return fmt.Errorf("%w: %d (must be a number between 1 and %d)", ErrorInvalidIndex, i, len(ls))
	}

	*cl = append(ls[:i-1], ls[i:]...)
	return nil
}

func NewColleagues() *ColleagueList {
	return &ColleagueList{}
}
