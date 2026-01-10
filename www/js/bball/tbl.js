// replaces /www/js/bball/dynamic_table.js
// dynamic replacement for lg_ldg_scorers.js, rg_ldg_scorer.js, reamrecs.js
class Tbl {
    constructor(elName, ttl, rows, hdrs, searchFor = null, url) {
        this.el = null;
        this.data = [];
        this.btnFunc = null;
        this.elName = elName;
        this.title = ttl;
        this.rowNum = rows;
        this.hdrs = hdrs;
        this.searchFor = searchFor;
        this.url = url;
        this.getData();
        this.build();
    }
    async getData() {
        try {
            const r = await fetch(this.url);
            let js = await r.json();
            this.data = js;
        }
        catch (ex) {
            throw new Error(`failed getting data from ${this.url}: ${ex}`);
        }
    }
    build() {
        let tbl = document.getElementById(this.elName);
        if (!tbl)
            throw new Error(`failed to find table element with id: ${this.elName}:`);
        let tblCapt = this.makeTitle();
        if (!tblCapt)
            throw new Error(`failed creating caption: ${this.title}`);
        tbl.appendChild(tblCapt);
        let hdrRow = this.makeHdrRow();
        if (!hdrRow)
            throw new Error(`failed to create header row: ${this.hdrs}`);
        tbl.appendChild(hdrRow);
        let rows = this.makeRows();
        if (!rows)
            throw new Error(`failed to make data rows`);
        for (let row of rows) {
            tbl.appendChild(row);
        }
    }
    makeTitle() {
        let capt = document.createElement('caption');
        if (!capt)
            return null;
        capt.textContent = this.title;
        return capt;
    }
    makeHdrRow() {
        let thead = document.createElement('thead');
        for (let hdr of this.hdrs) {
            let el = document.createElement('td');
            if (!el)
                return null;
            el.textContent = hdr.txt;
            thead.appendChild(el);
        }
        return thead;
    }
    makeRows() {
        let rows = [];
        for (let i = 0; i < this.rowNum; i++) {
            let row = document.createElement('tr');
            for (let hdr of this.hdrs) {
                let cell = document.createElement('td');
                const val = hdr.val(this.data[i], i);
                if (hdr.btnFn) {
                    let btn = document.createElement('button');
                    btn.type = 'button';
                    btn.textContent = val;
                    if (this.searchFor) {
                        btn.addEventListener('click', async () => {
                            await hdr.btnFn(this.searchFor ?? val);
                        });
                    }
                    cell.appendChild(btn);
                }
                else {
                    cell.textContent = val;
                }
                row.appendChild(cell);
            }
            rows.push(row);
        }
        return rows;
    }
}
export {};
//# sourceMappingURL=tbl.js.map