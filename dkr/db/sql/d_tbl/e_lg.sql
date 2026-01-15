create table lg.league (
    lg_id int primary key,
    lg_cde varchar(10) unique,
    lg varchar(10) unique,
    lg_name varchar(255)
);

create index idx_lg on lg.league(lg);

insert into lg.league values 
	(99, 'nat', 'NAT', 'Not Assigned to a Team'),
    (0, 'nba', 'NBA', 'National Basketball Association'),
    (1, 'wnba', 'WNBA', 'Women''s National Basketball Association'),
    (9, 'misc', 'MISC', 'Temporary/Miscellaneous League');

create table lg.szn_type (
    sznt_id int primary key,
    sznt_cde varchar(2) not null,
    szn_type varchar(255)
);

insert into lg.szn_type values 
    (1, 'PS', 'Pre-Season'),
    (2, 'RG', 'Regular Season'),
    (3, 'PS', 'All-Star'),
    (4, 'PO', 'Playoffs'),
    (5, 'PI', 'Play-In Tournament'),
    (6, 'NC', 'NBA Cup'),
    (9, 'CA', 'Combined Aggregates');

create table lg.szn (
    szn_id int primary key,
    sznt_id int references lg.szn_type(sznt_id),
    szn varchar(255),
    szn_desc varchar(255),
    wszn varchar(255),
    wszn_desc varchar(255)
);

create index idx_sznt on lg.szn(sznt_id);
create index idx_szn on lg.szn(szn);
create index idx_wszn on lg.szn(wszn);

-- insert aggregate seasons for api stats tables
insert into lg.szn values
    (29999, 2, 'CRS', 'Career Regular Season', 'CRS', 'Career Regular Season'),
    (49999, 4, 'CPO', 'Career Playoffs', 'CPO', 'Career Playoffs'),
    (99999, 9, 'CC', 'Career Combined', 'CC', 'Career Combined');

-- update for jdeko.me agg season labels: 
update lg.szn set szn_desc = 'Career Regular Season Stats', 
    wszn_desc = 'Career Regular Season Stats'
where szn_id = 29999;

update lg.szn set szn_desc = 'Career Post Season Stats', 
    wszn_desc = 'Career Post Season Stats'
where szn_id = 49999;

create table lg.team (
    lg_id int references lg.league(lg_id),
    team_id bigint primary key,
    team varchar(10),
    team_cde varchar(255),
    team_long varchar(255),
    team_city varchar(255),
    team_shrt varchar(255)
);

create index idx_team on lg.team(team);
create index idx_tlg on lg.team(lg_id);
create index idx_team_cde on lg.team(team_cde);

insert into lg.team (lg_id, team_id, team, team_cde, team_long) values 
    (99, 0, 'NAT', 'noteam', 'Not Assigned');

create table lg.plr (
    lg_id int references lg.league(lg_id),
    player_id bigint primary key,
    plr_cde varchar(255),
    player varchar(255),
    last_first varchar(255),
    from_year varchar(4),
    to_year varchar(4)
);
create index idx_plg on lg.plr(lg_id);
create index idx_pfromyear on lg.plr(from_year);
create index idx_ptoyear on lg.plr(to_year);

create table lg.plr_crnt (
    lg_id int references lg.league(lg_id),
    team_id bigint references lg.team(team_id),
    player_id bigint references lg.plr(player_id),
    plr_cde varchar(255),
    primary key (player_id, team_id)
);

create index idx_cplg on lg.plr_crnt(lg_id);
create index idx_cpt on lg.plr_crnt(team_id);

