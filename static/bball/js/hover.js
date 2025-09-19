// show/hide hover messages
export async function showHideHvr(el, hvrName, msg) {
    const hvr = document.getElementById(hvrName);
    el.addEventListener('mouseover', async (event) => {
        // event.preventDefault();
        hvr.textContent = msg;
        hvr.style.display = 'block'; 
    })
    el.addEventListener('mouseleave', async (event) => {
        // event.preventDefault();
        hvr.textContent = '';
        hvr.style.display = 'none'; 
    })
}