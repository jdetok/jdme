export const NBA_WNBA_LOGO_IMGS = {
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
export function makeTmPlrImageDiv(el, opts) {
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
    };
}
//# sourceMappingURL=elements.js.map