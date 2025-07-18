// TODO: GLOBAL HOVER FUNCTION
export async function showHideHvr(el, hvrName) {
    const hvr = document.getElementById(hvrName);
    el.addEventListener('mouseover', async (event) => {
        event.preventDefault();
        hvr.style.display = 'block'; 
    })
    el.addEventListener('mouseleave', async (event) => {
        event.preventDefault();
        hvr.style.display = 'none'; 
    })
}