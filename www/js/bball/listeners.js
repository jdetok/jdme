import { searchPlayer } from "./player.js";
import { clearSearch } from "./inputs.js";
import { RED_BOLD, foldedLog } from "../global.js";
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
// make post + reg checkboxes exclusive (but allow neither checked)
export async function setupExclusiveCheckboxes(leftbox, rightbox) {
    let lbox = document.getElementById(leftbox);
    let rbox = document.getElementById(rightbox);
    if (!lbox || !rbox)
        throw new Error(`couldn't get ${lbox} or ${rbox}`);
    function handleCheck(e) {
        const target = e.target;
        if (!target)
            throw new Error(`can't find ${e.type}`);
        if (target.checked) {
            if (e.target === lbox)
                rbox.checked = false;
            if (e.target === rbox)
                lbox.checked = false;
        }
    }
    lbox.addEventListener("change", handleCheck);
    rbox.addEventListener("change", handleCheck);
}
export async function clearSearchBtn() {
    const btn = document.getElementById('clearS');
    if (!btn)
        return;
    btn.addEventListener('click', (event) => {
        event.preventDefault();
        clearSearch(true);
    });
}
// BUTTONS SECTION
export async function holdPlayerBtn(elId = 'holdP') {
    const btn = document.getElementById(elId);
    if (!btn)
        throw new Error(`couldn't get button element at ${elId}`);
    btn.addEventListener('click', async (event) => {
        event.preventDefault();
        // get player name held in pHold value, fill player search bar with it
        const holdElId = 'pHold';
        const hold = document.getElementById(holdElId);
        if (!hold)
            throw new Error(`couldn't get input element at ${holdElId}`);
        let player = hold.value;
        if (player === '') {
            foldedLog(`%chold button pressed, empty string in ${holdElId}`, RED_BOLD);
            return;
        }
        const searchElId = 'pSearch';
        let search = document.getElementById(searchElId);
        search.value = player;
    });
}
export async function setup_jump_btns() {
    const btns = [{ el: "jumptoresp", jumpTo: "player_title" }, { el: "jumptosearch", jumpTo: "ui" }];
    for (const btn of btns) {
        const btnEl = document.getElementById(btn.el);
        if (btnEl) {
            btnEl.addEventListener('click', async () => {
                const jmpEl = document.getElementById(btn.jumpTo);
                if (jmpEl) {
                    jmpEl.scrollIntoView({ behavior: "smooth", block: "start" });
                }
            });
        }
    }
}
//# sourceMappingURL=listeners.js.map