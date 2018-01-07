
function isIPAddress(addr) {
    if (addr.includes("/")) {
        addr = addr.split("/")[0];
    }

    var regexp = /((^\s*((([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]))\s*$)|(^\s*((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(%.+)?\s*$))/;

    return regexp.test(addr);
}

function LoadBlacklist() {
    $.ajax({
        url: "/api/v1/blacklist",
        dataType: "json",
        success: function(data) {
            var content = '<table class="table table-striped">';
            content += '<thead>';
            content += '<tr>';
            content += '<th>Address</th>';
            content += '<th>Reason</th>';
            content += '<th>Added at</th>';
            content += '<th>Expires on</th>';
            content += '</tr>';
            content += '</thead>';
            content += '<tbody>';
            $.each(data, function(idx, entry) {
                content += '<tr>';
                content += '<td>' + entry.Address + '</td>';
                content += '<td>' + entry.Reason + '</td>';
                content += '<td>' + entry.AddedAt + '</td>';
                content += '<td>' + entry.ExpireOn + '</td>';
                content += '</tr>';
            });
            content += '</tbody>';
            content += '</table>';

            console.log(content);
            $("#BlacklistView").html(content);
        },
        fail: function(data) {
            $("#BlacklistView").html("Failed to load blacklist " + data);
        }
    })
}

function LoadWhitelist() {
    $.ajax({
        url: "/api/v1/whitelist",
        dataType: "json",
        success: function(data) {
            var content = '<table class="table table-striped">';
            content += '<thead>';
            content += '<tr>';
            content += '<th>Address</th>';
            content += '<th>Description</th>';
            content += '<th>Actions</th>';
            content += '</tr>';
            content += '</thead>';
            content += '<tbody>';
            $.each(data, function(idx, entry) {
                content += '<tr>';
                content += '<td>' + entry.Address + '</td>';
                content += '<td>' + entry.Description + '</td>';
                content += '<td>';
                content += '<span onclick="WhitelistEditModal(\'' + entry.Address + '\',\'' + entry.Description + '\')" class="glyphicon glyphicon-edit" aria-hidden="true"></span>&nbsp;';
                content += '<span onclick="WhitelistRemoveConfirm(\'' + entry.Address + '\')" class="glyphicon glyphicon-remove" aria-hidden="true"></span>';
                content += '</td>';
                content += '</tr>';
            });
            content += '</tbody>';
            content += '</table>';

            $("#WhitelistContent").html(content);
        },
        fail: function(data) {
            $("#WhitelistContent").html("Failed to load blacklist " + data);
        }
    })
}

function WhitelistAdd() {
    var ipaddr = $("#WhitelistAddIpAddr").val();
    var descr = $("#WhitelistAddDescr").val();
    var foundErrors = false;

    if (!isIPAddress(ipaddr)) {
        $("#WhitelistAddIpAddrGroup").addClass("has-error");
        $("#WhitelistAddIpAddrGroup").addClass("has-feedback");
        foundErrors = true;
    } else {
        $("#WhitelistAddIpAddrGroup").removeClass("has-feedback");
        $("#WhitelistAddIpAddrGroup").removeClass("has-error");
    }

    if (descr === "") {
        $("#WhitelistAddDescrGroup").addClass("has-error");
        $("#WhitelistAddDescrGroup").addClass("has-feedback");
        foundErrors = true;
    } else {
        $("#WhitelistAddDescrGroup").removeClass("has-feedback");
        $("#WhitelistAddDescrGroup").removeClass("has-error");
    }

    if (foundErrors) {
        return
    }

    var newWhitelistEntry = {
      "ip_addr": ipaddr,
      "description": descr
    };

    $.ajax({
       url: "/api/v1/whitelist",
       method: "POST",
       dataType: "json",
       data: JSON.stringify(newWhitelistEntry),
       success: function(data) {
           $("#WhitelistAddIpAddr").val("");
           $("#WhitelistAddDescr").val("");
           LoadWhitelist();
       },
       fail: function(data) {
           console.log(data);
       }
    });
}

function WhitelistEditModal(addr, descr) {
    $("#WhitelistEditIpAddr").val(addr);
    $("#WhitelistEditDescr").val(descr);

    $("#WhitelistEditModal").modal('show');
}

function WhitelistEditConfirmed() {
    var newWhitelistEditRequest = {
        "ip_addr": $("#WhitelistEditIpAddr").val(),
        "description": $("#WhitelistEditDescr").val()
    };

    $.ajax({
        url: "/api/v1/whitelist",
        method: "PATCH",
        dataType: "json",
        data: JSON.stringify(newWhitelistEditRequest),
        success: function(data) {
            LoadWhitelist();
        },
        fail: function(data) {
            console.log(data);
        }
    });
}

function WhitelistRemoveConfirm(addr) {
    $("#WhitelistRemoveDoConfirm").on("click", function(ev) {
       ev.preventDefault();
       WhitelistRemoveConfirmed(addr);
    });
    $("#WhitelistRemoveConfirmIpAddr").text(addr);
    $("#WhitelistRemoveConfirmModal").modal('show');
}

function WhitelistRemoveConfirmed(addr) {
    var newWhitelistRemoveRequest = {
        "ip_addr": addr
    };

    $.ajax({
       url: "/api/v1/whitelist",
       method: "DELETE",
       dataType: "json",
       data: JSON.stringify(newWhitelistRemoveRequest),
       success: function(data) {
           LoadWhitelist();
       },
       fail: function(data) {
           console.log(data);
       }
    });
}

function ShowView(viewId) {
    allViews = ["#DashboardView", "#BlacklistView", "#WhitelistView", "#SettingsView"];
    allButtons = ["#NavPlayerView"];

    for (i = 0; i < allViews.length; i++) {
        if (allViews[i] == viewId) {
            continue;
        }
        $(allViews[i]).hide();
        $(allViews[i]).removeClass("active");
    }

    $(viewId).show();
    $(viewId).addClass("active");

    switch (viewId) {
        case "#BlacklistView":
            LoadBlacklist();
            break;
        case "#WhitelistView":
            LoadWhitelist();
            break;
    }
}


function main() {
    $(document).ready(function() {
        $("#DashboardView").hide();
        $("#BlacklistView").hide();
        $("#WhitelistView").hide();
        $("#SettingsView").hide();

        ShowView("#DashboardView");

        $("#DashboardViewNav").on("click", function(ev) {
            ev.preventDefault();
            ShowView("#DashboardView");
        });

        $("#BlacklistViewNav").on("click", function(ev) {
            ev.preventDefault();
            ShowView("#BlacklistView");
        });

        $("#WhitelistViewNav").on("click", function(ev) {
            ev.preventDefault();
            ShowView("#WhitelistView");
        });

        $("#SettingsViewNav").on("click", function(ev) {
            ev.preventDefault();
            ShowView("#SettingsView");
        });

        $("#WhitelistAdd").on("click", function(ev) {
            ev.preventDefault();
            WhitelistAdd();
        });

        $("#WhitelistEditDoConfirm").on("click", function(ev) {
            ev.preventDefault();
            WhitelistEditConfirmed();
        })
    });
}
main();