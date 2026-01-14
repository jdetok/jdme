// REPLACES /www/js/listen.js
import { foldedLog, MSG, RED_BOLD } from "./global.js";
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
const ftrsToChange = {
    home: {
        lg_txt: '<a href="https://dev.jdeko.me/">jdeko.me</a>',
        sm_txt: '<a href="https://dev.jdeko.me/">jdeko.me</a>',
    },
    me: {
        lg_txt: 'created & maintained by Justin DeKock',
        sm_txt: 'by Justin DeKock',
    },
    email: {
        lg_txt: 'contact me at: <a href="mailto:jdekock17@gmail.com">jdekock17@gmail.com</a>',
        sm_txt: '<a href="mailto:jdekock17@gmail.com">jdekock17@gmail.com</a>',
    },
    src: {
        lg_txt: '<a id="src" target="_blank" href="https://github.com/jdetok/go-api-jdeko.me">jdeko.me source code</a>',
        sm_txt: '<a id="src" target="_blank" href="https://github.com/jdetok/go-api-jdeko.me">source code</a>',
    },
};
async function mediaQueryMenuSizes(e) {
    foldedLog(`%csetting page headers...`, MSG);
    const matches = e.matches ?? mq.matches;
    for (const [elName, val] of Object.entries(toChange)) {
        const el = document.getElementById(elName);
        if (!el)
            continue;
        el.textContent = matches ? val.sm_txt : val.lg_txt;
    }
    foldedLog(`%cbuilding footer`, MSG);
    const ftr = document.getElementById('ftr');
    if (!ftr) {
        foldedLog(`%cftr el not found`, RED_BOLD);
        return;
    }
    ftr.innerHTML = '';
    for (const val of Object.values(ftrsToChange)) {
        const d = document.createElement('div');
        d.innerHTML = matches ? val.sm_txt : val.lg_txt;
        ftr.appendChild(d);
    }
}
document.addEventListener('DOMContentLoaded', async () => {
    await mediaQueryMenuSizes(mq);
    mq.addEventListener("change", async (e) => {
        await mediaQueryMenuSizes(e);
    });
});
// document.addEventListener('DOMContentLoaded', async () => {
//     foldedLog(`%cbuilding footer`, MSG);
//     const ftr = document.getElementById('ftr');
//     if (!ftr) {
//         foldedLog(`%cftr el not found`, RED_BOLD);
//         return;
//     }
//     for (const ftrHTML of ftrs) {
//         const d = document.createElement('div');
//         d.innerHTML = ftrHTML;
//         ftr.appendChild(d);
//     }
// });
//# sourceMappingURL=hdr_ftr.js.map