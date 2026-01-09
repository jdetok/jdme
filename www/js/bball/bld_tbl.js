// replaces /www/js/bball/dynamic_table.js
// dynamic replacement for lg_ldg_scorers.js, rg_ldg_scorer.js, reamrecs.js
import { foldedLog, MSG } from "../global.js";
// build full table
export async function table5f(data, element_id, title, fields, rows_to_display, rowfunc) {
    foldedLog(`%cbuilding table from object with keys: ${Object.keys(data)}...`, MSG);
    const tbl = document.getElementById(element_id);
    if (!tbl)
        throw new Error(`no table element found with element ${element_id}`);
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
        await rowfunc(i, tbl, data);
    }
}
function getTblCaption(title) {
    let capt = document.createElement('caption');
    if (!capt)
        return null;
    capt.textContent = title;
    return capt;
}
function getTblHdrRow(fields) {
    let thead = document.createElement('thead');
    for (let fld of fields) {
        let el = document.createElement('td');
        if (!el)
            return null;
        el.textContent = fld;
        thead.appendChild(el);
    }
    return thead;
}
//# sourceMappingURL=bld_tbl.js.map