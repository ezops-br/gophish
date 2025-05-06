// Presets contains the UI components for creating and managing campaign presets
var presets = []
var preset = {}

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
    
    // Setup multiple modals
    // Code based on http://miles-by-motorcycle.com/static/bootstrap-modal/index.html
    $('.modal').on('hidden.bs.modal', function(event) {
        $(this).removeClass('fv-modal-stack');
        $('body').data('fv_open_modals', $('body').data('fv_open_modals') - 1);
    });
    
    $('.modal').on('shown.bs.modal', function(event) {
        // Keep track of the number of open modals
        if (typeof($('body').data('fv_open_modals')) == 'undefined') {
            $('body').data('fv_open_modals', 0);
        }
        // if the z-index of this modal has been set, ignore.
        if ($(this).hasClass('fv-modal-stack')) {
            return;
        }
        $(this).addClass('fv-modal-stack');
        // Increment the number of open modals
        $('body').data('fv_open_modals', $('body').data('fv_open_modals') + 1);
        // Setup the appropriate z-index
        $(this).css('z-index', 1040 + (10 * $('body').data('fv_open_modals')));
        $('.modal-backdrop').not('.fv-modal-stack').css('z-index', 1039 + (10 * $('body').data('fv_open_modals')));
        $('.modal-backdrop').not('fv-modal-stack').addClass('fv-modal-stack');
    });
    
    // Scrollbar fix - https://stackoverflow.com/questions/19305821/multiple-modals-overlay
    $(document).on('hidden.bs.modal', '.modal', function() {
        $('.modal:visible').length && $(document.body).addClass('modal-open');
    });
    
    $('#modal').on('hidden.bs.modal', function(event) {
        dismiss();
    });
});

// load fetches the templates and groups data
var load = function() {
    $("#loading").show();
    
    setupOptions();
    loadPresets();
};

// setupOptions configures the select2 options
function setupOptions() {
    api.templates.get()
        .success(function(templates) {
            if (templates.length == 0) {
                $("#emptyMessage").show();
                $("#loading").hide();
                return;
            }
            
            var template_s2 = $.map(templates, function(obj) {
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
            
            var page_s2 = $.map(pages, function(obj) {
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
            
            var profile_s2 = $.map(profiles, function(obj) {
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
}

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
        $("#template").val("").change();
        $("#page").val("").change();
        $("#profile").val("").change();
        $("#modal").modal("show");
        return;
    }
    
    // load the preset data
    api.presets.get(id)
        .success(function(p) {
            preset = p;
            $("#name").val(preset.name);
            $("#url").val(preset.url);
            
            if (preset.template.id) {
                $("#template").val(preset.template.id.toString()).change();
            }
            
            if (preset.page.id) {
                $("#page").val(preset.page.id.toString()).change();
            }
            
            if (preset.smtp.id) {
                $("#profile").val(preset.smtp.id.toString()).change();
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
        }
        $('button:contains("OK")').on('click', function() {
            location.reload();
        });
    });
};

var dismiss = function() {
    $("#modal\\.flashes").empty();
    $("#name").val("");
    $("#url").val("");
    $("#template").val("").change();
    $("#page").val("").change();
    $("#profile").val("").change();
    $("#modal").modal('hide');
};

// save creates a new preset or updates an existing preset
var save = function() {
    Swal.fire({
        title: "Are you sure?",
        text: "This will save the campaign preset.",
        type: "question",
        animation: false,
        showCancelButton: true,
        confirmButtonText: "Save",
        confirmButtonColor: "#428bca",
        reverseButtons: true,
        allowOutsideClick: false,
        showLoaderOnConfirm: true,
        preConfirm: function() {
            return new Promise(function(resolve, reject) {
                // Validate our fields
                var templateData = $("#template").select2("data");
                var pageData = $("#page").select2("data");
                var profileData = $("#profile").select2("data");
                
                // Create the preset
                preset = {
                    name: $("#name").val(),
                    template: {
                        name: templateData[0].text
                    },
                    url: $("#url").val(),
                    page: {
                        name: pageData[0].text
                    },
                    smtp: {
                        name: profileData[0].text
                    }
                }
                
                // Submit the preset
                if (window.editId) {
                    preset.id = window.editId;
                    api.presets.put(preset)
                        .success(function(data) {
                            resolve();
                            preset = data;
                        })
                        .error(function(data) {
                            $("#modal\\.flashes").empty().append("<div style=\"text-align:center\" class=\"alert alert-danger\">\
                            <i class=\"fa fa-exclamation-circle\"></i> " + data.responseJSON.message + "</div>");
                            Swal.close();
                        });
                } else {
                    api.presets.post(preset)
                        .success(function(data) {
                            resolve();
                            preset = data;
                        })
                        .error(function(data) {
                            $("#modal\\.flashes").empty().append("<div style=\"text-align:center\" class=\"alert alert-danger\">\
                            <i class=\"fa fa-exclamation-circle\"></i> " + data.responseJSON.message + "</div>");
                            Swal.close();
                        });
                }
            });
        }
    }).then(function(result) {
        if (result.value) {
            Swal.fire(
                'Preset Saved!',
                'This preset has been saved!',
                'success'
            );
        }
        $('button:contains("OK")').on('click', function() {
            dismiss();
            load();
        });
    });
}; 