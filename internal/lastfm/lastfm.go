package lastfm

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const baseURL = "https://ws.audioscrobbler.com/2.0/"

type LatestSong struct {
	Title     string
	Artist    string
	Album     string
	Thumbnail string
}

type LastFM struct {
	latest *LatestSong

	mut    sync.Mutex
	cancel context.CancelFunc
	key    string
}

func NewLastFM(ctx context.Context, key string) *LastFM {
	ctx, cancel := context.WithCancel(ctx)

	l := &LastFM{
		cancel: cancel,
		key:    key,
	}

	// ToDo: Maybe this should dynamically scale? If same song has been displayed for a while, or no little/no visitors
	//       recently, no reason to query so often
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		l.updateLatestSong()

		for {
			select {
			case <-ticker.C:
				l.updateLatestSong()

			case <-ctx.Done():
				return
			}
		}
	}()

	return l
}

func (l *LastFM) GetLatestSong() *LatestSong {
	l.mut.Lock()
	defer l.mut.Unlock()

	return l.latest
}

func (l *LastFM) updateLatestSong() {
	l.mut.Lock()
	defer l.mut.Unlock()

	values := url.Values{}
	values.Set("method", "user.getrecenttracks")
	values.Set("limit", "1")
	values.Set("format", "json")
	values.Set("user", "Fluffy_Bean_")
	values.Set("api_key", l.key)

	res, err := http.Get(baseURL + "?" + values.Encode())
	if err != nil {
		slog.Error("query latest song", "error", err)
		return
	}
	defer res.Body.Close()

	var data struct {
		RecentTracks struct {
			Tracks []struct {
				Name   string `json:"name"`
				URL    string `json:"url"`
				Artist struct {
					Text string `json:"#text"`
				} `json:"artist"`
				Images []struct {
					Text string `json:"#text"`
				} `json:"image"`
				Album struct {
					Text string `json:"#text"`
				} `json:"album"`
			} `json:"track"`
		} `json:"recenttracks"`
	}

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		slog.Error("query latest song", "error", err)

		return
	}

	if len(data.RecentTracks.Tracks) == 0 {
		return
	}

	track := data.RecentTracks.Tracks[0]
	image := track.Images[2]

	l.latest = &LatestSong{
		Title:     track.Name,
		Artist:    track.Artist.Text,
		Album:     track.Album.Text,
		Thumbnail: image.Text,
	}
}
