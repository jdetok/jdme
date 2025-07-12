import { getRandP } from "./pdash.js";
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
        // call randomplayer 
        // mdiv.textContent = '';

        await getRandP(dev, 'random');
    })
}

export async function search() {

    // const btn = document.getElementById('search');
    // btn.addEventListener('click', async (event) => {
    const frm = document.getElementById('ui');
    frm.addEventListener('submit', async (event) => {
        event.preventDefault();

        const player = encodeURIComponent(
            document.getElementById('pSearch').value.trim()
            ).toLowerCase();
        await getRandP(dev, player)
    }) 
}

// function searchListener() {
//     const btn = document.getElementById('searchBtn');
//     btn.addEventListener('click', async (event) => {
//         event.preventDefault();
//         await search();
//     });
// };
