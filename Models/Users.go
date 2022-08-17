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

type Users struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username       string             `json:"username"`
	Email          string             `json:"email"`
	Password       string             `json:"password"`
	PasswordHash   string             `json:"passwordbase64"`
	Fullname       string             `json:"fullname"`
	PhoneNumber    string             `json:"phonenumber"`
	Address        string             `json:"address"`
	City           string             `json:"city"`
	Country        string             `json:"country"`
	PhotoBase64    string             `json:"photobase64"`
	PhotoPath      string             `json:"photopath"`
	PhotoExtension string             `json:"photoextension"`
	Registered     bool               `json:"registered"`
}
type UsersSearch struct {
	ID                primitive.ObjectID `json:"id"`
	IDIsUsed          bool               `json:"idisused"`
	Username          string             `json:"username"`
	UsernameIsUsed    bool               `json:"usernameisused"`
	Fullname          string             `json:"fullname"`
	FullnameIsUsed    bool               `json:"fullnameisused"`
	PhoneNumber       string             `json:"phonenumber"`
	PhoneNumberIsUsed bool               `json:"phonenumberisused"`
	Address           string             `json:"address"`
	AddressIsUsed     bool               `json:"addressisused"`
	City              string             `json:"city"`
	CityIsUsed        bool               `json:"cityisused"`
	Country           string             `json:"country"`
	CountryIsUsed     bool               `json:"countryisused"`
	Registered        bool               `json:"registered"`
	RegisteredIsUsed  bool               `json:"registeredisused"`
}
type UsersLogin struct {
	EmailOrUsername string `json:"usernameoremail"`
	Password        string `json:"password"`
}

func (obj Users) ValidateSignUp() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Username, validation.Required),
		validation.Field(&obj.Email, validation.Required),
		validation.Field(&obj.Password, validation.Required),
	)
}
func (obj Users) ValidateRegister() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Email, validation.Required),
		validation.Field(&obj.Fullname, validation.Required),
		validation.Field(&obj.PhoneNumber, validation.Required),
		validation.Field(&obj.Address, validation.Required),
		validation.Field(&obj.City, validation.Required),
		validation.Field(&obj.Country, validation.Required),
	)
}
func (obj UsersLogin) ValidateLogin() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.EmailOrUsername, validation.Required),
	)
}
func (obj Users) GetModifcationBSONObj() bson.M {
	self := bson.M{}
	valueOfObj := reflect.ValueOf(obj)
	typeOfObj := valueOfObj.Type()
	invalidFieldNames := []string{"ID", "Username", "PasswordHash"}

	for i := 0; i < valueOfObj.NumField(); i++ {
		if Utils.ArrayStringContains(invalidFieldNames, typeOfObj.Field(i).Name) {
			continue
		}
		self[strings.ToLower(typeOfObj.Field(i).Name)] = valueOfObj.Field(i).Interface()
	}
	return self
}
func (obj UsersSearch) GetUsersSearchBSONObj() bson.M {
	self := bson.M{}
	if obj.IDIsUsed {
		self["_id"] = obj.ID
	}

	if obj.FullnameIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Fullname)
		self["fullname"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}
	if obj.PhoneNumberIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.PhoneNumber)
		self["phonenumber"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}
	if obj.AddressIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Address)
		self["address"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}
	if obj.CityIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.City)
		self["city"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}
	if obj.CountryIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Country)
		self["country"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.UsernameIsUsed {
		self["username"] = obj.Username
	}
	if obj.RegisteredIsUsed {
		self["registered"] = obj.Registered
	}

	return self
}
