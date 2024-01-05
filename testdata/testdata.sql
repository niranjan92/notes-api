INSERT INTO users (id, name, password, created_at, updated_at)
VALUES ('1', 'demo1', 'pass',  '2019-10-11 19:43:18'::timestamp, '2019-10-11 19:43:18'::timestamp),
       ('2', 'demo2', 'pass',  '2019-10-01 15:36:38'::timestamp, '2019-10-01 15:36:38'::timestamp);

INSERT INTO notes (id, title, text, text_searchable, user_id, created_at, updated_at)
VALUES ('asdf', 'note title', 'apple a day keeps doctor away. brown fox jumped',  'apple a day keeps doctor away. brown fox jumped', '1', '2019-10-11 19:43:18'::timestamp, '2019-10-11 19:43:18'::timestamp),
      ('asdfsds', 'note title 2', 'quick brown fox',  'quick brown fox', '1', '2019-10-01 15:36:38'::timestamp, '2019-10-01 15:36:38'::timestamp),
      ('erter', 'note title 3', 'striver like striver', 'striver like striver', '2', '2019-10-01 15:36:38'::timestamp, '2019-10-01 15:36:38'::timestamp),
      ('erterer', 'note title 4', 'sun rises in the east','sun rises in the east',  '2', '2019-10-01 15:36:38'::timestamp, '2019-10-01 15:36:38'::timestamp),
      ('ertererer', 'note title 5', 'AI is the future', 'AI is the future',  '2', '2019-10-01 15:36:38'::timestamp, '2019-10-01 15:36:38'::timestamp);

INSERT INTO shared_notes (id, note_id, shared_user_id)
VALUES ('3', 'erter', '1'),
         ('4', 'erterer', '1');


