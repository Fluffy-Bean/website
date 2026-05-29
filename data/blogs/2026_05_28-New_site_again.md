# New site.... again!

As seems to be a tradition with my sites, I write a blog for when I rewrite it, but then never again. I'll keep this one
short, as tbh there is not much to talk about.

Anyways, Go!

## Le Website

It's all Go now, as being employed made me realize that there is not enough time on earth for me to want to use it
anymore. Also, not doing maintenance on my website for about 2 years lead to my emails getting spammed by GitHub
alerting me of new issues with various npm packages...

![GitHub Security and quality issues tab, displaying 56 issues](/assets/blog/2026/05/Screenshot_20260529_000345.png)

There is now a chat on the site, as of now no chats are retained, more or less behaving as a broadcast. Once I have time
again, I hope to add custom usernames and message history. It uses SSE and JWT, as barebones as you can get! But it's
very reliable, yay for more free time.

Other than that, it's a pretty standard Go site, no fancy templates and no fancy CSS.

## Other projects

Been working on Urchin Engine still, other than I rewrote it again, this time going for C++ which actually has been
pretty enjoyable so far, surprisingly!

[Here's a little demo video of what I got so far](https://www.youtube.com/watch?v=RmOUC3TEfls)

The base of the engine uses a custom file format for defining how `Object`s behave and interact. It's just JSON with
fewer quotes, here's an example crate object:

```
name: "crate";

attributes: (
    { name: "draggable"; data: true; },
);

handlers: (
    { name: "click"; data: ( "core::objectCrateOpen", "debug::objectPrintHello" ); },
);

body: {
    textures: (
        { name: "open";   path: "resources/textures/generated/OpenCrate.png"; },
        { name: "closed"; path: "resources/textures/generated/Crate.png";     },
    );
    initial: {
        texture: "closed";
        animation: "default";
    };
    offset: ( -64, -94 );
};

bounds: (
    (  29,  22 ),
    (  29, -60 ),
    ( -29, -60 ),
    ( -29,  22 ),
);

collider: {
    kind: "active";
    shape: "rect";
    data: ( 28, 22 );
};
```

The language supports maps, but my C interface from rust doesn't actually support iterating over keys, so I went the
lazy array route...

This was my first rust project, but I'm pretty happy with how much I have gotten working! [Here is the git link](https://git.leggy.dev/Fluffy/UrchinEnginePlusPlus/src/branch/master/lib/config)
to the source if you, the reader, are interested in seeing the internals.

The `handlers` are called on certain events, like `update`, `fixed_update`, `click` and whatnot. It's a list of
functions to call within C++, that itself is just a map of pointers to functions. I originally wanted to write the
handlers within the configs themselves, but I struggled getting Lua working exactly how I wanted; also tried making my
own scripting language, but then got into a tangle of Rust lifetimes...

Here's how the `core::objectCrateOpen` handler looks like:

```cpp
void
core_object_crate_open(Object& object, HandlerContext& ctx)
{
    if (object.attributes["_is_open"].get_boolean())
    {
        object.texture_animation.play(object.texture_handles["closed"].get(), "default");
        object.attributes["_is_open"].set_boolean(false);
    }
    else
    {
        object.texture_animation.play(object.texture_handles["open"].get(), "default", true);
        object.attributes["_is_open"].set_boolean(true);
    }
}
```

`attributes` is just key-value storage on the `Object`. As `Object`s don't inherit and are built at runtime though these
config files, they need to have some way to apply object-specific values beyond what's available on the `Object` class.

In the crate example, it's also used to tell the engine that the `Object` is draggable by the player, it's quite funny
applying the draggable value on trees and houses.

## NFC 2026

Finally went to my first furry con this year, it was great! Was also my second time traveling by myself, double the
anxiety!

It was pretty surreal seeing the entire town of Malmö being overtaken by fursuiters, for once in my life I really did
not need to worry about getting judged or being perceived as _weird_.

Hoping to come back next year, also looking at going to Confuzzled next year - or trying at least, with their retched
lottery system.

![NFC sign covered in a thick layer of stickers](/assets/blog/2026/05/PXL_20260219_085424591.MP.jpg)

See if you can spot me.

## What's next

Nothing.

Anyway, see you maybe in a few years when I remember I have this site :p
