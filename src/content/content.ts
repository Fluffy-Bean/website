import { z, defineCollection } from "astro:content";

const postsCollection = defineCollection({
    type: "content",
    schema: z.object({
        draft: z.boolean().optional(),
        title: z.string(),
        description: z.string(),
        pubDate: z.date(),
        image: z.string().optional(),
        tags: z.array(z.string()),
    }),
});

const projectsCollection = defineCollection({
    type: "content",
    schema: z.object({
        draft: z.boolean().optional(),
        title: z.string(),
        description: z.string(),
        image: z.string().optional(),
        tags: z.array(z.string()),
    }),
});

const certificatesCollection = defineCollection({
    type: "data",
    schema: z.object({
        title: z.string(),
        provider: z.string(),
        achieved: z.date(),
        skills: z.array(z.string()).optional(),
        link: z.string().optional(),
    }),
});

export const collections = {
    posts: postsCollection,
    projects: projectsCollection,
    certificates: certificatesCollection,
};
