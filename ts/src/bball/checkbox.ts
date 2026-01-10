export async function checkBoxGroupValue(lgrp, rgrp, dflt) {
    const l = await checkBoxes(lgrp.box, lgrp.slct);
    const r = await checkBoxes(rgrp.box, rgrp.slct);

    if (l) return l;
    if (r) return r;
    
    return dflt;
}

export async function checkBoxes(box: string, sel: string) {
    const b = document.getElementById(box) as HTMLInputElement;
    const s = document.getElementById(sel) as HTMLInputElement;
    if (!b || !s) {
        throw new Error(`couldn't get element with id ${box} or ${sel}`);
    }
    if (b.checked) {
        return s.value
    }
}

export async function clearCheckBoxes(boxes: string[]) {
    for (let i = 0; i < boxes.length; i++) {
        let b = document.getElementById(boxes[i]) as HTMLInputElement;
        if (!b) throw new Error(`couldn't find input element with id ${boxes[i]}`);
        b.checked = false;
    }
}

// make post + reg checkboxes exclusive (but allow neither checked)
export async function setupExclusiveCheckboxes(leftbox, rightbox) {
    let lbox = document.getElementById(leftbox) as HTMLInputElement;
    let rbox = document.getElementById(rightbox) as HTMLInputElement;
    if (!lbox || !rbox) throw new Error(`couldn't get ${lbox} or ${rbox}`);
    function handleCheck(e) {
        if (e.target.checked) {
            if (e.target === lbox) rbox.checked = false;
            if (e.target === rbox) lbox.checked = false;
        }
    }
    lbox.addEventListener("change", handleCheck);
    rbox.addEventListener("change", handleCheck);
}