package api

// Task
type Task struct {
	tableName struct{} `pg:"tasks,alias:task"`
	ID        int      `pg:",pk"`
	Title     string   `pg:",notnull"`
	Done      bool     `pg:"default:FALSE"`
	TodoID    int
	Todo      *Todo `pg:"rel:has-one,notnull"`
}
