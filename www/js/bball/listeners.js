import { searchPlayer } from "./player.js";
import { clearSearch } from "./inputs.js";
import { RED_BOLD } from "../global.js";
import { makeLgTopScorersTbl, makeRgTopScorersTbl, makeTeamRecordsTbl } from "./tbls_onload.js";
const WINDOWSIZE = 700;
let exBtnsInitComplete = false;
// counter class for number of rows displayed per table
class rowNum {
    val;
    min;
    constructor(val, min = 2) {
        this.val = val;
        this.min = min;
    }
    ;
    get value() {
        return this.val;
    }
    increase = (by = 5) => {
        this.val += by;
        return this.val;
    };
    decrease = (by = 5) => {
        this.val = Math.max(this.min, this.val - by);
        return this.val;
    };
    reset = (to) => {
        if (to) {
            this.val = Math.max(this.min, to);
        }
        else {
            this.val = window.innerWidth <= WINDOWSIZE ? 5 : 10;
        }
        return this.val;
    };
}
;
// tracks row counters for both top scorer tables
export class rowsState {
    lgRowNum;
    rgRowNum;
    startRows;
    constructor(winSize = WINDOWSIZE) {
        this.startRows = window.innerWidth <= winSize ? 5 : 10;
        this.lgRowNum = new rowNum(this.startRows);
        this.rgRowNum = new rowNum(this.startRows);
        this.listenForResize();
    }
    ;
    resetRows(winSize = WINDOWSIZE) {
        const rows = window.innerWidth <= winSize ? 5 : 10;
        this.lgRowNum.reset(rows);
        this.rgRowNum.reset(rows);
        return rows;
    }
    listenForResize() {
        const mq = window.matchMedia(`(max-width: ${WINDOWSIZE}px)`);
        mq.addEventListener('change', async () => {
            const newRows = this.resetRows();
            await Promise.all([
                makeTeamRecordsTbl(newRows),
                makeLgTopScorersTbl(this.lgRowNum.value),
                makeRgTopScorersTbl(this.rgRowNum.value)
            ]);
        });
    }
}
;
export async function expandedListBtns(rs, btns = [
    { elId: "seemoreplayers", rows: rs.lgRowNum, pm: '+', build: makeLgTopScorersTbl },
    { elId: "seelessplayers", rows: rs.lgRowNum, pm: '-', build: makeLgTopScorersTbl },
    { elId: "resetplayers", rows: rs.lgRowNum, pm: 'rst', build: makeLgTopScorersTbl },
    { elId: "seemoreRGplayers", rows: rs.rgRowNum, pm: '+', build: makeRgTopScorersTbl },
    { elId: "seelessRGplayers", rows: rs.rgRowNum, pm: '-', build: makeRgTopScorersTbl },
    { elId: "resetRGplayers", rows: rs.rgRowNum, pm: 'rst', build: makeRgTopScorersTbl },
]) {
    if (exBtnsInitComplete)
        return;
    exBtnsInitComplete = true;
    for (let btnObj of btns) {
        const btn = document.getElementById(btnObj.elId);
        if (!btn)
            continue;
        btn.addEventListener('click', async () => {
            let newNum;
            switch (btnObj.pm) {
                case '+':
                    newNum = btnObj.rows.increase();
                    break;
                case '-':
                    newNum = btnObj.rows.decrease();
                    break;
                case 'rst':
                    newNum = btnObj.rows.reset();
                    break;
            }
            // const newNum = btnObj.pm === '+' ? btnObj.rows.increase() : btnObj.rows.decrease();
            await btnObj.build(newNum);
        });
    }
}
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
            console.error(`%chold button pressed, empty string in ${holdElId}`, RED_BOLD);
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