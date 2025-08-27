package repository

import (
	"context"
	"github.com/saleh-ghazimoradi/GraphQLCRUD/graph/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type JobListingMongo struct {
	ID          bson.ObjectID `bson:"_id"`
	Title       string        `bson:"title"`
	Description string        `bson:"description"`
	Company     string        `bson:"company"`
	URL         string        `bson:"url"`
}

type JobRepository interface {
	GetJob(ctx context.Context, id string) (*model.JobListing, error)
	GetJobs(ctx context.Context) ([]*model.JobListing, error)
	CreateJobListing(ctx context.Context, input model.CreateJobListingInput) (*model.JobListing, error)
	UpdateJobListing(ctx context.Context, id string, input model.UpdateJobListingInput) (*model.JobListing, error)
	DeleteJobListing(ctx context.Context, id string) (*model.DeleteJobResponse, error)
}

type jobRepository struct {
	collection *mongo.Collection
}

func (j *jobRepository) GetJob(ctx context.Context, id string) (*model.JobListing, error) {
	_id, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var jobMongo JobListingMongo
	if err := j.collection.FindOne(ctx, bson.M{"_id": _id}).Decode(&jobMongo); err != nil {
		return nil, err
	}
	return &model.JobListing{
		ID:          jobMongo.ID.Hex(),
		Title:       jobMongo.Title,
		Description: jobMongo.Description,
		Company:     jobMongo.Company,
		URL:         jobMongo.URL,
	}, nil
}

func (j *jobRepository) GetJobs(ctx context.Context) ([]*model.JobListing, error) {
	var jobsMongo []JobListingMongo
	cursor, err := j.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &jobsMongo); err != nil {
		return nil, err
	}
	jobs := make([]*model.JobListing, len(jobsMongo))
	for i, jobMongo := range jobsMongo {
		jobs[i] = &model.JobListing{
			ID:          jobMongo.ID.Hex(),
			Title:       jobMongo.Title,
			Description: jobMongo.Description,
			Company:     jobMongo.Company,
			URL:         jobMongo.URL,
		}
	}
	return jobs, nil
}

func (j *jobRepository) CreateJobListing(ctx context.Context, input model.CreateJobListingInput) (*model.JobListing, error) {
	_id := bson.NewObjectID()
	_, err := j.collection.InsertOne(ctx, bson.M{
		"_id":         _id,
		"title":       input.Title,
		"description": input.Description,
		"company":     input.Company,
		"url":         input.URL,
	})
	if err != nil {
		return nil, err
	}
	return &model.JobListing{
		ID:          _id.Hex(),
		Title:       input.Title,
		Description: input.Description,
		Company:     input.Company,
		URL:         input.URL,
	}, nil
}

func (j *jobRepository) UpdateJobListing(ctx context.Context, id string, input model.UpdateJobListingInput) (*model.JobListing, error) {
	_id, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	updatedJobInfo := bson.M{}
	if input.Title != nil {
		updatedJobInfo["title"] = *input.Title
	}
	if input.Description != nil {
		updatedJobInfo["description"] = *input.Description
	}
	if input.URL != nil {
		updatedJobInfo["url"] = *input.URL
	}
	var jobMongo JobListingMongo
	err = j.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": _id},
		bson.M{"$set": updatedJobInfo},
	).Decode(&jobMongo)
	if err != nil {
		return nil, err
	}
	return &model.JobListing{
		ID:          jobMongo.ID.Hex(),
		Title:       jobMongo.Title,
		Description: jobMongo.Description,
		Company:     jobMongo.Company,
		URL:         jobMongo.URL,
	}, nil
}

func (j *jobRepository) DeleteJobListing(ctx context.Context, id string) (*model.DeleteJobResponse, error) {
	_id, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	_, err = j.collection.DeleteOne(ctx, bson.M{"_id": _id})
	if err != nil {
		return nil, err
	}
	return &model.DeleteJobResponse{DeleteJobID: id}, nil
}

func NewJobRepository(db *mongo.Database, collectionName string) JobRepository {
	return &jobRepository{
		collection: db.Collection(collectionName),
	}
}
