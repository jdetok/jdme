import { base } from "./listen.js";

async function makeOption(slct, txt, val) {
    let opt = document.createElement('option');
    opt.textContent = txt;
    opt.value = val;
    slct.appendChild(opt);
}

export async function loadSznOptions() {
    const r = await fetch(base + '/seasons');
    if (!r.ok) { 
        throw new Error(`HTTP Error: ${s.status}`);
    } 
    const data = await r.json();
    await buildSznSelects(data);
}

async function buildSznSelects(data) {
    const rs = document.getElementById('rs_slct');
    const ps = document.getElementById('ps_slct');
    const cr = document.getElementById('cr_slct');
    for (let s of data) {
        if (s.season_id > 99990) {
            await makeOption(cr, s.season, s.season_id);
        } else if (s.season_id.substring(0, 1) === '4') {
            await makeOption(ps, s.season, s.season_id);
        } else if (s.season_id.substring(0, 1) === '2') {
            await makeOption(rs, s.season, s.season_id);
        }
    }
}

export async function loadAllTeamOpts() {
    const r = await fetch(base + '/teams');
    if (!r.ok) { 
            throw new Error(`HTTP Error: ${r.status}`);
        } // CONVERT SUCCESSFUL RESPONSE TO JSON
    const data = await r.json();
    if (data[0] == '') {
        console.log('empty json');
    }
    const w = document.getElementById('wnba_teams');//.value.trim()
    const n = document.getElementById('nba_teams');//.value.trim()
    const def = document.createElement('option');
    const defW= document.createElement('option');

    def.textContent = `All NBA Teams`
    def.value = 0;
    n.appendChild(def);
    defW.textContent = `All WNBA Teams`
    defW.value = 0;
    w.appendChild(defW);
    let i;
    for (i=0; i<data.length; i++){
        let opt = document.createElement('option');
        let optW = document.createElement('option');
        if (data[i].league === "NBA") {
            opt.textContent = data[i].team_long;
            opt.value = data[i].team_id;
            n.appendChild(opt);
        }
        if (data[i].league === "WNBA") {
            optW.textContent = data[i].team_long;
            optW.value = data[i].team_id;
            w.appendChild(optW);
        }
    }   
};