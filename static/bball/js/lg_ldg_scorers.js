import { tblCaption } from "./table.js"
import { base } from "./listen.js"
import { playerBtnListener } from "./ui.js"

// build top x players table
export async function makeScoringLeaders(numPl) {
    const r = await fetch(`${base}/league/scoring-leaders?num=${numPl}`);
    if (!r.ok) {
        throw new Error(`${r.status}: error calling /league/scoring-leaders`)
    }
    const data = await r.json();
    console.log(data);
    await buildLeadingScorersTbl(data, 'top_lg_players', numPl)
}

/*
build table with top numPl leading scorers in the nba and wnba for their current
respective seasons 
*/
export async function buildLeadingScorersTbl(data, elName, numPl) {
    const tblcont = document.getElementById(elName);
    const tbl = document.getElementById('nba_tstbl');
    const thead = document.createElement('thead');
    const nbaH = document.createElement('td');
    const ptsH = document.createElement('td');
    const wnbaH = document.createElement('td');
    const wptsH = document.createElement('td');

    const caption = `Scoring Leaders | Current NBA/WNBA Top ${numPl}`
    tblCaption(tbl, caption);
    
    nbaH.textContent = `nba | ${data.nba[0].season}`;
    ptsH.textContent = 'points';
    wnbaH.textContent = `wnba | ${data.wnba[0].season}`;
    wptsH.textContent = 'points';

    thead.appendChild(nbaH);
    thead.appendChild(ptsH);
    thead.appendChild(wnbaH);
    thead.appendChild(wptsH);

    tbl.appendChild(thead);
    
    for (let i = 0; i < numPl; i++) {
        await lgTopScorerRow(tbl, data.nba[i], data.wnba[i]);
    }

    tblcont.appendChild(tbl);
}

/* 
add a row to the league top scorers table. called within a loop
adds nba player with button, their points, wnba player with button, their points
*/
export async function lgTopScorerRow(tbl, nba, wnba) {
    let r = document.createElement('tr');

    let pName = document.createElement('td');
    let pts = document.createElement('td');
    let wpName = document.createElement('td');
    let wpts = document.createElement('td');

    let btn = document.createElement('button');
    btn.textContent = `${nba.player} | ${nba.team}`;
    btn.type = 'button';
    btn.addEventListener('click', async () => {
        await playerBtnListener(nba.player);
    }); 

    let wbtn = document.createElement('button');
    wbtn.textContent = `${wnba.player} | ${wnba.team}`;
    wbtn.type = 'button';

    wbtn.addEventListener('click', async () => {
        await playerBtnListener(wnba.player);
    }); 

    pName.appendChild(btn);
    wpName.appendChild(wbtn);

    pts.textContent = nba.points;
    wpts.textContent = wnba.points;
    r.appendChild(pName);
    r.appendChild(pts);
    r.appendChild(wpName);
    r.appendChild(wpts);

    tbl.appendChild(r);
}