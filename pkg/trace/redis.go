package trace

type OptionalType string

func (o OptionalType) String() string {
	return string(o)
}

const (
	GetOptional      OptionalType = "Get"
	SetOptional      OptionalType = "Set"
	DeleteOptional   OptionalType = "Delete"
	IncreaseOptional OptionalType = "Increase"
)

type Redis struct {
	Timestamp   int64        `json:"timestamp"`       // 时间戳
	Optional    OptionalType `json:"handle"`          // 操作
	Key         string       `json:"key"`             // Key
	Value       string       `json:"value,omitempty"` // Value
	TTL         float64      `json:"ttl,omitempty"`   // 超时时长(单位秒)
	CostSeconds float64      `json:"cost_seconds"`    // 执行时间(单位秒)
}
