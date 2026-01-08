import { base } from "./listen.js"
import { playerBtnListener } from "./player_search.js"
import { table5f } from "./dynamic_table.js"; 
import { bytes_in_resp } from "./util.js";

// build top x players table
export async function makeScoringLeaders(numPl) {
    const url = `${base}/league/scoring-leaders?num=${numPl}`;
    const r = await fetch(url);
    if (!r.ok) {
        console.error(` error calling ${url}`);
    }
    console.trace(`%c ${await bytes_in_resp(r)} bytes received from ${url}}`, 'color: green; font-weight: bold;')
    const data = await r.json();
    console.log(data);
    // await buildLeadingScorersTbl(data, 'top_lg_players', numPl);
    await table5f(data, 'nba_tstbl', 
        `Scoring Leaders | NBA/WNBA Top ${numPl}`, 
        ["rank", `nba | ${data.nba[0].season}`, "points", 
         `wnba | ${(data.wnba[0].season).substring(0, 4)}`, "points"],
        numPl, lgTopScorerRow,
    )
}

// add a row to the league top scorers table. called within a loop
// adds nba player with button, their points, wnba player with button, their points
export async function lgTopScorerRow(tbl, data, i) {
    let r = document.createElement('tr');

    let rank = document.createElement('td');
    let pName = document.createElement('td');
    let pts = document.createElement('td');
    let wpName = document.createElement('td');
    let wpts = document.createElement('td');
    let btn = document.createElement('button');
    let wbtn = document.createElement('button');

    rank.textContent = i + 1;
    btn.textContent = `${data.nba[i].player} | ${data.nba[i].team}`;
    btn.type = 'button';
    btn.addEventListener('click', async () => {
        await playerBtnListener(data.nba[i].player);
    }); 
    pName.appendChild(btn);

    wbtn.textContent = `${data.wnba[i].player} | ${data.wnba[i].team}`;
    wbtn.type = 'button';
    wbtn.addEventListener('click', async () => {
        await playerBtnListener(data.wnba[i].player);
    }); 
    wpName.appendChild(wbtn);

    pts.textContent = data.nba[i].points;
    wpts.textContent = data.wnba[i].points;

    r.appendChild(rank);
    r.appendChild(pName);
    r.appendChild(pts);
    r.appendChild(wpName);
    r.appendChild(wpts);

    tbl.appendChild(r);
}