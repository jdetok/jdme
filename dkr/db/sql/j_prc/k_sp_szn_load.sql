/* 
functions and procedures to take a season id & generate the values for lg.season
relies on lg.szn_type being populated to create the concatenated szn_desc and wszn_desc fields
*/

create or replace function lg.fn_szn(s_id int)
returns varchar(255)
language plpgsql
as $$
declare
    season_concat varchar(255);
begin
    select
        right(cast(s_id as varchar(10)), 4) || '-' ||
        right(cast((s_id+1) as varchar(10)), 2)
    into season_concat
    from lg.szn_type
    where cast(left(cast(s_id as varchar(10)), 1) as int) = sznt_id;
    return season_concat;
end; $$; 

create or replace function lg.fn_szn_desc(s_id int)
returns varchar(255)
language plpgsql
as $$
declare
    season_concat varchar(255);
begin
    select 
        right(cast(s_id as varchar(10)), 4) || '-' ||
        right(cast((s_id+1) as varchar(10)), 2) || ' ' || szn_type
    into season_concat
    from lg.szn_type
    where cast(left(cast(s_id as varchar(10)), 1) as int) = sznt_id;
    return season_concat;
end; $$; 

create or replace function lg.fn_wszn(s_id int)
returns varchar(255)
language plpgsql
as $$
declare
    season_concat varchar(255);
begin
    select
        right(cast(s_id as varchar(10)), 4)
    into season_concat
    from lg.szn_type
    where cast(left(cast(s_id as varchar(10)), 1) as int) = sznt_id;
    return season_concat;
end; $$; 

create or replace function lg.fn_wszn_desc(s_id int)
returns varchar(255)
language plpgsql
as $$
declare
    season_concat varchar(255);
begin
    select 
        right(cast(s_id as varchar(10)), 4) || ' WNBA ' || szn_type
    into season_concat
    from lg.szn_type
    where cast(left(cast(s_id as varchar(10)), 1) as int) = sznt_id;
    return season_concat;
end; $$; 

create or replace procedure lg.sp_szn_load()
language plpgsql
as $$
begin
    insert into lg.szn
        select 
            a.season_id, 
            b.sznt_id, 
            lg.fn_szn(a.season_id),
            lg.fn_szn_desc(a.season_id),
            lg.fn_wszn(a.season_id),
            lg.fn_wszn_desc(a.season_id)
        from intake.gm_team a
        inner join lg.szn_type b
            on cast(left(cast(a.season_id as varchar(10)), 1) as int) = b.sznt_id
        group by a.season_id, b.sznt_id
    on conflict (szn_id) do nothing;
end; $$;
-- call lg.sp_szn_load();
-- select * from lg.szn;