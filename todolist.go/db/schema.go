package db

// schema.go provides data models in DB
import (
	"time"
)

// Task corresponds to a row in `tasks` table
type Task struct {
	ID        uint64    `db:"id"`
	Title     string    `db:"title"`
	CreatedAt time.Time `db:"created_at"`
	IsDone    bool      `db:"is_done"`
}
// Account corresponds to a row in `accounts` table
type Account struct{
    ID        uint64     `db:"id"`
    Name      string     `db:"name"`
	Password  string     `db:"password"`
	Birth     string     `db:"birth"`
}

// Task_owner corresponds to a row in `task_owners` table
type Task_owner struct{
	User_ID        uint64 `db:"user_id"`
	Task_ID        uint64 `db:"task_id"`
}

