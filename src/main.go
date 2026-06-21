package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	texttemplate "text/template"
)

const baseURL = "https://fantasticbeastsandwheretofindthem.xyz"

var slugRe = regexp.MustCompile(`[^a-z0-9]+`)

type Beast struct {
	Name           string
	Classification int
	Description    string
	Slug           string
	ClassLabel     string
	ClassCSS       string
}

type PageData struct {
	BaseURL       string
	ActiveNav     string
	Search        string
	SearchParam   string
	Sort          string
}

type IndexData struct {
	PageData
	Beasts []Beast
}

type BeastDetailData struct {
	PageData
	Beast Beast
}

func main() {
	beasts, err := loadBeasts(beastsJSONPath())
	if err != nil {
		log.Fatalf("loading beasts: %v", err)
	}

	beastMap := make(map[string]Beast, len(beasts))
	beastByName := make(map[string]Beast, len(beasts))
	for _, b := range beasts {
		beastMap[b.Slug] = b
		beastByName[strings.ToLower(b.Name)] = b
	}

	htmlFuncs := template.FuncMap{
		"raw":       func(s string) template.HTML { return template.HTML(s) },
		"truncate":  truncateText,
	}
	partials := []string{
		"templates/partials/layout_styles.html.tmpl",
		"templates/partials/header.html.tmpl",
		"templates/partials/footer.html.tmpl",
	}

	tmplRobots, err := template.ParseFiles("templates/robots.txt.tmpl")
	if err != nil {
		log.Fatalf("parsing robots template: %v", err)
	}
	tmplSitemap, err := texttemplate.ParseFiles("templates/sitemap.xml.tmpl")
	if err != nil {
		log.Fatalf("parsing sitemap template: %v", err)
	}

	tmplIndex, err := template.New("index.html.tmpl").Funcs(htmlFuncs).ParseFiles(append([]string{"templates/index.html.tmpl"}, partials...)...)
	if err != nil {
		log.Fatalf("parsing index template: %v", err)
	}
	tmplBeast, err := template.New("beast.html.tmpl").Funcs(htmlFuncs).ParseFiles(append([]string{"templates/beast.html.tmpl"}, partials...)...)
	if err != nil {
		log.Fatalf("parsing beast template: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		search := strings.TrimSpace(r.URL.Query().Get("search"))
		sortKey := r.URL.Query().Get("sort")
		if sortKey == "" {
			sortKey = "asc"
		}

		items := beasts
		if search != "" {
			items = semanticSearchBeasts(beasts, search, 100)
		}
		items = sortBeasts(items, sortKey)

		searchParam := ""
		if search != "" {
			searchParam = "&search=" + url.QueryEscape(search)
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmplIndex.Execute(w, IndexData{
			PageData: PageData{BaseURL: baseURL, ActiveNav: "home", Search: search, SearchParam: searchParam, Sort: sortKey},
			Beasts:   items,
		}); err != nil {
			log.Printf("executing index template: %v", err)
			http.Error(w, "Internal Server Error", 500)
		}
	})

	http.HandleFunc("/beasts/", func(w http.ResponseWriter, r *http.Request) {
		slug := strings.TrimPrefix(r.URL.Path, "/beasts/")
		if slug == "" {
			http.Redirect(w, r, "/", 302)
			return
		}
		beast, ok := beastMap[slug]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmplBeast.Execute(w, BeastDetailData{
			PageData: PageData{BaseURL: baseURL, ActiveNav: "home", Sort: "asc"},
			Beast:    beast,
		}); err != nil {
			log.Printf("executing beast template: %v", err)
			http.Error(w, "Internal Server Error", 500)
		}
	})

	http.HandleFunc("/beast-details.php", func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimSpace(r.URL.Query().Get("name"))
		if name == "" {
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
			return
		}
		beast, ok := beastByName[strings.ToLower(name)]
		if !ok {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/beasts/"+beast.Slug, http.StatusMovedPermanently)
	})

	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if err := tmplRobots.Execute(w, struct{ BaseURL string }{BaseURL: baseURL}); err != nil {
			log.Printf("executing robots template: %v", err)
			http.Error(w, "Internal Server Error", 500)
		}
	})

	http.HandleFunc("/sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		data := struct {
			BaseURL string
			Beasts  []Beast
		}{BaseURL: baseURL, Beasts: beasts}
		if err := tmplSitemap.Execute(&buf, data); err != nil {
			log.Printf("executing sitemap template: %v", err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.Write(buf.Bytes())
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Serving at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func beastsJSONPath() string {
	if p := os.Getenv("BEASTS_JSON"); p != "" {
		return p
	}
	for _, candidate := range []string{"assets/beasts.json", "../assets/beasts.json"} {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return "assets/beasts.json"
}

func loadBeasts(path string) ([]Beast, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw []struct {
		Name           string `json:"name"`
		Classification int    `json:"classification"`
		Description    string `json:"description"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	beasts := make([]Beast, 0, len(raw))
	for _, item := range raw {
		beasts = append(beasts, enrichBeast(item.Name, item.Classification, item.Description))
	}

	sort.Slice(beasts, func(i, j int) bool {
		return strings.ToLower(beasts[i].Name) < strings.ToLower(beasts[j].Name)
	})

	return beasts, nil
}

func enrichBeast(name string, classification int, description string) Beast {
	label := strings.Repeat("X", classification)
	return Beast{
		Name:           name,
		Classification: classification,
		Description:    description,
		Slug:           beastSlug(name),
		ClassLabel:     label,
		ClassCSS:       "classification-" + strings.ToLower(label),
	}
}

func truncateText(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max]) + "..."
}

func beastSlug(name string) string {
	s := strings.ToLower(name)
	s = slugRe.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

func sortBeasts(beasts []Beast, sortKey string) []Beast {
	sorted := make([]Beast, len(beasts))
	copy(sorted, beasts)

	switch sortKey {
	case "desc":
		sort.Slice(sorted, func(i, j int) bool {
			return strings.ToLower(sorted[i].Name) > strings.ToLower(sorted[j].Name)
		})
	case "1to5":
		sort.Slice(sorted, func(i, j int) bool {
			if sorted[i].Classification == sorted[j].Classification {
				return strings.ToLower(sorted[i].Name) < strings.ToLower(sorted[j].Name)
			}
			return sorted[i].Classification < sorted[j].Classification
		})
	case "5to1":
		sort.Slice(sorted, func(i, j int) bool {
			if sorted[i].Classification == sorted[j].Classification {
				return strings.ToLower(sorted[i].Name) < strings.ToLower(sorted[j].Name)
			}
			return sorted[i].Classification > sorted[j].Classification
		})
	default:
		sort.Slice(sorted, func(i, j int) bool {
			return strings.ToLower(sorted[i].Name) < strings.ToLower(sorted[j].Name)
		})
	}

	return sorted
}
