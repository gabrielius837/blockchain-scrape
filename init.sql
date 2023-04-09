PRAGMA foreign_keys = ON;

CREATE TABLE block (
    number INTEGER NOT NULL PRIMARY KEY
);

-- last block without any transactions
INSERT INTO block VALUES(46146);

CREATE TABLE address (
    id BLOB NOT NULL PRIMARY KEY CHECK(length(id) = 20 or length(id) = 0),
    block INTEGER NOT NULL,
    flag INTEGER NOT NULL
);

INSERT INTO address (id, block, flag) VALUES(X'', 0, -1);

CREATE TABLE tx (
    id BLOB NOT NULL PRIMARY KEY CHECK(length(id) = 32),
    "from" BLOB NOT NULL,
    "to" BLOB NOT NULL,
    FOREIGN KEY("from") REFERENCES address(id)
    	ON UPDATE CASCADE
    	ON DELETE RESTRICT,
    FOREIGN KEY("to") REFERENCES address(id)
    	ON UPDATE CASCADE
    	ON DELETE RESTRICT
);

CREATE INDEX idx_tx_from_to 
ON tx("from", "to");
