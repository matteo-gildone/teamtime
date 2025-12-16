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

func (c Colleague) Validate() error {

	if c.Name == "" {
		return ErrMissingName
	}

	if c.City == "" {
		return ErrMissingCity
	}

	if c.Timezone == "" {
		return ErrMissingTimezone
	}

	if _, err := time.LoadLocation(c.Timezone); err != nil {
		return err
	}

	return nil
}

func NewColleague(name, city, tz string) (Colleague, error) {
	name = strings.TrimSpace(name)
	city = strings.TrimSpace(city)
	tz = strings.TrimSpace(tz)
	newColleague := Colleague{
		Name:     name,
		City:     city,
		Timezone: tz,
	}

	if err := newColleague.Validate(); err != nil {
		return Colleague{}, err
	}

	return newColleague, nil
}

type ColleagueList []Colleague

func (cl *ColleagueList) Add(newColleague Colleague) {
	*cl = append(*cl, newColleague)
}

func (cl *ColleagueList) Remove(idx int) (Colleague, error) {
	if len(*cl) == 0 {
		return Colleague{}, ErrEmptyList
	}

	if idx <= 0 || idx > len(*cl) {
		return Colleague{}, fmt.Errorf("%w: %d (must be a number between 1 and %d)", ErrorInvalidIndex, idx, len(*cl))
	}
	deleted := (*cl)[idx-1]

	copy((*cl)[idx-1:], (*cl)[idx:])
	(*cl)[len(*cl)-1] = Colleague{}
	*cl = (*cl)[:len(*cl)-1]
	return deleted, nil
}

func NewColleagues() *ColleagueList {
	return &ColleagueList{}
}
