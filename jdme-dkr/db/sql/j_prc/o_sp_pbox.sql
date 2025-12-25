-- work in progress, need to make season load first

create or replace procedure stats.sp_pbox()
language plpgsql
as $$
begin
    insert into stats.pbox
        select
            season_id,
            team_id, 
            game_id,
            player_id,
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
        from intake.gm_player
    on conflict (team_id, game_id, player_id) do nothing;
end; $$;
-- call stats.sp_pbox();
-- select * from stats.pbox;
create or replace procedure stats.sp_pbox()
language plpgsql
as $$
begin
    insert into stats.pbox
        select
            a.season_id,
            a.team_id, 
            a.game_id,
            a.player_id,
            a.game_date,
            a.matchup,
            a.wl,
            a.min,
            a.pts,
            a.ast,
            a.reb,
            a.stl,
            a.blk,
            a.tov,
            a.oreb,
            a.dreb,
            a.pf,
            a.plus_minus,
            a.fgm, 
            a.fga,
            a.fg_pct,
            a.fg3m,
            a.fg3a,
            a.fg3_pct,
            a.ftm,
            a.fta,
            a.ft_pct
        from intake.gm_player a 
        inner join lg.plr b on b.player_id = a.player_id
    on conflict (team_id, game_id, player_id) do nothing;
end; $$;