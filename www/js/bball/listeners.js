import { searchPlayer } from "./player.js";
export async function submitPlayerSearch(elId = 'ui') {
    const frm = document.getElementById(elId);
    if (!frm)
        throw new Error(`couldn't get element at Id ${elId}`);
    frm.addEventListener('submit', async (e) => {
        e.preventDefault();
        await searchPlayer();
    });
}
// get a random player from the API and getPlayerStats
export async function randPlayerBtn(elId = 'randP') {
    const btn = document.getElementById(elId);
    if (!btn)
        throw new Error(`couldn't get button element at id ${elId}`);
    btn.addEventListener('click', async (e) => {
        e.preventDefault();
        await searchPlayer('random');
    });
}
//# sourceMappingURL=listeners.js.map