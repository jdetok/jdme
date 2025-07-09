const base = "https://jdeko.me/bball";

import * as q from "./q.js";
import * as tbl from "./table.js";

// recent games table on home page
export async function gamesRecent() {
    try {
        let data = await q.getAPIResp(base + '/games/recent');
        console.log(data);
        await tbl.tableJSON(data.recentGames, 'recent-games-tbl');
    } catch (error) {
        console.error('error getting recent games: ', error);
    }
}

export async function topScorer() {
    try {
        let data = await q.getAPIResp(base + '/games/recent/top-scorer');
        console.log(data);
        await tbl.tableJSON(data, 'top-scorer-tbl');
    } catch (error) {
        console.error('error getting top scorer: ', error);
    }
}

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