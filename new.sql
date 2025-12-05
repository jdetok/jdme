select distinct team_id 
from stats.tbox
where;

select exists (select team_id from stats.tbox where szn_id = $1 and team_id = $2);
select distinct team_id from stats.tbox where szn_id = $1;
select distinct a.team_id, b.lg_id 
from stats.tbox a
inner join lg.team b on b.team_id = a.team_id 
where szn_id = $1;

select szn_id, string_agg(distinct team_id::text, ',')
from stats.pbox
where player_id = $1 
and substr(szn_id::text, 2, 4)::int between substr($2::text, 2, 4)::int and substr($3::text, 2, 4)::int
group by player_id, szn_id;


select substr(szn_id::text, 2, 4) from lg.szn;
create index idx_szn_no_prefix on lg.szn(substr(szn_id::text, 2, 4)); 