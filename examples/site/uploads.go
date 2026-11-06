package site

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Tend-Stack/GothUI/v2026/components"
)

const (
	uploadFileLimit   = 64 << 10
	uploadBodyLimit   = 72 << 10
	uploadRowLimit    = 50
	sampleProjectsCSV = "project,owner,status\nStudio website,Alex,active\nBrand guidelines,Sam,complete\nCommunity portal,Jordan,planned\n"
)

type uploadData struct {
	Rows  []components.TableRow
	Error string
}

func uploadDemo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	d := uploadData{}
	status := http.StatusOK
	switch r.Method {
	case http.MethodGet, http.MethodHead:
	case http.MethodPost:
		var data []byte
		data, d.Error, status = readProjectUpload(w, r)
		if d.Error == "" {
			d.Rows, d.Error = parseProjectCSV(data)
			if d.Error != "" {
				status = http.StatusUnprocessableEntity
			}
		}
	default:
		w.Header().Set("Allow", "GET, HEAD, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	renderPage(w, r, "CSV preview · GothUI", uploadPage(d), status)
}

// Bound the entire request before multipart parsing, including headers, extra parts
// and trailing bytes. The small payload stays in memory; no temporary files or
// client-provided paths are ever opened. Browser MIME/extension hints are not trusted.
func readProjectUpload(w http.ResponseWriter, r *http.Request) ([]byte, string, int) {
	r.Body = http.MaxBytesReader(w, r.Body, uploadBodyLimit)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			return nil, "The upload request exceeds 72 KiB. Choose a CSV file of 64 KiB or less.", http.StatusRequestEntityTooLarge
		}
		return nil, "The upload could not be read. Choose your file and try again.", http.StatusBadRequest
	}
	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "multipart/form-data" || params["boundary"] == "" {
		return nil, "Use the file chooser to submit one CSV file.", http.StatusBadRequest
	}
	reader := multipart.NewReader(bytes.NewReader(body), params["boundary"])
	part, err := reader.NextRawPart()
	if err != nil {
		return nil, "Choose one CSV file to preview.", http.StatusBadRequest
	}
	disposition, fields, err := mime.ParseMediaType(part.Header.Get("Content-Disposition"))
	if err != nil || len(part.Header.Values("Content-Disposition")) != 1 || disposition != "form-data" || fields["name"] != "projects" || fields["filename"] == "" || len(part.Header.Values("Content-Transfer-Encoding")) != 0 {
		return nil, "Submit one unencoded file using the project CSV field.", http.StatusBadRequest
	}
	data, err := io.ReadAll(io.LimitReader(part, uploadFileLimit+1))
	if len(data) > uploadFileLimit {
		return nil, "The CSV file exceeds 64 KiB. Choose a smaller file.", http.StatusRequestEntityTooLarge
	}
	if err != nil {
		return nil, "The uploaded file is incomplete. Choose your file and try again.", http.StatusBadRequest
	}
	if _, err := reader.NextRawPart(); err != io.EOF {
		return nil, "Submit exactly one CSV file with no additional form fields.", http.StatusBadRequest
	}
	return data, "", http.StatusOK
}

func parseProjectCSV(data []byte) ([]components.TableRow, string) {
	if !utf8.Valid(data) {
		return nil, "Save the CSV as UTF-8 text and try again."
	}
	data = bytes.TrimPrefix(data, []byte{0xef, 0xbb, 0xbf})
	reader := csv.NewReader(bytes.NewReader(data))
	reader.FieldsPerRecord = 3
	header, err := reader.Read()
	if err != nil || header[0] != "project" || header[1] != "owner" || header[2] != "status" {
		return nil, "Use the exact header project,owner,status. Download the sample to get started."
	}
	var rows []components.TableRow
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			var csvError *csv.ParseError
			if errors.As(err, &csvError) {
				return nil, fmt.Sprintf("Line %d: use three comma-separated fields and properly quoted values.", csvError.Line)
			}
			return nil, "The CSV could not be parsed."
		}
		if len(rows) == uploadRowLimit {
			return nil, "Preview up to 50 project rows at a time."
		}
		line, _ := reader.FieldPos(0)
		for i, value := range record {
			// Disallow embedded controls before trimming; edge spaces are fine.
			if strings.IndexFunc(value, unicode.IsControl) >= 0 {
				return nil, fmt.Sprintf("Line %d: fields must be single-line text without control characters.", line)
			}
			record[i] = strings.TrimSpace(value)
			if n := utf8.RuneCountInString(record[i]); n == 0 || n > 80 {
				return nil, fmt.Sprintf("Line %d: each field needs between 1 and 80 characters.", line)
			}
		}
		if record[2] != "planned" && record[2] != "active" && record[2] != "complete" {
			return nil, fmt.Sprintf("Line %d: status must be planned, active or complete.", line)
		}
		rows = append(rows, components.TableRow{Cells: []components.TableCell{{Text: record[0]}, {Text: record[1]}, {Text: record[2]}}})
	}
	if len(rows) == 0 {
		return nil, "Add at least one project row after the header."
	}
	return rows, ""
}

func downloadProjectSample(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="projects.csv"`)
	if r.Method != http.MethodHead {
		_, _ = io.WriteString(w, sampleProjectsCSV)
	}
}
