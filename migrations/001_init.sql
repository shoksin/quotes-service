CREATE TABLE IF NOT EXISTS quotes
(
    id SERIAL PRIMARY KEY,
    author VARCHAR(255) NOT NULL,
    quote TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO quotes (author, quote)
VALUES ('Confucius', 'Life is simple, but we insist on making it complicated.'),
       ('Albert Einstein', 'Imagination is more important than knowledge.'),
       ('Maya Angelou', 'If you don''t like something, change it. If you can''t change it, change your attitude.'),
       ('Steve Jobs', 'Innovation distinguishes between a leader and a follower.'),
       ('Winston Churchill', 'Success is not final, failure is not fatal: it is the courage to continue that counts.');
