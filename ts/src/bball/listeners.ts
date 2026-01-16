
import { searchPlayer } from "./player.js";
import { clearSearch } from "./inputs.js";
import { RED_BOLD } from "../global.js";
import { makeLgTopScorersTbl, makeRgTopScorersTbl, makeTeamRecordsTbl } from "./tbls_onload.js";

const WINDOWSIZE = 700;
const BIGWINDOW = 2000;
const LARGEROWS = 25;
let exBtnsInitComplete = false;

// counter class for number of rows displayed per table
class rowNum {
    constructor(private val: number, private min = 2) {};
    get value(): number {
        return this.val;
    }
    increase = (by = 5): number => {
        this.val += by;
        return this.val;
    }

    decrease = (by = 5): number => {
        this.val = Math.max(this.min, this.val - by);
        return this.val;
    }

    reset = (to?: number): number => {
        if (to) {
            this.val = Math.max(this.min, to);
        } else {
            this.val = window.innerWidth <= WINDOWSIZE ? 5 : 10;
        }
        return this.val;
    }
};

// tracks row counters for both top scorer tables
export class rowsState {
    lgRowNum: rowNum;
    lgRowLarge: number = LARGEROWS;
    rgRowNum: rowNum;
    trRowNum: rowNum;
    startRows: number;
    constructor(winSize: number = WINDOWSIZE, bigWinSize: number = BIGWINDOW) {
        this.startRows = window.innerWidth <= winSize ? 5 : 10;
        if (window.innerWidth >= bigWinSize) { 
            this.lgRowNum = new rowNum(this.lgRowLarge);
        } else {
            this.lgRowNum = new rowNum(this.startRows);
        }
        this.trRowNum = new rowNum(this.startRows);
        this.rgRowNum = new rowNum(this.startRows);
        this.listenForResize();
    };
    resetRows(winSize: number = WINDOWSIZE, bigWinSize: number = BIGWINDOW): number {
        let rows: number = window.innerWidth <= winSize ? 5 : 10;
        this.rgRowNum.reset(rows);
        this.trRowNum.reset(rows);
        if (window.innerWidth >= bigWinSize) {
            this.lgRowNum.reset(this.lgRowLarge);
        } else {
            this.lgRowNum.reset(rows);
        }
        return rows;
    }
    listenForResize(winSize: number = WINDOWSIZE, bigWinSize: number = BIGWINDOW) { // change row nums and rebuild tables when window size changes
        const small_mq = window.matchMedia(`(max-width: ${winSize}px)`);
        const large_mq = window.matchMedia(`(min-width: ${bigWinSize}px)`);
        small_mq.addEventListener('change', async () => { await this.handleMediaQueries(); });
        large_mq.addEventListener('change', async () => { await this.handleMediaQueries(); });
    }
    async handleMediaQueries() {
        const newRows = this.resetRows();
        await Promise.all([
            makeTeamRecordsTbl(this.trRowNum.value),
            makeLgTopScorersTbl(this.lgRowNum.value),
            makeRgTopScorersTbl(this.rgRowNum.value)
        ])
    }
};

export type expandTblBtns = {
    elId: string, 
    rows: rowNum, 
    build: (numRows: number) => Promise<void>
}

export async function makeExpandTblBtns(rs: rowsState, tblBtns: expandTblBtns[] = [
    {elId: "seemorelessLGbtns", rows: rs.lgRowNum, build: makeLgTopScorersTbl},
    {elId: "seemorelessRGbtns", rows: rs.rgRowNum, build: makeRgTopScorersTbl},
    {elId: "seemorelessTRbtns", rows: rs.trRowNum, build: makeTeamRecordsTbl},
]) {
    if (exBtnsInitComplete) return;
    exBtnsInitComplete = true;
    for (let etb of tblBtns) {
        const d = document.getElementById(etb.elId);
        if (!d) continue;
        
        let to_append: HTMLButtonElement[] = [];
        for (const obj of [{ op: '+', lbl: 'see more' }, {op:  '-', lbl:  'see less' }, { op: 'rst', lbl: 'reset'}]) {
            let newNum: number;
            
            const btn = document.createElement('button');
            btn.textContent = obj.lbl;
            btn.addEventListener('click', async () => {
                switch (obj.op) {
                    case '+':
                        newNum = etb.rows.increase();
                        break;
                    case '-':
                        newNum = etb.rows.decrease();
                        break;
                    case 'rst':
                        if (etb.elId === 'seemorelessLGbtns' && window.innerWidth >= BIGWINDOW) {
                            newNum = etb.rows.reset(LARGEROWS);
                        } else {
                            newNum = etb.rows.reset();
                        }
                        break;
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

export type expandableTbl = {
    elId: string, 
    rows: rowNum, 
    pm: '+' | '-' | 'rst',
    build: (numRows: number) => Promise<void>
}
export async function expandedListBtns(rs: rowsState, btns: expandableTbl[] = [
    {elId: "seemoreplayers", rows: rs.lgRowNum, pm: '+', build: makeLgTopScorersTbl}, 
    {elId: "seelessplayers", rows: rs.lgRowNum, pm: '-', build: makeLgTopScorersTbl},
    {elId: "resetplayers", rows: rs.lgRowNum, pm: 'rst', build: makeLgTopScorersTbl},
    {elId: "seemoreRGplayers", rows: rs.rgRowNum, pm: '+', build: makeRgTopScorersTbl}, 
    {elId: "seelessRGplayers", rows: rs.rgRowNum, pm: '-', build: makeRgTopScorersTbl},
    {elId: "resetRGplayers", rows: rs.rgRowNum, pm: 'rst', build: makeRgTopScorersTbl},
]) {
    if (exBtnsInitComplete) return;
    exBtnsInitComplete = true;
    for (let btnObj of btns) {
        const btn = document.getElementById(btnObj.elId);
        if (!btn) continue;
        btn.addEventListener('click', async() => {
            let newNum: number;
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
    const frm = document.getElementById(elId) as HTMLFormElement;
    if (!frm) throw new Error(`couldn't get element at Id ${elId}`);
    frm.addEventListener('submit', async (e: SubmitEvent) => {
        e.preventDefault();
        await searchPlayer();
    });
}

// get a random player from the API and getPlayerStats
export async function randPlayerBtn(elId = 'randP') {
    const btn = document.getElementById(elId) as HTMLInputElement;
    if (!btn) throw new Error(`couldn't get button element at id ${elId}`);
    btn.addEventListener('click', async (e: Event) => {
        e.preventDefault();
        await searchPlayer('random');
    });
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

export async function clearSearchBtn(): Promise<void> {
    const btn = document.getElementById('clearS');
    if (!btn) return;
    btn.addEventListener('click', (event) => {
        event.preventDefault();
        clearSearch(true);
    });
}


// BUTTONS SECTION
export async function holdPlayerBtn(elId = 'holdP') {
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
            console.error(`%chold button pressed, empty string in ${holdElId}`, RED_BOLD);
            return;
        }

        const searchElId = 'pSearch';
        let search = document.getElementById(searchElId) as HTMLInputElement;
        search.value = player;
    })
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