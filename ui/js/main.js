String.prototype.hashCode = function(){
    if (Array.prototype.reduce){
        return this.split("").reduce(function(a,b){a=((a<<5)-a)+b.charCodeAt(0);return a&a},0);
    }
    var hash = 0;
    if (this.length === 0) return hash;
    for (var i = 0; i < this.length; i++) {
        var character  = this.charCodeAt(i);
        hash  = ((hash<<5)-hash)+character;
        hash = hash & hash; // Convert to 32bit integer
    }
    return hash;
};

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
            var content = '<div class="container">';
            $.each(data, function(idx, entry) {
                var entryString = entry.Address + "-" + entry.FlowId + "-" + entry.Reason + "-" + entry.AddedAt;
                var entryHash = entryString.hashCode();
                content += '<div class="row navitem" data-toggle="collapse" data-target="#BlacklistDetails-' + entryHash + '" onclick="BlacklistFetchDetails(\'' + entryHash + '\',\'' + entry.FlowId + '\',\'' + entry.AddedAt +'\',\'' + entry.ExpireOn + '\')">';
                content += '  <div class="col col-md-1">' + entry.Address + '</div>';
                content += '  <div class="col col-md-offset-4">' + entry.Reason + '</div>';
                content += '</div>';
                content += '<div id="BlacklistDetails-' + entryHash + '" class="row collapse">';
                content += '  <div class="col panel panel-default">';
                content += '    <div class="panel-body">';
                content += '      <div class="row">';
                content += '        <div class="col col-md-1">';
                content += '          Flow';
                content += '        </div>';
                content += '        <div id="Flow-' + entryHash + '" class="col col-md-6">';
                content += '        </div>';
                content += '        <div class="col col-md-1">';
                content += '          Country';
                content += '        </div>';
                content += '        <div id="Country-' + entryHash + '" class="col col-md-2">';
                content += '        </div>';
                content += '      </div>';
                content += '      <div class="row">';
                content += '        <div class="col col-md-1">';
                content += '          Bytes';
                content += '        </div>';
                content += '        <div class="col col-md-6">';
                content += '          Received: <span id="BytesReceived-' + entryHash + '"></span> Sent: <span id="BytesSent-' + entryHash + '"></span>';
                content += '        </div>';
                content += '        <div class="col col-md-1">';
                content += '          City';
                content += '        </div>';
                content += '        <div id="City-' + entryHash + '" class="col col-md-2">';
                content += '        </div>';
                content += '      </div>';
                content += '      <div class="row">';
                content += '        <div class="col col-md-1">';
                content += '          Packets';
                content += '        </div>';
                content += '        <div class="col col-md-6">';
                content += '          Received: <span id="PktsReceived-' + entryHash + '"></span> Sent: <span id="PktsSent-' + entryHash + '"></span>';
                content += '        </div>';
                content += '        <div class="col col-md-1">';
                content += '          PCode';
                content += '        </div>';
                content += '        <div id="PCode-' + entryHash + '" class="col col-md-2">';
                content += '        </div>';
                content += '      </div>';
                content += '      <div class="row">';
                content += '        <div class="col col-md-1">';
                content += '          Signature';
                content += '        </div>';
                content += '        <div id="Signature-' + entryHash +'" class="col col-md-6">';
                content += '        </div>';
                content += '        <div class="col col-md-1">';
                content += '          Timestamp';
                content += '        </div>';
                content += '        <div id="Timestamp-' + entryHash + '" class="col col-md-3">';
                content += '        </div>';
                content += '      </div>';
                content += '      <div class="row">';
                content += '        <div class="col col-md-1">';
                content += '          Category';
                content += '        </div>';
                content += '        <div id="Category-' + entryHash + '" class="col col-md-6">';
                content += '        </div>';
                content += '        <div class="col col-md-1">';
                content += '          Expires';
                content += '        </div>';
                content += '        <div id="Expires-' + entryHash + '" class="col col-md-3">';
                content += '        </div>';
                content += '      </div>';
                content += '    </div>';
                content += '    <div class="row">';
                content += '      <hr/>';
                content += '      <div class="col col-md-6">';
                content += '        <button type="button" class="btn btn-default">Make permanent</button>';
                content += '        <button type="button" class="btn btn-default">Unblock</button>';
                content += '        <button type="button" class="btn btn-default">Add to whitelist</button>';
                content += '      </div>';
                content += '    </div>';
                content += '  </div>';
                content += '</div>';
            });
            content += '</div>';

            $("#BlacklistView").html(content);
        },
        fail: function(data) {
            $("#BlacklistView").html("Failed to load blacklist " + data);
        }
    })
}

