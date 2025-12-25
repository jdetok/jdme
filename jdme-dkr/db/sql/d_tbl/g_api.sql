/*
TABLES IN API SCHEMA
NEEDS TO MATCH `api_player_stats` from legacy mariadb to work with website
*/

-- EXPLORE USING A VIEW OR NORMAL QUERY WITH THE NEW INDEXING

create table if not exists api.plr_stats (
    player_id int references lg.plr(player_id),	
	team_id int references lg.team(team_id),
	league varchar(5) references lg.league(lg),
	season_id int references lg.szn(szn_id),
	stat_type varchar(50),
	player varchar(255),
	team varchar(5),
	team_name varchar(255),
	games_played int,
	minutes decimal(10,2),
	points decimal(10,2),
	assists decimal(10,2),
	rebounds decimal(10,2),
	steals decimal(10,2),
	blocks decimal(10,2),
	fgm decimal(10,2),
	fga decimal(10,2),
	fgp varchar(10),
	fg3m decimal(10,2),
	fg3a decimal(10,2),
	fg3p varchar(10),
	ftm decimal(10,2),
	fta decimal(10,2),
	ftp varchar(10),
    primary key (player_id, team_id, stat_type)
);

-- same as above with season & wseason desc between season_id and stat_type
create table if not exists api.plr_agg (
    player_id int references lg.plr(player_id),	
	team_id int references lg.team(team_id),
	league varchar(5) references lg.league(lg),
	season_id int references lg.szn(szn_id),
	season_desc varchar(255),
	wseason_desc varchar(255),
	stat_type varchar(50),
	player varchar(255),
	team varchar(5),
	team_name varchar(255),
	games_played int,
	minutes decimal(10,2),
	points decimal(10,2),
	assists decimal(10,2),
	rebounds decimal(10,2),
	steals decimal(10,2),
	blocks decimal(10,2),
	fgm decimal(10,2),
	fga decimal(10,2),
	fgp varchar(10),
	fg3m decimal(10,2),
	fg3a decimal(10,2),
	fg3p varchar(10),
	ftm decimal(10,2),
	fta decimal(10,2),
	ftp varchar(10),
    primary key (player_id, season_id, stat_type)
);