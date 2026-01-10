// replace some of /www/js/bball/ui.js
// import { base, checkBoxEls, foldedLog, MSG, scrollIntoBySize } from "../global.js";
export function clearSearch() {
    const btn = document.getElementById('clearS');
    if (!btn)
        return;
    btn.addEventListener('click', (event) => {
        event.preventDefault();
        let pSearch = document.getElementById('pSearch');
        if (!pSearch)
            return;
        pSearch.textContent = '';
        pSearch.focus();
    });
}
export async function setPHold(player) {
    const ph = document.getElementById('pHold');
    ph.value = player;
}
export async function lgRadioBtns() {
    const selected = document.querySelector('input[name="leagues"]:checked');
    if (selected) {
        return selected.value;
    }
    else {
        return 'all';
    }
}
//# sourceMappingURL=btns.js.map