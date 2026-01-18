// @vitest-environment jsdom
import { toTop } from '../../src/global.js';
import { vi } from 'vitest';

test('toTop scrolls header into view', () => {
    document.body.innerHTML = `<div id="hdr"></div>`;

    // jsdom does not implement this — define it
    Element.prototype.scrollIntoView = vi.fn();

    toTop('hdr');

    expect(Element.prototype.scrollIntoView).toHaveBeenCalled();
});