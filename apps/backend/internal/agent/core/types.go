package core

type StrategyOutput struct {
	Action      string  `json:"action"`
	DelayHours  int     `json:"delay_hours"`
	Reasoning   string  `json:"reasoning"`
	Confidence  float32 `json:"confidence"`
}
