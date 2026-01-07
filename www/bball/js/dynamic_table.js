// dynamic replacement for lg_ldg_scorers.js, rg_ldg_scorer.js, reamrecs.js

// build full table
export async function table5f(data, element_id, title, fields_hdrs, rows_to_display, rowfunc) {
    const tbl = document.getElementById(element_id);
    tbl.textContent = "";
    
    // first make caption
    (function(tbl, title) {
        let capt = document.createElement('caption');
        capt.textContent = title;
        tbl.appendChild(capt);
    })(tbl, title);

    // build header
    (function(tbl, fields_hdrs) {
        let thead = document.createElement('thead');
        for (let hdr of fields_hdrs) {
            let el = document.createElement('td');
            el.textContent = hdr;
            thead.appendChild(el);
        }
        tbl.appendChild(thead);
    })(tbl, fields_hdrs);
    

    // build each row
    for (let i = 0; i < rows_to_display; i++) {       
        await rowfunc(tbl, data, i);
    }
}