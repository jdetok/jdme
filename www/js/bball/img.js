import { foldedLog, MSG } from "../global.js";
export async function makeLogoImgs() {
    const imgs = await fillImageDiv({
        el: 'lg_imgs', imgs: [
            {
                url: 'https://cdn.nba.com/logos/leagues/logo-wnba.svg',
                alt: 'could not load WNBA logo',
            },
            {
                url: 'https://cdn.nba.com/logos/leagues/logo-nba.svg',
                alt: 'could not load NBA logo',
            },
        ]
    });
    await normalizeImgHeights(imgs);
}
export async function normalizeImgHeights(imgs) {
    const minHeight = Math.min(...imgs.map(img => img.getBoundingClientRect().height));
    foldedLog(`%cusing image height ${minHeight}px`, MSG);
    const pxBoostImg = window.innerWidth <= 900 ? 50 : 100;
    await changeImgPxBoost(pxBoostImg, minHeight, imgs);
}
export async function changeImgPxBoost(px, minHeight, imgs = []) {
    for (const img of imgs) {
        img.style.height = `${minHeight + px}px`;
    }
}
export async function fillImageDiv(idiv) {
    const d = document.getElementById(idiv.el);
    if (!d)
        throw new Error(`couldnt' get response title element at ${idiv.el}`);
    d.innerHTML = '';
    let imgs = [];
    for (const im of idiv.imgs) {
        const img = document.createElement('img');
        img.src = im.url;
        img.alt = im.alt ?? 'image not found';
        imgs.push(img);
        d.appendChild(img);
    }
    return imgs;
}
//# sourceMappingURL=img.js.map