export async function tableJSON(data, element) {
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