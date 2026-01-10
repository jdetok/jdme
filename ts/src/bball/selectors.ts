import { base, checkBoxEls, foldedLog, MSG, scrollIntoBySize } from "../global.js";

// append options to select
async function makeOption(slct, txt, val) {
    let opt = document.createElement('option');
    opt.textContent = txt;
    opt.value = val;
    opt.style.width = '100%';
    slct.appendChild(opt);
}

// call seaons endpoint for the opts
export async function loadSznOptions() {
    const url = base + '/seasons';
    const r = await fetch(url);
    if (!r.ok) throw new Error(`HTTP Error from ${url}`);
     
    const data = await r.json();
    await buildSznSelects(data);
}

// accept seasons in data object and make an option for each
async function buildSznSelects(data) {
    const rs = document.getElementById('rs_slct');
    const ps = document.getElementById('ps_slct');
    for (let s of data) {
        if (s.season_id.substring(0, 1) === '4') {
            await makeOption(ps, s.season, s.season_id);
        } else if (s.season_id.substring(0, 1) === '2') {
            await makeOption(rs, s.season, s.season_id);
        }
    }
}

// call seaons endpoint for the opts
export async function loadTeamOptions() {
    const url = base + '/teams';
    const r = await fetch(url);
    if (!r.ok) throw new Error(`HTTP Error from ${url}`);
    const data = await r.json();
    await buildTeamSelects(data);
}

// accept seasons in data object and make an option for each
async function buildTeamSelects(data) {
    const nba = document.getElementById('tm_slct');
    const wnba = document.getElementById('wTm_slct');
    for (let t of data) {
        let txt = `${t.team} | ${t.team_long}`
        if (t.league === 'NBA') {
            await makeOption(nba, txt, t.team_id);
        } else if (t.league === 'WNBA') {
            await makeOption(wnba, txt, t.team_id);
        }
    }
}