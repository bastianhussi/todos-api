package api

import "time"

// Todo
type Todo struct {
	tableName struct{}  `pg:"todos,alias:todo"`
	ID        int       `pg:",pk"`
	Title     string    `pg:",notnull"`
	CreatedAt time.Time `pg:"default:now()"`
	ProfileID int
	Profile   *Profile `pg:"rel:has-one,notnull"`
}
