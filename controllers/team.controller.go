package controllers

import (
	"api/app/pb"
	"api/models"
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
