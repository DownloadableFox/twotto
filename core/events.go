package core

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type EventFunc[T any] func(ctx context.Context, s *discordgo.Session, e *T) error

func HandleEvent[T any](fn EventFunc[T]) interface{} {
	return func(s *discordgo.Session, e *T) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Error().Any("panic", rec).Msg("[EventHandler] Recovered from fatal error while executing event!")
			}
		}()

		if err := fn(context.Background(), s, e); err != nil {
			log.Error().Err(err).Msg("[EventHandler] Error executing event not handled!")
		}
	}
}

type MiddlewareFunc[T any] func(next EventFunc[T]) EventFunc[T]

func ApplyMiddlewares[T any](fn EventFunc[T], middlewares ...MiddlewareFunc[T]) EventFunc[T] {
	// Create wrapper function
	next := fn
	for i := len(middlewares) - 1; i >= 0; i-- {
		next = middlewares[i](next)
	}

	// Return execute to wrapper
	return func(c context.Context, s *discordgo.Session, e *T) error {
		return next(c, s, e)
	}
}
