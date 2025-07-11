export async function getRandP(base) {
    const r = await fetch(base + '/player?player=random&season=88888');
    if (!r.ok) {
        throw new Error(`HTTP Error: ${r.status}`);
    }
    const js = await r.json();
    const data = js.player[0];
    console.log(data);
    
    const home = document.getElementById('main');
    home.textContent = ''

    await appendImg(data.player_meta.headshot_url, 'imgs', 'pl_img')
    await appendImg(data.player_meta.team_logo_url, 'imgs', 'tm_img')
    await playerTitle(data.player_meta, 'player_title');
    
}

async function appendImg(url, pElName, cElName) {
    const pEl = document.getElementById(pElName);
    const cEl = document.getElementById(cElName);
    // pEl.textContent = '';
    cEl.textContent = '';
    const img = document.createElement('img');
    img.src = url;
    img.alt = "image not found";
    cEl.appendChild(img);
    pEl.append(cEl);
}
async function playerImg(url, elName) {
    const imgEl = document.getElementById(elName);
    imgEl.textContent = ''
    let img = document.createElement('img')
    img.src = url
    img.alt = "player's headshot not found"
    imgEl.append(img)
}


async function playerTitle(meta, elName) {
    let cont = document.getElementById(elName);
    cont.textContent = '';

    let d = document.createElement('div');
    // let h = document.createElement('h3');
    let t = document.createElement('h1');
    
    t.textContent = meta.caption;
    // h.textContent = meta.team_name;
    d.append(t);
    

    
    cont.append(d);
}
