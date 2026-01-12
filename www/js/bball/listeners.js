import { handlePlayerSearch, handleRandomPlayer } from "./player.js";
export async function submitPlayerSearch(elId = 'ui') {
    const frm = document.getElementById(elId);
    if (!frm)
        throw new Error(`couldn't get element at Id ${elId}`);
    frm.addEventListener('submit', handlePlayerSearch);
}
// get a random player from the API and getPlayerStats
export async function randPlayerBtn(elId = 'randP') {
    const btn = document.getElementById(elId);
    if (!btn)
        throw new Error(`couldn't get button element at id ${elId}`);
    btn.addEventListener('click', handleRandomPlayer);
}
//# sourceMappingURL=listeners.js.map