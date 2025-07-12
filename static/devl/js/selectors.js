import { base } from "./listen.js";

export async function sznChangeListener() {
    const rs = document.getElementById('rs_slct');
    const ps = document.getElementById('ps_slct');
    const cr = document.getElementById('cr_slct');
    rs.addEventListener('change', async (event) => {
        event.preventDefault();
        // await ;
    });
};

export async function loadSznOptions() {
    const r = await fetch(base + '/seasons');
    if (!r.ok) { 
        throw new Error(`HTTP Error: ${s.status}`);
    } 
    const data = await r.json();
    console.log(data);
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

async function makeOption(slct, txt, val) {
    let opt = document.createElement('option');
    opt.textContent = txt;
    opt.value = val;
    slct.appendChild(opt);
}
