-- name: InsertMaterial :copyfrom
INSERT INTO materials(id, session_id, index, type, is_shared, name, data)
VALUES($1, $2, $3, $4, $5, $6, $7);

-- name: DeleteSessionMaterials :exec
DELETE FROM materials WHERE session_id = sqlc.arg(session_id);

-- name: GetMaterialsBySessionID :many
SELECT * from materials WHERE session_id = sqlc.arg(session_id);
