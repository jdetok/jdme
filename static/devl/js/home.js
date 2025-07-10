const base = "https://jdeko.me/bball";

import * as q from "./q.js";
// import * as tbl from "./table.js";

// recent games table on home page
export async function gamesRecent() {
    try {
        let data = await q.getAPIResp(base + '/games/recent');
        console.log(data.recent_games);
        let cont = document.getElementById('recent-games');
        let d = document.createElement('div')
        let t = document.createElement('h1')
        
        cont.textContent = ''
        t.textContent = `Game(s) played on ${data.recent_games[0].game_date}:`
        d.append(t);
        for (const g of data.recent_games)  {
            let h = document.createElement('h3');
            h.textContent = g.final;
            d.append(h);
        }
        cont.append(d);
        // await tbl.tableJSON(data.recentGames, 'recent-games-tbl');
    } catch (error) {
        console.error('error getting recent games: ', error);
    }
}

export async function topScorer() {
    try {
        let data = await q.getAPIResp(base + '/games/recent/top-scorer');
        
        console.log(data.players[0].player_meta);
        console.log(data.players[0].stats);

        let cont = document.getElementById('top-scorer')
        cont.textContent = '';
        makeTopScorerText(data.players[0]);

        let imgEl = document.getElementById('ts-img');
        let img = document.createElement('img')
        img.src = data.players[0].player_meta.headshot_url
        img.alt = "player's headshot not found"
        imgEl.append(img)
        // await makeShtgTable(data.players[0].shooting_stats)
        makeTable(data.players[0].box_stats, 'ts-box', 'Box-Stats');
        await makeShtgTable(data.players[0].shooting_stats, 'Shooting Stats')
    } catch (error) {
        console.error('error getting top scorer: ', error);
    }
}

async function makeShtgTable(data, caption) {
    let cont = document.getElementById('ts-shtg');
    if (!cont) {
        console.log("null ts-shtg")
    }
    let shtgTbl = document.createElement('table');
    let capt = document.createElement('caption');
    let thead = document.createElement('thead');

    let typ = document.createElement('th');
    typ.setAttribute('scope', 'col');
    typ.textContent = "Shot Type";
    thead.appendChild(typ);

    let m = document.createElement('th');
    m.setAttribute('scope', 'col');
    m.textContent = "Makes";
    thead.appendChild(m);

    let a = document.createElement('th');
    a.setAttribute('scope', 'col');
    a.textContent = "Attempts";
    thead.appendChild(a);

        let p = document.createElement('th');
    p.setAttribute('scope', 'col');
    p.textContent = "Percent";
    thead.appendChild(p);

    let f2 = document.createElement('th');
    f2.setAttribute('scope', 'row');
    f2.textContent = "Field Goals";

    let f3 = document.createElement('th');
    f3.setAttribute('scope', 'row');
    f3.textContent = "3-Pointers";

    let ft = document.createElement('th')
    ft.setAttribute('scope', 'row');
    ft.textContent = "Free Throws";

    let fgRow = document.createElement('tr');
    let fg3Row = document.createElement('tr');
    let ftRow = document.createElement('tr');
    fgRow.appendChild(f2);
    fg3Row.appendChild(f3);
    ftRow.appendChild(ft);

    let keys = Object.keys(data);
    for (let i=0; i<keys.length; i++) {
        // idx of first _ - gets the left e.g. fg, fg3, ft
        let _idx = keys[i].indexOf('_');
        if (keys[i].substring(0, _idx) === "fg3") {
            // console.log(`threes ${keys[i]}: ${data[keys[i]]}`);
            if (keys[i].substring(_idx + 1, _idx + 2) === "m") {
                let tdm = document.createElement('td');
                tdm.textContent = data[keys[i]];
                fg3Row.appendChild(tdm);
            } else if (keys[i].substring(_idx + 1, _idx + 2) === "a") {
                let tda = document.createElement('td');
                tda.textContent = data[keys[i]];
                fg3Row.appendChild(tda);
            } else if (keys[i].substring(_idx + 1, _idx + 2) === "p") {
                let tdp = document.createElement('td');
                tdp.textContent = data[keys[i]];
                fg3Row.appendChild(tdp);
            }
        } else if (keys[i].substring(0, _idx) === "fg") {
            // console.log(`field goals ${keys[i]}: ${data[keys[i]]}`);
            if (keys[i].substring(_idx + 1, _idx + 2) === "m") {
                let tdm = document.createElement('td');
                tdm.textContent = data[keys[i]];
                fgRow.appendChild(tdm);
            } else if (keys[i].substring(_idx + 1, _idx + 2) === "a") {
                let tda = document.createElement('td');
                tda.textContent = data[keys[i]];
                fgRow.appendChild(tda);
            } else if (keys[i].substring(_idx + 1, _idx + 2) === "p") {
                let tdp = document.createElement('td');
                tdp.textContent = data[keys[i]];
                fgRow.appendChild(tdp);
            }
        } else if (keys[i].substring(0, _idx) === "ft"){
            // console.log(`free throws ${keys[i]}: ${data[keys[i]]}`);
            if (keys[i].substring(_idx + 1, _idx + 2) === "m") {
                let tdm = document.createElement('td');
                tdm.textContent = data[keys[i]];
                ftRow.appendChild(tdm);
            } else if (keys[i].substring(_idx + 1, _idx + 2) === "a") {
                let tda = document.createElement('td');
                tda.textContent = data[keys[i]];
                ftRow.appendChild(tda);
            } else if (keys[i].substring(_idx + 1, _idx + 2) === "p") {
                let tdp = document.createElement('td');
                tdp.textContent = data[keys[i]];
                ftRow.appendChild(tdp);
            }
        }
    }
    capt.textContent = caption;
    shtgTbl.appendChild(capt);
    shtgTbl.appendChild(thead);
    shtgTbl.appendChild(fg3Row);
    shtgTbl.appendChild(fgRow);
    shtgTbl.appendChild(ftRow);
    cont.textContent = '';
    cont.appendChild(shtgTbl);
}

