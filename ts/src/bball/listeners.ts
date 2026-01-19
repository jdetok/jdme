import { RED_BOLD, errMsg, foldedErr, foldedLog } from "../global.js";
import { fetchAndBuildPlayerDash } from "./player_dash.js";
import { clearCheckBoxes, clearSearch } from "./inputs.js";

export async function listenForInput(): Promise<void> {
    await clearSearchBtn();
    await submitPlayerSearch();
    await randPlayerBtn();
    await holdPlayerBtn(); 
}

export async function setup_jump_btns() {
    const btns = [{el: "jumptoresp", jumpTo: "player_title"}, {el: "jumptosearch", jumpTo: "ui"}]
    for (const btn of btns) {
        const btnEl = document.getElementById(btn.el);
        if (btnEl) {
            btnEl.addEventListener('click', async() => {
                const jmpEl = document.getElementById(btn.jumpTo);
                if (jmpEl) {
                    jmpEl.scrollIntoView({behavior: "smooth", block: "start"});
                }
            })
        }    
    }
}

export type SznTmSelectors = {
    szn_el: string,
    szn: {
        lbox: string, 
        rbox: string,
    },
    tm_el: string,
    tm: {
        lbox: string,
        rbox: string,
    }
}

export async function setupExclusiveSelectorGroups(st: SznTmSelectors = {
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
}): Promise<void> {
    // setup internal exclusive listeners
    await setupExclusiveCheckboxes(st.szn.lbox, st.szn.rbox);
    await setupExclusiveCheckboxes(st.tm.lbox, st.tm.rbox);
    setupExclusiveGroups(
        [st.szn.lbox, st.szn.rbox],
        [st.tm.lbox, st.tm.rbox]
    );
}

// make post + reg checkboxes exclusive (but allow neither checked)
export async function setupExclusiveCheckboxes(leftbox: string, rightbox: string) {
    let lbox = document.getElementById(leftbox) as HTMLInputElement;
    let rbox = document.getElementById(rightbox) as HTMLInputElement;
    if (!lbox || !rbox) throw new Error(`couldn't get ${lbox} or ${rbox}`);
    function handleCheck(e: Event) {
        const target = e.target as HTMLInputElement;
        if (!target) throw new Error(`can't find ${e.type}`);
        if (target.checked) { 
            if (e.target === lbox) rbox.checked = false;
            if (e.target === rbox) lbox.checked = false;
        }
    }
    lbox.addEventListener("change", handleCheck);
    rbox.addEventListener("change", handleCheck);
}

function setupExclusiveGroups(
    groupA: string[],
    groupB: string[]
) {
    const aEls = groupA
        .map(id => document.getElementById(id) as HTMLInputElement | null)
        .filter(Boolean) as HTMLInputElement[];

    const bEls = groupB
        .map(id => document.getElementById(id) as HTMLInputElement | null)
        .filter(Boolean) as HTMLInputElement[];

    if (!aEls.length || !bEls.length) {
        throw new Error("could not resolve checkbox groups");
    }

    function clearGroup(group: HTMLInputElement[]) {
        group.forEach(cb => (cb.checked = false));
    }

    aEls.forEach(cb =>
        cb.addEventListener("change", () => {
            if (cb.checked) clearGroup(bEls);
        })
    );

    bEls.forEach(cb =>
        cb.addEventListener("change", () => {
            if (cb.checked) clearGroup(aEls);
        })
    );
}

async function submitPlayerSearch(elId = 'ui') {
    const frm = document.getElementById(elId) as HTMLFormElement;
    if (!frm) throw new Error(`couldn't get element at Id ${elId}`);
    frm.addEventListener('submit', async (e: SubmitEvent) => {
        e.preventDefault();
        await fetchAndBuildPlayerDash();
    });
}

// get a random player from the API and getPlayerStats
async function randPlayerBtn(elId = 'randP') {
    const btn = document.getElementById(elId) as HTMLInputElement;
    if (!btn) throw new Error(`couldn't get button element at id ${elId}`);
    btn.addEventListener('click', async (e: Event) => {
        e.preventDefault();
        try {
            await fetchAndBuildPlayerDash('random');
        } catch (e) {
            foldedErr(`error getting random player: ${e}`);
        }
    });
}

async function clearSearchBtn(): Promise<void> {
    const btn = document.getElementById('clearS');
    if (!btn) return;
    btn.addEventListener('click', (event) => {
        event.preventDefault();
        clearSearch(true);
    });
}

// BUTTONS SECTION
async function holdPlayerBtn(elId = 'holdP') {
    const btn = document.getElementById(elId);
    if (!btn) throw new Error(`couldn't get button element at ${elId}`);
    btn.addEventListener('click', async (event) => {
        event.preventDefault();
        // get player name held in pHold value, fill player search bar with it
        const holdElId = 'pHold';
        const hold = document.getElementById(holdElId) as HTMLInputElement;
        if (!hold) throw new Error(`couldn't get input element at ${holdElId}`);
        let player = hold.value;
        if (player === '') {
            foldedLog(`%chold button pressed, empty string in ${holdElId}`, RED_BOLD);
            return;
        }

        const searchElId = 'pSearch';
        let search = document.getElementById(searchElId) as HTMLInputElement;
        search.value = player;
    })
}
