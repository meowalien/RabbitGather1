module rabbitgather_database_init

go 1.16

require (
	github.com/meowalien/rabbitgather-lib v0.2.5
	gorm.io/gorm v1.21.15 // indirect
)

replace github.com/meowalien/rabbitgather-lib => ../lib
