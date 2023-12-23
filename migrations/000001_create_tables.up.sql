CREATE TABLE IF NOT EXISTS signatures
(
    id         INTEGER      NOT NULL PRIMARY KEY AUTO_INCREMENT,
    user_jwt   VARCHAR(512) NOT NULL,
    created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    signature  VARCHAR(512) NOT NULL,
    questions  json         NOT NULL,
    CONSTRAINT signature_unique UNIQUE (signature)
);