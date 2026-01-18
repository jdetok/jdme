// REPLACES /www/js/bball/util.js
export const base = "https://jdeko.me/bball";
export const homeurl = "https://jdeko.me/";
export const checkBoxEls = ['post', 'reg', 'nbaTm', 'wnbaTm'];
export const mediaQueryBreak = 850;
export const MSG = `color: mediumseagreen;`;
export const SBL = `color: skyblue;`;
export const MSG_BOLD = `color: mediumseagreen; font-weight: bold;`;
export const RED_BOLD = 'color: red; font-weight: bold;';
export async function bytes_in_resp(r) {
    const buf = await r.clone().arrayBuffer();
    return buf.byteLength;
}
export function foldedErr(...args) {
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
export function foldedLog(...args) {
    console.groupCollapsed(...args);
    console.trace();
    console.groupEnd();
}
export async function logResp(url, r) {
    console.groupCollapsed(`%crequesting ${url}...`, SBL);
    console.trace();
    console.log(`%c${await bytes_in_resp(r)} bytes received from ${url}}`, MSG);
    console.groupEnd();
}
export function scrollIntoBySize(wpx, hpx, el) {
    if (window.innerWidth <= wpx || window.innerHeight <= hpx) {
        const res = document.getElementById(el);
        if (!res) {
            foldedErr(`couldnt' find elementent with id=${el}`);
            return;
        }
        res.scrollIntoView({ behavior: "smooth", block: "start" });
    }
}
export function toTop(el = 'hdr') {
    const hdr = document.getElementById(el);
    if (!hdr) {
        foldedErr(`couldnt' find elementent with id=${el}`);
        return;
    }
    hdr.scrollIntoView({ behavior: "smooth", block: "start" });
}
export async function fetchJSON(url) {
    let r;
    try {
        r = await fetch(url);
    }
    catch (e) {
        foldedErr(`fetch error for ${url}: ${e}`);
        return;
    }
    await logResp(url, r);
    return await r.json();
}
//# sourceMappingURL=global.js.map