// replaces /www/js/bball/dynamic_table.js
// dynamic replacement for lg_ldg_scorers.js, rg_ldg_scorer.js, reamrecs.js

import { foldedLog, MSG } from "../global.js";

type CellValue<T> = (data: T | any, idx: number) => string;

type Column<T> = {
    header: string;
    value: CellValue<T>;
    button?: {
        onClick: (value: string, data: T, idx: number) => Promise<void>;
    };
};

export class Tbl<T> {
    private data!: T;

    constructor(
        private elId: string,
        private title: string,
        private rowCount: number,
        private url: string,
        private columns: Column<T>[],
    ) {}

    async init(): Promise<void> {
        const r = await fetch(this.url);
        this.data = await r.json();
        this.build();
    }

    private build(): void {
        let tbl = document.getElementById(this.elId) as HTMLTableElement;
        tbl.innerHTML = '';

        tbl.appendChild(this.makeTitle());

        tbl.appendChild(this.makeHdrRow());

        for (let i = 0; i < this.rowCount; i++) {
            tbl.appendChild(this.makeRow(i));
        }
    }

    makeTitle(): HTMLTableCaptionElement {
    let capt = document.createElement('caption') as HTMLTableCaptionElement;
        capt.textContent = this.title;
        return capt;
    }

    makeHdrRow(): HTMLTableSectionElement {
        const thead = document.createElement('thead');
        const tr = document.createElement('tr');
        for (const col of this.columns) {
            const td = document.createElement('td');
            td.textContent = col.header;
            tr.appendChild(td);
        }
        thead.appendChild(tr);
        return thead;
    }

    makeRow(idx: number): HTMLTableRowElement {
        const tr = document.createElement("tr");

        this.columns.forEach(col => {
            const td = document.createElement('td');
            const val = col.value(this.data, idx);
            if (col.button) {
                const btn = document.createElement("button");
                btn.type = "button";
                btn.textContent = val;
                btn.onclick = () => col.button!.onClick(val, this.data, idx);
                td.appendChild(btn);
            } else {
                td.textContent = val;
            }

            tr.appendChild(td);
        })
        return tr;
    }
}