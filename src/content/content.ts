import { z, defineCollection, reference } from "astro:content";

const posts = defineCollection({
    type: "content",
    schema: z.object({
        draft: z.boolean().optional().default(false),
        title: z.string(),
        description: z.string(),
        pubDate: z.string().transform((str) => new Date(str)),
        tags: z.array(reference("tags")),
    }),
});

const certificates = defineCollection({
    type: "data",
    schema: z.object({
        title: z.string(),
        provider: z.string(),
        achieved: z.date(),
        skills: z.array(z.string()).optional(),
        link: z.string().optional(),
    }),
});

const tags = defineCollection({
    type: "content",
    schema: z.object({
        name: z.string(),
    })
})

export const collections = {
    posts,
    certificates,
    tags,
};
