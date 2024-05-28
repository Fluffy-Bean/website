import { defineConfig } from "astro/config";
import syntaxTheme from "./syntax-theme.json";

import mdx from "@astrojs/mdx";

// https://astro.build/config
export default defineConfig({
    output: "static",
    markdown: {
        syntaxHighlight: "shiki",
        shikiConfig: {
            theme: syntaxTheme,
        },
    },
    integrations: [mdx()],
});
