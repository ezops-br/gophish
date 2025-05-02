package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// CampaignPreset represents a saved preset for campaign fields
type CampaignPreset struct {
	Id          int64     `json:"id"`
	UserId      int64     `json:"-"`
	Name        string    `json:"name" sql:"not null"`
	CreatedDate time.Time `json:"created_date"`
	TemplateId  int64     `json:"-"`
	Template    Template  `json:"template"`
	PageId      int64     `json:"-"`
	Page        Page      `json:"page"`
	SMTPId      int64     `json:"-"`
	SMTP        SMTP      `json:"smtp"`
	URL         string    `json:"url"`
}

// GetCampaignPresets returns all campaign presets for a user
func GetCampaignPresets(uid int64) ([]CampaignPreset, error) {
	presets := []CampaignPreset{}
	err := db.Where("user_id=?", uid).Find(&presets).Error
	if err != nil {
		return presets, err
	}
	for i := range presets {
		err = presets[i].getDetails()
		if err != nil {
			return presets, err
		}
	}
	return presets, err
}

// GetCampaignPreset returns the campaign preset specified by the id
func GetCampaignPreset(id int64, uid int64) (CampaignPreset, error) {
	preset := CampaignPreset{}
	err := db.Where("id=? and user_id=?", id, uid).First(&preset).Error
	if err != nil {
		return preset, err
	}
	err = preset.getDetails()
	return preset, err
}

// PostCampaignPreset creates a new campaign preset
func PostCampaignPreset(p *CampaignPreset) error {
	p.CreatedDate = time.Now().UTC()
	p.UserId = p.UserId
	return db.Save(p).Error
}

// PutCampaignPreset edits an existing campaign preset
func PutCampaignPreset(p *CampaignPreset) error {
	return db.Save(p).Error
}

// DeleteCampaignPreset deletes a campaign preset
func DeleteCampaignPreset(id int64) error {
	return db.Delete(&CampaignPreset{Id: id}).Error
}

// getDetails retrieves the related attributes of the preset
func (p *CampaignPreset) getDetails() error {
	err := db.Table("templates").Where("id=?", p.TemplateId).Find(&p.Template).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
		p.Template = Template{Name: "[Deleted]"}
	}
	err = db.Table("pages").Where("id=?", p.PageId).Find(&p.Page).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
		p.Page = Page{Name: "[Deleted]"}
	}
	err = db.Table("smtp").Where("id=?", p.SMTPId).Find(&p.SMTP).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
		p.SMTP = SMTP{Name: "[Deleted]"}
	}
	return nil
} 