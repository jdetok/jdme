export async function checkBoxGroupValue(lgrp, rgrp, dflt) {
    const l = await checkBoxes(lgrp.box, lgrp.slct);
    const r = await checkBoxes(rgrp.box, rgrp.slct);
    if (l)
        return l;
    if (r)
        return r;
    // 88888 for season, 0 for team
    return dflt;
    // return `2${new Date().getFullYear()}`;
}
export async function checkBoxes(box, sel) {
    const b = document.getElementById(box);
    const s = document.getElementById(sel);
    if (!b || !s) {
        throw new Error(`couldn't get element with id ${box} or ${sel}`);
    }
    if (b.checked) {
        return s.value;
    }
}
export async function clearCheckBoxes(boxes) {
    for (let i = 0; i < boxes.length; i++) {
        let b = document.getElementById(boxes[i]);
        if (!b)
            throw new Error(`couldn't find input element with id ${boxes[i]}`);
        b.checked = false;
    }
}
//# sourceMappingURL=checkbox.js.map