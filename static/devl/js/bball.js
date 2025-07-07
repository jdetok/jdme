// BASE URLS
const base = "https://jdeko.me/bball";
const nbaHsBase = "https://cdn.nba.com/headshots/nba/latest/1040x760";
const wnbaHsBase = "https://cdn.wnba.com/headshots/wnba/latest/1040x760";

document.addEventListener('DOMContentLoaded', () => {
    loadSeasonOpts();
    loadTeamOpts();
    lgChangeListener();
    gamesRecent();
    topScorer();
});


async function getAPIResp(url) {
    try { // WAIT FOR API RESPONSE
        const response = await fetch(url);
        if (!response.ok) { 
            throw new Error(`HTTP Error: ${response.status}`)
        } // CONVERT SUCCESSFUL RESPONSE TO JSON & CLEAR LOADMSG
        const data = await response.json();
        return data
    }
    catch(error) {
        console.log(error);
    };
};

async function tableJSON(data, element) {
    // DIV TO CREATE STATS ELEMENTS
    const contEl = document.getElementById(element);
    contEl.innerHTML = ''; 

    const keys = Object.keys(data[0]);
    for (const obj of data) { 
        const objTbl = document.createElement('table');
       
        
        // LOOP THROUGH FIELDS > numCapFlds, EACH LOOP APPENDS A ROW TO TABLE
        for (let i = 0; i < keys.length; i++) {
            const row = document.createElement('tr');
            const label = document.createElement('th');
            const val = document.createElement('td');

            // FIELD NAME IN LEFT COLUMN OF TABLE (RIGHT ALIGNED)
            label.textContent = keys[i];
            label.style.textAlign = 'right';

            // VALUE IN RIGHT COLUMN OF TABLE (LEFT ALIGNED)
            val.textContent = obj[keys[i]];
            val.style.textAlign = 'left';
            
            row.appendChild(label); // APPEND LABEL TO ROW
            row.appendChild(val); // APPEND VALUE TO ROW
            objTbl.appendChild(row); // APPEND ROW TO TABLE
        };
        
        contEl.append(objTbl);
        // div.append(objTbl); // APPEND TABLE TO DIV
    };
};


async function gamesRecent() {
    let data = await getAPIResp(base + '/games/recent');
    await tableJSON(data, 'recent-games-tbl');
}

async function topScorer() {
    let data = await getAPIResp(base + '/games/recent/top-scorer');
    await tableJSON(data, 'top-scorer-tbl');
}


function resetListener() {
    const btn = document.getElementById('reset');
    btn.addEventListener('click', async (event) => {
        event.preventDefault();
        document.getElementById('playerForm').reset();
    });
};

async function lgChangeListener() {
    const slct = document.getElementById('lg-slct');
    slct.addEventListener('change', async (event) => {
        event.preventDefault();
        await loadTeamOpts();
    });
};

// LOAD OPTIONS FOR SEASON SELECTOR
async function loadSeasonOpts() {
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
async function loadTeamOpts() {
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

async function teamsToOpts(data) {
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
