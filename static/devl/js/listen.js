/* INTENT: 
This file will be the main "script" file tagged in the HTML.
Event listeners are called here
*/
// import { loadSeasonOpts, loadTeamOpts, lgChangeListener } from "./test.js"
import * as home from "./home.js";
import * as buttons from "./buttons.js"
document.addEventListener('DOMContentLoaded', async () => {
    await home.loadSeasonOpts();
    await home.loadTeamOpts();
    await home.gamesRecent();
    await home.topScorer();
    await home.lgChangeListener();
    await buttons.clear();
});



// document.addEventListener('DOMContentLoaded', () => {
//     test.loadSeasonOpts();
//     test.loadTeamOpts();
//     test.lgChangeListener();
//     home.gamesRecent();
//     test.topScorer();
// });


