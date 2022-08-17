package Controllers

import (
	"context"
	"encoding/json"
	"errors"

	"example.com/example/DBManager"
	"example.com/example/Models"
	"example.com/example/Utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ElectionsNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Elections
	var self Models.Elections
	c.BodyParser(&self)
	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}
	for i, val := range self.Positions {
		err := val.Validate()
		if err != nil {
			c.Status(500)
			return err
		}
		if val.CandidatesNumber != len(val.Candidates) {
			return errors.New("All Candidates Are Not Added For " + val.PositionTitle + " Position")
		}
		if val.CandidatesNumber != len(val.Symbols) {
			return errors.New("All Symbols Are Not Added For " + val.PositionTitle + " Position")
		}
		for j, val2 := range val.Candidates {
			userObj, err := UsersGetByIDFunction(val2)
			if err != nil {
				return errors.New("Invalid Candidate For " + val.PositionTitle + " Position")
			}
			if userObj.Username != val.Symbols[j].Username {
				return errors.New("Invalid Symbol For " + userObj.Username + " Candidate")
			}
		}
		for k, val3 := range val.Symbols {
			if val3.PhotoBase64 == "" && val3.PhotoPath == "" {
				return errors.New("Invalid Symbol For " + val3.Username + " Candidate")
			}
			// upload photo
			filePath, err := Utils.UploadImageBase64(val3.PhotoBase64, val3.PhotoExtension)
			if err != nil {
				c.Status(500)
				return err
			}
			if filePath != "" {
				self.Positions[i].Symbols[k].PhotoPath = filePath
			}
			self.Positions[i].Symbols[k].PhotoBase64 = ""
		}
	}
	_, err = collection.InsertOne(context.Background(), self)
	if err != nil {
		c.Status(500)
		return err
	}
	c.Set("Content-Type", "application/json")
	c.Status(200).Send([]byte("Election Created Successfully"))
	return nil
}
func ElectionsModify(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Elections
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.M{
		"_id": objID,
	}
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		c.Status(404)
		return errors.New("id is not found")
	}
	var self Models.Elections
	c.BodyParser(&self)
	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}
	for i, val := range self.Positions {
		err := val.Validate()
		if err != nil {
			c.Status(500)
			return err
		}
		if val.CandidatesNumber != len(val.Candidates) {
			return errors.New("All Candidates Are Not Added For " + val.PositionTitle + " Position")
		}
		if val.CandidatesNumber != len(val.Symbols) {
			return errors.New("All Symbols Are Not Added For " + val.PositionTitle + " Position")
		}
		for j, val2 := range val.Candidates {
			userObj, err := UsersGetByIDFunction(val2)
			if err != nil {
				return errors.New("Invalid Candidate For " + val.PositionTitle + " Position")
			}
			if userObj.Username != val.Symbols[j].Username {
				return errors.New("Invalid Symbol For " + userObj.Username + " Candidate")
			}
		}
		for k, val3 := range val.Symbols {
			if val3.PhotoBase64 == "" && val3.PhotoPath == "" {
				return errors.New("Invalid Symbol For " + val3.Username + " Candidate")
			}
			// upload photo
			filePath, err := Utils.UploadImageBase64(val3.PhotoBase64, val3.PhotoExtension)
			if err != nil {
				c.Status(500)
				return err
			}
			if filePath != "" {
				self.Positions[i].Symbols[k].PhotoPath = filePath
			}
			self.Positions[i].Symbols[k].PhotoBase64 = ""
		}
	}
	updateData := bson.M{
		"$set": self.GetModifcationBSONObj(),
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateData)
	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred when modifying Branch Document")
	} else {
		c.Status(200).Send([]byte("Modified Successfully"))
		return nil
	}
}
func ElectionsGetByID(id primitive.ObjectID) (Models.Elections, error) {
	collection := DBManager.SystemCollections.Elections
	filter := bson.M{"_id": id}
	var self Models.Elections
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) <= 0 {
		return self, errors.New("Object Not Found")
	}
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return self, nil
}
func ElectionsGetAll(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Elections

	results := []bson.M{}

	var searchParams Models.ElectionsSearch
	c.BodyParser(&searchParams)

	cur, err := collection.Find(context.Background(), searchParams.GetElectionsSearchBSONObj())
	if err != nil {
		c.Status(500)
		return err
	}
	defer cur.Close(context.Background())
	cur.All(context.Background(), &results)
	response, _ := json.Marshal(bson.M{
		"results": results,
	})
	c.Set("content-type", "application/json")
	c.Status(200).Send(response)
	return nil
}
func ElectionsGetByIDPopulated(objID primitive.ObjectID, ptr *Models.Elections) (Models.ElectionsPopulated, error) {
	var currentDoc Models.Elections
	if ptr == nil {
		currentDoc, _ = ElectionsGetByID(objID)
	} else {
		currentDoc = *ptr
	}
	populatedResult := Models.ElectionsPopulated{}
	populatedResult.CloneFrom(currentDoc)
	allPositionsPopulated := make([]Models.PositionsPopulated, len(currentDoc.Positions))
	for i, content := range currentDoc.Positions {
		var positionPopulated Models.PositionsPopulated
		positionPopulated.CloneFrom(content)
		for _, candidate := range content.Candidates {
			candidateObj, err := UsersGetByIDFunction(candidate)
			if err != nil {
				return populatedResult, err
			}
			positionPopulated.Candidates = append(positionPopulated.Candidates, candidateObj)
		}
		allPositionsPopulated[i] = positionPopulated
	}
	populatedResult.CreatedBy, _ = UsersGetByIDFunction(currentDoc.CreatedBy)
	populatedResult.Positions = allPositionsPopulated
	return populatedResult, nil
}
func ElectionsGetAllPopulated(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Elections

	var results []bson.M
	var searchRequests Models.ElectionsSearch
	c.BodyParser(&searchRequests)

	b, results := Utils.FindByFilter(collection, searchRequests.GetElectionsSearchBSONObj())
	if !b {
		c.Status(500)
		return errors.New("Obj Not Found")
	}

	// Convert
	var allRequestsDocuments []Models.Elections
	byteArr, _ := json.Marshal(results)
	json.Unmarshal(byteArr, &allRequestsDocuments)

	populatetedResults := make([]Models.ElectionsPopulated, len(allRequestsDocuments))

	for i, v := range allRequestsDocuments {
		populatetedResults[i], _ = ElectionsGetByIDPopulated(v.ID, &v)
	}

	allpopulated, _ := json.Marshal(bson.M{"results": populatetedResults})

	c.Set("Content-Type", "application/json")
	c.Status(200).Send(allpopulated)
	return nil
}
