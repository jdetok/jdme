// replace some of /www/js/bball/ui.js
// import { base, checkBoxEls, foldedLog, MSG, scrollIntoBySize } from "../global.js";


export function clearSearch(): void {
    const btn = document.getElementById('clearS');
    if (!btn) return;
    btn.addEventListener('click', (event) => {
        event.preventDefault();
        let pSearch = document.getElementById('pSearch');
        if (!pSearch) return;
        pSearch.textContent = '';
        pSearch.focus();
    });
}
