package business

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"irj/internal/catalogs"
	"irj/internal/jwt"
	queries "irj/internal/postgres/_generated"
	"irj/pkg/api"
	_http "irj/pkg/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/julienschmidt/httprouter"
)

const MaxFileSize = 5 << 20 // 5 MB

type (
	NocoMedia struct {
		Title string `json:"titre"`
		ID    int32  `json:"id"`
	}

	MediaPath struct {
		Path     string `json:"path"`
		Title    string `json:"title"`
		Mimetype string `json:"mimetype"`
		Size     int    `json:"size"`
		ID       string `json:"id"`
	}
)

func parseMediaPath(path string) (MediaPath, error) {
	var media []MediaPath

	err := json.Unmarshal([]byte(path), &media)
	if err != nil {
		return MediaPath{}, err
	}

	return media[0], err
}

func (b *BusinessService) parseMedias(rawMedias interface{}) ([]NocoMedia, error) {
	data, err := json.Marshal(rawMedias)
	if err != nil {
		return nil, err
	}

	var medias []NocoMedia
	if err := json.Unmarshal(data, &medias); err != nil {
		return nil, err
	}

	return medias, nil
}

func (b *BusinessService) GetMediaByID(w http.ResponseWriter, r *http.Request) {
	subCtx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	params := httprouter.ParamsFromContext(subCtx)

	id, err := strconv.ParseInt(params.ByName("id"), 10, 32)
	if err != nil {
		http.Error(w, "id path param is invalid", http.StatusBadRequest)

		return
	}

	rawMedia, err := b.postgresService.Queries.FindRawMediaCheminByID(subCtx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "media not found", http.StatusNotFound)

			return
		}

		http.Error(w, "error while fetching data", http.StatusInternalServerError)

		return
	}

	media, err := parseMediaPath(rawMedia.String)
	if err != nil {
		http.Error(w, "error while fetching media", http.StatusInternalServerError)

		return
	}

	filename := filepath.Base(media.Path)

	w.Header().Set("Content-Type", media.Mimetype)

	http.ServeFile(w, r, filepath.Join(b.config.FileSystem.UploadDir, filename))
}

func (b *BusinessService) UploadImage(w http.ResponseWriter, r *http.Request) *_http.APIError {
	_, ok := r.Context().Value(catalogs.AccessToken).(jwt.SessionInfo)
	if !ok {
		return _http.ErrUnauthorized.Msg("invalid token")
	}

	r.Body = http.MaxBytesReader(w, r.Body, MaxFileSize+1)

	err := r.ParseMultipartForm(MaxFileSize + 1)
	if err != nil {
		return _http.ErrBadRequest.Msg("Formulaire invalide ou fichier trop gros").Err(err)
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		return _http.ErrBadRequest.Msg("Aucune image trouvée : " + err.Error()).Err(err)
	}

	defer func(file multipart.File) {
		_ = file.Close()
	}(file)

	if fileHeader.Size > MaxFileSize {
		return _http.ErrBadRequest.Msg(fmt.Sprintf("Fichier %s trop volumineux (max 5Mo)", fileHeader.Filename))
	}

	// Générer un nom unique
	filename := uuid.New().String() + filepath.Ext(fileHeader.Filename)
	savePath := filepath.Join(b.config.FileSystem.UploadDirForDB, filename)
	realPath := filepath.Join(b.config.FileSystem.UploadDir, filename)

	dst, err := os.Create(realPath)
	if err != nil {
		return _http.ErrInternalError.Msg("Erreur écriture fichier").Err(err)
	}

	// Copier avec une limite de taille
	size, err := io.Copy(dst, io.LimitReader(file, MaxFileSize+1))
	_ = dst.Close()

	if err != nil {
		return _http.ErrInternalError.Msg("Erreur copie fichier").Err(err)
	}

	mimeType := fileHeader.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	title := r.FormValue("title")

	// Insérer en DB
	var id int32

	id, err = b.postgresService.Queries.CreateNewMedia(r.Context(), queries.CreateNewMediaParams{
		Title: pgtype.Text{String: title, Valid: title != ""},
		CheminMedia: pgtype.Text{
			String: fmt.Sprintf("[{\"title\":%q,\"mimetype\":%q,\"size\":%d,\"path\":%q}]", title, mimeType, size, savePath),
			Valid:  true,
		},
		DateCreation: pgtype.Text{String: time.Now().String(), Valid: true},
	})
	if err != nil {
		return _http.ErrInternalError.Msg("Erreur DB").Err(err)
	}

	// Réponse JSON
	return _http.WriteJSONResponse(w, http.StatusCreated, api.Media{
		ID:    &id,
		Title: &title,
	})
}
