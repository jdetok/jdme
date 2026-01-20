import { RED_BOLD, foldedErr, foldedLog } from "../global.js";
import { fetchAndBuildPlayerDash } from "./player_dash.js";
import { clearSearch } from "./inputs.js";
export async function listenForInput() {
    await clearSearchBtn();
    await submitPlayerSearch();
    await randPlayerBtn();
    await holdPlayerBtn();
}
export async function setup_jump_btns() {
    const btns = [{ el: "jumptoresp", jumpTo: "resp_subttl" }, { el: "jumptosearch", jumpTo: "ui" }];
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
export async function setupExclusiveSelectorGroups(st = {
    szn_el: 'select_szns',
    szn: {
        lbox: 'post',
        rbox: 'reg',
    },
    tm_el: 'select_teams',
    tm: {
        lbox: 'nbaTm',
        rbox: 'wnbaTm',
    }
}) {
    // setup internal exclusive listeners
    await setupExclusiveCheckboxes(st.szn.lbox, st.szn.rbox);
    await setupExclusiveCheckboxes(st.tm.lbox, st.tm.rbox);
    setupExclusiveGroups([st.szn.lbox, st.szn.rbox], [st.tm.lbox, st.tm.rbox]);
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
function setupExclusiveGroups(groupA, groupB) {
    const aEls = groupA
        .map(id => document.getElementById(id))
        .filter(Boolean);
    const bEls = groupB
        .map(id => document.getElementById(id))
        .filter(Boolean);
    if (!aEls.length || !bEls.length) {
        throw new Error("could not resolve checkbox groups");
    }
    function clearGroup(group) {
        group.forEach(cb => (cb.checked = false));
    }
    aEls.forEach(cb => cb.addEventListener("change", () => {
        if (cb.checked)
            clearGroup(bEls);
    }));
    bEls.forEach(cb => cb.addEventListener("change", () => {
        if (cb.checked)
            clearGroup(aEls);
    }));
}
async function submitPlayerSearch(elId = 'ui') {
    const frm = document.getElementById(elId);
    if (!frm)
        throw new Error(`couldn't get element at Id ${elId}`);
    frm.addEventListener('submit', async (e) => {
        e.preventDefault();
        await fetchAndBuildPlayerDash();
    });
}
// get a random player from the API and getPlayerStats
async function randPlayerBtn(elId = 'randP') {
    const btn = document.getElementById(elId);
    if (!btn)
        throw new Error(`couldn't get button element at id ${elId}`);
    btn.addEventListener('click', async (e) => {
        e.preventDefault();
        try {
            await fetchAndBuildPlayerDash('random');
        }
        catch (e) {
            foldedErr(`error getting random player: ${e}`);
        }
    });
}
async function clearSearchBtn() {
    const btn = document.getElementById('clearS');
    if (!btn)
        return;
    btn.addEventListener('click', (event) => {
        event.preventDefault();
        clearSearch(true);
    });
}
// BUTTONS SECTION
async function holdPlayerBtn(elId = 'holdP') {
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
//# sourceMappingURL=listeners.js.map