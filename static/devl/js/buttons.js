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

        await getRandP(dev);
    })
}