import { foldedLog, MSG } from "../global.js";
import { imageDiv, NBA_WNBA_LOGO_IMGS } from "./elements.js";

export async function makeLogoImgs() {
    const imgs = await fillImageDiv(NBA_WNBA_LOGO_IMGS);
    await normalizeImgHeights(imgs);
}

export async function normalizeImgHeights(imgs: HTMLImageElement[]): Promise<void> {
    const minHeight = Math.min(...imgs.map(img => img.getBoundingClientRect().height));
    foldedLog(`%cusing image height ${minHeight}px`, MSG);
    
    const pxBoostImg = window.innerWidth <= 900 ? 50 : 100;
    await changeImgPxBoost(pxBoostImg, minHeight, imgs);
}

export async function changeImgPxBoost(px: number, minHeight: number, imgs: HTMLImageElement[] = []) {
    for (const img of imgs) {
        img.style.height = `${minHeight + px}px`;
    }
}

export async function fillImageDiv(idiv: imageDiv): Promise<HTMLImageElement[]> {
    const d = document.getElementById(idiv.el);
    if (!d) throw new Error(`couldnt' get response title element at ${idiv.el}`);
    d.innerHTML = '';
    let imgs: HTMLImageElement[] = [];
    for (const im of idiv.imgs) {
        const img = document.createElement('img');
        img.src = im.url;
        img.alt = im.alt ?? 'image not found';
        imgs.push(img);
        d.appendChild(img)
    }
    return imgs;
}