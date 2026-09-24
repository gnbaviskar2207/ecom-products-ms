package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	errs "github.com/gnbaviskar2207/ecom-common/pkg/err"
	"github.com/gnbaviskar2207/ecom-common/pkg/utils"
	"github.com/gnbaviskar2207/ecom-products-ms/internal/domain"
	"github.com/gnbaviskar2207/ecom-products-ms/internal/dto"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoProductRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
	logger     *slog.Logger
}

func New(ctx context.Context, url, database, collection string, logger *slog.Logger) (*MongoProductRepository, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		_ = client.Disconnect(ctx)
		return nil, fmt.Errorf("mongodb error while ping %w", err)
	}
	productCollection := client.Database(database).Collection(collection)
	repo := &MongoProductRepository{
		client:     client,
		collection: productCollection,
		logger:     logger,
	}

	if err := repo.ensureIndexes(ctx); err != nil {
		_ = client.Disconnect(ctx)
		return nil, fmt.Errorf("failed to ensure mongodb indexes: %w", err)
	}
	return repo, nil
}

func (r *MongoProductRepository) ensureIndexes(ctx context.Context) error {
	indexes := r.collection.Indexes()
	cursor, err := indexes.List(ctx)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var index struct {
			Name string `bson:"name"`
		}
		if err := cursor.Decode(&index); err != nil {
			return err
		}
		if index.Name == "payload_pid" {
			r.logger.Info("index payload_pid already exists")
			return nil
		}
	}
	if err := cursor.Err(); err != nil {
		return err
	}
	_, err = indexes.CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "payload.pid", Value: 1},
		},
		Options: options.Index().SetName("payload_pid"),
	})
	if err != nil {
		return fmt.Errorf("failed to create index payload_pid: %w", err)
	}
	r.logger.Info("index payload_pid created")
	return nil
}

func (m *MongoProductRepository) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

func (m *MongoProductRepository) Ping(ctx context.Context) error {
	return m.client.Ping(ctx, nil)
}

func (m *MongoProductRepository) FindOneByPid(ctx context.Context, pid string) (*domain.Product, error) {
	var err error
	var result struct {
		Payload domain.Product `bson:"payload"`
	}
	cursor := m.collection.FindOne(ctx, bson.M{"payload.pid": pid}, options.FindOne().SetProjection(bson.D{{Key: "payload", Value: 1}}))
	err = cursor.Decode(&result)
	if err != nil {
		m.logger.ErrorContext(ctx, "product retrieval failed",
			slog.String("method", "repository.FindOneByPid"),
			slog.String("pid", pid),
			slog.String("error", err.Error()),
		)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%w: pid=%s", errs.ErrNotFound, pid)
		}
		return nil, err
	}
	m.logger.DebugContext(ctx, "product retrieved", slog.String("method", "repository.FindOneByPid"), slog.String("pid", pid))
	return &result.Payload, nil
}

func (m *MongoProductRepository) ListProducts(ctx context.Context, req *dto.ListProductsRequestDTO) (*domain.PaginatedResult[*domain.Product], error) {
	nextCursor := ""
	filter := bson.M{}
	if req != nil && req.NextCursor != "" {
		cursor, err := utils.DecodeBase64Str(req.NextCursor)
		if err != nil {
			return nil, fmt.Errorf("%w %s", errs.ErrInvalidArgument, err)
		}
		filter["payload.pid"] = bson.M{"$gt": cursor}
	}

	findOptions := options.Find().
		SetSort(bson.D{{Key: "payload.pids", Value: 1}}).
		SetLimit(req.Limit).
		SetProjection(bson.D{{Key: "payload", Value: 1}})

	cursor, err := m.collection.Find(ctx, filter, findOptions)
	if err != nil {
		m.logger.ErrorContext(ctx, "product retrieval failed",
			slog.String("method", "repository.ListProducts"),
			slog.String("next_cursor", req.NextCursor),
			slog.String("error", err.Error()),
		)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%w: next_cursor=%s", errs.ErrNotFound, req.NextCursor)
		}
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*domain.Product
	for cursor.Next(ctx) {
		var result struct {
			Payload domain.Product `bson:"payload"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		products = append(products, &result.Payload)
	}
	if len(products) > 0 {
		nextCursor = utils.EncodeBase64Str(products[len(products)-1].Pid)
	}
	return &domain.PaginatedResult[*domain.Product]{
		Data:       products,
		NextCursor: nextCursor,
		HasMore:    len(products) == int(req.Limit),
	}, nil
}
