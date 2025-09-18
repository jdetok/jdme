// reusable functions for dynamically creating tables

// create caption and append to table
export async function tblCaption(tbl, caption) {
    const capt = document.createElement('caption');
    capt.textContent = caption;
    tbl.appendChild(capt);
} 

// call in top scorer table loop for each player
export async function topPlayerRow(tbl, scorer, data, mode) {
    let game = data.recent_games.find(g => g.player_id === scorer.player_id);
    let r = document.createElement('tr');

    let pName = document.createElement('td');
    let pTeam = document.createElement('td');
    let pts = document.createElement('td');

    let btn = document.createElement('button');
    btn.textContent = scorer.player;
    btn.type = 'button';
    btn.addEventListener('click', async () => {
        await playerLinkSearch(scorer.player, data);
    }); 

    pName.appendChild(btn);
    if (mode == 'recent') {
        pTeam.textContent = `${game ? game.team_name : ""} | \
        ${game ? game.matchup : ""} | ${game ? game.wl : ""}`;
    } else if (mode == 'league') {
        pTeam.textContent = scorer.team;
    }
    pts.textContent = scorer.points;
    r.appendChild(pName);
    r.appendChild(pTeam);
    r.appendChild(pts);

    tbl.appendChild(r);
}

export async function recGamesTbl(tbl) {
    // table headers
    const thead = document.createElement('thead');
    const nameH = document.createElement('td');
    const ptsH = document.createElement('td');
    const teamH = document.createElement('td');

    nameH.textContent = 'name';
    ptsH.textContent = 'points';
    teamH.textContent = 'team | matchup | win-loss';

    thead.appendChild(nameH);
    thead.appendChild(teamH);
    thead.appendChild(ptsH);

    tbl.appendChild(thead);
    
}


// FIRST ROW CONTAINS HEADERS. ALL COLUMNS CONTAIN A HEADER AND DATA
export async function basicTable(data, caption, pElName) {
    // parent element
    const pEl = document.getElementById(pElName);
    const tbl = document.createElement('table');
    const thead = document.createElement('thead');
    const tr = document.createElement('tr');
    const cols = Object.keys(data);
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
export async function rowHdrTable(data, caption, rowHdrLabel, pElName) {
    const pEl = document.getElementById(pElName);
    const tbl = document.createElement('table');
    const thead = document.createElement('thead');
    const keys = Object.keys(data);
    const cols = Object.keys(data[keys[0]])
    
    pEl.textContent = '' // clear parent element
    await tblCaption(tbl, caption); // create & append caption
    
    // append header row to table
    for (let i=0; i<(cols.length + 1); i++) {
        const th = document.createElement('th');
        if (i === 0) { // append first table cell (col header for row headers)
            const rowHdr = document.createElement('th');
            rowHdr.textContent = rowHdrLabel;  
            thead.appendChild(rowHdr);
            tbl.appendChild(thead);  
        } // append colunn headers to thead
        th.textContent = cols[i];
        thead.appendChild(th);
    } // append thead to table
    tbl.appendChild(thead);
    
    // outer loop: set row header for each row
    for (let i=0; i<keys.length; i++) {
        const tr = document.createElement('tr');
        const tdh = document.createElement('td');
        tdh.setAttribute('scope', 'row');
        tdh.textContent = keys[i]
        tr.appendChild(tdh);
    // inner loop: append each data point for each row
        for (let c=0; c<cols.length; c++) {
            const td = document.createElement('td');
            td.textContent = data[keys[i]][cols[c]]; 
            tr.appendChild(td);
            tbl.appendChild(tr);         
        }    
    } // append table to parent element
    pEl.appendChild(tbl);   
}