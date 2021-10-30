
-- +migrate Up
CREATE TABLE `node_agent_statuses` (
    `node_name` varchar(128),
    `agent_name` varchar(128),
    `status` varchar(128)
);

-- +migrate Down
DROP TABLE `node_agent_statuses`;

