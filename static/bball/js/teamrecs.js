import { tblCaption } from "./table.js"
import { base } from "./listen.js"

// build top x players table
export async function getTeamRecords() {
    try {
        const r = await fetch(`${base}/teamrecs`);
        const data = await r.json();
        return data;
    } catch(err) {
        throw new Error(`${r.status}: error calling /teamrecs`);
    }
}

/*
build table with top numPl leading scorers in the nba and wnba for their current
respective seasons 
*/
export async function buildTeamRecsTbl(data, elName) {
    console.log("build team recs table");
    const tblcont = document.getElementById(elName);
    const tbl = document.getElementById('trtbl');
    tbl.textContent = '';

    const thead = document.createElement('thead');
    const rankH = document.createElement('td');
    const nbaH = document.createElement('td');
    const recH = document.createElement('td');
    const wnbaH = document.createElement('td');
    const wrecH = document.createElement('td');

    const capMsg = `NBA/WNBA Regular Season Team Records`;
    tblCaption(tbl, capMsg);

    rankH.textContent = 'rank';
    nbaH.textContent = `nba | ${data.nba_team_records[0].season}`;
    recH.textContent = 'record';
    wnbaH.textContent = `wnba | ${data.wnba_team_records[0].season}`;
    wrecH.textContent = 'record';
    console.log(data.wnba_team_records[0].season);

    thead.appendChild(rankH);
    thead.appendChild(nbaH);
    thead.appendChild(recH);
    thead.appendChild(wnbaH);
    thead.appendChild(wrecH);
    tbl.appendChild(thead);
    
    for (let i = 0; i < 5; i++) {
        await teamRecsRow(tbl, data, i);
    }
    tblcont.appendChild(tbl);
}

/* 
add a row to the league top scorers table. called within a loop
adds nba player with button, their points, wnba player with button, their points
*/
export async function teamRecsRow(tbl, data, i) {
    const nba = data.nba_team_records[i];
    const wnba = data.wnba_team_records[i];

    let r = document.createElement('tr');
    let rank = document.createElement('td');
    let tName = document.createElement('td');
    let rec = document.createElement('td');
    let wtName = document.createElement('td');
    let wrec = document.createElement('td');

    rank.textContent = i + 1;
    tName.textContent = nba.team_long;
    wtName.textContent = wnba.team_long;
    rec.textContent = `${nba.wins}-${nba.losses}`;
    wrec.textContent = `${wnba.wins}-${wnba.losses}`;

    r.appendChild(rank);
    r.appendChild(tName);
    r.appendChild(rec);
    r.appendChild(wtName);
    r.appendChild(wrec);

    tbl.appendChild(r);
}