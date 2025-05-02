package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	ctx "github.com/gophish/gophish/context"
	"github.com/gophish/gophish/models"
	"github.com/gorilla/mux"
)

func (as *Server) GetCampaignPresets(w http.ResponseWriter, r *http.Request) {
	u := ctx.Get(r, "user").(models.User)

	presets, err := models.GetCampaignPresets(u.Id)
	if err != nil {
		http.Error(w, "Error getting presets", http.StatusInternalServerError)
		return
	}

	JSONResponse(w, presets, http.StatusOK)
}

func (as *Server) GetCampaignPreset(w http.ResponseWriter, r *http.Request) {
	u := ctx.Get(r, "user").(models.User)

	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 0, 64)
	preset, err := models.GetCampaignPreset(id, u.Id)
	if err != nil {
		http.Error(w, "Error getting preset", http.StatusInternalServerError)
		return
	}

	JSONResponse(w, preset, http.StatusOK)
}

func (as *Server) PostCampaignPreset(w http.ResponseWriter, r *http.Request) {
	u := ctx.Get(r, "user").(models.User)

	preset := &models.CampaignPreset{}
	err := json.NewDecoder(r.Body).Decode(preset)
	if err != nil {
		http.Error(w, "Error decoding preset", http.StatusBadRequest)
		return
	}

	preset.UserId = u.Id

	// Get template ID
	template, err := models.GetTemplateByName(preset.Template.Name, u.Id)
	if err != nil {
		http.Error(w, "Template not found", http.StatusBadRequest)
		return
	}
	preset.TemplateId = template.Id

	// Get page ID
	page, err := models.GetPageByName(preset.Page.Name, u.Id)
	if err != nil {
		http.Error(w, "Page not found", http.StatusBadRequest)
		return
	}
	preset.PageId = page.Id

	// Get SMTP ID
	smtp, err := models.GetSMTPByName(preset.SMTP.Name, u.Id)
	if err != nil {
		http.Error(w, "SMTP profile not found", http.StatusBadRequest)
		return
	}
	preset.SMTPId = smtp.Id

	err = models.PostCampaignPreset(preset)
	if err != nil {
		http.Error(w, "Error creating preset", http.StatusInternalServerError)
		return
	}

	JSONResponse(w, preset, http.StatusCreated)
}

func (as *Server) PutCampaignPreset(w http.ResponseWriter, r *http.Request) {
	u := ctx.Get(r, "user").(models.User)

	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 0, 64)
	preset, err := models.GetCampaignPreset(id, u.Id)
	if err != nil {
		http.Error(w, "Preset not found", http.StatusNotFound)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&preset)
	if err != nil {
		http.Error(w, "Error decoding preset", http.StatusBadRequest)
		return
	}

	// Get template ID
	template, err := models.GetTemplateByName(preset.Template.Name, u.Id)
	if err != nil {
		http.Error(w, "Template not found", http.StatusBadRequest)
		return
	}
	preset.TemplateId = template.Id

	// Get page ID
	page, err := models.GetPageByName(preset.Page.Name, u.Id)
	if err != nil {
		http.Error(w, "Page not found", http.StatusBadRequest)
		return
	}
	preset.PageId = page.Id

	// Get SMTP ID
	smtp, err := models.GetSMTPByName(preset.SMTP.Name, u.Id)
	if err != nil {
		http.Error(w, "SMTP profile not found", http.StatusBadRequest)
		return
	}
	preset.SMTPId = smtp.Id

	err = models.PutCampaignPreset(&preset)
	if err != nil {
		http.Error(w, "Error updating preset", http.StatusInternalServerError)
		return
	}

	JSONResponse(w, preset, http.StatusOK)
}

// DeleteCampaignPreset deletes a preset
func (as *Server) DeleteCampaignPreset(w http.ResponseWriter, r *http.Request) {
	u := ctx.Get(r, "user").(models.User)

	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 0, 64)
	_, err := models.GetCampaignPreset(id, u.Id)
	if err != nil {
		http.Error(w, "Preset not found", http.StatusNotFound)
		return
	}

	err = models.DeleteCampaignPreset(id)
	if err != nil {
		http.Error(w, "Error deleting preset", http.StatusInternalServerError)
		return
	}

	JSONResponse(w, nil, http.StatusOK)
} 