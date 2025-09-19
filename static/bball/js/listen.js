// script to load in HTML -- all listener functions are called here 
import { loadSznOptions, selHvr, randPlayerBtn, search, clearSearch, 
    holdPlayerBtn } from "./ui.js"
import { makeScoringLeaders } from "./lg_ldg_scorers.js"
import { makeRGTopScorers } from "./rg_ldg_scorers.js";

export const base = "http://localhost:8080/bball";

document.addEventListener('DOMContentLoaded', async () => {
    await loadSznOptions();
    await selHvr();
    await randPlayerBtn();
    await search();
    await clearSearch();
    await holdPlayerBtn();
    await makeRGTopScorers();
    await makeScoringLeaders(10);
});