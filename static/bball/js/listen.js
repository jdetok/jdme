// script to load in HTML -- all listener functions are called here 
import * as ui from "./ui.js"
import * as resp from "./resp.js"
import { makeScoringLeaders } from "./lg_ldg_scorers.js"
import { makeRGTopScorers } from "./rg_ldg_scorers.js";

export const base = "http://localhost:8080/bball";

document.addEventListener('DOMContentLoaded', async () => {
    await ui.loadSznOptions();
    await ui.selHvr();
    await ui.randPlayerBtn();
    await ui.search();
    await ui.clearSearch();
    await ui.holdPlayerBtn();
    await makeRGTopScorers();
    await makeScoringLeaders(10);
});