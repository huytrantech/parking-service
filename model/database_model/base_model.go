package database_model

import "time"

type BaseEntity struct {
	Id          int        `db:"id"`
	CreatedDate time.Time  `db:"created_date"`
	CreatedName string     `db:"created_name"`
	UpdatedDate *time.Time `db:"updated_date"`
	UpdatedName *string    `db:"updated_name"`
}
