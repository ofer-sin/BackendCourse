-- name: CreateEntry :one
INSERT INTO entries (
  account_id,
  amount
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: ListEntries :many
-- This query retrieves a list of entries from the "entries" table that belong to a specific account (filtered by account_id).
-- The results are ordered by the "id" column in ascending order.
-- The "LIMIT $2" clause restricts the number of rows returned to the value specified by the second parameter.
-- The "OFFSET $3" clause skips the first $3 rows, allowing for pagination of results.

-- Key Points:
-- account_id = $1: Filters rows based on the provided account ID.
-- ORDER BY id: Ensures the results are sorted by the id column in ascending order.
-- LIMIT $2: Limits the number of rows returned to the value of $2.
-- OFFSET $3: Skips the first $3 rows, useful for implementing pagination.
-- This query is commonly used in applications to fetch a subset of data for a specific account, often for displaying paginated results in a UI.

-- name: ListEntries :many
SELECT * FROM entries
WHERE account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;