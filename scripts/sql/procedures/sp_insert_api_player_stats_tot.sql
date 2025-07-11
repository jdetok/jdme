DELIMITER $$
CREATE OR REPLACE PROCEDURE sp_insert_api_player_stats_tot()
BEGIN
	
	set session foreign_key_checks=0;
	DELETE FROM api_player_stats where stat_type = 'tot';
	
	INSERT INTO api_player_stats 
		SELECT 
			a.player_id,
			a.team_id,
			c.lg,
			a.season_id,
			"tot",
			c.player,
			e.team,
			e.team_name,
			count(a.game_id),
			sum(a.mins),
			sum(a.pts),
			sum(a.ast),
			sum(a.reb),
			sum(a.stl),
			sum(a.blk),
			sum(b.fgm),
			sum(b.fga),
			case when sum(b.fgm) > 0
				then concat(round((sum(b.fgm) / sum(b.fga)) * 100, 2), "%")
				else "0%"
			end as fg_pct,
			sum(b.fg3m),
			sum(b.fg3a),
			case when sum(b.fg3m) > 0
				then concat(round((sum(b.fg3m) / sum(b.fg3a)) * 100, 2), "%")
				else "0%"
			end as fg3_pct,
			sum(b.ftm),
			sum(b.fta),
			case when sum(b.ftm) > 0
				then concat(round((sum(b.ftm) / sum(b.fta)) * 100, 2), "%")
				else "0%"
			end as ft_pct
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
		GROUP BY a.player_id, a.season_id;
		
	INSERT INTO api_player_stats
		SELECT 
			a.player_id,
			c.team_id,
			c.lg,
			99999,
			"tot",
			c.player,
			e.team,
			e.team_name,
			count(a.game_id),
			sum(a.mins),
			sum(a.pts),
			sum(a.ast),
			sum(a.reb),
			sum(a.stl),
			sum(a.blk),
			sum(b.fgm),
			sum(b.fga),
			case when sum(b.fgm) > 0
				then concat(round((sum(b.fgm) / sum(b.fga)) * 100, 2), "%")
				else "0%"
			end as fg_pct,
			sum(b.fg3m),
			sum(b.fg3a),
			case when sum(b.fg3m) > 0
				then concat(round((sum(b.fg3m) / sum(b.fg3a)) * 100, 2), "%")
				else "0%"
			end as fg3_pct,
			sum(b.ftm),
			sum(b.fta),
			case when sum(b.ftm) > 0
				then concat(round((sum(b.ftm) / sum(b.fta)) * 100, 2), "%")
				else "0%"
			end as ft_pct
		FROM p_box a
		INNER JOIN p_shtg b
			ON b.player_id = a.player_id
			AND b.game_id = a.game_id
		INNER JOIN player c
			ON c.player_id = a.player_id
		INNER JOIN team e
			ON e.team_id = c.team_id
		-- joined to player instead of p_box to get recent team
		WHERE c.lg <> "GNBA"
		AND LEFT(a.season_id, 1) in ('2', '4')
		GROUP BY a.player_id;
	
	INSERT INTO api_player_stats
		SELECT 
			a.player_id,
			c.team_id,
			c.lg,
			case when LEFT(a.season_id, 1) = '2'
				then 99998
				else 99997
			end,
			"tot",
			c.player,
			e.team,
			e.team_name,
			count(a.game_id),
			sum(a.mins),
			sum(a.pts),
			sum(a.ast),
			sum(a.reb),
			sum(a.stl),
			sum(a.blk),
			sum(b.fgm),
			sum(b.fga),
			case when sum(b.fgm) > 0
				then concat(round((sum(b.fgm) / sum(b.fga)) * 100, 2), "%")
				else "0%"
			end as fg_pct,
			sum(b.fg3m),
			sum(b.fg3a),
			case when sum(b.fg3m) > 0
				then concat(round((sum(b.fg3m) / sum(b.fg3a)) * 100, 2), "%")
				else "0%"
			end as fg3_pct,
			sum(b.ftm),
			sum(b.fta),
			case when sum(b.ftm) > 0
				then concat(round((sum(b.ftm) / sum(b.fta)) * 100, 2), "%")
				else "0%"
			end as ft_pct
		FROM p_box a
		INNER JOIN p_shtg b
			ON b.player_id = a.player_id
			AND b.game_id = a.game_id
		INNER JOIN player c
			ON c.player_id = a.player_id
		INNER JOIN team e
			ON e.team_id = c.team_id
		-- joined to player instead of p_box to get recent team
		WHERE c.lg <> "GNBA"
		AND LEFT(a.season_id, 1) in ('2', '4')
		GROUP BY a.player_id, LEFT(a.season_id, 1);
	set session foreign_key_checks=1;
END$$
DELIMITER ;

call sp_insert_api_player_stats_tot();
select * from api_player_stats where player = "LeBron James";
select * from api_player_stats where season_id = 22024 and player_id = 2544;
