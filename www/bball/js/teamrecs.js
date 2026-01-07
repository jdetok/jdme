import { base } from "./listen.js"
import { table5f } from "./dynamic_table.js";

// build top x players table
export async function getTeamRecords() {
    try {
        const r = await fetch(`${base}/teamrecs`);
        const data = await r.json();
        return data;
    } catch(err) {
        throw new Error(`error calling /teamrecs`);
    }
}

export async function makeTeamRecsTable(numRecs) {
    const data = await getTeamRecords();
    await table5f(data, 'trtbl', `NBA/WNBA Regular Season Team Records`,
        ["rank", `nba | ${data.nba_team_records[0].season}`, "record", 
            `wnba | ${data.wnba_team_records[0].season}`, "record"],
        numRecs, teamRecsRow)
}

// add a row to the league top scorers table. called within a loop
// adds nba player with button, their points, wnba player with button, their points
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