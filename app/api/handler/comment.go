package handler

import (
	schemas "echo-postgres-sample/api/schemas"
	models "echo-postgres-sample/db"
	crud "echo-postgres-sample/db/crud"

	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// -- handle routes
func (h *Handler) CreateComment(c echo.Context) (err error) {
	var req schemas.CommentRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error)
	}
	comment, err := crud.CreateComment(h.DB, &models.Comment{
		Comment: req.Comment,
		UserID:  req.UserID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error)
	}
	return c.JSON(http.StatusOK, comment)
}

func (h *Handler) GetComments(c echo.Context) (err error) {
	comments, err := crud.GetComments(h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error)
	}
	return c.JSON(http.StatusOK, comments)
}

func (h *Handler) GetCommentByID(c echo.Context) (err error) {
	commentID := c.Param("commentID")
	comment, err := crud.GetComment(h.DB, commentID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error)
	}
	return c.JSON(http.StatusOK, comment)
}

func (h *Handler) UpdateCommentByID(c echo.Context) (err error) {
	req := schemas.CommentRequest{}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error)
	}
	commentID := c.Param("commentID")
	intCommentID, err := strconv.ParseInt(commentID, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error)
	}
	comment, err := crud.UpdateComment(h.DB, &models.Comment{
		ID:      intCommentID,
		Comment: req.Comment,
		UserID:  req.UserID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error)
	}
	return c.JSON(http.StatusOK, comment)
}

func (h *Handler) DeleteCommentByID(c echo.Context) (err error) {
	commentID := c.Param("commentID")
	intCommentID, err := strconv.ParseInt(commentID, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error)
	}
	err = crud.DeleteComment(h.DB, intCommentID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error)
	}
	return c.JSON(http.StatusOK, "OK")
}
