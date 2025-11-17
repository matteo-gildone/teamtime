package types

import (
	"errors"
	"fmt"
)

var (
	ErrorInvalidIndex = errors.New("invalid index")
	ErrEmptyList      = errors.New("colleagues list is empty")
)

type colleague struct {
	Name     string `json:"name"`
	City     string `json:"city"`
	Timezone string `json:"timezone"`
}

type ColleagueList []colleague

func (cl *ColleagueList) Add(name, city, tz string) {
	newCol := colleague{
		Name:     name,
		City:     city,
		Timezone: tz,
	}

	*cl = append(*cl, newCol)
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
