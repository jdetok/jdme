// REPLACES /www/js/listen.js

import { foldedLog, MSG } from "./global.js";

const mq = window.matchMedia("(max-width: 1350px)"); 
const toChange = {
        abt: {
            lg_txt: "about jdeko.me",
            sm_txt: "about",
        },
        src: {
            lg_txt: "jdeko.me source code",
            sm_txt: "source code",
        },
        resume: {
            lg_txt: "professional resume",
            sm_txt: "resume",
        },
        bball: {
            lg_txt: "nba/wnba stats api",
            sm_txt: "stats api",
        },
        tech: {
            lg_txt: "supporting technologies",
            sm_txt: "technologies",
        },
        github: {
            lg_txt: "jdetok on github",
            sm_txt: "github",
        },
        linkedin: {
            lg_txt: "linkedin",
            sm_txt: "linkedin"
        }
    };

document.addEventListener('DOMContentLoaded', async () => {
    await mediaQueryMenuSizes(mq);
    mq.addEventListener("change",  async (e: MediaQueryListEvent) => {
        await mediaQueryMenuSizes(e);
    });
});

async function mediaQueryMenuSizes(e: MediaQueryListEvent | MediaQueryList): Promise<void> {
    foldedLog(`%csetting page headers...`, MSG);
    const matches = e.matches ?? mq.matches; 

    for (const [elName, val] of Object.entries(toChange)) {
        const el = document.getElementById(elName);
        if (!el) continue;
        el.textContent = matches ? val.sm_txt : val.lg_txt;
    }

}


