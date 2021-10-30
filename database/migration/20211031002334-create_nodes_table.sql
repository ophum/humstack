
-- +migrate Up
CREATE TABLE `nodes` (
    `name` varchar(128) PRIMARY KEY,
    `hostname` varchar(255),
    `status` varchar(128),
    `created_at` datetime,
    `updated_at` datetime
);

-- +migrate Down
DROP TABLE `nodes`;
