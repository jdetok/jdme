export const BLU_BOLD = 'color: blue; font-weight: bold;'
export const AQUA_BOLD = 'color: aqua; font-weight: bold;'
export const GRN_BOLD = 'color: green; font-weight: bold;'
export const PRP_BOLD = 'color: purple; font-weight: bold;'
export const FUSC_BOLD = 'color: fuchsia; font-weight: bold;'
export const YLW_BOLD = 'color: yellow; font-weight: bold;'

export const BLU = 'color: blue;'
export const AQUA = 'color: aqua;'
export const GRN = 'color: green;'
export const PRP = 'color: purple;'
export const FUSC = 'color: fuchsia;'
export const YLW = 'color: yellow; font-weight: bold;'

export async function bytes_in_resp(r) {
    const buf = await r.clone().arrayBuffer();
    return buf.byteLength;
}