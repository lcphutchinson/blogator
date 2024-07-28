module github.com/lcphutchinson/blogator

go 1.22.4

replace github.com/lcphutchinson/database v0.0.0 => ./internal/database

require github.com/lcphutchinson/database v0.0.0

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
)
