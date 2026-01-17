import { searchPlayer } from "./player.js";
import { clearSearch } from "./inputs.js";
import { MSG, SBL, RED_BOLD, foldedLog } from "../global.js";
import { makeLgTopScorersTbl, makeRgTopScorersTbl, makeTeamRecordsTbl } from "./tbls_onload.js";
import { makeLogoImgs } from "./img.js";
const WINDOWSIZE = 700;
const BIGWINDOW = 2000;
const LARGEROWS = 25;
let exBtnsInitComplete = false;
// counter class for number of rows displayed per table
class rowNum {
    val;
    minR;
    maxR;
    constructor(val, minR = 1, maxR = 100) {
        this.val = val;
        this.minR = minR;
        this.maxR = maxR;
    }
    ;
    get value() {
        return this.val;
    }
    snap(val) {
        if (val <= 1)
            return 1;
        return Math.ceil(val / 5) * 5;
    }
    increase = (by = 5) => {
        const now = this.val;
        this.val = this.snap(Math.min(this.maxR, this.val + by));
        foldedLog(`%cincrease from ${now} to ${this.val}`, SBL);
        return this.val;
    };
    decrease = (by = 5) => {
        const now = this.val;
        this.val = this.snap(Math.max(this.minR, this.val - by));
        foldedLog(`%cdecrease from ${now} to ${this.val}`, SBL);
        return this.val;
    };
    reset = (to) => {
        const now = this.val;
        if (to) {
            this.val = Math.max(this.minR, to);
        }
        else {
            this.val = window.innerWidth <= WINDOWSIZE ? 5 : 10;
        }
        foldedLog(`%creset from ${now} to ${this.val}`, SBL);
        return this.val;
    };
    max = () => {
        foldedLog(`%cmax rows requested: current ${this.val} to ${this.maxR}`, SBL);
        this.val = this.maxR;
        return this.val;
    };
    min = () => {
        foldedLog(`%cmin rows requested: current ${this.val} to ${this.minR}`, SBL);
        this.val = this.minR;
        return this.val;
    };
}
;
// tracks row counters for both top scorer tables
export class rowsState {
    lgRowNum;
    lgRowLarge = LARGEROWS;
    rgRowNum;
    trRowNum;
    startRows;
    constructor(winSize = WINDOWSIZE, bigWinSize = BIGWINDOW) {
        // this.max = max;
        this.startRows = window.innerWidth <= winSize ? 5 : 10;
        if (window.innerWidth >= bigWinSize) {
            this.lgRowNum = new rowNum(this.lgRowLarge);
        }
        else {
            this.lgRowNum = new rowNum(this.startRows);
        }
        this.trRowNum = new rowNum(this.startRows);
        this.rgRowNum = new rowNum(this.startRows);
        this.listenForResize();
    }
    ;
    resetRows(winSize = WINDOWSIZE, bigWinSize = BIGWINDOW) {
        let rows = window.innerWidth <= winSize ? 5 : 10;
        this.rgRowNum.reset(rows);
        this.trRowNum.reset(rows);
        if (window.innerWidth >= bigWinSize) {
            foldedLog(`%cresetting rows for big screen: ${this.lgRowLarge}`, MSG);
            this.lgRowNum.reset(this.lgRowLarge);
        }
        else {
            this.lgRowNum.reset(rows);
        }
        return rows;
    }
    listenForResize(winSize = WINDOWSIZE, bigWinSize = BIGWINDOW) {
        const small_mq = window.matchMedia(`(max-width: ${winSize}px)`);
        const large_mq = window.matchMedia(`(min-width: ${bigWinSize}px)`);
        small_mq.addEventListener('change', async () => { await this.handleMediaQueries(); });
        large_mq.addEventListener('change', async () => { await this.handleMediaQueries(); });
    }
    async handleMediaQueries() {
        this.resetRows();
        await Promise.all([
            makeLogoImgs(),
            makeTeamRecordsTbl(this.trRowNum.value),
            makeLgTopScorersTbl(this.lgRowNum.value),
            makeRgTopScorersTbl(this.rgRowNum.value)
        ]);
    }
}
;
export async function makeExpandTblBtns(rs, tblBtns = [
    { elId: "seemorelessLGbtns", rows: rs.lgRowNum, build: makeLgTopScorersTbl },
    { elId: "seemorelessRGbtns", rows: rs.rgRowNum, build: makeRgTopScorersTbl },
    { elId: "seemorelessTRbtns", rows: rs.trRowNum, build: makeTeamRecordsTbl },
]) {
    if (exBtnsInitComplete)
        return;
    exBtnsInitComplete = true;
    for (let etb of tblBtns) {
        const d = document.getElementById(etb.elId);
        if (!d)
            continue;
        let to_append = [];
        for (const obj of [
            { op: 'all', lbl: 'see all' },
            { op: '+', lbl: 'see more' },
            { op: '-', lbl: 'see less' },
            { op: 'rst', lbl: 'reset' },
            { op: 'min', lbl: 'minimize' }
        ]) {
            let newNum;
            const btn = document.createElement('button');
            btn.textContent = obj.lbl;
            btn.addEventListener('click', async () => {
                switch (obj.op) {
                    case 'all':
                        newNum = etb.rows.max();
                        break;
                    case 'min':
                        newNum = etb.rows.min();
                        break;
                    case '+':
                        newNum = etb.rows.increase();
                        break;
                    case '-':
                        newNum = etb.rows.decrease();
                        break;
                    case 'rst':
                        if (etb.elId === 'seemorelessLGbtns' && window.innerWidth >= BIGWINDOW) {
                            newNum = etb.rows.reset(LARGEROWS);
                        }
                        else {
                            newNum = etb.rows.reset();
                        }
                        break;
                    default:
                        throw new Error(`invalid case: ${obj.op} | ${obj.lbl}`);
                }
                await etb.build(newNum);
            });
            to_append.push(btn);
        }
        for (const b of to_append) {
            d.appendChild(b);
        }
    }
}
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