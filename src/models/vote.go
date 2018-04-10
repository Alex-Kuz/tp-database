package models


import (
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)


type Vote struct {

	Nickname string `json:"nickname"`

	Voice    int32  `json:"voice"`
}


func (m *Vote) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateNickname(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateVoice(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Vote) validateNickname(formats strfmt.Registry) error {

	if err := validate.RequiredString("nickname", "body", string(m.Nickname)); err != nil {
		return err
	}

	return nil
}

var voteTypeVoicePropEnum []interface{}

func init() {
	var res []int32
	if err := json.Unmarshal([]byte(`[-1,1]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		voteTypeVoicePropEnum = append(voteTypeVoicePropEnum, v)
	}
}

// prop value enum
func (m *Vote) validateVoiceEnum(path, location string, value int32) error {
	if err := validate.Enum(path, location, value, voteTypeVoicePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *Vote) validateVoice(formats strfmt.Registry) error {

	if err := validate.Required("voice", "body", int32(m.Voice)); err != nil {
		return err
	}

	// value enum
	if err := m.validateVoiceEnum("voice", "body", m.Voice); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Vote) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Vote) UnmarshalBinary(b []byte) error {
	var res Vote
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
