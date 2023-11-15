package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/itslychee/mpdrp/mpd"
)

type ReleaseGroups struct {
	Count         int `json:"count"`
	ReleaseGroups []struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		PrimaryType string `json:"primary-type"`
	} `json:"release-groups"`
}

func GetCoverArt(r mpd.Response) (cover string, err error) {
	if *noAlbumCovers {
		log(Debug, "not retrieving album cover as configured")
		cover = NoAlbumAsset
		return
	}

	v, ok := r.Records["musicbrainz_releasegroupid"]
	if !ok {
		var query strings.Builder
		// FIXME: Possible index error
		album, ok := r.Records["Album"]
		if ok {
		  query.WriteString(fmt.Sprintf("releasegroup:%s ", strconv.Quote(strings.TrimSpace(album[0]))))
		}
		artist, ok := r.Records["Artist"]
		if ok {
		  query.WriteString(fmt.Sprintf("artist:%s ", strconv.Quote(strings.TrimSpace(artist[0]))))
		}
		albumArtist, ok := r.Records["AlbumArtist"]
		if ok {
		  query.WriteString(fmt.Sprintf("albumartist:%s ", strconv.Quote(strings.TrimSpace(albumArtist[0]))))
		}
		title, ok := r.Records["Title"]
		if ok {
		  query.WriteString(fmt.Sprintf("title:%s ", strconv.Quote(strings.TrimSpace(title[0]))))
		}

		if query.String() == "" {
			log(Normal, "not enough metadata to use in order to search for song's album cover")
			cover = NoAlbumAsset
			return
		}
		req := &http.Request{
			Method: "GET",
			URL: &url.URL{
				Scheme: "https",
				Host:   "musicbrainz.org",
				Path:   "/ws/2/release-group",
				RawQuery: url.Values{
					"query": []string{query.String()},
					"limit": []string{"1"},
				}.Encode(),
			},
			Header: http.Header{
				"Accept":     []string{"application/json"},
				"User-Agent": []string{"mpdrp (https://github.com/itslychee/mpdrp)"},
			},
		}

		var response *http.Response
		var body []byte
		var d ReleaseGroups
		response, err = http.DefaultClient.Do(req)
		if err != nil {
			log(Normal, "encountered error during http req")
			return
		}
		body, err = io.ReadAll(response.Body)
		if err != nil {
			return
		}

		logjson(Network, "MusicBrainz", json.RawMessage(body))

		err = json.Unmarshal(body, &d)
		if err != nil {
			return
		}
		if len(d.ReleaseGroups) == 0 {
			log(Normal, "no releases")
			return
		}
		rel := d.ReleaseGroups[0].ID
		resp, err := http.Get(fmt.Sprintf("https://coverartarchive.org/release-group/%s/front-250", rel))
		if err != nil {
			return "", err
		}
		if resp.StatusCode != 200 {
			log(Network, "bad asset, returning default album asset")
			v = []string{NoAlbumAsset}
		} else {
			v = []string{fmt.Sprintf("https://coverartarchive.org/release-group/%s/front-250", rel)}
		}
	}
	cover = v[0]
	return 
}




