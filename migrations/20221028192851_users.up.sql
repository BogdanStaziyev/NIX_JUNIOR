CREATE TABLE IF NOT EXISTS users (
    id              serial PRIMARY KEY,
    email           varchar(50) NOT NULL,
    password        varchar(100) NOT NULL,
    name            varchar(50) NOT NULL
);
