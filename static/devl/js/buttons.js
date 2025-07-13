import * as pdash from "./pdash.js";
import { base, dev, updateCrnt } from "./listen.js";
const mdiv = document.getElementById('main');
export async function clear() {
    const btn = document.getElementById('clear');
    btn.addEventListener('click', async (event) => {
        event.preventDefault();
        mdiv.textContent = '';
    })
}

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
        // const team = document.getElementById('rs_slct').value;

        await pdash.getP(dev, 'random', season, 0);
    })
}

export async function holdPlayerBtn() {
    const btn = document.getElementById('holdP');
    btn.addEventListener('click', async (event) => {
        event.preventDefault();
        // console.log(`pHold: ${document.getElementById('pHold').value}`)
        let player = document.getElementById('pHold').value;
        document.getElementById('pSearch').value = player;


        
        // crnt = 
        // const player = document.getElementById('pSearch');
        // await pdash.getP(dev, 'random', season);
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
        
        // const season = szn.value.trim(); // pass before encoding for error msg
        
        console.log(team)
        await pdash.getP(dev, player, season, team);
        input.value = ''; // clear input box after searching
        
    }) 
}