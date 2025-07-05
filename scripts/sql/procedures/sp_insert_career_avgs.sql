-- select and insert into season_totals 
DELIMITER $$
CREATE OR REPLACE PROCEDURE sp_insert_career_avgs()
BEGIN
	DELETE FROM career_avgs;
	
	INSERT INTO career_avgs
		SELECT 
			a.player_id,
			c.team_id,
			c.lg,
			case
				when LEFT(a.season_id, 1) = '2'
				then "Regular Season"
				else "Playoffs"
			end as season_type,
			c.player,
			e.team,
			e.team_name,
			c.active,
			count(a.game_id),
			round(avg(a.pts), 2),
			round(avg(a.ast), 2),
			round(avg(a.reb), 2),
			round(avg(a.stl), 2),
			round(avg(a.blk), 2),
			
			round(avg(b.fgm), 2),
			round(avg(b.fga), 2),
			concat(round(avg(b.fg_pct) * 100, 2), "%"),

			round(avg(b.fg3m), 2),
			round(avg(b.fg3a), 2),
			concat(round(avg(b.fg3_pct) * 100, 2), "%"),
			
			round(avg(b.ftm), 2),
			round(avg(b.fta), 2),
			concat(round(avg(b.ft_pct) * 100, 2), "%")
			
		FROM p_box a
		INNER JOIN p_shtg b
			ON b.player_id = a.player_id
			AND b.game_id = a.game_id
		INNER JOIN player c
			ON c.player_id = a.player_id
		INNER JOIN season d
		 	ON d.season_id = a.season_id
		INNER JOIN team e
			ON e.team_id = c.team_id
		-- joined to player instead of p_box to get recent team
		WHERE c.lg <> "GNBA"
		AND LEFT(a.season_id, 1) in ('2', '4')
		GROUP BY a.player_id, LEFT(a.season_id, 1);
END$$
DELIMITER ;

CALL sp_insert_career_avgs();

select count(*) from career_avgs;
select * from career_avgs;
