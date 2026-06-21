package main

import (
	"math"
	"regexp"
	"sort"
	"strings"
)

var tokenRe = regexp.MustCompile(`[^\p{L}\p{N}\s]+`)

var stopwords = map[string]struct{}{
	"the": {}, "a": {}, "an": {}, "and": {}, "or": {}, "but": {}, "if": {}, "then": {}, "else": {},
	"for": {}, "to": {}, "of": {}, "in": {}, "on": {}, "with": {}, "at": {}, "by": {}, "from": {},
	"as": {}, "is": {}, "are": {}, "was": {}, "were": {}, "it": {}, "its": {}, "be": {}, "this": {},
	"that": {}, "these": {}, "those": {}, "their": {}, "his": {}, "her": {}, "they": {}, "them": {},
	"we": {}, "you": {}, "your": {}, "our": {}, "i": {}, "my": {}, "me": {},
}

func normalizeText(text string) string {
	text = strings.ToLower(text)
	text = tokenRe.ReplaceAllString(text, " ")
	text = strings.Join(strings.Fields(text), " ")
	return strings.TrimSpace(text)
}

func singularizeBasic(token string) string {
	switch {
	case strings.HasSuffix(token, "ches"), strings.HasSuffix(token, "shes"),
		strings.HasSuffix(token, "xes"), strings.HasSuffix(token, "ses"):
		return strings.TrimSuffix(token, "es")
	case strings.HasSuffix(token, "ies"):
		return strings.TrimSuffix(token, "ies") + "y"
	case strings.HasSuffix(token, "s") && !strings.HasSuffix(token, "ss"):
		return strings.TrimSuffix(token, "s")
	default:
		return token
	}
}

func tokenize(text string) []string {
	normalized := normalizeText(text)
	if normalized == "" {
		return nil
	}
	tokens := strings.Fields(normalized)
	filtered := make([]string, 0, len(tokens))
	for _, t := range tokens {
		if t == "" {
			continue
		}
		if _, ok := stopwords[t]; ok {
			continue
		}
		filtered = append(filtered, singularizeBasic(t))
	}
	return filtered
}

func fieldWeightedTokens(beast Beast) []string {
	nameTokens := tokenize(beast.Name)
	descTokens := tokenize(beast.Description)

	tokens := make([]string, 0, len(nameTokens)*6+len(descTokens))
	for _, t := range nameTokens {
		for i := 0; i < 6; i++ {
			tokens = append(tokens, t)
		}
	}
	tokens = append(tokens, descTokens...)
	return tokens
}

type searchDocument struct {
	beast  Beast
	tokens []string
}

func buildCorpusDocuments(beasts []Beast) []searchDocument {
	docs := make([]searchDocument, 0, len(beasts))
	for _, beast := range beasts {
		docs = append(docs, searchDocument{
			beast:  beast,
			tokens: fieldWeightedTokens(beast),
		})
	}
	return docs
}

func computeIDF(documents []searchDocument) map[string]float64 {
	numDocuments := max(1, len(documents))
	documentFrequency := map[string]int{}
	for _, doc := range documents {
		seen := map[string]struct{}{}
		for _, token := range uniqueStrings(doc.tokens) {
			seen[token] = struct{}{}
		}
		for token := range seen {
			documentFrequency[token]++
		}
	}
	idf := make(map[string]float64, len(documentFrequency))
	for token, df := range documentFrequency {
		idf[token] = math.Log(float64(numDocuments+1)/float64(df+1)) + 1.0
	}
	return idf
}

func computeTF(tokens []string) map[string]float64 {
	tf := map[string]float64{}
	for _, token := range tokens {
		tf[token]++
	}
	length := max(1, len(tokens))
	for token, count := range tf {
		tf[token] = count / float64(length)
	}
	return tf
}

func computeVector(tf, idf map[string]float64) map[string]float64 {
	vector := map[string]float64{}
	for token, tfv := range tf {
		idfv := idf[token]
		weight := tfv * idfv
		if weight != 0 {
			vector[token] = weight
		}
	}
	return vector
}

func vectorNorm(vector map[string]float64) float64 {
	sum := 0.0
	for _, w := range vector {
		sum += w * w
	}
	return math.Sqrt(sum)
}

func cosineSimilarity(a, b map[string]float64) float64 {
	lenA := vectorNorm(a)
	lenB := vectorNorm(b)
	if lenA == 0 || lenB == 0 {
		return 0
	}

	small, large := a, b
	if len(a) >= len(b) {
		small, large = b, a
	}

	dot := 0.0
	for token, wa := range small {
		if wb, ok := large[token]; ok {
			dot += wa * wb
		}
	}
	return dot / (lenA * lenB)
}

func levenshtein(a, b string) int {
	if a == b {
		return 0
	}
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}

	prev := make([]int, len(b)+1)
	curr := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}

	for i := 1; i <= len(a); i++ {
		curr[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			curr[j] = min(
				curr[j-1]+1,
				prev[j]+1,
				prev[j-1]+cost,
			)
		}
		prev, curr = curr, prev
	}
	return prev[len(b)]
}

func nameBonus(query, name string) float64 {
	q := normalizeText(query)
	n := normalizeText(name)
	if q == "" || n == "" {
		return 0
	}

	score := 0.0
	if strings.Contains(n, q) {
		score += 0.5
	}

	lev := levenshtein(q, n)
	maxLen := max(1, max(len(q), len(n)))
	sim := 1.0 - min(1.0, float64(lev)/float64(maxLen))
	score += 0.35 * sim

	for _, qt := range tokenize(query) {
		if qt != "" && strings.HasPrefix(n, qt) {
			score += 0.25
		}
	}
	return score
}

func termMatchBonus(query string, beast Beast) float64 {
	qTokens := tokenize(query)
	desc := normalizeText(beast.Description)
	name := normalizeText(beast.Name)

	score := 0.0
	for _, qt := range qTokens {
		if qt != "" && strings.Contains(name, qt) {
			score += 1.0
		}
	}
	for _, qt := range qTokens {
		if qt != "" && strings.Contains(desc, qt) {
			score += 0.3
		}
	}
	return score
}

func semanticSearchBeasts(beasts []Beast, query string, limit int) []Beast {
	query = strings.TrimSpace(query)
	if query == "" {
		return beasts
	}

	documents := buildCorpusDocuments(beasts)
	idf := computeIDF(documents)

	queryTokens := tokenize(query)
	expandedQueryTokens := make([]string, 0, len(queryTokens)*2)
	for _, qt := range queryTokens {
		expandedQueryTokens = append(expandedQueryTokens, qt, qt)
	}
	queryVector := computeVector(computeTF(expandedQueryTokens), idf)

	type scoredBeast struct {
		beast Beast
		score float64
	}
	scored := make([]scoredBeast, 0, len(documents))
	for _, doc := range documents {
		docVector := computeVector(computeTF(doc.tokens), idf)
		score := cosineSimilarity(queryVector, docVector)
		score += nameBonus(query, doc.beast.Name)
		score += termMatchBonus(query, doc.beast)
		if score > 0 {
			scored = append(scored, scoredBeast{beast: doc.beast, score: score})
		}
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	result := make([]Beast, 0, limit)
	for _, row := range scored {
		result = append(result, row.beast)
		if len(result) >= limit {
			break
		}
	}

	if len(result) == 0 {
		q := normalizeText(query)
		for _, beast := range beasts {
			if strings.Contains(normalizeText(beast.Name), q) {
				result = append(result, beast)
			}
		}
	}

	return result
}

func uniqueStrings(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, v := range values {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

