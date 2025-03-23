-- name: CreateSession :exec
INSERT INTO
    sessions (token, user_id, expire_at, created_at)
VALUES
    (?, ?, ?, unixepoch ("now"));

-- name: GetSession :one
SELECT
    *
FROM
    sessions
WHERE
    token = ?
    AND expire_at > unixepoch ("now");

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE
    token = ?;

-- name: GetUserFromSession :one
SELECT
    users.*
FROM
    users
    JOIN sessions ON users.id = sessions.user_id
WHERE
    sessions.token = ?;
