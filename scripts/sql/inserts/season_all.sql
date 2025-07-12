select * from season;
insert into season values (99999, "all", "all", "all", "all");
insert into season values (99998, "rs", "rs", "rs", "rs");
insert into season values (99997, "ps", "ps", "ps", "ps");

insert into season values (99999, "rs", "Career Regular Season Statistics", 
    "all", "Career Regular Season Statistics");
insert into season values (99998, "ps", "Career Playoff Statistics", "rs", 
    "Career Playoff Statistics");
insert into season values (99997, "all", "Career Combined Reg-Season/Playoff Statistics", 
    "ps", "Career Combined Reg-Season/Playoff Statistics");

update season 
set season_desc = "Career Combined Reg-Season/Playoff Statistics", 
    wseason_desc = "Career Combined Reg-Season/Playoff Statistics"
where season_id = 99997;

update season 
set season_desc = "Career Regular Season Statistics", 
    wseason_desc = "Career Regular Season Statistics"
where season_id = 99999;

update season 
set season_desc = "Career Playoff Statistics", 
    wseason_desc = "Career Playoff Statistics"
where season_id = 99998;