/*INTENT: 
take in and query a URL
used for API calls, making links for images, etc
*/

export async function getAPIResp(url) {
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

// renamed from getStats
export async function FetchURL(url) {
    const loadMsg = document.getElementById('loadmsg');
    loadMsg.textContent = 'Requesting data from API...';
    try { // WAIT FOR API RESPONSE
        const response = await fetch(url);
        if (!response.ok) { 
            throw new Error(`HTTP Error: ${response.status}`)
        } // CONVERT SUCCESSFUL RESPONSE TO JSON & CLEAR LOADMSG
        const data = await response.json();
        loadMsg.textContent = ''; 
        // CONVERT JSON RESPONSE TO HTML TABLE ELEMENTS
        return data
    }
    catch(error) {
        console.log(error);
        loadMsg.textContent = "Failed to load player data";
    };
};

export async function getPlayerId(url, player) {
    const idUrl = url + `/players/id?player=${player}`;
    const response = await fetch(idUrl);
    if (!response.ok) {
        throw new Error(`HTTP Error getting player id: ${response.status}`)
    }
    const jsonResp = await response.json();
    const playerId = jsonResp.playerId;
    return String(playerId);
};

// get player's headshot
export async function getHeadshot(lg, player) {
    let playerId = await getPlayerId(base, player);
    let url = `https://cdn.${lg}.com/headshots/${lg}/latest/1040x760/${playerId}.png`
    return makeImg(url);
}