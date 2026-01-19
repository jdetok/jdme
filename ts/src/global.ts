// REPLACES /www/js/bball/util.js

export const base = "https://jdeko.me/bball";
export const homeurl = "https://jdeko.me/";
export const checkBoxEls = ['post', 'reg', 'nbaTm', 'wnbaTm'] as string[];


export const WINDOWSIZE = 700;
export const BIGWINDOW = 1400;
export const LARGEROWS = 25;

export const mediaQueryBreak = 850;

export const MSG = `color: mediumseagreen;`
export const SBL = `color: skyblue;`
export const MSG_BOLD = `color: mediumseagreen; font-weight: bold;`
export const RED_BOLD = 'color: red; font-weight: bold;'

export const wsize = (): string => { return `W:${window.innerWidth}px X H:${window.innerHeight}px` }

export async function bytes_in_resp(r: Response): Promise<number> {
    const buf = await r.clone().arrayBuffer();
    return buf.byteLength;
}

export function foldedErr(...args: any[]): void {
    console.groupCollapsed(`%c** ERROR **`, RED_BOLD);
    if (args.length == 1) {
        if (args[0].substring(0, 2) !== '%c') {
            args[0] = `%c${args[0]}`;
        }
        args.push(RED_BOLD);
    }
    console.error(...args);
    console.groupEnd();
}

export function foldedLog(...args: any[]): void {
    console.groupCollapsed(...args);
    console.trace();
    console.groupEnd();
}

export async function logResp(url: string, r: Response) {
    console.groupCollapsed(`%crequesting '${url}'...`, SBL);
    console.trace();
    console.log(`%c${await bytes_in_resp(r)} bytes received from ${url}}`, MSG);
    console.groupEnd();
}

export function scrollIntoBySize(wpx: number, hpx: number, el: string): void {
    if (window.innerWidth <= wpx || window.innerHeight <= hpx) {
        const res = document.getElementById(el);
        if (!res) {
            foldedErr(`couldnt' find elementent with id=${el}`); return;
        }
        res.scrollIntoView({behavior: "smooth", block: "start"});
    }
}

export function toTop(el = 'hdr'): void {
    const hdr = document.getElementById(el);
    if (!hdr) {
        foldedErr(`couldnt' find elementent with id=${el}`); return;
    }
    hdr.scrollIntoView({behavior: "smooth", block: "start"});
}

export async function fetchJSON(url: string): Promise<any> {
    let r: Response;
    try {
        r = await fetch(url);
    } catch (e) {
        foldedErr(`fetch error for ${url}: ${e}`);
        return;
    }
    
    await logResp(url, r);
    return await r.json()
}

export async function errMsg(msg: string, el_id = 'errmsg'): Promise<void> { 
    const el = document.getElementById(el_id);
    if (!el) throw new Error(`can't find element with id ${el_id}`)
    el.textContent = msg;
    el.style.display = 'block';
}

export async function hideErr(el_id = 'errmsg'): Promise<void> {
    const el = document.getElementById(el_id);
    if (!el) throw new Error(`can't find element with id ${el_id}`)
    el.textContent = '';
    el.style.display = 'none';
}