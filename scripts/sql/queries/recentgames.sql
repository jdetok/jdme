select a.game_id, a.team_id, e.player_id, f.player, b.lg, c.team, c.team_name,
b.game_date, b.matchup, b.final, b.ot, a.pts
from t_box a
inner join game b on b.game_id = a.game_id
inner join team c on c.team_id = a.team_id
inner join p_box d on d.game_id = a.game_id
	and d.team_id = a.team_id
inner join (
	select game_id, player_id, team_id
	from p_box
	group by game_id, team_id
	order by pts desc
	limit 1
) e on e.game_id = a.game_id and e.team_id = a.team_id
inner join player f on f.player_id = e.player_id and f.team_id = a.team_id
where b.game_date = (select max(game_date) from game)
group by a.game_id, a.team_id;