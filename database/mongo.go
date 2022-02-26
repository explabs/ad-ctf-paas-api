package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/explabs/ad-ctf-paas-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var collection, flags, scoreboard, configurations *mongo.Collection

var ctx = context.TODO()

func InitMongo() {
	adminPass := os.Getenv("ADMIN_PASS")
	if adminPass == "" {
		adminPass = "admin"
	}
	credential := options.Credential{
		Username: "admin",
		Password: adminPass,
	}

	mongoAddr := os.Getenv("MONGODB")
	if mongoAddr == "" {
		mongoAddr = "localhost:27017"
	}
	mongoURI := fmt.Sprintf("mongodb://%s", mongoAddr)
	clientOptions := options.Client().ApplyURI(mongoURI).SetAuth(credential)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("ad").Collection("teams")
	flags = client.Database("ad").Collection("flags")
	scoreboard = client.Database("ad").Collection("scoreboard")
	configurations = client.Database("ad").Collection("config")
}

func CreateTeam(team *models.Team) error {
	_, err := collection.InsertOne(ctx, team)
	return err
}
func GetTeams() ([]*models.TeamInfo, error) {
	// passing bson.D{{}} matches all documents in the collection
	filter := bson.M{"login": bson.M{"$ne": "admin"}}
	return FilterTeams(filter)
}
func GetTeam(login string) (*models.TeamInfo, error) {
	// passing bson.D{{}} matches all documents in the collection
	var team models.TeamInfo
	filter := bson.M{"login": login}
	err := collection.FindOne(ctx, filter).Decode(&team)
	if err != nil {
		return nil, err
	}
	return &team, err
}
func GetUsers() ([]*models.TeamInfo, error) {
	// passing bson.D{{}} matches all documents in the collection
	filter := bson.D{{}}
	return FilterTeams(filter)
}
func GetAuthTeam(login string) (team models.Team, err error) {
	cur := collection.FindOne(ctx, bson.M{"login": login})
	cur.Decode(&team)
	if err != nil {
		return team, err
	}
	return team, nil
}

func FilterTeams(filter interface{}) ([]*models.TeamInfo, error) {
	// A slice of teams for storing the decoded documents
	var teams []*models.TeamInfo

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return teams, err
	}

	for cur.Next(ctx) {
		var t models.TeamInfo
		err := cur.Decode(&t)
		if err != nil {
			return teams, err
		}

		teams = append(teams, &t)
	}

	if err := cur.Err(); err != nil {
		return teams, err
	}

	// once exhausted, close the cursor
	cur.Close(ctx)

	if len(teams) == 0 {
		return teams, mongo.ErrNoDocuments
	}

	return teams, nil
}

func DeleteTeam(name string) error {
	filter := bson.D{primitive.E{Key: "login", Value: name}}

	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("No teams were deleted")
	}

	filter = bson.D{primitive.E{Key: "name", Value: name}}

	res, err = scoreboard.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func AddAttackFlag(team, service string) error {
	return AddFlag(team, service, "gained")
}

func AddDefenceFlag(team, service string) error {
	return AddFlag(team, service, "lost")
}

func AddFlag(team, service, field string) error {
	_, err := flags.UpdateOne(ctx, bson.M{
		"team":    team,
		"service": service,
	}, bson.D{
		{"$inc", bson.D{{field, 1}}},
	}, options.Update().SetUpsert(true))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

type ServiceFlagsStats struct {
	Gained float64 `bson:"gained"`
	Lost   float64 `bson:"lost"`
}

func GetServiceFlagsStats(team, service string) (f ServiceFlagsStats) {
	res := flags.FindOne(ctx, bson.M{
		"team":    team,
		"service": service,
	})
	res.Decode(&f)
	return f
}

func GetScoreboard() ([]models.Score, error) {
	cur, err := scoreboard.Find(context.TODO(), bson.D{})
	defer cur.Close(context.TODO())
	if err != nil {
		return []models.Score{}, err
	}

	var scoreboard []models.Score

	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var score models.Score
		err := cur.Decode(&score)
		if err != nil {
			log.Fatal(err)
		}
		scoreboard = append(scoreboard, score)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return scoreboard, nil
}
func GetTeamsScoreboard(teamName string) (models.Score, error) {
	var teamsScore models.Score
	err := scoreboard.FindOne(context.TODO(), bson.M{"name": teamName}).Decode(&teamsScore)
	if err == mongo.ErrNoDocuments {
		return models.Score{
			Name:         teamName,
			Round:        0,
			Services:     map[string]models.ScoreService{},
			LastServices: map[string]models.ScoreService{},
			Score:        0,
			LastScore:    0,
		}, nil
	}
	if err != nil {
		return models.Score{}, err
	}
	return teamsScore, nil
}
