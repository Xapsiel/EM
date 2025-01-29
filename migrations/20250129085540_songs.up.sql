CREATE TABLE groups (
                        id SERIAL PRIMARY KEY,
                        name TEXT UNIQUE NOT NULL
);
CREATE TABLE songs (
                       id SERIAL PRIMARY KEY,
                       group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
                       song_name TEXT NOT NULL,
                       release_date DATE NOT NULL,
                       text TEXT NOT NULL,
                       link TEXT NOT NULL,
                       created_at TIMESTAMP DEFAULT NOW(),
                       updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_songs_group_id ON songs (group_id);
CREATE INDEX idx_songs_song_name ON songs (song_name);