import { bytes_in_resp } from "../../src/global.js"

test('bytes_in_resp returns correct byte size', async () => {
    const r = new Response(new Uint8Array(10));
    expect(await bytes_in_resp(r)).toBe(10);
});