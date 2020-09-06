package model

type Interaction struct {
	ID string
}

func NewInteraction(id string) *Interaction {
	return &Interaction{
		ID: id,
	}
}

func (Interaction) TableName() string { return "interactions" }

// func (Interaction) CreateTableDDL() string {
// 	return `CREATE TABLE IF NOT EXISTS interactions(` +
// 		`"id"   INTEGER  NOT NULL PRIMARY KEY` +
// 		`);`
// }
