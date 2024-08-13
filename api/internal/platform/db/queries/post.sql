-- name: GetPostsByUserAndTags :one
WITH post_likes AS (
    SELECT p.id, COUNT(ulp.post_id) AS total_likes
    FROM "post" p
             LEFT JOIN "user_liked_post" ulp ON p.id = ulp.post_id
    GROUP BY p.id
),
     filtered_posts AS (
         SELECT p.id, p.title, p.description, p.created_at AS "createdAt", ap.price, p.views,
                pl.total_likes AS "totalLikes",
                json_build_object('userId', u.id,
                                  'username', u.username,
                                  'name', u.name,
                                  'avatar', u.avatar_link) AS "user",
                json_build_object('artworkId', a.id,
                                  'image', a.processed_image_link,
                                  'type', a.type,
                                  'isBuyable', a.is_buyable) AS "artwork",
                EXISTS (
                    SELECT 1
                    FROM "user_liked_post" ulp
                    WHERE ulp.post_id = p.id AND ulp.user_id = NULLIF(sqlc.arg(user_id)::unknown, '00000000-0000-0000-0000-000000000000')::UUID
                ) AS "isLiked",
                ROW_NUMBER() OVER (ORDER BY
                    CASE WHEN sqlc.arg(sort_by)::text = 'title' AND sqlc.arg(sort_order)::text = 'ASC' THEN p.title END,
                    CASE WHEN sqlc.arg(sort_by)::text = 'title' AND sqlc.arg(sort_order)::text = 'DESC' THEN p.title END DESC,
                    CASE WHEN sqlc.arg(sort_by)::text = 'price' AND sqlc.arg(sort_order)::text = 'ASC' THEN ap.price END,
                    CASE WHEN sqlc.arg(sort_by)::text = 'price' AND sqlc.arg(sort_order)::text = 'DESC' THEN ap.price END DESC,
                    CASE WHEN sqlc.arg(sort_by)::text = 'createdAt' AND sqlc.arg(sort_order)::text = 'ASC' THEN p.created_at END,
                    CASE WHEN sqlc.arg(sort_by)::text = 'createdAt' AND sqlc.arg(sort_order)::text = 'DESC' THEN p.created_at END DESC,
                    CASE WHEN sqlc.arg(sort_by)::text = 'totalLikes' AND sqlc.arg(sort_order)::text = 'ASC' THEN pl.total_likes END,
                    CASE WHEN sqlc.arg(sort_by)::text = 'totalLikes' AND sqlc.arg(sort_order)::text = 'DESC' THEN pl.total_likes END DESC,
                    p.id
                    ) AS row_num
         FROM "post" p
                  INNER JOIN "user" u ON p.user_id = u.id
                  INNER JOIN "artwork" a ON p.id = a.post_id
                  LEFT JOIN "artwork_price" ap ON a.id = ap.artwork_id AND ap.to_date IS NULL
                  LEFT JOIN post_likes pl ON p.id = pl.id
         WHERE
             (u.username ILIKE '%' || sqlc.arg(search_term)::text || '%')
           AND ((sqlc.arg(tag_ids)::uuid[] IS NULL) OR (p.id IN (
             SELECT pt.post_id
             FROM "post_tags" pt
             WHERE pt.tag_id = ANY(sqlc.arg(tag_ids)::uuid[])
         ))) AND a.is_buyable = sqlc.arg(is_buyable)::boolean
           AND p.is_ban = false
           AND p.is_deleted = false
           AND a.is_deleted = false
     ),
     total_count AS (
         SELECT COUNT(*) AS count FROM filtered_posts
     )
SELECT
    json_build_object(
            'items', json_agg(filtered_posts.* ),
            'totalCount', (SELECT count FROM total_count),
            'hasNextPage', (SELECT count FROM total_count) > (sqlc.arg(size)::INT * sqlc.arg(page)::INT),
            'hasPrevPage', sqlc.arg(page) > 1,
            'totalPages', CEIL(CAST((SELECT count FROM total_count) AS FLOAT)/sqlc.arg(size)),
            'pageNumber', sqlc.arg(page),
            'pageSize', sqlc.arg(size),
            'sortBy', sqlc.arg(sort_by),
            'sortOrder', sqlc.arg(sort_order)
    ) AS result
FROM
    filtered_posts
WHERE
    filtered_posts.row_num > ((sqlc.arg(page) - 1) * sqlc.arg(size))
  AND filtered_posts.row_num <= (sqlc.arg(page) * sqlc.arg(size));