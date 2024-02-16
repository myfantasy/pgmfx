CREATE SCHEMA config;

CREATE TABLE config.dynamic_config (
    data_center text NOT NULL,
	app_name text NOT NULL,
	"name" text NOT NULL,
	"version" text NOT NULL,
	description text NOT NULL,
	cfg json NOT NULL,
	CONSTRAINT dynamic_config_pkey PRIMARY KEY (app_name, name)
);

CREATE TABLE config.dynamic_config_logs (
	"_log_id" int8 DEFAULT nextval('config.dynamic_config_logs_sq'::regclass) NOT NULL,
	"_ts" timestamptz DEFAULT now() NOT NULL,
	"_create_month" int4 DEFAULT EXTRACT(month FROM now()) NOT NULL,
	"_action" text NOT NULL,
    data_center text NOT NULL,
	app_name text NOT NULL,
	"name" text NOT NULL,
	"version" text NOT NULL,
	description text NOT NULL,
	cfg json NOT NULL,
	CONSTRAINT dynamic_config_logs_pkey PRIMARY KEY (_log_id)
);

CREATE OR REPLACE FUNCTION config.dynamic_config_get_one_by_name_and_app_name(_data_center text, _app_name text, _name text)
 RETURNS json
 LANGUAGE plpgsql
 SECURITY DEFINER
AS $function$
    declare _result json;
begin
	with cte as (
		select 	data_center,
				app_name,
				name,
				version,
				description,
				cfg
			from 	config.dynamic_config dc
			where	_name like (dc.name||'%')
				and _app_name like (dc.app_name||'%')
				and _data_center like (dc.data_center||'%')
			order by  
                char_length(dc.name) desc,
                char_length(dc.data_center) desc,
                char_length(dc.app_name) desc
			limit 1000
	)
	select 
			jsonb_agg(to_json(r.*))
		into	
			_result
		from cte r;
	
	return jsonb_pretty(coalesce(_result ,'null'::json)::jsonb);
END;
$function$
;
