
-- +migrate Up
CREATE TABLE `disks` (
    name varchar(128) PRIMARY KEY,
    type varchar(64),
    request_bytes int,
    limit_bytes int,
    status varchar(64),
    created_at datetime,
    updated_at datetime
);

-- +migrate Down
DROP TABLE `disks`;