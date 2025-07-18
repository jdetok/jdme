import { getP } from "./pdash.js";
import { base, showHideHvr } from "./listen.js";

// read pHold invisible val to add on-screen player's name to search bar
export async function holdPlayerBtn() {
    const btn = document.getElementById('holdP');
    
    btn.addEventListener('click', async (event) => {
        event.preventDefault();
        let player = document.getElementById('pHold').value;
        document.getElementById('pSearch').value = player;
    })
    await showHideHvr(
        btn, 
        'hvrmsg',
        `fill the input box with the current player's name`
    )
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
export async function handleSeasonBoxes() {
    const c = await checkBoxes('career', 'cr_slct')
    const p = await checkBoxes('post', 'ps_slct')
    const r = await checkBoxes('reg', 'rs_slct')
    if (c) {
        return c;
    }
    if (p) {
        return p;
    }
    if (r) {
        return r;
    }
    return 88888;
}

export async function randPlayerBtn() {
    const btn = document.getElementById('randP');
    btn.addEventListener('click', async (event) => {        
        event.preventDefault();
        const season = await handleSeasonBoxes();
        console.log(season);
        await getP(base, 'random', season, 0);
    })
    await showHideHvr(
        btn, 
        'hvrmsg',
        `get the stats for a random player in the selected season. if no season 
        is specified, the current/most recent season will be used. if the 
        random player did not play in the selected season, their most 
        recent (or first, whichever is closer) season will be used`
    )
}

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