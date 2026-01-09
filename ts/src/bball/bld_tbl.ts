// replaces /www/js/bball/dynamic_table.js
// dynamic replacement for lg_ldg_scorers.js, rg_ldg_scorer.js, reamrecs.js

import { foldedLog, MSG } from "../global.js";

type RowFn = (tbl: HTMLTableElement, data: any, idx: number, ) => Promise<void>;
type FldSearchFn = (searchFor: string) => Promise<void>;

type fld<T> = {
    txt: string;
    val: (row: T, idx: number) => string;
    btnFn: FldSearchFn | null;
}

class Tbl<T> {
    el: HTMLTableElement | null = null;
    elName: string;
    title: string;
    hdrs: fld<T>[];
    rowNum: number;
    searchFor: string | null;
    url: string;
    data: T[] = [];
    btnFunc: FldSearchFn | null = null;

    constructor(elName: string, ttl: string, rows: number, hdrs: fld<T>[], searchFor: string | null = null, url: string) {
        this.elName = elName;
        this.title = ttl;
        this.rowNum = rows;
        this.hdrs = hdrs;
        this.searchFor = searchFor;
        this.url = url;
        this.getData();
        this.build();
    }

    async getData(): Promise<void>{
        try {
            const r = await fetch(this.url);
            let js = await r.json();
            this.data = js;
        } catch (ex) {
            throw new Error(`failed getting data from ${this.url}: ${ex}`)
        }
    }

    build() {
        let tbl = document.getElementById(this.elName) as HTMLTableElement | null;
        if (!tbl) throw new Error(`failed to find table element with id: ${this.elName}:`);
        let tblCapt = this.makeTitle();
        if (!tblCapt) throw new Error(`failed creating caption: ${this.title}`);
        tbl.appendChild(tblCapt);

        let hdrRow = this.makeHdrRow();
        if (!hdrRow) throw new Error(`failed to create header row: ${this.hdrs}`);
        tbl.appendChild(hdrRow);

        let rows = this.makeRows();
        if (!rows) throw new Error(`failed to make data rows`);
        for (let row of rows) {
            tbl.appendChild(row);
        }
    }

    makeTitle(): HTMLTableCaptionElement | null {
    let capt = document.createElement('caption') as HTMLTableCaptionElement | null;
        if (!capt) return null
        capt.textContent = this.title;
        return capt;
    }

    makeHdrRow(): HTMLTableSectionElement | null {
        let thead = document.createElement('thead');
        for (let hdr of this.hdrs) {
            let el = document.createElement('td') as HTMLTableCellElement | null;
            if (!el) return null;
            el.textContent = hdr.txt;
            thead.appendChild(el);
        }
        return thead;
    }

    makeRows(): HTMLTableRowElement[] | null {
        let rows: HTMLTableRowElement[] = [];
        for (let i = 0; i < this.rowNum; i++) {
            let row = document.createElement('tr') as HTMLTableRowElement;
            for (let hdr of this.hdrs) {
                let cell = document.createElement('td') as HTMLTableCellElement;
                if (hdr.btnFn) {
                    let btn = document.createElement('button') as HTMLButtonElement;
                    btn.textContent = hdr.vals[i];
                    btn.type = 'button';
                    if (this.searchFor) {
                        btn.addEventListener('click', async () => {
                            await hdr.btnFn!(this.searchFor ?? hdr.vals[i]);
                        });
                    }
                    cell.appendChild(btn);
                } else {
                    cell.textContent = hdr.vals[i];
                }
                row.appendChild(cell);
            }
            rows.push(row);
        }
        return rows;
    }
}


export async function buildTbl(ttl: string) {
    const tbl = new Tbl(
        ttl, [
            new fld( )
        ]
    );
}

// build full table
export async function buildTableWithHdr(data: any, element_id: string, title: string,
    fields: string[], rows_to_display: number, rowfunc: RowFn | null
): Promise<void> {
    foldedLog(`%cbuilding table from object with keys: ${Object.keys(data)}...`, MSG);

    const tbl = document.getElementById(element_id) as HTMLTableElement | null;
    if (!tbl) throw new Error(`no table element found with element ${element_id}`);
    tbl.textContent = "";    

    let capt = getTblCaption(title);
    if (!capt) {
        throw new Error(`failed making table caption from passed title {${title}}`);
    }
    tbl.appendChild(capt);

    let hdr = getTblHdrRow(fields);
    if (!hdr) {
        throw new Error(`failed making header row passed fields: {${console.table(fields)}}`);
    }
    tbl.appendChild(hdr);
    
    // build each row
    for (let i = 0; i < rows_to_display; i++) {
        if (rowfunc) {
            await rowfunc(tbl, data, i);
        } else {
            let row = getTblDataRow();
        }
        
    }
}

function getTblCaption(title: string): HTMLTableCaptionElement | null {
    let capt = document.createElement('caption') as HTMLTableCaptionElement | null;
    if (!capt) return null
    capt.textContent = title;
    return capt;
}

function getTblHdrRow(fields: string[]): HTMLTableSectionElement | null {
    let thead = document.createElement('thead');
    for (let fld of fields) {
        let el = document.createElement('td') as HTMLTableCellElement | null;
        if (!el) return null;
        el.textContent = fld;
        thead.appendChild(el);
    }
    return thead;
}
type CellData = {hdr: string, data: string};

function getTblDataRow(idx: number, cells: CellData[]): HTMLTableRowElement | null {
    let row = document.createElement('tr') as HTMLTableRowElement | null;
    if (!row) return null;
    for (let cell of cells) {
        let el = document.createElement('td') as HTMLTableCellElement | null 
        if (!el) return null;
        el.textContent = cell.data;
        row.appendChild(el);
    }
    return row;
}

