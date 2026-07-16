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
	"strings"
	"time"

	"irj/internal/catalogs"
	"irj/internal/jwt"
	queries "irj/internal/postgres/_generated"
	"irj/pkg/api"
	_http "irj/pkg/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

const MaxFileSize = 1 << 20 // 1 MB

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

func (b *BusinessService) parseMedias(rawMedias any) ([]NocoMedia, error) {
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

	logger := zerolog.Ctx(subCtx)

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

	path := strings.Replace(media.Path, "download", b.config.FileSystem.BaseDir, 1)

	logger.Info().Msgf("media: %s", path)

	w.Header().Set("Content-Type", media.Mimetype)

	http.ServeFile(w, r, path)
}

func (b *BusinessService) UploadImage(w http.ResponseWriter, r *http.Request) *_http.APIError {
	logger := zerolog.Ctx(r.Context())

	_, ok := r.Context().Value(catalogs.AccessToken).(jwt.SessionInfo)
	if !ok {
		return _http.ErrUnauthorized.Msg("invalid token")
	}

	r.Body = http.MaxBytesReader(w, r.Body, MaxFileSize+1)

	err := r.ParseMultipartForm(MaxFileSize + 1) //nolint:gosec // G120: Unbound form data is limited to MaxFileSize + 1 bytes, which is safe
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

	logger.Info().Str("db_path", savePath).Str("path", realPath).Msg("Uploading image")

	dst, err := os.Create(realPath) //nolint:gosec // G703: path is constructed server-side from UUID, not user input
	if err != nil {
		return _http.ErrInternalError.Msg("Erreur écriture fichier").Err(err)
	}

	// Copier avec une limite de taille
	size, err := io.Copy(dst, io.LimitReader(file, MaxFileSize+1))
	_ = dst.Close()

	if err != nil {
		return _http.ErrInternalError.Msg("Erreur copie fichier").Err(err)
	}

	mimeType, apiError := extractMimeType(fileHeader)
	if apiError != nil {
		return apiError
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
		DateCreation: pgtype.Timestamptz{Time: time.Now(), Valid: true},
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

func (b *BusinessService) UpdateMediaTitle(w http.ResponseWriter, r *http.Request) *_http.APIError {
	subCtx, cancel := context.WithTimeout(r.Context(), defaultTimeOut)
	defer cancel()

	_, ok := r.Context().Value(catalogs.AccessToken).(jwt.SessionInfo)
	if !ok {
		return _http.ErrUnauthorized.Msg("invalid token")
	}

	params := httprouter.ParamsFromContext(subCtx)

	id, err := strconv.ParseInt(params.ByName("id"), 10, 32)
	if err != nil {
		return _http.ErrBadRequest.Msg("id path param is invalid")
	}

	var body struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return _http.ErrBadRequest.Msg("invalid request body")
	}

	err = b.postgresService.Queries.UpdateMediaTitle(subCtx, queries.UpdateMediaTitleParams{
		Title:     pgtype.Text{String: body.Title, Valid: true},
		JsonTitle: body.Title,
		ID:        int32(id),
	})
	if err != nil {
		return _http.ErrInternalError.Msg("error while updating media title").Err(err)
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (b *BusinessService) duplicateMediasIfShared(ctx context.Context, logger *zerolog.Logger, mediaIds, parentMediaIds []int32) []int32 {
	parentSet := make(map[int32]bool, len(parentMediaIds))
	for _, id := range parentMediaIds {
		parentSet[id] = true
	}

	result := make([]int32, 0, len(mediaIds))

	for _, id := range mediaIds {
		if !parentSet[id] {
			result = append(result, id)

			continue
		}

		newID, err := b.postgresService.Queries.DuplicateMedia(ctx, id)
		if err != nil {
			logger.Error().Err(err).Int32("source_id", id).Msg("failed to duplicate media")
			result = append(result, id) // fallback: use original

			continue
		}

		result = append(result, newID)
	}

	return result
}

func (b *BusinessService) applyMediaTitleChanges(
	ctx context.Context, logger *zerolog.Logger, origIDs, newIDs []int32, mediaTitles map[string]string,
) {
	if len(mediaTitles) == 0 {
		return
	}

	for i, origID := range origIDs {
		if i >= len(newIDs) {
			break
		}

		title, ok := mediaTitles[strconv.FormatInt(int64(origID), 10)]
		if !ok {
			continue
		}

		err := b.postgresService.Queries.UpdateMediaTitle(ctx, queries.UpdateMediaTitleParams{
			Title:     pgtype.Text{String: title, Valid: true},
			JsonTitle: title,
			ID:        newIDs[i],
		})
		if err != nil {
			logger.Error().Err(err).Int32("media_id", newIDs[i]).Msg("failed to update media title after duplication")
		}
	}
}

func extractMimeType(fileHeader *multipart.FileHeader) (string, *_http.APIError) {
	mimeType := fileHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(mimeType, "image/") && mimeType != "" {
		return "", _http.ErrBadRequest.Msg("Type de fichier non autorisé (jpeg, png, webp uniquement)")
	}

	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	return mimeType, nil
}
