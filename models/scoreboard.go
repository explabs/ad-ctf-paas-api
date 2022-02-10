package models

type Score struct {
	Name         string                  `bson:"name" json:"name"`
	Round        int                     `bson:"round" json:"round"`
	Services     map[string]ScoreService `bson:"services" json:"services"`
	LastServices map[string]ScoreService `bson:"last_services" json:"last_services"`
	SLA          float64                 `bson:"sla" json:"sla"`
	LastSLA      float64                 `bson:"last_sla" json:"last_sla"`
	Score        float64                 `bson:"score" json:"score"`
	LastScore    float64                 `bson:"last_score" json:"last_score"`
}

type ScoreService struct {
	SLA   float64 `bson:"sla" json:"sla"`
	State int     `bson:"state" json:"state"`
}
