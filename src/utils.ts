import { type ContentEntryMap, getCollection } from "astro:content";

// https://github.com/hellotham/hello-astro/blob/e05706cf488bcec6e4c5494a622eedfc4e47d763/src/config.ts#L55C1-L62C2
export async function getPosts(collection: keyof ContentEntryMap) {
    const posts = await getCollection(collection, ({ data }) => {
        return data.draft !== true;
    });
    return posts.sort((a, b) =>
        a.data.pubDate && b.data.pubDate
            ? Number(b.data.pubDate) - Number(a.data.pubDate)
            : 0,
    );
}
