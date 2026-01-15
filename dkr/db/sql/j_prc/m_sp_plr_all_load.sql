create or replace procedure lg.sp_plr_all_load()
language plpgsql
as $$
begin
    insert into lg.plr (
        lg_id, player_id, plr_cde, player, last_first, from_year, to_year)
        select 
            0, 
            person_id,
            playercode,
            display_first_last,
            display_last_comma_first,
            from_year,
            to_year
        from intake.player
		union
        select 
            1, 
            person_id,
            playercode,
            display_first_last,
            display_last_comma_first,
            from_year,
            to_year
        from intake.wplayer
    on conflict (player_id) do nothing;
end; $$;