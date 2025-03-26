package app

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log/slog"
	"time"

	"{{SOSO_REPO}}/internal/database"
)

type Storage struct {
    q *database.Queries
    log *slog.Logger
}

func NewStorage(q *database.Queries, log *slog.Logger) *Storage {
    return &Storage{ q: q }
}

func (s *Storage) ValidateNewUser(token string) (error, bool, string) {
    return s.ValidateNewUserWithContext(context.Background(), token)
}

func (s *Storage) NewUser(email string, username string, hash string) (error, string) {
    return s.NewUserWithContext(context.Background(), email, username, hash)
}

func (s *Storage) StoreBugReport(username string, problem string, result string, steps string) (error, database.GetUserRow) {
    return s.StoreBugReportWithContext(context.Background(), username, problem, result, steps)
}

func (s *Storage) StoreBugReportWithContext(ctx context.Context, username string, problem string, result string, steps string) (error, database.GetUserRow) {
    user, err := s.q.GetUser(ctx, username)
    if err != nil {
        s.log.Error("getting user", "err", err)
        return err, database.GetUserRow{}
    }

    err = s.q.ReportBug(ctx, database.ReportBugParams{
        Problem: problem,
        Result: result,
        Steps: steps,
        Uid: user.ID,
    })

    if err != nil {
        s.log.Error("reporting bug", "err", err)
        return err, database.GetUserRow{}
    }

    return nil, user
}

func (s *Storage) NewUserWithContext(ctx context.Context, email string, username string, hash string) (error, string) {
    existingUser, err := s.q.GetUser(ctx, username)
    if err != nil && err != sql.ErrNoRows {
        s.log.Error("sql err", "err", err)
        return fmt.Errorf(INTERNAL_ERROR), ""
    }

    if existingUser.Username != "" {
        return fmt.Errorf(USERNAME_EXISTS_ERROR), ""
    }

    validbytes := make([]byte, 32)
    rand.Read(validbytes)
    validToken := base64.URLEncoding.EncodeToString(validbytes)[:16]

    err = s.q.SaveUser(ctx, database.SaveUserParams{
        Username: username,
        Email: email,
        Password: hash,
        ValidToken: sql.NullString{ String: validToken, Valid: true },
    })

    return nil, validToken
}

func (s *Storage) ValidateNewUserWithContext(ctx context.Context, token string) (error, bool, string) {
    user, err := s.q.GetUserByValidToken(ctx, sql.NullString{
        String: token,
        Valid: true,
    })

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, false, ""
        } else {
            s.log.Error("checking valid token", "token", token, "err", err)
            return fmt.Errorf(INTERNAL_ERROR), false, ""
        }
    }

    if user.Username != "" {
        err = s.q.ValidateUser(ctx, user.Username)

        if err != nil {
            s.log.Error("validating user", "user", user.Username, "err", err)
            return fmt.Errorf(INTERNAL_ERROR), false, ""
        }

        return nil, true, user.Username
    }

    return nil, false, ""
}
