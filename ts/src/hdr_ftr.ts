import { foldedLog, homeurl, MSG, RED_BOLD, toTop } from "./global.js";

const CHANGE_AT_WIDTH = 1350;

const media_query = window.matchMedia(`(max-width: ${CHANGE_AT_WIDTH}px)`); 

document.addEventListener('DOMContentLoaded', async () => {
    await mediaQueryMenuSizes(media_query);
    media_query.addEventListener("change",  async (e: MediaQueryListEvent) => {
        await mediaQueryMenuSizes(e);
    });
});

async function mediaQueryMenuSizes(e: MediaQueryListEvent | MediaQueryList): Promise<void> {
    foldedLog(`%csetting page headers for ${window.location.pathname}...`, MSG);
    const matches = e.matches ?? media_query.matches; 

    foldedLog(`%cbuilding header`, MSG);
    const vals = Object.values(headers);
    const hdr = document.getElementById('hdr');
    if (!hdr) {
        foldedLog(`%cftr el not found`, RED_BOLD);
        return;
    }
    hdr.innerHTML = ''; 
    let to_append: HTMLDivElement[] = [];
    for (const val of vals) {
        if (val.path === window.location.pathname) continue;
        const label_txt = matches ? val.sm_txt : val.lg_txt;
        const el_txt = (val.link ? 
            `<a href="${val.link}"${val.blank ? ' target="_blank"' : ''}>${label_txt}</a>` 
            : label_txt
        );
        const d = document.createElement('div');
        d.innerHTML = el_txt;
        to_append.push(d);
    }
    hdr.style.gridTemplateColumns = `repeat(${to_append.length}, 1fr)`;
    for (const d of to_append) {
        hdr.appendChild(d);
    }
    

    foldedLog(`%cbuilding footer`, MSG);
    const ftr_vals = Object.values(footers);
    const ftr = document.getElementById('ftr');
    if (!ftr) {
        foldedLog(`%cftr el not found`, RED_BOLD);
        return;
    }
    ftr.innerHTML = ''; 
    let ftrs_to_append: HTMLDivElement[] = [];
    for (const val of ftr_vals) {
        const label_txt = matches ? val.sm_txt : val.lg_txt;
        const d = document.createElement('div');
        if (val.link === "top") {
            const btn = document.createElement('button');
            btn.textContent = label_txt ?? '';
            btn.addEventListener('click', () => {
                toTop();
            });
            d.appendChild(btn);
        } else if (val.link) {
            const a = document.createElement('a');
            a.href = val.link;
            if (val.blank) a.target = "_blank";
            a.textContent = label_txt ?? '';
            d.appendChild(a);
        } else {
            d.textContent = label_txt ?? '';
        }
        ftrs_to_append.push(d);
    }
    ftr.style.gridTemplateColumns = `repeat(${ftrs_to_append.length}, 1fr)`;
    for (const d of ftrs_to_append) {
        ftr.appendChild(d);
    }
}

type hdrftr = {
    path: string | false,
    link: string | "top" | false,
    blank: boolean,
    lg_txt: string,
    sm_txt: string | null,
};

const headers = {
    home: {
        path: "/",
        link: `${homeurl}`,
        blank: false,
        lg_txt: "jdeko.me",
        sm_txt: "jdeko.me",
    },
    abt: {
        path: "/about/",
        link: `${homeurl}about/`,
        blank: false,
        lg_txt: "about jdeko.me",
        sm_txt: "about",
    },
    src: {
        path: false,
        link: "https://github.com/jdetok/jdme",
        blank: true,
        lg_txt: "jdeko.me source code",
        sm_txt: "source code",
    },
    bball: {
        path: "/bball/",
        link: `${homeurl}bball/`,
        blank: false,
        lg_txt: "nba/wnba stats api",
        sm_txt: "stats api",
    },
    tech: {
        path: "/tech/",
        link: `${homeurl}tech/`,
        blank: false,
        lg_txt: "supporting technologies",
        sm_txt: "technologies",
    },
    github: {
        path: false,
        link: "https://github.com/jdetok",
        blank: true,
        lg_txt: "jdetok on github",
        sm_txt: "github",
    },
    resume: {
        path: "/resume/cv/",
        link: `${homeurl}resume/cv/`,
        blank: false,
        lg_txt: "professional resume",
        sm_txt: "resume",
    },
    linkedin: {
        path: false,
        link: "https://www.linkedin.com/in/justin-dekock-257879185",
        blank: true,
        lg_txt: "linkedin",
        sm_txt: "linkedin"
    },
} satisfies Record<string, hdrftr>;

const footers = {
    home: {
        path: "/",
        link: `${homeurl}`,
        blank: false,
        lg_txt: "jdeko.me",
        sm_txt: "jdeko.me",
    },
    me: {
        path: false,
        link: false,
        blank: false,
        lg_txt: 'created & maintained by Justin DeKock',
        sm_txt: 'by Justin DeKock',
    },
    email: {
        path: false,
        link: "mailto:jdekock17@gmail.com",
        blank: true,
        lg_txt: 'email: jdekock17@gmail.com',
        sm_txt: 'jdekock17@gmail.com',
    },
    src: {
        path: false,
        link: "https://www.linkedin.com/in/justin-dekock-257879185",
        blank: true,
        lg_txt: 'jdeko.me source code',
        sm_txt: 'source code',
    },
    top: {
        path: false,
        link: "top",
        blank: true,
        lg_txt: 'top of this page',
        sm_txt: 'top of page',
    },
} satisfies Record<string, hdrftr>;

