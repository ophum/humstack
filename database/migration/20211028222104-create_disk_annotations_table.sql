
-- +migrate Up
CREATE TABLE `disk_annotations` (
    disk_name varchar(128),
    key varchar(128),
    value text
);

-- +migrate Down
DROP TABLE `disk_annotations`;