// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// SiteProcessingSettings site processing settings
//
// swagger:model siteProcessingSettings
type SiteProcessingSettings struct {

	// css
	CSS *MinifyOptions `json:"css,omitempty"`

	// html
	HTML *SiteProcessingSettingsHTML `json:"html,omitempty"`

	// images
	Images *SiteProcessingSettingsImages `json:"images,omitempty"`

	// js
	Js *MinifyOptions `json:"js,omitempty"`

	// skip
	Skip bool `json:"skip,omitempty"`
}

// Validate validates this site processing settings
func (m *SiteProcessingSettings) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCSS(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateHTML(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateImages(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateJs(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SiteProcessingSettings) validateCSS(formats strfmt.Registry) error {

	if swag.IsZero(m.CSS) { // not required
		return nil
	}

	if m.CSS != nil {
		if err := m.CSS.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("css")
			}
			return err
		}
	}

	return nil
}

func (m *SiteProcessingSettings) validateHTML(formats strfmt.Registry) error {

	if swag.IsZero(m.HTML) { // not required
		return nil
	}

	if m.HTML != nil {
		if err := m.HTML.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("html")
			}
			return err
		}
	}

	return nil
}

func (m *SiteProcessingSettings) validateImages(formats strfmt.Registry) error {

	if swag.IsZero(m.Images) { // not required
		return nil
	}

	if m.Images != nil {
		if err := m.Images.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("images")
			}
			return err
		}
	}

	return nil
}

func (m *SiteProcessingSettings) validateJs(formats strfmt.Registry) error {

	if swag.IsZero(m.Js) { // not required
		return nil
	}

	if m.Js != nil {
		if err := m.Js.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("js")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *SiteProcessingSettings) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SiteProcessingSettings) UnmarshalBinary(b []byte) error {
	var res SiteProcessingSettings
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
