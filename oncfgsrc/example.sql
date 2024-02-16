INSERT INTO config.dynamic_config
(data_center, app_name, "name", "version", description, cfg)
VALUES('dc', 'a', 'n1', 'v0', 'd', '{}'::json);
INSERT INTO config.dynamic_config
(data_center, app_name, "name", "version", description, cfg)
VALUES('dc', 'a.2', 'n1', 'v0', 'd2', '{}'::json);
INSERT INTO config.dynamic_config
(data_center, app_name, "name", "version", description, cfg)
VALUES('dc', 'b', 'n1', 'v0', 'd4', '{}'::json);
INSERT INTO config.dynamic_config
(data_center, app_name, "name", "version", description, cfg)
VALUES('dc', 'c', 'n1', 'v0', 'd7', '{}'::json);
INSERT INTO config.dynamic_config
(data_center, app_name, "name", "version", description, cfg)
VALUES('dc', 'c', 'n2', 'v0', 'd10', '{}'::json);
INSERT INTO config.dynamic_config
(data_center, app_name, "name", "version", description, cfg)
VALUES('dc', 'd.2.5', 'n1', 'v0', 'd15', '{"a":7}'::json);
INSERT INTO config.dynamic_config
(data_center, app_name, "name", "version", description, cfg)
VALUES('dc', 'd', 'n1', 'v0', 'd12', '{"a":5, "c":98}'::json);
INSERT INTO config.dynamic_config
(data_center, app_name, "name", "version", description, cfg)
VALUES('dc', 'd.2', 'n1', 'v0', 'd19', '{"a":6, "b":99}'::json);
INSERT INTO config.dynamic_config
(data_center, app_name, "name", "version", description, cfg)
VALUES('dc', 'd.2', 'n1.fun', 'v0', 'd5', '{"f":"f"}'::json);