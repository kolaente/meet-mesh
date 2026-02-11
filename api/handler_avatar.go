package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
	"gorm.io/gorm"
)

const (
	maxAvatarSize = 2 << 20 // 2MB
	avatarMaxDim  = 256     // Max width/height in pixels
	avatarQuality = 85      // JPEG quality
)

// AvatarHandler handles avatar upload, delete, and serving.
// These are plain HTTP handlers (not ogen) because ogen does not support multipart uploads.
type AvatarHandler struct {
	db          *gorm.DB
	auth        *AuthService
	avatarsPath string
}

func NewAvatarHandler(db *gorm.DB, auth *AuthService, avatarsPath string) *AvatarHandler {
	return &AvatarHandler{
		db:          db,
		auth:        auth,
		avatarsPath: avatarsPath,
	}
}

// Upload handles POST /api/avatars
func (h *AvatarHandler) Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Authenticate the user
	userID, ok := h.authenticateRequest(r)
	if !ok {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// Limit request body size
	r.Body = http.MaxBytesReader(w, r.Body, maxAvatarSize)

	// Parse multipart form
	if err := r.ParseMultipartForm(maxAvatarSize); err != nil {
		http.Error(w, "File too large (max 2MB)", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, "Missing avatar file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate content type
	contentType := header.Header.Get("Content-Type")
	if !isAllowedImageType(contentType) {
		http.Error(w, "Invalid file type. Allowed: JPEG, PNG, WebP", http.StatusBadRequest)
		return
	}

	// Read the file into memory for processing
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Decode the image
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		http.Error(w, "Invalid image file", http.StatusBadRequest)
		return
	}

	// Resize if needed
	resized := resizeImage(img, avatarMaxDim)

	// Encode as JPEG
	jpegData, err := encodeJPEG(resized, avatarQuality)
	if err != nil {
		http.Error(w, "Failed to process image", http.StatusInternalServerError)
		return
	}

	// Generate filename from content hash for cache busting
	hash := sha256.Sum256(jpegData)
	filename := hex.EncodeToString(hash[:8]) + ".jpg"

	// Save to disk
	destPath := filepath.Join(h.avatarsPath, filename)
	if err := os.WriteFile(destPath, jpegData, 0644); err != nil {
		http.Error(w, "Failed to save avatar", http.StatusInternalServerError)
		return
	}

	// Get current user to find old avatar
	var user User
	if err := h.db.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	oldFilename := user.AvatarFilename

	// Update user record
	if err := h.db.Model(&user).Update("avatar_filename", filename).Error; err != nil {
		_ = os.Remove(destPath)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	// Clean up old avatar file if different
	if oldFilename != "" && oldFilename != filename {
		_ = os.Remove(filepath.Join(h.avatarsPath, oldFilename))
	}

	// Return the avatar URL as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"avatar_url":"/api/avatars/%s"}`, filename)
}

// Delete handles DELETE /api/avatars
func (h *AvatarHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Authenticate the user
	userID, ok := h.authenticateRequest(r)
	if !ok {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// Get current user
	var user User
	if err := h.db.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// Nothing to delete
	if user.AvatarFilename == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"avatar_url":""}`)
		return
	}

	oldFilename := user.AvatarFilename

	// Clear avatar in database
	if err := h.db.Model(&user).Update("avatar_filename", "").Error; err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	// Remove file from disk (best-effort, don't fail if file is already gone)
	_ = os.Remove(filepath.Join(h.avatarsPath, oldFilename))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"avatar_url":""}`)
}

// authenticateRequest validates the session cookie and returns the user ID.
func (h *AvatarHandler) authenticateRequest(r *http.Request) (uint, bool) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return 0, false
	}

	session, err := h.auth.ParseSessionCookie(cookie)
	if err != nil {
		return 0, false
	}

	return session.UserID, true
}

func isAllowedImageType(contentType string) bool {
	switch contentType {
	case "image/jpeg", "image/png", "image/webp":
		return true
	}
	return false
}

// resizeImage resizes the image so its longest side is at most maxDim pixels.
func resizeImage(img image.Image, maxDim int) image.Image {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	if w <= maxDim && h <= maxDim {
		return img
	}

	var newW, newH int
	if w > h {
		newW = maxDim
		newH = h * maxDim / w
	} else {
		newH = maxDim
		newW = w * maxDim / h
	}

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)
	return dst
}

// encodeJPEG encodes an image as JPEG bytes.
func encodeJPEG(img image.Image, quality int) ([]byte, error) {
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Serve handles GET /api/avatars/{filename}
// This endpoint is public (no auth) since avatars appear on public pages.
func (h *AvatarHandler) Serve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract filename from path: /api/avatars/{filename}
	filename := filepath.Base(r.URL.Path)

	// Sanitize: only allow valid filenames
	if filename == "" || filename == "." || filename == ".." {
		http.NotFound(w, r)
		return
	}

	// Ensure filename doesn't contain path separators
	if filepath.Dir(filename) != "." {
		http.NotFound(w, r)
		return
	}

	filePath := filepath.Join(h.avatarsPath, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	// Set cache headers - filenames are content-hashed so we can cache forever
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.Header().Set("Content-Type", "image/jpeg")

	http.ServeFile(w, r, filePath)
}

// HandleAvatars is the main handler for /api/avatars/ that dispatches by method.
func (h *AvatarHandler) HandleAvatars(w http.ResponseWriter, r *http.Request) {
	// /api/avatars/ with no filename = upload or delete
	// /api/avatars/{filename} = serve
	filename := strings.TrimPrefix(r.URL.Path, "/api/avatars/")
	filename = strings.TrimPrefix(filename, "/api/avatars")

	switch {
	case r.Method == http.MethodPost:
		h.Upload(w, r)
	case r.Method == http.MethodDelete:
		h.Delete(w, r)
	case r.Method == http.MethodGet && filename != "" && filename != "/":
		h.Serve(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
