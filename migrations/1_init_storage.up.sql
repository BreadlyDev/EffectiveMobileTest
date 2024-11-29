CREATE TABLE IF NOT EXISTS groups(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
); 

CREATE INDEX idx_name ON groups(name);

CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    group_id INT REFERENCES groups(id) ON DELETE CASCADE, 
    release_date DATE,
    text TEXT DEFAULT '',
    link TEXT DEFAULT ''
);
