package main

import (
	"git.leggy.dev/Fluffy/Website/internal/routes"
)

var PagesConfig = routes.PagesConfig{
	GeneratedSizes: []int{256, 512, 1024},
	Artists: []routes.AristConfig{
		{
			Name: "Anchee",
			NSFW: true,
			Socials: map[string]string{
				"Twitter":     "https://twitter.com/AnChee950249",
				"Furaffinity": "https://www.furaffinity.net/user/anchee",
			},
		},
		{
			Name: "Jaylacine",
			NSFW: false,
			Socials: map[string]string{
				"Bluesky":     "https://bsky.app/profile/jaylacine.bsky.social",
				"Furaffinity": "https://www.furaffinity.net/user/galeb",
			},
		},
		{
			Name: "Easyvock",
			NSFW: false,
			Socials: map[string]string{
				"Bluesky":     "https://bsky.app/profile/easyvock.bsky.social",
				"Furaffinity": "https://www.furaffinity.net/user/sinser115",
			},
		},
		{
			Name: "Oggy",
			NSFW: false,
			Socials: map[string]string{
				"Twitter":     "https://twitter.com/OggyOsbourne",
				"Bluesky":     "https://bsky.app/profile/oggythefox.bsky.social",
				"Furaffinity": "https://www.furaffinity.net/user/oggythefox",
			},
		},
		{
			Name: "Zadok",
			NSFW: false,
			Socials: map[string]string{
				"Twitter": "https://twitter.com/Zadoktater",
			},
		},
		{
			Name: "Pulex",
			NSFW: false,
			Socials: map[string]string{
				"Bluesky":     "https://bsky.app/profile/pulex.bsky.social",
				"Furaffinity": "https://www.furaffinity.net/user/pulex",
			},
		},
		{
			Name: "Shep",
			NSFW: false,
			Socials: map[string]string{
				"Twitter": "https://twitter.com/ShepGoesBlep",
				"Bluesky": "https://bsky.app/profile/sheppybleppy.bsky.social",
			},
		},
	},
	Art: []routes.ArtConfig{
		{Path: "/static/images/art/anchee/wink.jpg", Artist: "Anchee"},
		{Path: "/static/images/art/jaylacine/mr_wolp.png", Artist: "Jaylacine"},
		{Path: "/static/images/art/jaylacine/bint.png", Artist: "Jaylacine"},
		{Path: "/static/images/art/easyvock/autistic_stare.png", Artist: "Easyvock"},
		{Path: "/static/images/art/oggy/67.png", Artist: "Oggy"},
		{Path: "/static/images/art/oggy/mood.png", Artist: "Oggy"},
		{Path: "/static/images/art/zadok/taidum.png", Artist: "Zadok"},
		{Path: "/static/images/art/zadok/cross_eye.png", Artist: "Zadok"},
		{Path: "/static/images/art/pulex/please.png", Artist: "Pulex"},
		{Path: "/static/images/art/pulex/kiss.png", Artist: "Pulex"},
		{Path: "/static/images/art/pulex/eeee.png", Artist: "Pulex"},
		{Path: "/static/images/art/shep/sneak.png", Artist: "Shep"},
		{Path: "/static/images/art/shep/heart.png", Artist: "Shep"},
	},
}
