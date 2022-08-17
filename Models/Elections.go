package Models

import (
	"fmt"
	"reflect"
	"strings"

	"example.com/example/Utils"
	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Elections struct {
	ID                    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CreatedBy             primitive.ObjectID `json:"createdby"`
	ElectionTitle         string             `json:"electiontitle"`
	ElectionType          string             `json:"electiontype"`
	ElectionStartDate     primitive.DateTime `json:"electionstartdate"`
	ElectionEndDate       primitive.DateTime `json:"electionenddate"`
	RegistrationStartDate primitive.DateTime `json:"registrationstartdate"`
	RegistrationEndDate   primitive.DateTime `json:"registrationenddate"`
	Positions             []Positions        `json:"positions"`
}
type ElectionsPopulated struct {
	ID                    primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	CreatedBy             Users                `json:"createdby"`
	ElectionTitle         string               `json:"electiontitle"`
	ElectionType          string               `json:"electiontype"`
	ElectionStartDate     primitive.DateTime   `json:"electionstartdate"`
	ElectionEndDate       primitive.DateTime   `json:"electionenddate"`
	RegistrationStartDate primitive.DateTime   `json:"registrationstartdate"`
	RegistrationEndDate   primitive.DateTime   `json:"registrationenddate"`
	Positions             []PositionsPopulated `json:"positions"`
}
type ElectionsSearch struct {
	ID                          primitive.ObjectID `json:"_id"`
	IDIsUsed                    bool               `json:"idisused"`
	CreatedBy                   primitive.ObjectID `json:"createdby"`
	CreatedByIsUsed             bool               `json:"createdbyisused"`
	ElectionTitle               string             `json:"electiontitle"`
	ElectionTitleIsUsed         bool               `json:"electiontitleisused"`
	ElectionType                string             `json:"electiontype"`
	ElectionTypeIsUsed          bool               `json:"electiontypeisused"`
	ElectionDateRangeFrom       primitive.DateTime `json:"electiondaterangefrom"`
	ElectionDateRangeTo         primitive.DateTime `json:"electiondaterangeto"`
	ElectionDateRangeIsUsed     bool               `json:"electiondaterangeisused"`
	RegistrationDateRangeFrom   primitive.DateTime `json:"registrationdaterangefrom"`
	RegistrationDateRangeTo     primitive.DateTime `json:"registrationdaterangeto"`
	RegistrationDateRangeIsUsed bool               `json:"registrationdaterangeisused"`
}
type Positions struct {
	PositionTitle    string               `json:"positiontitle"`
	CandidatesNumber int                  `json:"candidatesnumber"`
	Candidates       []primitive.ObjectID `json:"candidates"`
	Symbols          []Symbols            `json:"symbols"`
}
type PositionsPopulated struct {
	PositionTitle    string    `json:"positiontitle"`
	CandidatesNumber int       `json:"candidatesnumber"`
	Candidates       []Users   `json:"candidates"`
	Symbols          []Symbols `json:"symbols"`
}
type Symbols struct {
	Username       string `json:"username"`
	PhotoBase64    string `json:"photobase64"`
	PhotoPath      string `json:"photopath"`
	PhotoExtension string `json:"photoextension"`
}

func (obj Elections) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.ElectionTitle, validation.Required),
		validation.Field(&obj.ElectionType, validation.Required),
		validation.Field(&obj.ElectionStartDate, validation.Required),
		validation.Field(&obj.ElectionEndDate, validation.Required),
		validation.Field(&obj.RegistrationStartDate, validation.Required),
		validation.Field(&obj.RegistrationEndDate, validation.Required),
		validation.Field(&obj.Positions, validation.Required, validation.Length(1, 0)),
	)
}
func (obj Positions) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.PositionTitle, validation.Required),
		validation.Field(&obj.CandidatesNumber, validation.Required),
		validation.Field(&obj.Candidates, validation.Required, validation.Length(1, 0)),
		validation.Field(&obj.Symbols, validation.Required, validation.Length(1, 0)),
	)
}
func (obj Elections) GetModifcationBSONObj() bson.M {
	self := bson.M{}
	valueOfObj := reflect.ValueOf(obj)
	typeOfObj := valueOfObj.Type()
	invalidFieldNames := []string{"ID"}

	for i := 0; i < valueOfObj.NumField(); i++ {
		if Utils.ArrayStringContains(invalidFieldNames, typeOfObj.Field(i).Name) {
			continue
		}
		self[strings.ToLower(typeOfObj.Field(i).Name)] = valueOfObj.Field(i).Interface()
	}
	return self
}
func (obj ElectionsSearch) GetElectionsSearchBSONObj() bson.M {
	self := bson.M{}
	if obj.IDIsUsed {
		self["_id"] = obj.ID
	}

	if obj.CreatedByIsUsed {
		self["createdby"] = obj.CreatedBy
	}
	if obj.ElectionTitleIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.ElectionTitleIsUsed)
		self["electiontitle"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}
	if obj.ElectionTypeIsUsed {
		self["electiontype"] = obj.ElectionType
	}
	if obj.ElectionDateRangeIsUsed {
		self["electionstartdate"] = bson.M{
			"$gte": obj.ElectionDateRangeFrom,
		}
		self["electionenddate"] = bson.M{
			"$lte": obj.ElectionDateRangeTo,
		}
	}
	if obj.RegistrationDateRangeIsUsed {
		self["registrationstartdate"] = bson.M{
			"$gte": obj.RegistrationDateRangeFrom,
		}
		self["registrationenddate"] = bson.M{
			"$lte": obj.RegistrationDateRangeTo,
		}
	}

	return self
}
func (obj *PositionsPopulated) CloneFrom(other Positions) {
	obj.PositionTitle = other.PositionTitle
	obj.CandidatesNumber = other.CandidatesNumber
	obj.Candidates = []Users{}
	obj.Symbols = other.Symbols
}
func (obj *ElectionsPopulated) CloneFrom(other Elections) {
	obj.ID = other.ID
	obj.CreatedBy = Users{}
	obj.ElectionTitle = other.ElectionTitle
	obj.ElectionType = other.ElectionType
	obj.ElectionStartDate = other.ElectionStartDate
	obj.ElectionEndDate = other.ElectionEndDate
	obj.RegistrationStartDate = other.RegistrationStartDate
	obj.RegistrationEndDate = other.RegistrationEndDate
	obj.Positions = []PositionsPopulated{}
}
