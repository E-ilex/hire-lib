SELECT
    q.body AS question_body,
	json_group_array(
        json_object(
            'body', o.body,
            'correct', o.correct
        )
    ) AS options
FROM (
    SELECT * FROM question
    WHERE (id > $1 OR $1 IS NULL)
    ORDER BY id DESC
    LIMIT COALESCE($2, 10000000)
) AS q
JOIN option o ON q.id = o.question_id
GROUP BY q.id
ORDER BY q.id, o.id DESC;
