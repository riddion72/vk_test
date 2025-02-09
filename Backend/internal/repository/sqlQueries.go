package repository

const (
	queryGetPing   = `SELECT * FROM addresses LIMIT $1 OFFSET $2`
	queryGetNumber = `SELECT COUNT(*) FROM addresses`
	querySetPing   = `INSERT INTO addresses (ip, response_time, last_successful_ping) VALUES ($1, $2, $3)
	ON CONFLICT (ip) DO UPDATE SET response_time = $2, last_successful_ping = $3 RETURNING *`
	querySetNoAnswer = `INSERT INTO addresses (ip, response_time) VALUES ($1, $2)
	ON CONFLICT (ip) DO UPDATE SET response_time = $2 RETURNING *`
)
