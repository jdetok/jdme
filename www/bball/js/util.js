export async function bytes_in_resp(r) {
    const buf = await r.clone().arrayBuffer();
    return buf.byteLength;
}

export async function logresp(r, url) {
    const style = 'color: green; font-weight: bold;';
    console.trace(
        `%c request status ${r.status}: \
        ${await bytes_in_resp(r)} bytes received from ${url}}`, style);
}