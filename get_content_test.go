package main

import (
	"net/url"
	"reflect"
	"testing"
)

// Tests for getH1FromHTML
func TestGetH1FromHTMLBasic(t *testing.T) {
	inputBody := "<html><body><h1>Test Title</h1></body></html>"
	actual := getH1FromHTML(inputBody)
	expected := "Test Title"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestNoH1InHTML(t *testing.T) {
	inputBody := "<html><body><p>No H1 here</p></body></html>"
	actual := getH1FromHTML(inputBody)
	expected := ""

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestMultipleH1InHTML(t *testing.T) {
	inputBody := "<html><body><h1>First Title</h1><h1>Second Title</h1></body></html>"
	actual := getH1FromHTML(inputBody)
	expected := "First Title"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

// Tests for getFirstParagraphFromHTML
func TestGetFirstParagraphFromHTMLMainPriority(t *testing.T) {
	inputBody := `<html><body>
		<p>Outside paragraph.</p>
		<main>
			<p>Main paragraph.</p>
		</main>
	</body></html>`
	actual := getFirstParagraphFromHTML(inputBody)
	expected := "Main paragraph."

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestNoParagraphinHTML(t *testing.T) {
	inputBody := "<html><body><h1>No paragraphs here</h1></body></html>"
	actual := getFirstParagraphFromHTML(inputBody)
	expected := ""

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetFirstParagraphFromHTMLNoMain(t *testing.T) {
	inputBody := `<html><body>
		<h1>Title</h1>
		<p>First paragraph.</p>
		<p>Second paragraph.</p>
	</body></html>`
	actual := getFirstParagraphFromHTML(inputBody)
	expected := "First paragraph."

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}
func TestGetFirstParagraphFromHTMLMainEmpty(t *testing.T) {
	inputBody := `<html><body>
		<main></main>
		<p>First paragraph outside main.</p>
	</body></html>`
	actual := getFirstParagraphFromHTML(inputBody)
	expected := "First paragraph outside main."

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

// Tests for getURLsFromHTML
func TestGetURLsFromHTMLAbsolute(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	inputBody := `<html><body><a href="https://blog.boot.dev"><span>Boot.dev</span></a></body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getURLsFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://blog.boot.dev"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetURLsFromHTMLRelative(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	inputBody := `<html><body><a href="/path/to/page"><span>Boot.dev Page</span></a></body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getURLsFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://blog.boot.dev/path/to/page"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetURLsFromHTMLMultiple(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	inputBody := `<html><body>
		<a href="https://blog.boot.dev/page1">Page 1</a>
		<a href="/page2">Page 2</a>
	</body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getURLsFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{
		"https://blog.boot.dev/page1",
		"https://blog.boot.dev/page2",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

// Tests for getImagesFromHTML
func TestGetImagesFromHTMLRelative(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	inputBody := `<html><body><img src="/logo.png" alt="Logo"></body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getImagesFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://blog.boot.dev/logo.png"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetImagesFromHTMLAbsolute(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	inputBody := `<html><body><img src="https://blog.boot.dev/logo.png" alt="Logo"></body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getImagesFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://blog.boot.dev/logo.png"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetImagesFromHTMLMultiple(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	inputBody := `<html><body>
		<img src="https://blog.boot.dev/image1.png" alt="Image 1">
		<img src="/image2.png" alt="Image 2">
	</body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getImagesFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{
		"https://blog.boot.dev/image1.png",
		"https://blog.boot.dev/image2.png",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetImagesFromHTMLMissingSrc(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	inputBody := `<html><body>
		<img alt="no src here">
		<img src="" alt="empty src">
		<img src="/logo.png" alt="Logo">
	</body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getImagesFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://blog.boot.dev/logo.png"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

// Tests for extractPageData
func TestExtractPageData(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	inputBody := `<html><body>
        <h1>Test Title</h1>
        <p>This is the first paragraph.</p>
        <a href="/link1">Link 1</a>
        <img src="/image1.jpg" alt="Image 1">
    </body></html>`

	actual := extractPageData(inputBody, inputURL)

	expected := PageData{
		URL:            "https://blog.boot.dev",
		H1:             "Test Title",
		FirstParagraph: "This is the first paragraph.",
		OutgoingLinks:  []string{"https://blog.boot.dev/link1"},
		ImageURLs:      []string{"https://blog.boot.dev/image1.jpg"},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func TestExtractPageData_NoH1(t *testing.T) {
	inputURL := "https://example.com"
	inputBody := `<html><body>
        <p>First para</p>
    </body></html>`

	actual := extractPageData(inputBody, inputURL)

	expected := PageData{
		URL:            "https://example.com",
		H1:             "",
		FirstParagraph: "First para",
		OutgoingLinks:  nil,
		ImageURLs:      nil,
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %#v, got %#v", expected, actual)
	}

	if actual.H1 != expected.H1 {
		t.Errorf("H1 mismatch: expected %q, got %q", expected.H1, actual.H1)
	}
	if actual.FirstParagraph != expected.FirstParagraph {
		t.Errorf("FirstParagraph mismatch: expected %q, got %q", expected.FirstParagraph, actual.FirstParagraph)
	}
}

func TestExtractPageData_NoParagraph(t *testing.T) {
	inputURL := "https://example.com"
	inputBody := `<html><body>
        <h1>Title</h1>
    </body></html>`

	actual := extractPageData(inputBody, inputURL)

	expected := PageData{
		URL:            "https://example.com",
		H1:             "Title",
		FirstParagraph: "",
		OutgoingLinks:  nil,
		ImageURLs:      nil,
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %#v, got %#v", expected, actual)
	}
}

func TestExtractPageData_MultipleLinksAndImages(t *testing.T) {
	inputURL := "https://example.com"
	inputBody := `<html><body>
        <h1>Title</h1>
        <p>First para</p>
        <a href="/rel">Rel</a>
        <a href="https://other.com/">Abs</a>
        <img src="/pic.png">
        <img src="https://cdn.com/img.png">
    </body></html>`

	actual := extractPageData(inputBody, inputURL)

	expected := PageData{
		URL:            "https://example.com",
		H1:             "Title",
		FirstParagraph: "First para",
		OutgoingLinks: []string{
			"https://example.com/rel",
			"https://other.com/",
		},
		ImageURLs: []string{
			"https://example.com/pic.png",
			"https://cdn.com/img.png",
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func TestExtractPageData_EmptyHTML(t *testing.T) {
	inputURL := "https://example.com"
	inputBody := `<html><body></body></html>`

	actual := extractPageData(inputBody, inputURL)

	expected := PageData{
		URL:            "https://example.com",
		H1:             "",
		FirstParagraph: "",
		OutgoingLinks:  nil,
		ImageURLs:      nil,
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %#v, got %#v", expected, actual)
	}
}

func TestExtractPageData_URLField(t *testing.T) {
	inputURL := "https://another.example/path"
	inputBody := `<html><body>
        <h1>Title</h1>
        <p>Para</p>
    </body></html>`

	actual := extractPageData(inputBody, inputURL)

	if actual.URL != inputURL {
		t.Errorf("expected URL %q, got %q", inputURL, actual.URL)
	}
}
