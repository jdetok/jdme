DELIMITER $$
CREATE OR REPLACE PROCEDURE sp_insert_api_player_stats_avg()
BEGIN
	
	set session foreign_key_checks=0;
	DELETE FROM api_player_stats where stat_type = 'avg';
	-- BY SEASON
	INSERT INTO api_player_stats 
		SELECT 
			a.player_id,
			a.team_id,
			c.lg,
			a.season_id,
			"avg",
			c.player,
			e.team,
			e.team_name,
			count(distinct a.game_id),
			avg(a.mins),
			avg(a.pts),
			avg(a.ast),
			avg(a.reb),
			avg(a.stl),
			avg(a.blk),
			avg(b.fgm),
			avg(b.fga),
			case when sum(b.fgm) > 0
				then concat(round((avg(b.fgm) / avg(b.fga)) * 100, 2), "%")
				else "0%"
			end as fg_pct,
			avg(b.fg3m),
			avg(b.fg3a),
			case when sum(b.fg3m) > 0
				then concat(round((avg(b.fg3m) / avg(b.fg3a)) * 100, 2), "%")
				else "0%"
			end as fg3_pct,
			avg(b.ftm),
			avg(b.fta),
			case when sum(b.ftm) > 0
				then concat(round((avg(b.ftm) / avg(b.fta)) * 100, 2), "%")
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
	
	-- ALL STATS REGARDLESS OF SEASON TYPE
	INSERT INTO api_player_stats
		SELECT 
			a.player_id,
			c.team_id,
			c.lg,
			99999,
			"avg",
			c.player,
			e.team,
			e.team_name,
			count(distinct a.game_id),
			avg(a.mins),
			avg(a.pts),
			avg(a.ast),
			avg(a.reb),
			avg(a.stl),
			avg(a.blk),
			avg(b.fgm),
			avg(b.fga),
			case when sum(b.fgm) > 0
				then concat(round((avg(b.fgm) / avg(b.fga)) * 100, 2), "%")
				else "0%"
			end as fg_pct,
			avg(b.fg3m),
			avg(b.fg3a),
			case when sum(b.fg3m) > 0
				then concat(round((avg(b.fg3m) / avg(b.fg3a)) * 100, 2), "%")
				else "0%"
			end as fg3_pct,
			avg(b.ftm),
			avg(b.fta),
			case when sum(b.ftm) > 0
				then concat(round((avg(b.ftm) / avg(b.fta)) * 100, 2), "%")
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
			ON e.team_id = c.team_id
		-- joined to player instead of p_box to get recent team
		WHERE c.lg <> "GNBA"
		AND LEFT(a.season_id, 1) in ('2', '4')
		GROUP BY a.player_id;
	
	-- REG/PS TOTALS
	INSERT INTO api_player_stats
		SELECT 
			a.player_id,
			c.team_id,
			c.lg,
			case when LEFT(a.season_id, 1) = '2'
				then 99998
				else 99997
			end,
			"avg",
			c.player,
			e.team,
			e.team_name,
			count(distinct a.game_id),
			avg(a.mins),
			avg(a.pts),
			avg(a.ast),
			avg(a.reb),
			avg(a.stl),
			avg(a.blk),
			avg(b.fgm),
			avg(b.fga),
			case when sum(b.fgm) > 0
				then concat(round((avg(b.fgm) / avg(b.fga)) * 100, 2), "%")
				else "0%"
			end as fg_pct,
			avg(b.fg3m),
			avg(b.fg3a),
			case when sum(b.fg3m) > 0
				then concat(round((avg(b.fg3m) / avg(b.fg3a)) * 100, 2), "%")
				else "0%"
			end as fg3_pct,
			avg(b.ftm),
			avg(b.fta),
			case when sum(b.ftm) > 0
				then concat(round((avg(b.ftm) / avg(b.fta)) * 100, 2), "%")
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
			ON e.team_id = c.team_id
		-- joined to player instead of p_box to get recent team
		WHERE c.lg <> "GNBA"
		AND LEFT(a.season_id, 1) in ('2', '4')
		GROUP BY a.player_id, LEFT(a.season_id, 1);
		
	set session foreign_key_checks=1;
END$$
DELIMITER ;

-- CALL PROC
call sp_insert_api_player_stats_avg();

-- TESTS
select * from api_player_stats;