campaigns: {
    get: function (id) {
        if (id) {
            return $.get("/api/campaigns/" + id + "/");
        }
        return $.get("/api/campaigns/");
    },
    post: function (data) {
        return $.post("/api/campaigns/", data);
    },
    put: function (id, data) {
        return $.ajax({
            url: "/api/campaigns/" + id + "/",
            type: "PUT",
            data: data
        });
    },
    delete: function (id) {
        return $.ajax({
            url: "/api/campaigns/" + id + "/",
            type: "DELETE"
        });
    },
    complete: function (id) {
        return $.post("/api/campaigns/" + id + "/complete");
    },
    summary: function () {
        return $.get("/api/campaigns/summary");
    },
    results: function (id) {
        return $.get("/api/campaigns/" + id + "/results");
    },
    campaignSummary: function (id) {
        return $.get("/api/campaigns/" + id + "/summary");
    },
    import: function (data) {
        return $.post("/api/campaigns/import", data);
    },
    export: function (id) {
        return $.get("/api/campaigns/" + id + "/export");
    }
},
presets: {
    get: function (id) {
        if (id) {
            return $.get("/api/presets/" + id + "/");
        }
        return $.get("/api/presets/");
    },
    post: function (data) {
        return $.post("/api/presets/", data);
    },
    put: function (id, data) {
        return $.ajax({
            url: "/api/presets/" + id + "/",
            type: "PUT",
            data: data
        });
    },
    delete: function (id) {
        return $.ajax({
            url: "/api/presets/" + id + "/",
            type: "DELETE"
        });
    }
},
templates: {
    // ... existing code ...
} 