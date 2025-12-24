-- tables that match stats.nba.com json responses

create table intake.player (
    person_id bigint primary key,
    display_last_comma_first varchar(255),
    display_first_last varchar(255),
    rosterstatus boolean,
    from_year varchar(4),
    to_year varchar(4),
    playercode varchar(255),
    player_slug varchar(255),
    team_id bigint,
    team_city varchar(255),
    team_name varchar(255),
    team_abbreviation varchar(10),
    team_slug varchar(255),
    team_code varchar(255),
    games_played_flag varchar(1),
    otherleague_experience_ch varchar(10)
);

create index idx_inpl_team on intake.player(team_id);

create table intake.wplayer (
    person_id bigint primary key,
    display_last_comma_first varchar(255),
    display_first_last varchar(255),
    rosterstatus boolean,
    from_year varchar(4),
    to_year varchar(4),
    playercode varchar(255),
    player_slug varchar(255),
    team_id bigint,
    team_city varchar(255),
    team_name varchar(255),
    team_abbreviation varchar(10),
    team_code varchar(255),
    team_slug varchar(255),
    is_nba_assigned boolean,
    nba_assigned_team_id boolean,
    games_played_flag varchar(1)
);

create index idx_inwpl_team on intake.wplayer(team_id);

create table intake.gm_player (
    season_id int not null,
    player_id bigint not null,
    player_name varchar(255),
    team_id bigint not null,
    team_abbreviation varchar(3),
    team_name varchar(255),
    game_id bigint not null,
    game_date date,
    matchup varchar(50),
    wl varchar(1),
    min int,
    fgm int,
    fga int,
    fg_pct numeric(5, 4),
    fg3m int,
    fg3a int,
    fg3_pct numeric(5, 4),
    ftm int,
    fta int,
    ft_pct numeric(5, 4),
    oreb int,
    dreb int,
    reb int,
    ast int,
    stl int,
    blk int,
    tov int,
    pf int,
    pts int,
    plus_minus int,
    fantasy_pts numeric(5, 2),
    video_available smallint,
    primary key (game_id, player_id)
);

create index idx_ingmpl_team_id on intake.gm_player(team_id);
create index idx_ingmpl_team_abbr on intake.gm_player(team_abbreviation);
create index idx_ingmpl_szn on intake.gm_player(season_id);
create index idx_ingmpl_gdate on intake.gm_player(game_date);

create table intake.gm_team (
    season_id int not null,
    team_id bigint not null,
    team_abbreviation varchar(3),
    team_name varchar(255),
    game_id bigint not null,
    game_date date,
    matchup varchar(50),
    wl varchar(1),
    min int,
    fgm int,
    fga int,
    fg_pct numeric(5, 4),
    fg3m int,
    fg3a int,
    fg3_pct numeric(5, 4),
    ftm int,
    fta int,
    ft_pct numeric(5, 4),
    oreb int,
    dreb int,
    reb int,
    ast int,
    stl int,
    blk int,
    tov int,
    pf int,
    pts int,
    plus_minus int,
    video_available smallint,
    primary key (game_id, team_id)
);

create index idx_ingmtm_team_abbr on intake.gm_team(team_abbreviation);
create index idx_ingmtm_szn on intake.gm_team(season_id);
create index idx_ingmtm_gdate on intake.gm_team(game_date);