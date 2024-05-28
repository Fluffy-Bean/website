import {
    type CollectionEntry,
    type ContentEntryMap,
    getCollection,
} from "astro:content";

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

export function getMonth(date: Date): string {
    const months = [
        "January",
        "February",
        "March",
        "April",
        "May",
        "June",
        "July",
        "August",
        "September",
        "October",
        "November",
        "December",
    ];

    return months[date.getMonth()];
}

export function getDaySuffix(date: Date): string {
    let suffix = "th";
    if (date.getDate() % 10 === 1 && date.getDate() !== 11) {
        suffix = "st";
    } else if (date.getDate() % 10 === 2 && date.getDate() !== 12) {
        suffix = "nd";
    } else if (date.getDate() % 10 === 3 && date.getDate() !== 13) {
        suffix = "rd";
    }

    return suffix;
}

export async function getTagsBySlug(
    postTags: string[],
): Promise<CollectionEntry<"tags">[]> {
    const allTags: CollectionEntry<"tags">[] = await getCollection("tags");
    const tags: CollectionEntry<"tags">[] = [];
    postTags.forEach((postTag) => {
        allTags.forEach((allTag) => {
            if (allTag.slug === postTag) tags.push(allTag);
        });
    });
    return tags;
}
