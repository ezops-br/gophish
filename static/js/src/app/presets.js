// Presets contains the UI components for creating and managing campaign presets
$(document).ready(function() {
    // Setup datatables
    $("#presetsTable").DataTable({
        columnDefs: [
            {
                orderable: false,
                targets: "no-sort"
            }
        ],
        order: [[2, "desc"]]
    });
    
    // Setup the "New Preset" button to show the modal
    $("#new_button").on("click", function() {
        edit("new");
    });
    
    // Load the initial data
    load();
    
    // Setup multiple targets selection
    $("#template").select2({
        placeholder: "Select a Template",
        allowClear: true
    });
    
    $("#page").select2({
        placeholder: "Select a Landing Page",
        allowClear: true
    });
    
    $("#profile").select2({
        placeholder: "Select a Sending Profile",
        allowClear: true
    });
});

// load fetches the templates and groups data
var load = function() {
    $("#loading").show();
    api.templates.get()
        .success(function(templates) {
            if (templates.length == 0) {
                $("#emptyMessage").show();
                $("#loading").hide();
                return;
            }
            $("#loading").hide();
            template_s2 = $.map(templates, function(obj) {
                obj.text = obj.name;
                return obj;
            });
            $("#template").select2({
                placeholder: "Select a Template",
                data: template_s2
            });
        })
        .error(function() {
            $("#loading").hide();
            errorFlash("Error loading templates");
        });
        
    api.pages.get()
        .success(function(pages) {
            if (pages.length == 0) {
                $("#emptyMessage").show();
                $("#loading").hide();
                return;
            }
            $("#loading").hide();
            page_s2 = $.map(pages, function(obj) {
                obj.text = obj.name;
                return obj;
            });
            $("#page").select2({
                placeholder: "Select a Landing Page",
                data: page_s2
            });
        })
        .error(function() {
            $("#loading").hide();
            errorFlash("Error loading landing pages");
        });
        
    api.SMTP.get()
        .success(function(profiles) {
            if (profiles.length == 0) {
                $("#emptyMessage").show();
                $("#loading").hide();
                return;
            }
            $("#loading").hide();
            profile_s2 = $.map(profiles, function(obj) {
                obj.text = obj.name;
                return obj;
            });
            $("#profile").select2({
                placeholder: "Select a Sending Profile",
                data: profile_s2
            });
        })
        .error(function() {
            $("#loading").hide();
            errorFlash("Error loading sending profiles");
        });
        
    loadPresets();
};

// loadPresets fetches the existing presets and loads them into the view
var loadPresets = function() {
    $("#presetsTable").hide();
    $("#emptyMessage").hide();
    $("#loading").show();
    api.presets.get()
        .success(function(ps) {
            presets = ps;
            $("#loading").hide();
            if (presets.length == 0) {
                $("#emptyMessage").show();
            } else {
                presetsTable = $("#presetsTable").DataTable();
                presetsTable.clear();
                $.each(presets, function(i, preset) {
                    presetsTable.row.add([
                        escapeHtml(preset.name),
                        escapeHtml(preset.template.name),
                        moment(preset.created_date).format('MMMM Do YYYY, h:mm:ss a'),
                        "<div class='pull-right'>" +
                        "<button class='btn btn-primary' data-toggle='modal' data-backdrop='static' data-target='#modal' onclick='edit(" + preset.id + ")'><i class='fa fa-pencil'></i></button>" +
                        "<button class='btn btn-danger' onclick='deletePreset(" + preset.id + ")'><i class='fa fa-trash-o'></i></button></div>"
                    ]).draw();
                });
                $("#presetsTable").show();
            }
        })
        .error(function() {
            $("#loading").hide();
            errorFlash("Error loading presets");
        });
};

// edit opens the preset edit modal
var edit = function(id) {
    $("#presetModalLabel").text("Edit Preset");
    if (id == "new") {
        $("#presetModalLabel").text("New Preset");
        $("#name").val("");
        $("#url").val("");
        $("#template").select2("val", "");
        $("#page").select2("val", "");
        $("#profile").select2("val", "");
        $("#modal").modal("show");
        return;
    }
    
    // load the preset data
    api.presets.get(id)
        .success(function(preset) {
            $("#name").val(preset.name);
            $("#url").val(preset.url);
            
            if (preset.template.id) {
                $("#template").val(preset.template.id.toString());
                $("#template").trigger("change.select2");
            }
            
            if (preset.page.id) {
                $("#page").val(preset.page.id.toString());
                $("#page").trigger("change.select2");
            }
            
            if (preset.smtp.id) {
                $("#profile").val(preset.smtp.id.toString());
                $("#profile").trigger("change.select2");
            }
            
            // Store the preset ID
            window.editId = id;
            $("#modal").modal("show");
        })
        .error(function() {
            errorFlash("Error loading preset");
        });
};

// deletePreset deletes the selected preset
var deletePreset = function(id) {
    Swal.fire({
        title: "Are you sure?",
        text: "This will delete the preset. This can't be undone!",
        type: "warning",
        animation: false,
        showCancelButton: true,
        confirmButtonText: "Delete Preset",
        confirmButtonColor: "#428bca",
        reverseButtons: true,
        allowOutsideClick: false,
        showLoaderOnConfirm: true,
        preConfirm: function() {
            return new Promise(function(resolve, reject) {
                api.presets.delete(id)
                    .success(function(msg) {
                        resolve();
                    })
                    .error(function(data) {
                        reject(data.responseJSON.message);
                    });
            });
        }
    }).then(function(result) {
        if (result.value) {
            Swal.fire(
                'Preset Deleted!',
                'The preset has been deleted!',
                'success'
            );
            $('button[data-dismiss="modal"]').click();
            loadPresets();
        }
    });
};

// dismiss closes the new/edit modal
var dismiss = function() {
    $("#name").val("");
    $("#url").val("");
    $("#template").select2("val", "");
    $("#page").select2("val", "");
    $("#profile").select2("val", "");
    $("#modal\\.flashes").empty();
    $("#modal").modal("hide");
};

// save creates a new preset or updates an existing preset
var save = function() {
    var preset = {
        name: $("#name").val(),
        url: $("#url").val(),
        template: {
            name: $("#template").select2("data")[0].text
        },
        page: {
            name: $("#page").select2("data")[0].text
        },
        smtp: {
            name: $("#profile").select2("data")[0].text
        }
    };
    
    // Validate data
    if (!preset.name) {
        errorFlash("Preset name is required");
        return;
    }
    
    if (!preset.template.name) {
        errorFlash("Template is required");
        return;
    }
    
    if (!preset.page.name) {
        errorFlash("Landing page is required");
        return;
    }
    
    if (!preset.smtp.name) {
        errorFlash("Sending profile is required");
        return;
    }
    
    if (!preset.url) {
        errorFlash("URL is required");
        return;
    }
    
    // Submit the data
    if (window.editId) {
        // Update existing preset
        preset.id = window.editId;
        api.presets.put(preset)
            .success(function(data) {
                successFlash("Preset updated successfully!");
                dismiss();
                loadPresets();
            })
            .error(function(data) {
                errorFlash(data.responseJSON.message);
            });
    } else {
        // Create new preset
        api.presets.post(preset)
            .success(function(data) {
                successFlash("Preset created successfully!");
                dismiss();
                loadPresets();
            })
            .error(function(data) {
                errorFlash(data.responseJSON.message);
            });
    }
}; 