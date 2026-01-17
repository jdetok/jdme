import { foldedLog, MSG } from "../global.js";
export async function makeLogoImgs() {
    const imgs = await createImages();
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
export async function newImg(url, elId) {
    const d = document.getElementById(elId);
    if (!d)
        throw new Error(`couldnt' get response title element at ${elId}`);
    const img = document.createElement('img');
    d.textContent = '';
    img.src = url;
    img.alt = "image not found";
    d.append(img);
    return img;
}
export async function createImages(lgs = [
    { el: 'wnba_img', url: 'https://cdn.nba.com/logos/leagues/logo-wnba.svg' },
    { el: 'nba_img', url: 'https://cdn.nba.com/logos/leagues/logo-nba.svg' },
]) {
    const imgs = [];
    for (const lg of lgs) {
        const img = await newImg(lg.url, lg.el);
        imgs.push(img);
    }
    // wait for all images to have real sizes
    await Promise.all(imgs.map(img => img.complete ? Promise.resolve() : new Promise(r => img.onload = r)));
    return imgs;
}
//# sourceMappingURL=img.js.map