
-- +migrate Up
CREATE TABLE `node_annotations` (
    `node_name` varchar(128),
    `key` varchar(128),
    `value` text
);

-- +migrate Down
DROP TABLE `node_annotations`;
