package e621

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/downloadablefox/twotto/core"
	"github.com/rs/zerolog/log"
)

var (
	E621ServiceKey         = core.NewIdentifier("e621", "services/client")
	ErrE621ServiceNotFound = errors.New("e621 service not found")
)

const MAX_POST_SIZE = 25 * 1024 * 1024

type IE621Service interface {
	GetRandomPost() (*E621Post, error)
	GetPostByID(id int) (*E621Post, error)
	SearchPosts(tags string, limit, page int) ([]*E621Post, error)
	GetPopularPosts() ([]*E621Post, error)
}

type E621Service struct {
	httpClient *http.Client
	userAgent  string
}

func NewE621Service(userAgent string) *E621Service {
	return &E621Service{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		userAgent: userAgent,
	}
}

func (e *E621Service) GetRandomPost() (*E621Post, error) {
	const url = "https://e621.net/posts/random.json"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", e.userAgent)

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var post struct {
		Post E621PostResponse `json:"post"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&post); err != nil {
		return nil, err
	}

	return e.parsePost(&post.Post)
}

func (e *E621Service) GetPostByID(id int) (*E621Post, error) {
	url := fmt.Sprintf("https://e621.net/posts/%d.json", id)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", e.userAgent)

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var post struct {
		Post E621PostResponse `json:"post"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&post); err != nil {
		return nil, err
	}

	return e.parsePost(&post.Post)
}

func (e *E621Service) SearchPosts(tags string, limit, page int) ([]*E621Post, error) {
	// URL encode the query
	tags = url.QueryEscape(tags)

	// Append the limit and page
	url := "https://e621.net/posts.json?tags=%s&limit=%d&page=%d"
	url = fmt.Sprintf(url, tags, limit, page)

	log.Debug().Msgf("[E621Service] Search URL: %s", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", e.userAgent)

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var posts struct {
		Posts []*E621PostResponse `json:"posts"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		return nil, err
	}

	var result []*E621Post
	for _, post := range posts.Posts {
		parsed, err := e.parsePost(post)
		if err != nil {
			continue
		}

		result = append(result, parsed)
	}

	return result, nil
}

func (e *E621Service) parsePost(post *E621PostResponse) (*E621Post, error) {
	if post.ID == 0 {
		return nil, errors.New("post was not found")
	}

	useSample := post.File.Size > MAX_POST_SIZE

	// If the post is too large and we can use the sample, use it
	if useSample && !post.Sample.Has {
		return nil, errors.New("post is too large and has no samples")
	}

	isVideo := post.File.Ext == "webm" || post.File.Ext == "mp4"

	url := post.File.URL
	if useSample {
		if isVideo {
			url = ""

			// Attempt to get the best quality video
			for altname, alt := range post.Sample.Alts {
				for _, current := range alt.URLs {
					if current == nil {
						continue
					}

					// Download the video and check the size
					// If it's smaller than the max size, use it
					req, err := http.NewRequest(http.MethodGet, *current, nil)
					if err != nil {
						return nil, err
					}

					req.Header.Set("User-Agent", e.userAgent)

					resp, err := e.httpClient.Do(req)
					if err != nil {
						return nil, err
					}

					// Read all contents of resp.Body to ensure it's not too large
					byteCount := 0
					buf := make([]byte, 1024)
					for {
						n, err := resp.Body.Read(buf)
						byteCount += n
						if err != nil {
							break
						}
					}

					log.Debug().Msgf("[E621Service] Video size for alt (%s): %.2f MB", altname, float64(byteCount)/1024.0/1024.0)

					if byteCount < MAX_POST_SIZE {
						resp.Body.Close()
						url = *current
						break
					}

					resp.Body.Close()
				}

				if url != "" {
					break
				}
			}

			if url == "" {
				return nil, fmt.Errorf("no suitable video found")
			}
		} else {
			// Check if the sample is too large
			req, err := http.NewRequest(http.MethodGet, post.Sample.URL, nil)
			if err != nil {
				return nil, err
			}

			req.Header.Set("User-Agent", e.userAgent)

			resp, err := e.httpClient.Do(req)
			if err != nil {
				return nil, err
			}

			// Read all contents of resp.Body to ensure it's not too large
			byteCount := 0
			buf := make([]byte, 1024)
			for {
				n, err := resp.Body.Read(buf)
				byteCount += n
				if err != nil {
					break
				}
			}
			defer resp.Body.Close()

			if byteCount > MAX_POST_SIZE {
				return nil, errors.New("no suitable image found")
			}

			url = post.Sample.URL
		}
	}

	return &E621Post{
		ID:   post.ID,
		URL:  url,
		Ext:  post.File.Ext,
		Size: post.File.Size,
	}, nil
}

func (e *E621Service) GetPopularPosts() ([]*E621Post, error) {
	return nil, errors.New("not implemented")
}
