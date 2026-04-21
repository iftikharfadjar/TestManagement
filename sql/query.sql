-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = ? LIMIT 1;

-- name: CreateUser :exec
INSERT INTO users (id, username, password_hash) VALUES (?, ?, ?);
