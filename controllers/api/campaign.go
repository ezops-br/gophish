package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	ctx "github.com/gophish/gophish/context"
	log "github.com/gophish/gophish/logger"
	"github.com/gophish/gophish/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// Campaigns returns a list of campaigns if requested via GET.
// If requested via POST, APICampaigns creates a new campaign and returns a reference to it.
func (as *Server) Campaigns(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		cs, err := models.GetCampaigns(ctx.Get(r, "user_id").(int64))
		if err != nil {
			log.Error(err)
		}
		JSONResponse(w, cs, http.StatusOK)
	//POST: Create a new campaign and return it as JSON
	case r.Method == "POST":
		c := models.Campaign{}
		// Put the request into a campaign
		err := json.NewDecoder(r.Body).Decode(&c)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid JSON structure"}, http.StatusBadRequest)
			return
		}
		err = models.PostCampaign(&c, ctx.Get(r, "user_id").(int64))
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
			return
		}
		// If the campaign is scheduled to launch immediately, send it to the worker.
		// Otherwise, the worker will pick it up at the scheduled time
		if c.Status == models.CampaignInProgress {
			go as.worker.LaunchCampaign(c)
		}
		JSONResponse(w, c, http.StatusCreated)
	}
}

// CampaignsSummary returns the summary for the current user's campaigns
func (as *Server) CampaignsSummary(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		cs, err := models.GetCampaignSummaries(ctx.Get(r, "user_id").(int64))
		if err != nil {
			log.Error(err)
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, cs, http.StatusOK)
	}
}

// Campaign returns details about the requested campaign. If the campaign is not
// valid, APICampaign returns null.
func (as *Server) Campaign(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 0, 64)
	c, err := models.GetCampaign(id, ctx.Get(r, "user_id").(int64))
	if err != nil {
		log.Error(err)
		JSONResponse(w, models.Response{Success: false, Message: "Campaign not found"}, http.StatusNotFound)
		return
	}
	switch {
	case r.Method == "GET":
		JSONResponse(w, c, http.StatusOK)
	case r.Method == "DELETE":
		err = models.DeleteCampaign(id)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Error deleting campaign"}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, models.Response{Success: true, Message: "Campaign deleted successfully!"}, http.StatusOK)
	}
}

// CampaignResults returns just the results for a given campaign to
// significantly reduce the information returned.
func (as *Server) CampaignResults(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 0, 64)
	cr, err := models.GetCampaignResults(id, ctx.Get(r, "user_id").(int64))
	if err != nil {
		log.Error(err)
		JSONResponse(w, models.Response{Success: false, Message: "Campaign not found"}, http.StatusNotFound)
		return
	}
	if r.Method == "GET" {
		JSONResponse(w, cr, http.StatusOK)
		return
	}
}

// CampaignSummary returns the summary for a given campaign.
func (as *Server) CampaignSummary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 0, 64)
	switch {
	case r.Method == "GET":
		cs, err := models.GetCampaignSummary(id, ctx.Get(r, "user_id").(int64))
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				JSONResponse(w, models.Response{Success: false, Message: "Campaign not found"}, http.StatusNotFound)
			} else {
				JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			}
			log.Error(err)
			return
		}
		JSONResponse(w, cs, http.StatusOK)
	}
}

// CampaignComplete effectively "ends" a campaign.
// Future phishing emails clicked will return a simple "404" page.
func (as *Server) CampaignComplete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 0, 64)
	switch {
	case r.Method == "GET":
		err := models.CompleteCampaign(id, ctx.Get(r, "user_id").(int64))
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Error completing campaign"}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, models.Response{Success: true, Message: "Campaign completed successfully!"}, http.StatusOK)
	}
}

// CampaignReport generates and returns the campaign report document.
func (as *Server) CampaignReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 0, 64)
	uid := ctx.Get(r, "user_id").(int64)

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to parse form data: " + err.Error(),
		})
		return
	}

	lang := r.FormValue("lang")
	if lang == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing required field: lang",
		})
		return
	}

	file, _, err := r.FormFile("template_file")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to read template file: " + err.Error(),
		})
		return
	}
	defer file.Close()

	exePath, err := os.Executable()
	if err != nil {
		log.Errorf("Error getting executable path: %v", err)
		http.Error(w, "Failed to determine application path (executable error)", http.StatusInternalServerError)

		return
	}
	basePath := filepath.Dir(exePath)

	defaultTemplateFile := filepath.Join(basePath, "Goreport", "template.docx")
	currentTemplateFile := defaultTemplateFile

	if r.Method == "POST" && strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			log.Errorf("Error parsing multipart form: %v", err)
			http.Error(w, fmt.Sprintf("Failed to parse multipart form: %v", err), http.StatusBadRequest)

			return
		}

		lang := r.FormValue("lang")

		file, handler, err := r.FormFile("template_file")

		if err != nil && err != http.ErrMissingFile {
			log.Errorf("Error retrieving 'template_file': %v", err)
			http.Error(w, fmt.Sprintf("Error processing template file: %v", err), http.StatusInternalServerError)

			return
		}

		if handler != nil && err == nil {
			defer file.Close()
			tempDir := os.TempDir()
			tempPath := filepath.Join(tempDir, handler.Filename)
			out, errCreate := os.Create(tempPath)

			if errCreate != nil {
				log.Errorf("Error creating temporary file for upload [%s]: %v", tempPath, errCreate)
				http.Error(w, "Failed to save uploaded template (create)", http.StatusInternalServerError)

				return
			}

			_, errCopy := io.Copy(out, file)
			out.Close()

			if errCopy != nil {
				log.Errorf("Error copying uploaded file to temporary location [%s]: %v", tempPath, errCopy)
				http.Error(w, "Failed to save uploaded template (copy)", http.StatusInternalServerError)

				return
			}

			currentTemplateFile = tempPath

			defer func() {
				if err := os.Remove(tempPath); err != nil {
					log.Warnf("Failed to remove temporary template file [%s]: %v", tempPath, err)
				}
			}()

			fmt.Printf("Using uploaded templateFile: [%s]\n", currentTemplateFile)
		} else {
			fmt.Printf("No 'template_file' uploaded or error retrieving it (err: %v), using default: [%s]\n", err, currentTemplateFile)
		}

		report, modelErr := models.CampaignReport(id, uid, lang, currentTemplateFile)

		if modelErr != nil {
			log.Errorf("Error from models.CampaignReport: %v", modelErr)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Failed to generate report: " + modelErr.Error(),
			})
			return
		}

		w.Header().Set("Content-Disposition", "attachment; filename=campaign_report.docx")
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
		w.Write(report)
	}

	http.Error(w, "Invalid request method or Content-Type for report generation.", http.StatusMethodNotAllowed)
}
