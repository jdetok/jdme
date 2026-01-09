export const base = "https://dev.jdeko.me/bball";
export const checkBoxEls = ['post', 'reg', 'nbaTm', 'wnbaTm'];
export const MSG = `color: mediumseagreen;`;
export const RED_BOLD = 'color: red; font-weight: bold;';
export async function bytes_in_resp(r) {
    const buf = await r.clone().arrayBuffer();
    return buf.byteLength;
}
export async function foldedLog(...args) {
    console.groupCollapsed(...args);
    console.trace();
    console.groupEnd();
}
export async function scrollIntoBySize(wpx, hpx, el) {
    if (window.innerWidth <= wpx || window.innerHeight <= hpx) {
        let res = document.getElementById(el);
        if (res) {
            res.scrollIntoView({ behavior: "smooth", block: "start" });
        }
    }
}
