import { defineConfig } from "astro/config";

import syntaxTheme from "./syntax-theme.json";

// https://astro.build/config
export default defineConfig({
    output: "static",
    markdown: {
        syntaxHighlight: "shiki",
        shikiConfig: {
            theme: syntaxTheme,
        },
    },
});
