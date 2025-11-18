package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/pin-app/pin/internal/server"
)

const maxUploadSize = 10 << 20 // 10MB

type UploadHandler struct {
	uploadDir string
}

func NewUploadHandler(uploadDir string) *UploadHandler {
	if uploadDir == "" {
		uploadDir = "uploads"
	}

	if !filepath.IsAbs(uploadDir) {
		if abs, err := filepath.Abs(uploadDir); err == nil {
			uploadDir = abs
		}
	}

	_ = os.MkdirAll(uploadDir, 0o755)

	return &UploadHandler{uploadDir: uploadDir}
}

func (h *UploadHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "failed to parse upload"})
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "file field is required"})
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if ext == "" {
		ext = ".jpg"
	}

	filename := fmt.Sprintf("%s%s", uuid.NewString(), ext)
	destPath := filepath.Join(h.uploadDir, filename)

	dst, err := os.Create(destPath)
	if err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to store file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to save file"})
		return
	}

	server.WriteJSON(w, http.StatusCreated, map[string]string{
		"url": fmt.Sprintf("/uploads/%s", filename),
	})
}
