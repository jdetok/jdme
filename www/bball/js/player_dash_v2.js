
// ts indicates 'top scorer' - used when called on page refresh to get recent game
export async function getPlayerStatsV2(base, player, season, team, lg) { // add season & team
    const errmsg = document.getElementById('sErr');
    if (errmsg.style.display === "block") {
        errmsg.style.display = 'none';
    }

    // encode passed args to be ready for query string
    const s = encodeURIComponent(season)
    const p = encodeURIComponent(player).toLowerCase();

    const req = `${base}/v2/players?player=${p}&season=${s}&team=${team}&league=${lg}`;
    // attempt to fetch from /player endpoint with encoded params
    try {
        const r = await fetch(req);
        if (!r.ok) {
            throw new Error(`HTTP Error (${r.status}) attempting to fetch ${player}`);
        }
        
        // get json, set first object in player array as data var
        const js = await r.json()
        if (js) {
            if (js.error_string) {
                errmsg.textContent = js.error_string;
                errmsg.style.display = "block";
                return;
            } else {
                return js;
            }
        }
    } catch(err) {
        errmsg.textContent = `can't find ${player}` 
        console.log(`an error occured attempting to fetch ${player}\n${err}`);
        errmsg.style.display = "block";
        throw new Error(`HTTP Error: ${r.status}`);
    }
}

