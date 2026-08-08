package lastfm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

	"git.leggy.dev/Fluffy/Website/internal/broker"
	"git.leggy.dev/Fluffy/Website/internal/events"
)

const baseURL = "https://ws.audioscrobbler.com/2.0/"

var imageURLs = []string{"https://lastfm.freetls.fastly.net", "https://lastfm-img.freetls.fastly.net"}
var transparentPixel = []byte("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII=")

type LatestSong struct {
	Title     string
	Artist    string
	Album     string
	Thumbnail []byte
}

type LastFM struct {
	Events *broker.Broker
	mut    sync.Mutex
	cancel context.CancelFunc
	key    string
	latest *LatestSong
}

func NewLastFM(ctx context.Context, key string) *LastFM {
	ctx, cancel := context.WithCancel(ctx)

	l := &LastFM{
		Events: broker.NewBroker(),
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
	track, err := l.requestLatestSong()
	if err != nil {
		slog.Error("request latest song", "error", err)

		return
	}

	l.mut.Lock()
	defer l.mut.Unlock()

	if l.latest != nil && l.latest.Title != track.Title {
		go l.Events.BroadcastEvent(events.NewSong{
			Title:  track.Title,
			Artist: track.Artist,
		})
	}

	l.latest = track
}

func (l *LastFM) requestLatestSong() (*LatestSong, error) {
	values := url.Values{}
	values.Set("method", "user.getrecenttracks")
	values.Set("limit", "1")
	values.Set("format", "json")
	values.Set("user", "Fluffy_Bean_")
	values.Set("api_key", l.key)

	res, err := http.Get(baseURL + "?" + values.Encode())
	if err != nil {
		slog.Error("get latest song", "error", err)

		return nil, err
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
		slog.Error("get latest song", "error", err)

		return nil, err
	}

	if len(data.RecentTracks.Tracks) == 0 {
		return nil, fmt.Errorf("no recent tracks")
	}

	track := data.RecentTracks.Tracks[0]

	latest := &LatestSong{
		Title:     track.Name,
		Artist:    track.Artist.Text,
		Album:     track.Album.Text,
		Thumbnail: transparentPixel,
	}

	image := track.Images[2]

	safePrefix := slices.ContainsFunc(imageURLs, func(url string) bool {
		return strings.HasPrefix(image.Text, url)
	})

	if !safePrefix {
		slog.Debug("suspicious latest song thumbnail endpoint", "url", image.Text)

		return latest, nil
	}

	current := l.GetLatestSong()
	if current != nil && current.Title != latest.Title {
		return latest, nil
	}

	res, err = http.Get(image.Text)
	if err != nil {
		slog.Error("get latest song thumbnail", "url", image.Text, "error", err)

		return latest, nil
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		slog.Error("get latest song thumbnail", "url", image.Text, "error", res.StatusCode)

		return latest, nil
	}

	thumbnail, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("read latest song thumbnail", "url", image.Text, "error", err)

		return latest, nil
	}

	if len(thumbnail) > 0 {
		latest.Thumbnail = thumbnail
	}

	return latest, nil
}
