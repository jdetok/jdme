import { base, fetchJSON } from "../global.js";
export async function checkBoxGroupValue(lgrp, rgrp, dflt) {
    const l = await checkBoxes(lgrp);
    const r = await checkBoxes(rgrp);
    if (l)
        return l;
    if (r)
        return r;
    return String(dflt);
}
export async function checkBoxes(cg) {
    const b = document.getElementById(cg.box);
    const s = document.getElementById(cg.slct);
    if (!b || !s) {
        throw new Error(`couldn't get element with id ${cg.box} or ${cg.slct}`);
    }
    if (b.checked) {
        return s.value;
    }
}
export async function clearCheckBoxes(boxes) {
    for (let i = 0; i < boxes.length; i++) {
        let b = document.getElementById(boxes[i]);
        if (!b)
            throw new Error(`couldn't find input element with id ${boxes[i]}`);
        b.checked = false;
    }
}
export function clearSearch(focus = false, elId = 'pSearch') {
    const pSearch = document.getElementById(elId);
    if (!pSearch)
        throw new Error(`couldn't get search bar element at ${elId}`);
    pSearch.value = '';
    if (focus)
        pSearch.focus();
}
export async function setPHold(player) {
    const ph = document.getElementById('pHold');
    ph.value = player;
}
export async function lgRadioBtns() {
    const selected = document.querySelector('input[name="leagues"]:checked');
    if (selected) {
        return selected.value;
    }
    else {
        return 'all';
    }
}
export async function getInputVals() {
    const lg = await lgRadioBtns();
    const season = await checkBoxGroupValue({ box: 'post', slct: 'ps_slct' }, { box: 'reg', slct: 'rs_slct' }, 88888);
    const team = await checkBoxGroupValue({ box: 'nbaTm', slct: 'tm_slct' }, { box: 'wnbaTm', slct: 'wTm_slct' }, 0);
    return { lg, season, team };
}
// SELECTORS
// append options to select
async function makeOption(slct, txt, val) {
    let opt = document.createElement('option');
    opt.textContent = txt;
    opt.value = val;
    opt.style.width = '100%';
    slct.appendChild(opt);
}
export async function getSeasons() {
    return await fetchJSON(`${base}/seasons`);
}
// call seaons endpoint for the opts
export async function loadSznOptions() {
    const data = await getSeasons();
    await buildSznSelects(data);
}
// accept seasons in data object and make an option for each
async function buildSznSelects(data) {
    const rs = document.getElementById('rs_slct');
    const ps = document.getElementById('ps_slct');
    for (let s of data) {
        if (s.season_id.substring(0, 1) === '4') {
            await makeOption(ps, s.season, s.season_id);
        }
        else if (s.season_id.substring(0, 1) === '2') {
            await makeOption(rs, s.season, s.season_id);
        }
    }
}
export async function getTeams() {
    return await fetchJSON(`${base}/teams`);
}
// call seaons endpoint for the opts
export async function loadTeamOptions() {
    const data = await getTeams();
    await buildTeamSelects(data);
}
// accept seasons in data object and make an option for each
async function buildTeamSelects(data) {
    const nba = document.getElementById('tm_slct');
    const wnba = document.getElementById('wTm_slct');
    for (let t of data) {
        let txt = `${t.team} | ${t.team_long}`;
        if (t.league === 'NBA') {
            await makeOption(nba, txt, t.team_id);
        }
        else if (t.league === 'WNBA') {
            await makeOption(wnba, txt, t.team_id);
        }
    }
}
//# sourceMappingURL=inputs.js.map