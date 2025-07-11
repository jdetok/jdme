// /* INTENT:
// functions that create & return html elements
// */

// // make html image with built url
// export async function makeImg(url) {
//     return new Promise((resolve, reject) => {
//         const img = document.createElement('img');
//         img.src = url;
//         img.alt = "image not found"
//         img.onload = () => resolve(img);
//         img.onerror = reject;
//     });
// }

// // append image to container element
// export async function appendImg(img, el) {
//     const container = document.getElementById(el);
//     container.innerHTML = '';
//     container.appendChild(img);
// }