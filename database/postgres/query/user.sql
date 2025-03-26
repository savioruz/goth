-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (email, password, level, google_id, full_name, profile_image, is_verified) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: UpdateUser :one
UPDATE users SET email = $1, password = $2, google_id = $3, full_name = $4, profile_image = $5, is_verified = $6, updated_at = now()
    WHERE id = $7 AND deleted_at IS NULL RETURNING *;

-- name: UpdateLastLogin :one
UPDATE users SET last_login = now() WHERE id = $1 AND deleted_at IS NULL RETURNING id;

-- name: CreateEmailVerification :one
INSERT INTO email_verifications (user_id, token) VALUES ($1, $2) RETURNING *;

-- name: GetEmailVerificationByToken :one
SELECT * FROM email_verifications WHERE token = $1 AND expires_at > now() LIMIT 1;

-- name: VerifyEmail :one
UPDATE users SET is_verified = true WHERE id = $1 AND deleted_at IS NULL RETURNING *;

-- name: CreatePasswordReset :one
INSERT INTO password_resets (user_id, token) VALUES ($1, $2) RETURNING *;

-- name: GetPasswordResetByToken :one
SELECT * FROM password_resets WHERE token = $1 AND expires_at > now() LIMIT 1;

-- name: ResetPassword :one
UPDATE users SET password = $1 WHERE id = $2 AND deleted_at IS NULL RETURNING *;
