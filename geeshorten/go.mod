module geeshorten

go 1.15

replace ./db => ./db

require (
	github.com/go-sql-driver/mysql v1.6.0
	gorm.io/driver/mysql v1.0.5
	gorm.io/gorm v1.21.7
)
