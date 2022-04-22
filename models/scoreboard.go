package models

type Score struct {
	Name         string                  `bson:"name" json:"name"`
	Login        string                  `bson:"login"`
	Round        int                     `bson:"round" json:"round"`
	Services     map[string]ScoreService `bson:"services" json:"services"`
	LastServices map[string]ScoreService `bson:"last_services" json:"last_services"`
	SLA          float64                 `bson:"sla" json:"sla"`
	LastSLA      float64                 `bson:"last_sla" json:"last_sla"`
	Score        float64                 `bson:"score" json:"score"`
	LastScore    float64                 `bson:"last_score" json:"last_score"`
}

type ScoreService struct {
	SLA     float64 `bson:"sla" json:"sla"`
	State   int     `bson:"state" json:"state"`
	Gained  int     `bson:"gained" json:"gained"`
	Lost    int     `bson:"lost" json:"lost"`
	HP      int     `bson:"hp" json:"hp"`
	TotalHP int     `bson:"total_hp" json:"total_hp"`
	Cost    int     `bson:"cost" json:"cost"`
}

type OutputScoreboard struct {
	Name         string                  `json:"name"`
	Login        string                  `json:"login"`
	Place        int                     `json:"place"`
	ChangedPlace int                     `json:"changed_place"`
	Round        int                     `json:"round"`
	Services     map[string]ScoreService `json:"services"`
	SLA          float64                 `json:"sla"`
	Score        float64                 `json:"score"`
}
