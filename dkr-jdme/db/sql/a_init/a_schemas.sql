-- tables to insert the data from stats.nba.com
create schema intake;

-- core league, season, team, player data
create schema lg;

-- stats tables, all should reference league schema
create schema stats;

-- api view tables, will need to match the mariadb ones
create schema api;

