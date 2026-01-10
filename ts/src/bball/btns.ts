// replace some of /www/js/bball/ui.js
// import { base, checkBoxEls, foldedLog, MSG, scrollIntoBySize } from "../global.js";

export async function clearSearchBtn(): Promise<void> {
    const btn = document.getElementById('clearS');
    if (!btn) return;
    btn.addEventListener('click', (event) => {
        event.preventDefault();
        clearSearch(true);
    });
}

export function clearSearch(focus=false): void {
    const elId = 'pSearch';
    const pSearch = document.getElementById(elId) as HTMLInputElement;
    if (!pSearch) throw new Error(`couldn't get search bar element at ${elId}`);
    pSearch.value = '';
    if (focus) pSearch.focus();
}

export async function setPHold(player: string) {
    const ph = document.getElementById('pHold') as HTMLInputElement;
    ph.value = player;
}

export async function lgRadioBtns() {
    const selected = document.querySelector('input[name="leagues"]:checked') as HTMLInputElement;
    if (selected) {
        return selected.value;
    } else {
        return 'all';
    }
}
