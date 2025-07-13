import { getP } from "./pdash.js";
import { base } from "./listen.js";

export async function clearSearch() {
    const btn = document.getElementById('clearS');
    btn.addEventListener('click', async (event) => {
        event.preventDefault();
        document.getElementById('pSearch').value = '';
    })
}

export async function randPlayerBtn() {
    const btn = document.getElementById('randP');
    btn.addEventListener('click', async (event) => {
        event.preventDefault();
        const szn = document.getElementById('rs_slct');
        const season = szn.value;
        await getP(base, 'random', season, 0);
    })
}

// read pHold invisible val to add on-screen player's name to search bar
export async function holdPlayerBtn() {
    const btn = document.getElementById('holdP');
    btn.addEventListener('click', async (event) => {
        event.preventDefault();
        let player = document.getElementById('pHold').value;
        document.getElementById('pSearch').value = player;
    })
}

export async function search() {
    const frm = document.getElementById('ui');
    frm.addEventListener('submit', async (event) => {
        event.preventDefault();

        const input = document.getElementById('pSearch');
        const player = input.value.trim();
        const szn = document.getElementById('rs_slct');
        const season = szn.value;
        const team = document.getElementById('nba_teams').value;
        
        console.log(team)
        await getP(base, player, season, '0');
        input.value = ''; // clear input box after searching
        
    }) 
}


