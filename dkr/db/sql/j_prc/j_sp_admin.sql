-- returns count of passed table as int
create or replace function fn_cnt(_tbl regclass, out cnt bigint)
returns bigint
language plpgsql
as $$
begin
	execute format('select count(*) from %s', _tbl) into cnt;
end; $$;

-- returns count of passed table as varchar
create or replace function fn_cntstr(_tbl regclass, out cntstr varchar(255))
language plpgsql
as $$
declare cnt bigint;
begin
	execute format('select count(*) from %s', _tbl) into cnt;
	cntstr := cast(cnt as varchar(255)) || format(' rows in table <%s>', _tbl);
end; $$;

select fn_cntstr('lg.plr');