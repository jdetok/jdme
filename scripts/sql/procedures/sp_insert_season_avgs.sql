-- select and insert into season_totals 
DELIMITER $$
CREATE OR REPLACE PROCEDURE sp_insert_season_avgs()
BEGIN
	DELETE FROM season_avgs;
	
	INSERT INTO season_avgs
		SELECT 
			a.player_id,
			a.team_id,
			c.lg,
			a.season_id,
			case when c.lg = "WNBA"
	        	then d.wseason_desc
	        	else d.season_desc
	    	end as season,
			c.player,
			e.team,
			e.team_name,
			c.active,
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
			ON e.team_id = a.team_id
			
		WHERE c.lg <> "GNBA"
		AND LEFT(a.season_id, 1) in ('2', '4')
		GROUP BY a.player_id, a.team_id, a.season_id;
END$$
DELIMITER ;

CALL sp_insert_season_avgs();

select count(*) from season_avgs;
select * from season_avgs;
