package mng

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DbSet[T any] struct {
	CollectionName   string
	ConnectionString string

	conn     *mongo.Client
	feniksDb *mongo.Database
	coll     *mongo.Collection
}

func (d *DbSet[T]) Open() error {
	uri := d.ConnectionString
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	d.conn = client
	d.feniksDb = client.Database("feniks")
	d.coll = d.feniksDb.Collection(d.CollectionName)

	return nil
}

func (d *DbSet[T]) Close() error {
	if d.conn == nil {
		return nil
	}
	return d.conn.Disconnect(context.TODO())
}

func (d *DbSet[T]) GetByQuery(query bson.M, sorts bson.D) ([]*T, error) {
	ctx := context.TODO()
	opts := options.Find()
	if sorts != nil {
		opts.SetSort(sorts)
	}
	cursor, err := d.coll.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	ret := make([]*T, 0)
	err = cursor.All(ctx, &ret)
	return ret, nil
}

func (d *DbSet[T]) DeleteOneById(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	return d.coll.DeleteOne(context.TODO(), bson.M{"_id": id})
}
