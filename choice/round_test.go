package choice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var premierLeague = map[string]int{
	"Chelsea":   20,
	"Liverpool": 30,
	"ManUTD":    50,
}

type Team struct {
	Name   string
	Weight int
}

func (t Team) GetWeight() int {
	return t.Weight
}

func (t Team) IsBad() bool {
	return false
}

// 普通测试
func TestChoice(t *testing.T) {
	result := []string{"ManUTD", "Liverpool", "Chelsea", "ManUTD", "Liverpool",
		"ManUTD", "ManUTD", "Chelsea", "Liverpool", "ManUTD", "ManUTD",
		"Liverpool", "Chelsea", "ManUTD", "Liverpool", "ManUTD", "ManUTD",
		"Chelsea", "Liverpool", "ManUTD"}
	teams := new(RoundRobin)
	for name, weight := range premierLeague {
		teams.AddChoice(Team{Name: name, Weight: weight})
	}
	for i := 1; i <= 20; i++ {
		team := teams.GetBest().(Team)
		t.Log(i, team.Name, team.GetWeight())
		result = append(result, team.Name)
		assert.Equal(t, result[i-1], team.Name)
	}
}
