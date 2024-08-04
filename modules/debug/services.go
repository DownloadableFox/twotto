package debug

import (
	"context"
	"errors"

	"github.com/downloadablefox/twotto/core"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Feature is a feature of the debug module
var (
	FeatureServiceKey         = core.NewIdentifier("debug", "service/features")
	ErrFeatureServiceNotFound = errors.New("ledger manager not found in context (missing injection)")
	ErrFeatureNotRegistered   = errors.New("feature not registered for guild")
)

type Feature struct {
	Identifier   *core.Identifier
	DefaultState bool
}

type FeatureService interface {
	RegisterFeature(identifier *core.Identifier, defaultValue bool) error
	ListFeatures() ([]Feature, error)
	GetFeature(ctx context.Context, identifier *core.Identifier, guildId string) (bool, error)
	SetFeature(ctx context.Context, identifier *core.Identifier, guildId string, enabled bool) error
}

type PostgresFeatureService struct {
	pool               *pgxpool.Pool
	registeredFeatures map[*core.Identifier]Feature
}

func NewPostgresFeatureService(pool *pgxpool.Pool) FeatureService {
	return &PostgresFeatureService{
		pool:               pool,
		registeredFeatures: make(map[*core.Identifier]Feature),
	}
}

func (s *PostgresFeatureService) RegisterFeature(identifier *core.Identifier, defaultValue bool) error {
	s.registeredFeatures[identifier] = Feature{
		Identifier:   identifier,
		DefaultState: defaultValue,
	}

	return nil
}

func (s *PostgresFeatureService) ListFeatures() ([]Feature, error) {
	var features []Feature
	for _, f := range s.registeredFeatures {
		features = append(features, f)
	}

	return features, nil
}

func (s *PostgresFeatureService) GetFeature(ctx context.Context, identifier *core.Identifier, guildId string) (bool, error) {
	var enabled bool
	err := s.pool.QueryRow(ctx, "SELECT enabled FROM debug_features WHERE guild_id = $1 AND name = $2", guildId, identifier.String()).Scan(&enabled)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, ErrFeatureNotRegistered
		}

		return false, err
	}

	return enabled, nil
}

func (s *PostgresFeatureService) SetFeature(ctx context.Context, identifier *core.Identifier, guildId string, enabled bool) error {
	_, err := s.pool.Exec(ctx,
		"INSERT INTO debug_features (guild_id, name, enabled) VALUES ($1, $2, $3) ON CONFLICT (guild_id, name) DO UPDATE SET enabled = $3",
		guildId, identifier.String(), enabled,
	)

	return err
}
