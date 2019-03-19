package choice

type Choice interface {
	GetWeight() int
	IsBad() bool
}

// 权重轮询
type RoundRobin struct {
	choices   []Choice
	scores    []int
	weightSum int
}

func (rr *RoundRobin) Len() int {
	return len(rr.choices)
}

func (rr *RoundRobin) AddChoice(c Choice) {
	if c.GetWeight() <= 0 {
		return
	}
	rr.choices = append(rr.choices, c)
	rr.scores = append(rr.scores, 0)
	rr.weightSum = 0 // 强制初始化
}

// 获取weight之和
func (rr *RoundRobin) GetWeightSum() int {
	if rr.weightSum <= 0 { // 重新初始化
		rr.weightSum = 0
		for i, c := range rr.choices {
			rr.weightSum += c.GetWeight()
			rr.scores[i] = 0
		}
	}
	return rr.weightSum
}

// 请先判断 RoundRobin.Len() > 0 ，避免索引超出范围
// 获得结果后判断 Choice.IsBad() ，遇到无效 Choice 需丢弃，再次轮询
func (rr *RoundRobin) GetBest() Choice {
	var best_i, mas_score int
	if rr.Len() <= 1 {
		return rr.choices[0]
	}
	sum := rr.GetWeightSum()
	for i, c := range rr.choices {
		rr.scores[i] += c.GetWeight()
		if rr.scores[i] > mas_score {
			mas_score = rr.scores[i]
			best_i = i
		}
	}
	rr.scores[best_i] -= sum
	return rr.choices[best_i]
}
