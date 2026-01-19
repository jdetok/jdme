import { foldedLog, MSG } from "../global.js";

export type imageInDiv = {
    url: string,
    el?: string,
    alt?: string,
};

export type imageDiv = {
    el: string,
    imgs: imageInDiv[],
};

export const NBA_WNBA_LOGO_IMGS: imageDiv = {
    el: 'lg_imgs',
    imgs: [
        {
            url: 'https://cdn.nba.com/logos/leagues/logo-wnba.svg',
            alt: 'could not load WNBA logo',
        },
        {
            url: 'https://cdn.nba.com/logos/leagues/logo-nba.svg',
            alt: 'could not load NBA logo',
        },
    ]
};

// team logo and player headshot side by side
export function makeTmPlrImageDiv(el: string,
    opts: {
        tm_url: string,
        tm: string,
        plr_url: string,
        plr: string,
    }
): imageDiv {
    return {
        el: el,
        imgs: [
            {
                url: opts.tm_url,
                alt: `could not load team logo for ${opts.tm}`,
            },
            {
                url: opts.plr_url,
                alt: `could not load headshot for player ${opts.plr}`,
            },
        ]
    }
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

