/* INTENT: 
This file will be the main "script" file tagged in the HTML.
Event listeners are called here
*/

// import { loadSeasonOpts } from "./test";
import * as test from "./test.js";

// BASE URLS


document.addEventListener('DOMContentLoaded', () => {
    test.loadSeasonOpts();
    test.loadTeamOpts();
    test.lgChangeListener();
    test.gamesRecent();
    test.topScorer();
});


