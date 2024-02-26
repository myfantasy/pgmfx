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
INSERT INTO config.dynamic_config
(data_center, app_name, "name", "version", description, cfg)
VALUES('', 'd.2', 'db.pg.test_all.cfg', 'v1', 'all_databases', '
{
	"test1":{
		 "conn_string": "user=postgres host=localhost port=5432 dbname=test_db1",
		 "max_connections": 5,
		 "min_connections": 1,
		 "pool_name": "test1"
		},
	"test2":{
		 "conn_string": "user=postgres host=localhost port=5432 dbname=test_db2",
		 "max_connections": 5,
		 "min_connections": 1,
		 "pool_name": "test1"
		},
	"test3":{
		 "conn_string": "user=postgres host=localhost port=5432 dbname=test_db3",
		 "max_connections": 5,
		 "min_connections": 1,
		 "pool_name": "test1"
		},
	"test4":{
		 "conn_string": "user=postgres host=localhost port=5432 dbname=test_db4",
		 "max_connections": 5,
		 "min_connections": 1,
		 "pool_name": "test1"
		}
}
'::json);
INSERT INTO config.dynamic_config
(data_center, app_name, "name", "version", description, cfg)
VALUES('', 'd.2', 'db.pg.test_all.pwd', 'v1', 'all_databases', '
{
	"test1":{
		 "password": "postgres"
		},
	"test2":{
		 "password": "postgres"
		},
	"test3":{
		 "password": "postgres"
		},
	"test4":{
		 "password": "postgres"
		}
}
'::json);