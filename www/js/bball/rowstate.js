import { MSG, SBL, foldedLog, WINDOWSIZE, BIGWINDOW, LARGEROWS } from "../global.js";
import { rebuildContent } from "./onload.js";
// counter class for number of rows displayed per table
export class rowNum {
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
    rgData;
    constructor(rgData, winSize = WINDOWSIZE, bigWinSize = BIGWINDOW) {
        let rgRows = rgData?.top_scorers?.length ?? 30;
        this.startRows = window.innerWidth <= winSize ? 5 : 10;
        if (window.innerWidth >= bigWinSize) {
            this.lgRowNum = new rowNum(this.lgRowLarge);
        }
        else {
            this.lgRowNum = new rowNum(this.startRows);
        }
        this.trRowNum = new rowNum(this.startRows, 1, 30);
        this.rgRowNum = new rowNum(this.startRows, 1, rgRows);
        this.rgData = rgData;
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
        await rebuildContent(this.trRowNum.value, this.lgRowNum.value, this.rgRowNum.value, this.rgData);
    }
}
;
//# sourceMappingURL=rowstate.js.map