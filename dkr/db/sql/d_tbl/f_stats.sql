--DECIDED TO REMOVE LEAGE. AVAILABLE WITH SINGLE JOIN TO PLAYER OR TEAM TABLE
create table stats.pbox (
    szn_id int references lg.szn(szn_id),
    team_id bigint not null references lg.team(team_id),
    game_id bigint not null,
    player_id bigint not null references lg.plr(player_id),
    gdate date,
    matchup varchar(50),
    wl varchar(1),
    mins int,
    pts int,
    ast int,
    reb int,
    stl int,
    blk int,
    tov int,
    oreb int,
    dreb int,
    foul int,
    pm int,
    fgm int,
    fga int,
    fgp numeric(5, 4),
    f3m int,
    f3a int,
    f3p numeric(5, 4),
    ftm int,
    fta int,
    ftp numeric(5, 4),
    primary key (team_id, game_id, player_id)
);

create index idx_pbox_szn on stats.pbox(szn_id);
create index idx_pbox_team on stats.pbox(team_id);
create index idx_pbox_gdate on stats.pbox(gdate);
create index idx_pbox_matchup on stats.pbox(matchup);
create index idx_pbox_wl on stats.pbox(wl);

create table stats.tbox (
    szn_id int references lg.szn(szn_id),
    team_id bigint not null references lg.team(team_id),
    game_id bigint not null,
    gdate date,
    matchup varchar(50),
    wl varchar(1),
    mins int,
    pts int,
    ast int,
    reb int,
    stl int,
    blk int,
    tov int,
    oreb int,
    dreb int,
    foul int,
    pm int,
    fgm int,
    fga int,
    fgp numeric(5, 4),
    f3m int,
    f3a int,
    f3p numeric(5, 4),
    ftm int,
    fta int,
    ftp numeric(5, 4),
    primary key (game_id, team_id)
);

create index idx_tbox_szn on stats.tbox(szn_id);
create index idx_tbox_gdate on stats.tbox(gdate);
create index idx_tbox_matchup on stats.tbox(matchup);
create index idx_tbox_wl on stats.tbox(wl);