package internal

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"irj/internal/smtp"
	"irj/pkg/utils"
)

func (e *Env) Cron(ctx context.Context) error {
	if err := e.DeleteInactiveUsers(ctx); err != nil {
		return err
	}

	if err := e.DeleteExpiredPasswordReset(ctx); err != nil {
		return err
	}

	e.Scheduler.StartAsync()

	return nil
}

func (e *Env) DeleteInactiveUsers(ctx context.Context) error {
	logger := e.Logger.With().Str("task", "delete_inactive_users").Logger()

	return utils.KeepError(e.Scheduler.Every(1).Day().At("00:00").Name("deleteInactiveUsers").Do(func() {
		users, err := e.PostgresService.Queries.GetUsers(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("failed to get users during inactive user analysis")

			return
		}

		for _, u := range users {
			if u.LastLogin.Valid && u.LastLogin.Time.Before(time.Now().Add(24*365*-3*time.Hour)) {
				err := e.PostgresService.Queries.DeleteUserByID(ctx, u.ID)
				if err != nil && !errors.Is(err, sql.ErrNoRows) {
					logger.Error().Err(err).Str("id", u.ID).Msg("failed to delete inactive user")

					continue
				}

				logger.Info().Str("id", u.ID).Time("lastLogin", u.LastLogin.Time).Msg("delete inactive user")

				_ = e.SMTPService.SendDeletionMail(ctx, []smtp.EmailPerson{
					{
						Name:  u.Prenom + " " + u.Nom,
						Email: u.Email,
					},
				}, true)
			}
		}
	}))
}

func (e *Env) DeleteExpiredPasswordReset(ctx context.Context) error {
	logger := e.Logger.With().Str("task", "delete_reset_password").Logger()

	return utils.KeepError(e.Scheduler.Every(1).Minute().Name("deleteExpiredResetPassword").Do(func() {
		err := e.PostgresService.Queries.DeleteExpiredPasswordReset(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("failed to delete expired password resets")

			return
		}
	}))
}
