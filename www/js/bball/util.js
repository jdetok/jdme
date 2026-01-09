export const base = "https://dev.jdeko.me/bball";
export const checkBoxEls = ['post', 'reg', 'nbaTm', 'wnbaTm'];

export const BLU_BOLD = 'color: blue; font-weight: bold;'
export const RED_BOLD = 'color: red; font-weight: bold;'
export const AQUA_BOLD = 'color: aqua; font-weight: bold;'
export const GRN_BOLD = 'color: green; font-weight: bold;'
export const PRP_BOLD = 'color: purple; font-weight: bold;'
export const FUSC_BOLD = 'color: fuchsia; font-weight: bold;'
export const YLW_BOLD = 'color: yellow; font-weight: bold;'
export const MSG_BOLD = `color: mediumseagreen; font-weight: bold;`

export const BLU = 'color: blue;'
export const AQUA = 'color: aqua;'
export const GRN = 'color: green;'
export const PRP = 'color: purple;'
export const FUSC = 'color: fuchsia;'
export const YLW = 'color: yellow;'
export const MSG = `color: mediumseagreen;`

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
            res.scrollIntoView({behavior: "smooth", block: "start"});
        }
    }
}