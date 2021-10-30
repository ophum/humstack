
-- +migrate Up
CREATE TABLE `node_agent_daemons` (
    `node_name` varchar(128),
    `name` varchar(128),
    `command` varchar(128),
    `args` text,
    `envs` text,
    `restart_policy` varchar(128)
);


-- +migrate Down
DROP TABLE `node_agent_daemons`;
