import { playerBtnListener } from "./player_search.js"
import { table5f } from "./dynamic_table.js"; 

export async function makeRGTopScorersTbl(data, numPl) {
    await table5f(data, 'tstbl', `Scoring Leaders | ${data.recent_games[0].game_date} Games`, 
        ["rank", "name | team", "matchup", "wl | score", "points"], numPl, rgTopScorerRow
    );
}

export async function rgTopScorerRow(tbl, data, i) {
    let r = document.createElement('tr');

    let rank = document.createElement('td');
    let pName = document.createElement('td');
    let matchup = document.createElement('td');
    let wlscore = document.createElement('td');
    let pts = document.createElement('td');
    let btn = document.createElement('button');

    let player = data.top_scorers[i];
    let game = data.recent_games.find(g => g.player_id === player.player_id);
    
    rank.textContent = i + 1;

    btn.textContent = `${player.player} | ${game.team}`;
    btn.type = 'button';
    btn.addEventListener('click', async () => {
        await playerBtnListener(player.player);
    }); 
    pName.appendChild(btn);

    matchup.textContent = game.matchup;
    wlscore.textContent = `${game.wl} | ${game.points}-${game.opp_points}`
    pts.textContent = player.points;

    r.appendChild(rank);
    r.appendChild(pName);
    r.appendChild(matchup);
    r.appendChild(wlscore);
    r.appendChild(pts);

    tbl.appendChild(r);
}