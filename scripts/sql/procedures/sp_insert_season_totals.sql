-- select and insert into season_totals 
DELIMITER $$
CREATE OR REPLACE PROCEDURE sp_insert_season_totals()
BEGIN
	DELETE FROM season_totals;
	
	INSERT INTO season_totals
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
			sum(a.pts), 
			sum(a.ast),
			sum(a.reb),
			sum(a.stl),
			sum(a.blk),
			sum(b.fgm),
			sum(b.fga),
			sum(b.fg3m),
			sum(b.fg3a),
			sum(b.ftm),
			sum(b.fta)
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

CALL sp_insert_season_totals();

select count(*) from season_totals;
select * from season_totals;