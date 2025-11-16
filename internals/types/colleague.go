package types

import "fmt"

type colleague struct {
	Name     string `json:"file"`
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
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("item %d does not exist", i)
	}

	*cl = append(ls[:i-1], ls[i:]...)
	return nil
}

func NewColleagues() *ColleagueList {
	return &ColleagueList{}
}
