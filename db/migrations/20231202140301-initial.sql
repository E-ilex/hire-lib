-- +migrate Up
CREATE TABLE question (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    body VARCHAR(255) NOT NULL,
    ts_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE TABLE option (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    question_id INTEGER NOT NULL,
    body VARCHAR(255) NOT NULL,
    correct BOOLEAN NOT NULL,
    ts_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (question_id) REFERENCES question(id)
);
CREATE INDEX idx_options_question_id ON option (question_id);
-- +migrate Down
DROP TABLE IF EXISTS option;
DROP TABLE IF EXISTS question;