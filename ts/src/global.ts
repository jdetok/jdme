export const base = "https://dev.jdeko.me/bball";
export const checkBoxEls = ['post', 'reg', 'nbaTm', 'wnbaTm'] as const;

export const MSG = `color: mediumseagreen;`
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