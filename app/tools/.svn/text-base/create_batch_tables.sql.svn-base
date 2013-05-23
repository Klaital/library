
CREATE TABLE lookup_requests (
	request_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	user_id    INT NOT NULL,
	code_type  VARCHAR(15) DEFAULT 'UPC',
	code       VARCHAR(63) NOT NULL,
	created    DATETIME,
    location   VARCHAR(63),
	notes      VARCHAR(125)
);

CREATE TABLE lookup_steps (
	step_id    INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	request_id INT NOT NULL,
	source     VARCHAR(63) NOT NULL,
	created    DATETIME,
	started    DATETIME,
	completed  DATETIME
);

