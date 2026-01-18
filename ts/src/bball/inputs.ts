import { base, fetchJSON } from "../global.js";

export type checkGroup = {box: string, slct: string};
export async function checkBoxGroupValue(lgrp: checkGroup, rgrp: checkGroup, dflt: number | string
): Promise<string> {
    const l = await checkBoxes(lgrp);
    const r = await checkBoxes(rgrp);

    if (l) return l;
    if (r) return r;
    
    return String(dflt);
}

export async function checkBoxes(cg: checkGroup) {
    const b = document.getElementById(cg.box) as HTMLInputElement;
    const s = document.getElementById(cg.slct) as HTMLInputElement;
    if (!b || !s) {
        throw new Error(`couldn't get element with id ${cg.box} or ${cg.slct}`);
    }
    if (b.checked) {
        return s.value
    }
}

export async function clearCheckBoxes(boxes: string[]) {
    for (let i = 0; i < boxes.length; i++) {
        let b = document.getElementById(boxes[i]) as HTMLInputElement;
        if (!b) throw new Error(`couldn't find input element with id ${boxes[i]}`);
        b.checked = false;
    }
}

export function clearSearch(focus=false, elId = 'pSearch'): void {
    const pSearch = document.getElementById(elId) as HTMLInputElement;
    if (!pSearch) throw new Error(`couldn't get search bar element at ${elId}`);
    pSearch.value = '';
    if (focus) pSearch.focus();
}

export async function setPHold(player: string) {
    const ph = document.getElementById('pHold') as HTMLInputElement;
    ph.value = player;
}

export async function lgRadioBtns() {
    const selected = document.querySelector('input[name="leagues"]:checked') as HTMLInputElement;
    if (selected) {
        return selected.value;
    } else {
        return 'all';
    }
}

export async function getInputVals() {
    const lg = await lgRadioBtns();
    const season = await checkBoxGroupValue(
        {box: 'post', slct: 'ps_slct'}, 
        {box: 'reg', slct: 'rs_slct'}, 
        88888);
    const team = await checkBoxGroupValue(
        {box: 'nbaTm', slct: 'tm_slct'}, 
        {box: 'wnbaTm', slct: 'wTm_slct'}, 
        0);

    return { lg, season, team };
}

// SELECTORS
// append options to select
async function makeOption(slct: HTMLSelectElement, txt: string, val: string) {
    let opt = document.createElement('option');
    opt.textContent = txt;
    opt.value = val;
    opt.style.width = '100%';
    slct.appendChild(opt);
}

export type Season = {
    season_id: string,
    season: string,
    wseason: string,
};

export type Seasons = Season[];

export async function getSeasons(): Promise<Seasons> {
    return await fetchJSON(`${base}/seasons`);
}

// call seaons endpoint for the opts
export async function loadSznOptions() {
    const data = await getSeasons();
    await buildSznSelects(data);
}

// accept seasons in data object and make an option for each
async function buildSznSelects(data: Seasons) {
    const rs = document.getElementById('rs_slct') as HTMLSelectElement;
    const ps = document.getElementById('ps_slct') as HTMLSelectElement;
    for (let s of data) {
        if (s.season_id.substring(0, 1) === '4') {
            await makeOption(ps, s.season, s.season_id);
        } else if (s.season_id.substring(0, 1) === '2') {
            await makeOption(rs, s.season, s.season_id);
        }
    }
}

export type Team = {
    league: string,
    team_id: string,
    team: string,
    team_long: string,
};

export type Teams = Team[];

export async function getTeams(): Promise<Teams> {
    return await fetchJSON(`${base}/teams`);
}

// call seaons endpoint for the opts
export async function loadTeamOptions() {
    const data = await getTeams();
    await buildTeamSelects(data);
}

// accept seasons in data object and make an option for each
async function buildTeamSelects(data: Teams) {
    const nba = document.getElementById('tm_slct') as HTMLSelectElement;
    const wnba = document.getElementById('wTm_slct') as HTMLSelectElement;
    for (let t of data) {
        let txt = `${t.team} | ${t.team_long}`
        if (t.league === 'NBA') {
            await makeOption(nba, txt, t.team_id);
        } else if (t.league === 'WNBA') {
            await makeOption(wnba, txt, t.team_id);
        }
    }
}