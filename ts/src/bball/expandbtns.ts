export type expandTblBtns = {
    elId: string, 
    rows: rowNum, 
    build: (numRows: number) => Promise<void>
}

export async function makeExpandTblBtns(rs: rowsState, tblBtns: expandTblBtns[] = [
    {elId: "seemorelessLGbtns", rows: rs.lgRowNum, build: makeLgTopScorersTbl},
    {elId: "seemorelessRGbtns", rows: rs.rgRowNum, build: makeRgTopScorersTbl},
    {elId: "seemorelessTRbtns", rows: rs.trRowNum, build: makeTeamRecordsTbl},
]) {
    if (exBtnsInitComplete) return;
    exBtnsInitComplete = true;
    for (let etb of tblBtns) {
        const d = document.getElementById(etb.elId);
        if (!d) continue;
        
        let to_append: HTMLButtonElement[] = [];
        for (const obj of [
            { op: 'all', lbl: 'see all' },
            { op: '+', lbl: 'see more' },
            { op: '-', lbl: 'see less' },
            { op: 'rst', lbl: 'reset' },
            { op: 'min', lbl: 'minimize' }
        ]) {
            let newNum: number;
            
            const btn = document.createElement('button');
            btn.textContent = obj.lbl;
            btn.addEventListener('click', async () => {
                switch (obj.op) {
                    case 'all':
                        newNum = etb.rows.max();
                        break;
                    case 'min':
                        newNum = etb.rows.min();
                        break;
                    case '+':
                        newNum = etb.rows.increase();
                        break;
                    case '-':
                        newNum = etb.rows.decrease();
                        break;
                    case 'rst':
                        if (etb.elId === 'seemorelessLGbtns' && window.innerWidth >= BIGWINDOW) {
                            newNum = etb.rows.reset(LARGEROWS);
                        } else {
                            newNum = etb.rows.reset();
                        }
                        break;
                    default:
                        throw new Error(`invalid case: ${obj.op} | ${obj.lbl}`)
                }
                await etb.build(newNum);
            });
            to_append.push(btn);
        }
        for (const b of to_append) {
            d.appendChild(b);
        }
        
    }
};