function BlacklistFetchDetails(entryHash, flowId, addedAt, expiresOn) {
    var ts = addedAt.replace("+01:00", "+0100");

    var newESQuery = {
        "query": {
            "constant_score": {
                "filter": {
                    "bool": {
                        "must": [
                            { "term": { "flow_id": parseInt(flowId) } },
                            { "term": { "timestamp": ts } }
                        ]
                    }
                }
            }
        }
    };

    console.log(newESQuery);

    var flow = $("#Flow-" + entryHash);
    var bytesSent = $("#BytesSent-" + entryHash);
    var bytesReceived = $("#BytesReceived-" + entryHash);
    var pktsSent = $("#PktsSent-" + entryHash);
    var pktsReceived = $("#PktsReceived-" + entryHash);
    var signature = $("#Signature-" + entryHash);
    var category = $("#Category-" + entryHash);
    var country = $("#Country-" + entryHash);
    var city = $("#City-" + entryHash);
    var pcode = $("#PCode-" + entryHash);
    var timestamp = $("#Timestamp-" + entryHash);
    var expires = $("#Expires-" + entryHash);
    var sensor = $("#Sensor-" + entryHash);
    var action = $("#Action-" + entryHash);

    $.ajax({
        url: "/api/v1/esproxy",
        method: "POST",
        data: JSON.stringify(newESQuery),
        dataType: "json",
        success: function(data) {
            if (data.status) {
                var d = JSON.parse(atob(data.data));

                if (d.hits.total === 1) {
                    var proto = d.hits.hits[0]._source.proto;
                    var srcIp = d.hits.hits[0]._source.src_ip;
                    var srcPort = d.hits.hits[0]._source.src_port;
                    var dstIp = d.hits.hits[0]._source.dest_ip;
                    var dstPort = d.hits.hits[0]._source.dest_port;

                    var flowContent = proto + " " + srcIp + ":" + srcPort + " > " + dstIp + ":" + dstPort;
                    flow.text(flowContent);

                    bytesSent.text(d.hits.hits[0]._source.flow.bytes_toclient);
                    bytesReceived.text(d.hits.hits[0]._source.flow.bytes_toserver);
                    pktsSent.text(d.hits.hits[0]._source.flow.pkts_toclient);
                    pktsReceived.text(d.hits.hits[0]._source.flow.pkts_toserver);

                    signature.text(d.hits.hits[0]._source.alert.signature);
                    category.text(d.hits.hits[0]._source.alert.category);


                    country.text(d.hits.hits[0]._source.geoip_src_ip.country_name);
                    city.text(d.hits.hits[0]._source.geoip_src_ip.city_name);
                    pcode.text(d.hits.hits[0]._source.geoip_src_ip.postal_code);
                    timestamp.text(d.hits.hits[0]._source.timestamp.split("+")[0].replace("T", " "));
                    expires.text(expiresOn.split("+")[0].replace("T", " "));

                    sensor.text(d.hits.hits[0]._source.host);
                    action.text(d.hits.hits[0]._source.alert.action);
                } else {
                    console.log("no hits!");
                    console.log(d.hits);
                }
            } else {
                console.log(data.message);
            }
        },
        fail: function(data) {
            $("#WhitelistContent").html("Failed to load blacklist " + data);
        }
    });
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