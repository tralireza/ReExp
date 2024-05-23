module main

go 1.22.3

require (
	github.com/go-sql-driver/mysql v1.8.1
	github.com/lib/pq v1.10.9
	github.com/tralireza/ReExp/reexp v0.0.0
)

require filippo.io/edwards25519 v1.1.0 // indirect

replace github.com/tralireza/ReExp/reexp v0.0.0 => ./reexp
