// @vitest-environment jsdom
import { describe, it, expect, vi, beforeEach } from "vitest";
import { mediaQueryMenuSizes } from "../../src/hdr_ftr.js";

// Mock window.matchMedia
beforeEach(() => {
    Object.defineProperty(window, 'matchMedia', {
        writable: true,
        value: vi.fn().mockImplementation(query => ({
            matches: false,
            media: query,
            addEventListener: vi.fn(),
            removeEventListener: vi.fn(),
            onchange: null,
            dispatchEvent: vi.fn()
        }))
    });
});

describe('mediaQueryMenuSizes', () => {
    it('can call mediaQueryMenuSizes without error', async () => {
        const e = { matches: false } as MediaQueryList;
        await mediaQueryMenuSizes(e);

        // nothing to assert here yet, just ensure it runs
        expect(window.matchMedia).toHaveBeenCalled;
    });
});