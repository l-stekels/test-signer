CREATE TABLE IF NOT EXISTS signatures
(
    id         INTEGER      NOT NULL PRIMARY KEY AUTO_INCREMENT,
    user_jwt   VARCHAR(512) NOT NULL,
    created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT user_jwt_unique UNIQUE (user_jwt)
);

CREATE TABLE IF NOT EXISTS questions
(
    id         INTEGER  NOT NULL PRIMARY KEY AUTO_INCREMENT,
    body       TEXT     NOT NULL,
    signature_id  INTEGER  NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT signature_id_fk FOREIGN KEY (signature_id) REFERENCES signatures (id)
);

CREATE TABLE IF NOT EXISTS answers
(
    id          INTEGER  NOT NULL PRIMARY KEY AUTO_INCREMENT,
    question_id INTEGER  NOT NULL,
    body        TEXT     NOT NULL,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT question_id_fk FOREIGN KEY (question_id) REFERENCES questions (id)
);