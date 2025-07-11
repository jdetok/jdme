-- select and insert into season_totals 
DELIMITER $$
CREATE OR REPLACE PROCEDURE sp_insert_top_scorers()
BEGIN
	DELETE FROM top_scorers;
	
	INSERT INTO top_scorers
		select 
			a.player_id,
			a.team_id,
			e.lg,
			a.season_id,
			a.game_id, 
			b.game_date,
			d.player,
			e.team,
			e.team_name,
			a.mins,
			a.pts,
			a.ast,
			a.reb,
			a.stl,
			a.blk,
			c.fgm,
			c.fga,
			case when c.fg_pct is not null and c.fgm > 0
				then concat(round((c.fgm / c.fga) * 100, 2), "%")
				else "0%"
			end as fg_pct,
			c.fg3m,
			c.fg3a,
			case when c.fg3_pct is not null and c.fg3m > 0
				then concat(round((c.fg3m / c.fg3a) * 100, 2), "%")
				else "0%"
			end as fg3_pct,
			c.ftm,
			c.fta,

			case when c.ft_pct is not null and c.ftm > 0
				then concat(round((c.ftm / c.fta) * 100, 2), "%")
				else "0%"
			end as ft_pct
		from p_box a
		inner join game b on b.game_id = a.game_id
			and b.game_date = (select max(z.game_date) from game z)
-- "2025-01-21"
-- 			(select max(z.game_date) from game z)
		inner join p_shtg c on c.player_id = a.player_id
			and c.game_id = a.game_id
		inner join player d on d.player_id = a.player_id
		inner join team e on e.team_id = a.team_id
			and e.lg <> "GNBA"
		group by a.player_id, a.game_id
		order by a.pts desc;
END$$
DELIMITER ;

CALL sp_insert_top_scorers();

-- select count(*) from top_scorers;
select * 
from top_scorers 
order by points desc, (assists + rebounds + steals + blocks) desc 
limit 1;