-- PLAYER ID FROM NAME
USE bball;
CREATE OR REPLACE FUNCTION lg.get_player_id_by_name(p_name text)
RETURNS bigint
LANGUAGE sql
STABLE
AS $$
    SELECT player_id
    FROM lg.plr
    WHERE lower(trim(player)) = lower(trim(p_name))
    LIMIT 1;
$$;
