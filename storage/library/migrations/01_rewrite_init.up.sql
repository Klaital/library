-- noinspection SqlNoDataSourceInspectionForFile

CREATE TABLE locations (
    id INTEGER NOT NULL PRIMARY KEY,
    name VARCHAR(64) NOT NULL UNIQUE,
    notes TEXT
);

CREATE TABLE items (
    id INTEGER NOT NULL PRIMARY KEY,
    location_id INT NOT NULL,

    code VARCHAR(64) NOT NULL,
    code_type VARCHAR(16) NOT NULL DEFAULT 'UPC',
    code_source VARCHAR(64) NOT NULL DEFAULT 'manual',

    title VARCHAR(128) NOT NULL,
    title_translated VARCHAR(128),
    title_transliterated VARCHAR(128),

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
