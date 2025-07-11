export async function clear() {
    const mdiv = document.getElementById('main');
    const btn = document.getElementById('clear');
    btn.addEventListener('click', async (event) => {
        event.preventDefault();
        mdiv.textContent = '';
    })
}