async function makeTopScorerText(data) {
    let cont = document.getElementById('top-scorer');
    let d = document.createElement('div');
    let h = document.createElement('h3');
    let t = document.createElement('h1');
    
    t.textContent = `Top Scorer from ${data.game_meta.game_date}:`;
    h.textContent = data.player_meta.caption;
    d.append(t);
    d.append(h);

    // cont.textContent = '';
    cont.append(d);
}

async function makeTable(data, element, caption) {
    let cont = document.getElementById(element);
    const objTbl = document.createElement('table');
    const capt = document.createElement('caption');
    // const headerRow = document.createElement('tr');
    const dataRow = document.createElement('tr');
    let thead = document.createElement('thead');
    let keys = Object.keys(data);
    for (let key of keys) {
        const th = document.createElement('th');
        const td = document.createElement('td');
        
        th.setAttribute('scope', 'col');
        th.textContent = key;

        td.textContent = data[key];

        thead.appendChild(th);
        dataRow.appendChild(td);
    }
    // thead.appendChild
    capt.textContent = caption;
    objTbl.appendChild(capt);
    objTbl.appendChild(thead);
    // objTbl.appendChild(headerRow);
    objTbl.appendChild(dataRow);
    cont.append(objTbl);
}
    
    
    // for (let i = 0; i < keys.length; i++) {
    //     // const cols = document.createElement('tr');
    //     // const dataRow = document.createElement('tr');
    //     const row = document.createElement('tr');
    //     const label = document.createElement('th');
    //     const val = document.createElement('td');
        
    //     label.setAttribute('scope', 'col')
    //     label.textContent = keys[i];
    //     label.style.textAlign = 'right';


    //     // VALUE IN RIGHT COLUMN OF TABLE (LEFT ALIGNED)
    //     val.textContent = data[keys[i]];
    //     val.style.textAlign = 'left';
    //     row.append(label, val);
    //     objTbl.appendChild(row);
    // }
    

export async function lgChangeListener() {
    const slct = document.getElementById('lg-slct');
    slct.addEventListener('change', async (event) => {
        event.preventDefault();
        await loadTeamOpts();
    });
};

// LOAD OPTIONS FOR SEASON SELECTOR
export async function loadSeasonOpts() {
    try {
        const response = await fetch(base + '/seasons');
        if (!response.ok) { 
                throw new Error(`HTTP Error: ${response.status}`);
            } // CONVERT SUCCESSFUL RESPONSE TO JSON
        const data = await response.json();
        if (data[0] == '') {
            console.log('empty json');
        }

        const slct = document.getElementById('szn-slct');
        // each player
        let i;
        for (i=0; i<data.length; i++){
            let opt = document.createElement('option');
            opt.textContent = data[i].Season;
            opt.value = data[i].SeasonId;
            slct.appendChild(opt);
            // console.log(data[i].Season);
        }   
    } catch (error) {
        console.error("failed to load seasons")
    }
};

// LOAD OPTIONS FOR TEAM SELECTOR
export async function loadTeamOpts() {
    try {
        const response = await fetch(base + '/teams');
        if (!response.ok) { 
                throw new Error(`HTTP Error: ${response.status}`);
            } // CONVERT SUCCESSFUL RESPONSE TO JSON
        const data = await response.json();
        if (data[0] == '') {
            console.log('empty json');
        }

        let lg = document.getElementById('lg-slct').value.trim()
        
        const slct = document.getElementById('team-slct');
        slct.innerHTML = ``;

        // default all teams option
        const defaultOpt = document.createElement('option');
        if (lg != "all") {
            defaultOpt.textContent = `All ${lg.toUpperCase()} Teams`;    
        } else {
            defaultOpt.textContent = `All Teams`;    
        }
        // defaultOpt.textContent = `All ${lg.toUpperCase()} Teams`;
        slct.appendChild(defaultOpt);
        // each player
        let i;
        for (i=0; i<data.length; i++){
            // TODO: only if team matches league selector
            if ((data[i].League).toLowerCase() === 
                    document.getElementById('lg-slct').value.trim()) {
                let opt = document.createElement('option');
                opt.textContent = data[i].CityTeam;
                opt.value = data[i].TeamAbbr;
                slct.appendChild(opt);
            } else {
                console.log('team not in league');
            }
            // console.log(data[i].Season);
        }   
    } catch (error) {
        console.error("failed to load seasons")
    }
};

export async function teamsToOpts(data) {
    let i;
        for (i=0; i<data.length; i++){
            // TODO: only if team matches league selector
            if ((data[i].League).toLowerCase() === 
                    document.getElementById('lg-slct').value.trim()) {
                let opt = document.createElement('option');
                opt.textContent = data[i].CityTeam;
                opt.value = data[i].TeamAbbr;
                slct.appendChild(opt);
            } else {
                console.log('team not in league');
            }
            // console.log(data[i].Season);
        }   
}