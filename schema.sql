CREATE TABLE items (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE tags (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE item_tags (
    item_id TEXT NOT NULL,
    tag_id TEXT NOT NULL,
    PRIMARY KEY (item_id, tag_id),
    FOREIGN KEY (item_id) REFERENCES items(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id)
);

CREATE TABLE lists (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE list_items (
    list_id TEXT NOT NULL,
    item_id TEXT NOT NULL,
    PRIMARY KEY (list_id, item_id),
    FOREIGN KEY (list_id) REFERENCES lists(id),
    FOREIGN KEY (item_id) REFERENCES items(id)
);
