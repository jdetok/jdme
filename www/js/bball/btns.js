// replace some of /www/js/bball/ui.js
// import { base, checkBoxEls, foldedLog, MSG, scrollIntoBySize } from "../global.js";
export async function clearSearchBtn() {
    const btn = document.getElementById('clearS');
    if (!btn)
        return;
    btn.addEventListener('click', (event) => {
        event.preventDefault();
        clearSearch(true);
    });
}
export function clearSearch(focus = false) {
    const elId = 'pSearch';
    const pSearch = document.getElementById(elId);
    if (!pSearch)
        throw new Error(`couldn't get search bar element at ${elId}`);
    pSearch.value = '';
    if (focus)
        pSearch.focus();
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