// reusable functions for dynamically creating tables

// create caption and append to table
export async function tblCaption(tbl, caption) {
    const capt = document.createElement('caption');
    capt.innerHTML = caption;
    tbl.appendChild(capt);
} 

// FIRST ROW CONTAINS HEADERS. ALL COLUMNS CONTAIN A HEADER AND DATA
export async function basicTable(data, caption, pElName) {
    // parent element
    const pEl = document.getElementById(pElName);
    const tbl = document.createElement('table');
    const thead = document.createElement('thead');
    const tr = document.createElement('tr');

    // get keys of json object and set as cols
    const cols = Object.keys(data);
    // console.log(`cols: ${cols}`);
    pEl.textContent = '' // clear parent element
    await tblCaption(tbl, caption); // create & append caption

    // append the header and data for each column
    for (let i=0; i<cols.length; i++) {
        const th = document.createElement('th');
        const td = document.createElement('td');
        th.textContent = cols[i];
        td.textContent = data[cols[i]]
        thead.appendChild(th);
        tr.appendChild(td);
    }
    tbl.appendChild(thead);
    tbl.appendChild(tr);
    pEl.appendChild(tbl);
}

// FIRST CELL OF EACH ROW IS A HEADER
// primarily used with shooting stats tables
export async function rowHdrTable(data, caption, rowHdrLabel, pElName) {
    const pEl = document.getElementById(pElName);
    const tbl = document.createElement('table');
    const thead = document.createElement('thead');
    
    // use keys from json object as row headers - fg, fg3, ft
    const rHdrs = Object.keys(data);

    // use keys from the first data[keys] as column headers - made, attempt, pct 
    const cHdrs = Object.keys(data[rHdrs[0]]);
    
    // console.log(`ROW HEADERS: ${rHdrs}`);
    // console.log(`COLUMN HEADERS: ${cHdrs}`);

    // use cols, create & append caption
    pEl.textContent = ''; 
    await tblCaption(tbl, caption);
    
    // append header rows to table
    for (let i=0; i<(cHdrs.length + 1); i++) {
        const th = document.createElement('th');
        if (i === 0) { // append first table cell (col header for row headers)
            const rowHdr = document.createElement('th');
            rowHdr.textContent = rowHdrLabel;  
            thead.appendChild(rowHdr);
            tbl.appendChild(thead);  
        } // append colunn headers to thead
        th.textContent = cHdrs[i];
        thead.appendChild(th);
    } // append thead to table
    tbl.appendChild(thead);
    
    // outer loop: set row header for each row
    for (let i=0; i<rHdrs.length; i++) {
        const tr = document.createElement('tr');
        const tdh = document.createElement('td');
        tdh.setAttribute('scope', 'row');
        tdh.textContent = rHdrs[i]
        tr.appendChild(tdh);
    // inner loop: append each data point for each row
        for (let c=0; c<cHdrs.length; c++) {
            const td = document.createElement('td');
            td.textContent = data[rHdrs[i]][cHdrs[c]]; 
            tr.appendChild(td);
            tbl.appendChild(tr);         
        }    
    } // append table to parent element
    pEl.appendChild(tbl);   
}