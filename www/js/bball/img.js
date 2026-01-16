import { appendImg } from "./player.js";
// https://cdn.nba.com/logos/leagues/logo-nba.svg
// https://cdn.nba.com/logos/leagues/logo-nba.svg
export async function makeLogoImgs() {
    const lgs = [
        {
            el: 'nba_img',
            url: 'https://cdn.nba.com/logos/leagues/logo-nba.svg',
        },
        {
            el: 'wnba_img',
            url: 'https://cdn.nba.com/logos/leagues/logo-wnba.svg',
        },
    ];
    for (const lg of lgs) {
        await appendImg(lg.url, lg.el);
    }
}
//# sourceMappingURL=img.js.map