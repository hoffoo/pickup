package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/werkshy/pickup/model"
	"log"
	"net/http"
	"strings"
	//"time"
)

type AlbumHandler struct {
	Music model.Collection
}

// Return a list of albums or a specific album
func (h AlbumHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/albums/"):]
	parts := strings.SplitN(path, "/", 3)
	fmt.Printf("Path: %s  parts: %q   len(parts): %d\n", r.URL.Path, parts,
		len(parts))
	// If only one part, we'll search for it
	if len(parts) == 1 {
		query := parts[0]
		if query != "" {
			fmt.Printf("Showing album search results for '%s'\n", query)
			h.searchAlbums(w, query)
		} else {
			fmt.Printf("Showing all albums\n")
			h.listAllAlbums(w)
		}
		return
	} else if len(parts) == 2 {
		// category/album
		h.getAlbum(w, parts[0], "", parts[1])
	} else {
		// Otherwise we assume category/artist/album
		h.getAlbum(w, parts[0], parts[1], parts[2])
	}
}

func (h AlbumHandler) getAlbum(w http.ResponseWriter,
	categoryName string, artistName string, albumName string) {
	log.Printf("Looking up album '%s/%s/%s\n", categoryName,
		artistName, albumName)
	album, err := model.GetAlbum(h.Music, categoryName, artistName, albumName)

	if err == nil {
		log.Printf("Found album: %s/%s; %d tracks", album.Artist, album.Name, len(album.Tracks))
		summary := model.NewAlbumSummary(album)
		j, _ := json.Marshal(summary)
		w.Write(j)
		return
	}
	log.Printf("Did not find album: %s/%s/%s (%v)", categoryName, artistName,
		albumName, err)

	//fmt.Fprintf(w, "\n<h1>Hello</h1><div>world</div>\n")
	writeError(w, http.StatusNotFound, fmt.Sprintf("Album not found '%s'",
		albumName))
}

func (h AlbumHandler) listAllAlbums(w http.ResponseWriter) {
	// TODO: list all albums
	log.Printf("TOOD: list all albums\n")
	/*
		t0 := time.Now()
		fmt.Printf("All albums (%d)\n", len(h.Music.Albums))
		// Convert to Album Summary to save on info
		albumSummaries := make([]model.AlbumSummary, len(h.Music.Albums))
		for i := 0; i < len(h.Music.Albums); i++ {
			albumSummaries[i] = model.NewAlbumSummary(h.Music.Albums[i])
		}
		j, _ := json.Marshal(albumSummaries)
		fmt.Println("Time to marshall all albums:", time.Since(t0))
		t1 := time.Now()
		w.Write(j)
		fmt.Println("Time to send all albums:", time.Since(t1))
	*/
}

func (h AlbumHandler) searchAlbums(w http.ResponseWriter, query string) {
	matches := model.SearchAlbums(h.Music, query)
	fmt.Printf("Found %d results\n", len(matches))
	for _, item := range matches {
		fmt.Printf("%s - %s\n", item.Artist, item.Name)
	}
	j, _ := json.Marshal(matches)
	w.Write(j)
	return
}
