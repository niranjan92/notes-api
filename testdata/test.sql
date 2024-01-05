SELECT * FROM notes n
LEFT JOIN shared_notes sn ON sn.note_id = n.id
WHERE n.user_id = '1'
OR sn.shared_user_id = '1';


SELECT * FROM notes n
LEFT JOIN shared_notes sn ON sn.note_id = n.id
WHERE sn.shared_user_id = '1';


SELECT n.* FROM notes n
LEFT JOIN shared_notes sn ON n.id = sn.note_id
WHERE (n.user_id = '1' OR sn.shared_user_id = '1')
AND n.text @@ to_tsquery('english', 'quick');

SELECT n.* FROM notes n
LEFT JOIN shared_notes sn ON n.id = sn.note_id
WHERE (n.user_id = '1' OR sn.shared_user_id = '1')
AND n.text @@ to_tsquery('english', 'brown');


CREATE TABLE shared_notes
(
    id         VARCHAR NOT NULL,
    note_id    VARCHAR NOT NULL,
    shared_user_id    VARCHAR NOT NULL,
    PRIMARY KEY (note_id, shared_user_id),
    -- FOREIGN KEY (note_id) REFERENCES notes (id),
    -- FOREIGN KEY (shared_user_id) REFERENCES users (id)
);