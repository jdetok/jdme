// REPLACES /www/js/bball/util.js

export const base = "https://dev.jdeko.me/bball";
export const homeurl = "https://dev.jdeko.me/";
export const checkBoxEls = ['post', 'reg', 'nbaTm', 'wnbaTm'] as string[];

export const mediaQueryBreak = 850;

export const MSG = `color: mediumseagreen;`
export const SBL = `color: skyblue;`
export const MSG_BOLD = `color: mediumseagreen; font-weight: bold;`
export const RED_BOLD = 'color: red; font-weight: bold;'

export async function bytes_in_resp(r: Response): Promise<number> {
    const buf = await r.clone().arrayBuffer();
    return buf.byteLength;
}

export function foldedLog(...args: any[]): void {
    console.groupCollapsed(...args);
    console.trace();
    console.groupEnd();
}

export function scrollIntoBySize(wpx: number, hpx: number, el: string): void {
    if (window.innerWidth <= wpx || window.innerHeight <= hpx) {
        let res = document.getElementById(el);
        if (res) {
            res.scrollIntoView({behavior: "smooth", block: "start"});
        }
    }
}

export function toTop(): void {
    const hdr = document.getElementById('hdr');
    if (hdr) {
        hdr.scrollIntoView({behavior: "smooth", block: "start"});
    }
}

export async function fetchJSON(url: string): Promise<any> {
    const r = await fetch(url);
    foldedLog(`%c ${await bytes_in_resp(r)} bytes received from ${url}}`, MSG_BOLD);
    return await r.json()
}