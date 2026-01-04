const mq = window.matchMedia("(max-width: 1200px)"); 
const toChange = {
        home: {
            lg_txt: "jdeko.me",
            sm_txt: "home",
        },
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
    mediaQueryMenuSizes(mq);
    mq.addEventListener("change", mediaQueryMenuSizes);
});

function mediaQueryMenuSizes(e) {
    const matches = e.matches ?? mq.matches; 

    for (const [elName, val] of Object.entries(toChange)) {
        console.log(`changing ${elName}...`);
        const el = document.getElementById(elName);
        if (!el) continue;
        el.textContent = matches ? val.sm_txt : val.lg_txt;
    }

}


