import { getP } from "./resp.js";
import { base } from "./listen.js";

export async function search() {
    const frm = document.getElementById('ui');
    frm.addEventListener('submit', async (event) => {
        event.preventDefault();

        const input = document.getElementById('pSearch');
        const player = input.value.trim();
        const season = await handleSeasonBoxes();
        console.log(`searhing for season ${season}`)
        await getP(base, player, season, '0');
        input.value = ''; // clear input box after searching
    }) 
}

export async function showHideHvr(el, hvrName, msg) {
    const hvr = document.getElementById(hvrName);
    el.addEventListener('mouseover', async (event) => {
        event.preventDefault();
        hvr.textContent = msg;
        hvr.style.display = 'block'; 
    })
    el.addEventListener('mouseleave', async (event) => {
        event.preventDefault();
        hvr.textContent = '';
        hvr.style.display = 'none'; 
    })
}

export async function selHvr() {
    const rs = document.getElementById('hlpRs');
    const ps = document.getElementById('hlpPs');
    await showHideHvr(rs, 'selhvr',
        `search for a specific regular-season. if the player being searched didn't
        play in the selected season, their first or most recent season, whichever
        is closer to the selected, will be used`);
    await showHideHvr(ps, 'selhvr',
        `search for a specific post-season. if the player being searched didn't
        play in the selected season, their first or most recent season, whichever
        is closer to the selected, will be used`);
}

// read pHold invisible val to add on-screen player's name to search bar
export async function holdPlayerBtn() {
    const btn = document.getElementById('holdP');
    const hlp = document.getElementById('hlpHld');
    btn.addEventListener('click', async (event) => {
        event.preventDefault();
        let player = document.getElementById('pHold').value;
        document.getElementById('pSearch').value = player;
    })
    await showHideHvr(
        hlp, 
        'hvrmsg',
        `fill the input box with the current player's name`
    )
}

export async function randPlayerBtn() {
    const btn = document.getElementById('randP');
    const hlp = document.getElementById('hlpRnd');
    btn.addEventListener('click', async (event) => {        
        event.preventDefault();
        const season = await handleSeasonBoxes();
        console.log(season);
        await getP(base, 'random', season, 0);
    })
    await showHideHvr(
        hlp, 
        'hvrmsg',
        `get the stats for a random player in the selected season. if no season 
        is specified, the current/most recent season will be used. if the 
        random player did not play in the selected season, their most 
        recent (or first, whichever is closer) season will be used`
    )
}

export async function handleSeasonBoxes() {
    // const c = await checkBoxes('career', 'cr_slct')
    const p = await checkBoxes('post', 'ps_slct')
    const r = await checkBoxes('reg', 'rs_slct')

    if (p) {
        return p;
    }
    if (r) {
        return r;
    }
    // return 77777;
    return `2${new Date().getFullYear() - 1}`;
/* FIX THIS - PASS 77777 after figuring out how to handle on server
    if (new Date().getMonth() < 7) {
        return `2${new Date().getFullYear() - 1}`;
    } else {
        return `2${new Date().getFullYear()}`;
    }
    
    // return 88888; */
}

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
    // const cr = document.getElementById('cr_slct');
    for (let s of data) {
        // if (s.season_id.substring(1, 4) === '9999') {
        //     await makeOption(cr, s.season, s.season_id);
        if (s.season_id.substring(0, 1) === '4') {
            await makeOption(ps, s.season, s.season_id);
        } else if (s.season_id.substring(0, 1) === '2') {
            await makeOption(rs, s.season, s.season_id);
        }
    }

}

export async function clearSearch() {
    const btn = document.getElementById('clearS');
    btn.addEventListener('click', async (event) => {
        event.preventDefault();
        document.getElementById('pSearch').value = '';
    })
}

export async function checkBoxes(box, sel) {
    const b = document.getElementById(box);
    const s = document.getElementById(sel);
    if (b.checked) {
        return s.value
    }
}