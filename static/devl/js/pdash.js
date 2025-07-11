export async function getRandP(base) {
    const r = await fetch(base + '/player?player=random&season=88888');
    if (!r.ok) {
        throw new Error(`HTTP Error: ${r.status}`);
    }
    const js = await r.json();
    const data = js.player[0];
    // console.log(data);
    
    const home = document.getElementById('main');
    home.textContent = ''

    await appendImg(data.player_meta.headshot_url, 'imgs', 'pl_img')
    await appendImg(data.player_meta.team_logo_url, 'imgs', 'tm_img')
    await playerTitle(data.player_meta, 'player_title');
    await shtgTable(data.totals.shooting, 'shooting');
}

async function shtgTable(shtg, pElName) {
    const pEl = document.getElementById(pElName);
    const shtgTbl = document.createElement('table');
    const thead = document.createElement('thead');
    // let capt = document.createElement('caption');
    const sType = Object.keys(shtg);
    const keys = Object.keys(shtg);
    const cols = Object.keys(shtg[keys[0]])
    console.log(`COLUMNS: ${cols}`)
    // for (let r of Object.keys(shtg[keys[0]])) {
    for (let i=0; i<(cols.length + 1); i++) {
        console.log(`OUTER LOOP LENGTH: ${cols.length + 1}`)
        // console.log(r);
        let th = document.createElement('th');
        if (i === 0) {
            let th2 = document.createElement('th');
            th2.textContent = 'shooting';  
            thead.appendChild(th2);
            shtgTbl.appendChild(thead);  
        } 
        th.textContent = cols[i];
        thead.appendChild(th);
    }
    shtgTbl.appendChild(thead);
    
    for (let i=0; i<sType.length; i++) {
        console.log(`OUTER LOOP 2 LENGTH: ${i}/${sType.length}`)
        let tr = document.createElement('tr');
        for (let c=0; c<cols.length + 1; c++) {
            if (c === 0) {
                let tdh = document.createElement('td');
                tdh.setAttribute('scope', 'row');
                
                tdh.textContent = sType[c];
                tr.appendChild(tdh);
                

                console.log(`INNER LOOP i=0: ${i}: ${sType[i]} | ${c}:${cols[c]}`)
                console.log(`ROW HEADER: ${sType[i]}`)
            } else {
                console.log(`ELSE ${i}/${cols.length + 1}`)
                let td = document.createElement('td');
                td.textContent = shtg[sType[i]][cols[c-1]];
                tr.appendChild(td);
                shtgTbl.appendChild(tr);         
                console.log(`INNERMOST: ${i}: ${sType[i]} | ${c}:${cols[c-1]}`)
                console.log(shtg[sType[i]][cols[c-1]])
            }
        }    
    }
    pEl.appendChild(shtgTbl);
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

