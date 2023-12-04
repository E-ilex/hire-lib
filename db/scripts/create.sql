CREATE TABLE question (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    body VARCHAR(255) NOT NULL,
    ts_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    ts_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE TABLE options (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    question_id INTEGER NOT NULL,
    body VARCHAR(255) NOT NULL,
    rank INTEGER CHECK (rank >= 0) NOT NULL,
    correct BOOLEAN NOT NULL,
    ts_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    ts_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (question_id) REFERENCES question(id)
);
CREATE INDEX idx_options_question_id ON options (question_id);