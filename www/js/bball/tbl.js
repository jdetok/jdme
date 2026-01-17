// data is passed as the type T 
export class Tbl {
    elId;
    title;
    rowCount;
    data;
    columns;
    constructor(elId, title, rowCount, data, columns) {
        this.elId = elId;
        this.title = title;
        this.rowCount = rowCount;
        this.data = data;
        this.columns = columns;
    }
    init() {
        this.build();
    }
    build() {
        let tbl = document.getElementById(this.elId);
        tbl.innerHTML = '';
        tbl.appendChild(this.makeTitle());
        tbl.appendChild(this.makeHdrRow());
        for (let i = 0; i < this.rowCount; i++) {
            tbl.appendChild(this.makeRow(i));
        }
    }
    makeTitle() {
        let capt = document.createElement('caption');
        capt.textContent = this.title;
        return capt;
    }
    makeHdrRow() {
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
    makeRow(idx) {
        const tr = document.createElement("tr");
        this.columns.forEach(col => {
            const td = document.createElement('td');
            const val = col.value(this.data, idx);
            if (col.button) {
                const btn = document.createElement("button");
                btn.type = "button";
                btn.textContent = val;
                btn.onclick = () => col.button.onClick(val, this.data, idx);
                td.appendChild(btn);
            }
            else {
                td.textContent = val;
            }
            tr.appendChild(td);
        });
        return tr;
    }
}
//# sourceMappingURL=tbl.js.map