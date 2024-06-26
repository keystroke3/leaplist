CREATE TABLE user (
    id PRIMARY KEY,
    username TEXT NOT NULL,
    display_name TEXT NOT NULL,
    passphrase BLOB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE station(
    id PRIMARY KEY,
    user_id TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
);
CREATE TABLE tag (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    label TEXT NOT NULL UNIQUE,
    station_id TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (station_id) REFERENCES station(id) ON DELETE CASCADE
);
CREATE TABLE relay (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    alias TEXT NOT NULL,
    destination TEXT NOT NULL,
    note TEXT,
    station_id TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (station_id) REFERENCES station(id) ON DELETE CASCADE
);
CREATE TABLE relay_tag (
    relay_id INT NOT NULL,
    tag_id INT NOT NULL,
    station_id TEXT NOT NULL,
    PRIMARY KEY (relay_id, tag_id),
    FOREIGN KEY (relay_id) REFERENCES relays(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
    FOREIGN KEY (station_id) REFERENCES station(id) ON DELETE CASCADE
);