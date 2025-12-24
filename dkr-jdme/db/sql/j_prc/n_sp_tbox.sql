-- work in progress, need to make season load first

create or replace procedure stats.sp_tbox()
language plpgsql
as $$
begin
    insert into stats.tbox
        select
            season_id,
            team_id, 
            game_id,
            game_date,
            matchup,
            wl,
            min,
            pts,
            ast,
            reb,
            stl,
            blk,
            tov,
            oreb,
            dreb,
            pf,
            plus_minus,
            fgm, 
            fga,
            fg_pct,
            fg3m,
            fg3a,
            fg3_pct,
            ftm,
            fta,
            ft_pct
        from intake.gm_team -- join here to lg.team /lg.player
    on conflict (game_id, team_id) do nothing;
end; $$;
-- call stats.sp_tbox();
-- select * from stats.tbox;