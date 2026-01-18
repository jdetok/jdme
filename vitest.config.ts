import { defineConfig } from 'vitest/config';

export default defineConfig({
    test: {
        include: ['ts/tests/**/*.test.ts'],
        globals: true,
        environment: 'node',
        typecheck: {
            tsconfig: 'ts/tests/tsconfig.json'
        }
    }
});
