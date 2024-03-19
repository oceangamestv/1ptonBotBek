package valueobject

type MineManyResult struct {
	Balance   int64 `json:"balance"`
	Mined     int32 `json:"mined"`
	NewEnergy int32 `json:"newEnergy"`
}
