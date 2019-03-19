package choice

import "testing"

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
	teams := new(RoundRobin)
	for name, weight := range premierLeague {
		teams.AddChoice(Team{Name: name, Weight: weight})
	}
	for i := 1; i <= 20; i++ {
		team := teams.GetBest().(Team)
		t.Log(i, team.Name, team.GetWeight())
	}
}
