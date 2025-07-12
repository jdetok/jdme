import * as pdash from "./pdash.js";
import { base, dev } from "./listen.js";
const mdiv = document.getElementById('main');
export async function clear() {
    const btn = document.getElementById('clear');
    btn.addEventListener('click', async (event) => {
        event.preventDefault();
        mdiv.textContent = '';
    })
}

export async function randPlayerBtn() {
    const btn = document.getElementById('randP');
    btn.addEventListener('click', async (event) => {
        event.preventDefault();
        console.log('random')
        await pdash.getP(dev, 'random');
    })
}

export async function search() {
    const frm = document.getElementById('ui');
    frm.addEventListener('submit', async (event) => {
        event.preventDefault();
        const input = document.getElementById('pSearch');
        const player = input.value.trim(); // pass before encoding for error msg
        await pdash.getP(dev, player);
        input.value = ''; // clear input box after searching
    }) 
}