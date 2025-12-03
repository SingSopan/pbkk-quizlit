	package services

import (
	"crypto/sha1"
	"math"
	"sort"
	"strings"
)

// EmbeddingProvider provides text embeddings for retrieval
type EmbeddingProvider interface {
	// Embed returns a vector embedding for the given text
	Embed(text string) ([]float64, error)
}

// RAGService offers chunking, embedding, and retrieval over PDF/text content
type RAGService struct {
	embedder EmbeddingProvider
	store    *VectorStore
	// configuration
	chunkSize    int
	chunkOverlap int
}

// NewRAGService constructs a RAG service with the given embedder
func NewRAGService(embedder EmbeddingProvider) *RAGService {
	return &RAGService{
		embedder:     embedder,
		store:        NewVectorStore(),
		chunkSize:    800, // characters per chunk
		chunkOverlap: 150, // overlap between chunks
	}
}

// BuildIndex tokenizes content into chunks, embeds them, and stores in memory
func (r *RAGService) BuildIndex(docID string, content string) error {
	chunks := r.chunkText(content)
	for i, ch := range chunks {
		// compute embedding
		emb, err := r.embedder.Embed(ch)
		if err != nil {
			return err
		}
		// use stable id
		chunkID := docID + ":" + itoa(i)
		r.store.Upsert(VectorItem{ID: chunkID, Text: ch, Embedding: emb})
	}
	return nil
}

// Retrieve returns topK most similar chunks given a query
func (r *RAGService) Retrieve(query string, topK int) ([]VectorItem, error) {
	if topK <= 0 {
		topK = 5
	}
	qEmb, err := r.embedder.Embed(query)
	if err != nil {
		return nil, err
	}
	return r.store.TopK(qEmb, topK), nil
}

// chunkText splits text into overlapping chunks suitable for retrieval
func (r *RAGService) chunkText(text string) []string {
	cleaned := strings.TrimSpace(text)
	if cleaned == "" {
		return []string{}
	}
	// simple sentence split with fallbacks
	delimiters := []string{". ", "\n", "? ", "! "}
	parts := []string{cleaned}
	for _, d := range delimiters {
		var next []string
		for _, p := range parts {
			next = append(next, strings.Split(p, d)...)
		}
		parts = next
	}
	// rebuild chunks up to chunkSize
	var chunks []string
	var buf strings.Builder
	for _, s := range parts {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if buf.Len()+len(s)+1 > r.chunkSize {
			// flush current chunk
			ch := strings.TrimSpace(buf.String())
			if len(ch) > 0 {
				chunks = append(chunks, ch)
			}
			// start new buffer with overlap from previous chunk
			if r.chunkOverlap > 0 && len(ch) > r.chunkOverlap {
				overlap := ch[len(ch)-r.chunkOverlap:]
				buf.Reset()
				buf.WriteString(overlap)
				buf.WriteString(" ")
			} else {
				buf.Reset()
			}
		}
		if buf.Len() > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(s)
	}
	// final flush
	if buf.Len() > 0 {
		chunks = append(chunks, strings.TrimSpace(buf.String()))
	}
	return chunks
}

// VectorItem is a single chunk with its embedding
type VectorItem struct {
	ID        string
	Text      string
	Embedding []float64
}

// VectorStore is a minimal in-memory store for embeddings
type VectorStore struct {
	items []VectorItem
}

func NewVectorStore() *VectorStore { return &VectorStore{items: []VectorItem{}} }

func (vs *VectorStore) Upsert(item VectorItem) {
	for i := range vs.items {
		if vs.items[i].ID == item.ID {
			vs.items[i] = item
			return
		}
	}
	vs.items = append(vs.items, item)
}

// TopK returns top-k items by cosine similarity
func (vs *VectorStore) TopK(query []float64, k int) []VectorItem {
	type scored struct {
		item  VectorItem
		score float64
	}
	var arr []scored
	for _, it := range vs.items {
		arr = append(arr, scored{item: it, score: cosine(query, it.Embedding)})
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].score > arr[j].score })
	if k > len(arr) {
		k = len(arr)
	}
	out := make([]VectorItem, 0, k)
	for i := 0; i < k; i++ {
		out = append(out, arr[i].item)
	}
	return out
}

func cosine(a, b []float64) float64 {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	n := min(len(a), len(b))
	var dot, na, nb float64
	for i := 0; i < n; i++ {
		dot += a[i] * b[i]
		na += a[i] * a[i]
		nb += b[i] * b[i]
	}
	denom := math.Sqrt(na) * math.Sqrt(nb)
	if denom == 0 {
		return 0
	}
	return dot / denom
}

// Simple SHA1-based deterministic pseudo-embedding when real provider is not configured.
// This keeps the pipeline functional for development and testing, but should be
// replaced by a proper embedding model in production.
type HashEmbedding struct{}

func (h HashEmbedding) Embed(text string) ([]float64, error) {
	sum := sha1.Sum([]byte(text))
	// Expand to 32-dim vector deterministically
	vec := make([]float64, 32)
	// use bytes repeatedly
	by := sum[:]
	for i := 0; i < len(vec); i++ {
		b := by[i%len(by)]
		vec[i] = float64(int(b)) / 255.0
	}
	// normalize
	var norm float64
	for _, v := range vec {
		norm += v * v
	}
	norm = math.Sqrt(norm)
	if norm == 0 {
		norm = 1
	}
	for i := range vec {
		vec[i] = vec[i] / norm
	}
	return vec, nil
}

// helper: integer to string without fmt allocation
func itoa(i int) string {
	// minimal fast itoa for non-negative integers
	if i == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + (i % 10))
		i /= 10
	}
	return string(buf[pos:])
}
