package controllers

import (
	"api/app/pb"
	"api/models"
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NflApiServiceServer struct {
	pb.UnimplementedNflApiServiceServer
	Db  *mongo.Collection
	Ctx context.Context
}

func (nfl NflApiServiceServer) CreateTeam(ctx context.Context, req *pb.CreateTeamRequest) (*pb.CreateTeamResponse, error) {
	data := models.Team{
		Name:                strings.Title(req.GetName()),
		Conference:          strings.ToUpper(req.GetConference()),
		Divisional:          strings.Title(req.GetDivisional()),
		Stadium:             strings.Title(req.GetStadium()),
		State:               strings.Title(req.GetState()),
		Titles:              req.GetTitles(),
		SuperBowlAppearance: req.GetSuperBowlAppearance(),
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	result, err := nfl.Db.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %s", err.Error()))
	}

	data.Id = result.InsertedID.(primitive.ObjectID)

	return &pb.CreateTeamResponse{
		Success: true,
	}, nil
}

func (nfl NflApiServiceServer) GetTeamById(ctx context.Context, req *pb.GetTeamByIdRequest) (*pb.Team, error) {
	objectId, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Cannot convert string id to objectID: %s", err.Error()))
	}

	result := nfl.Db.FindOne(ctx, primitive.M{"_id": objectId})
	var data models.Team
	err = result.Decode(&data)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	return &pb.Team{
		Name:                data.Name,
		Conference:          data.Conference,
		Divisional:          data.Divisional,
		Stadium:             data.Stadium,
		State:               data.State,
		Titles:              data.Titles,
		SuperBowlAppearance: data.SuperBowlAppearance,
	}, nil
}

func (nfl NflApiServiceServer) GetTeams(req *pb.Empty, stream pb.NflApiService_GetTeamsServer) error {
	cursor, err := nfl.Db.Find(nfl.Ctx, primitive.M{})
	if err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unknown internal error: %s", err.Error()))
	}

	defer cursor.Close(nfl.Ctx)

	var data models.Team

	for cursor.Next(nfl.Ctx) {
		err = cursor.Decode(&data)
		if err != nil {
			return status.Errorf(codes.Unavailable, fmt.Sprintf("Could not decode data: %s", err.Error()))
		}

		stream.Send(&pb.Team{
			Name:                data.Name,
			Conference:          data.Conference,
			Divisional:          data.Divisional,
			Stadium:             data.Stadium,
			State:               data.State,
			Titles:              data.Titles,
			SuperBowlAppearance: data.SuperBowlAppearance,
		})
	}

	return nil
}

func (nfl NflApiServiceServer) UpdateTeam(ctx context.Context, req *pb.UpdateTeamRequest) (*pb.Team, error) {
	objectId, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Cannot convert string id to objectID: %s", err.Error()))
	}

	filter := primitive.M{"_id": objectId}
	update := primitive.M{
		"name":                strings.Title(req.Team.GetName()),
		"conference":          strings.ToUpper(req.Team.GetConference()),
		"divisional":          strings.ToUpper(req.Team.GetDivisional()),
		"stadium":             strings.Title(req.Team.GetStadium()),
		"state":               strings.Title(req.Team.GetState()),
		"titles":              req.Team.GetTitles(),
		"superBowlAppearance": req.Team.GetSuperBowlAppearance(),
		"updatedAt":           time.Now(),
	}

	result := nfl.Db.FindOneAndUpdate(ctx, filter, primitive.M{"$set": update}, options.FindOneAndUpdate().SetReturnDocument(1))
	var data models.Team
	err = result.Decode(&data)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, fmt.Sprintf("Could not decode data: %s", err.Error()))
	}

	return &pb.Team{
		Name:                data.Name,
		Conference:          data.Conference,
		Divisional:          data.Divisional,
		Stadium:             data.Stadium,
		State:               data.State,
		Titles:              data.Titles,
		SuperBowlAppearance: data.SuperBowlAppearance,
	}, nil
}

func (nfl NflApiServiceServer) DeleteTeam(ctx context.Context, req *pb.DeleteTeamRequest) (*pb.CreateTeamResponse, error) {
	objectId, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Cannot convert string id to objectID: %s", err.Error()))
	}

	_, err = nfl.Db.DeleteOne(ctx, primitive.M{"_id": objectId})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find/delete team with id %s: %s", req.GetId(), err.Error()))
	}

	return &pb.CreateTeamResponse{
		Success: true,
	}, nil
}
