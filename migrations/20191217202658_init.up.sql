CREATE TABLE notes
(
    id         VARCHAR PRIMARY KEY,
    title       VARCHAR NOT NULL,
    text       VARCHAR NOT NULL,
    text_searchable       TSVECTOR,
    user_id    VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
CREATE INDEX notes_tsv_idx ON notes USING gin(text_searchable);


CREATE TABLE users
(
    id         VARCHAR PRIMARY KEY ,
    name       VARCHAR NOT NULL,
    password   VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE shared_notes
(
    id         VARCHAR NOT NULL,
    note_id    VARCHAR NOT NULL,
    shared_user_id    VARCHAR NOT NULL,
    PRIMARY KEY (note_id, shared_user_id)
);