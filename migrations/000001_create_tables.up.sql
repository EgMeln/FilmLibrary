CREATE TABLE IF NOT EXISTS actors(
    id          UUID PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    gender      VARCHAR(10) NOT NULL,
    birth_date  TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS movies (
    id              UUID PRIMARY KEY,
    title           VARCHAR(150) NOT NULL CHECK (LENGTH(title) > 0 AND LENGTH(title) <= 150),
    description     TEXT NOT NULL CHECK (LENGTH(description) <= 1000),
    release_date    TIMESTAMP NOT NULL,
    rating          INTEGER NOT NULL CHECK (rating >= 0 AND rating <= 10)
);

CREATE TABLE IF NOT EXISTS movie_actor (
    movie_id UUID REFERENCES movies(id),
    actor_id UUID REFERENCES actors(id),
    PRIMARY KEY (movie_id, actor_id)
);

CREATE TABLE IF NOT EXISTS users (
    id        UUID PRIMARY KEY,
    username  VARCHAR(30) NOT NULL,
    password  VARCHAR(150) NOT NULL,
    role      VARCHAR(10) DEFAULT 'user'
);