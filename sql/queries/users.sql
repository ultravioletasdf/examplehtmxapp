-- name: CreateUser :exec
INSERT INTO
    users (id, email, password, created_at)
VALUES
    (?, ?, ?, unixepoch ("now"));

-- name: VerifyUser :exec
UPDATE users
SET
    verified = 1
WHERE
    id = ?;

-- name: SetVerification :exec
UPDATE users
SET
    verify_code = ?,
    verify_expire_at = strftime ('%s', 'now', '15 minutes')
WHERE
    id = ?;

-- name: GetUser :one
SELECT
    *
FROM
    users
WHERE
    id = ?;

-- name: GetUserByEmail :one
SELECT
    *
FROM
    users
WHERE
    email = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE
    id = ?;
