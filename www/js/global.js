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
export function foldedLog(...args) {
    console.groupCollapsed(...args);
    console.trace();
    console.groupEnd();
}
export function scrollIntoBySize(wpx, hpx, el) {
    if (window.innerWidth <= wpx || window.innerHeight <= hpx) {
        let res = document.getElementById(el);
        if (res) {
            res.scrollIntoView({ behavior: "smooth", block: "start" });
        }
    }
}
export function toTop() {
    const hdr = document.getElementById('hdr');
    if (hdr) {
        hdr.scrollIntoView({ behavior: "smooth", block: "start" });
    }
}
//# sourceMappingURL=global.js.map