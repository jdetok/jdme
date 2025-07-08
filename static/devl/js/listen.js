/* INTENT: 
This file will be the main "script" file tagged in the HTML.
Event listeners are called here
*/

import * as test from "./test.js";

document.addEventListener('DOMContentLoaded', () => {
    test.loadSeasonOpts();
    test.loadTeamOpts();
    test.lgChangeListener();
    test.gamesRecent();
    test.topScorer();
});


