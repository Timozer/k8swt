
export function GetQueryParams() {
    var params = {};
    var query = window.location.search.substring(1);
    var vars = query.split("&");
    for (var i = 0; i < vars.length; i++) {
        var pair = vars[i].split("=");
        // If first entry with this name
        if (typeof params[pair[0]] === "undefined") {
            params[pair[0]] = decodeURIComponent(pair[1]);
            // If second entry with this name
        } else if (typeof params[pair[0]] === "string") {
            var arr = [params[pair[0]], decodeURIComponent(pair[1])];
            params[pair[0]] = arr;
            // If third or later entry with this name
        } else {
            params[pair[0]].push(decodeURIComponent(pair[1]));
        }
    }
    return params;
}

export function GetEnvsAndCmdsFromQueryStr() {
    let data = "";
    let params = GetQueryParams();
    for (var key in params) {
        if (key === "envs" || key === "cmds") {
            let values = decodeURI(params[key]).split(";");
            values.forEach(val => {
                if (val.length > 0) {
                    switch (key) {
                        case "envs":
                            data += "export " + val + ";"
                            break;
                        case "cmds":
                            data += val + ";"
                            break;
                        default:
                            break;
                    }
                }
            });
        }
    }
    if (data.length > 0) {
        data += "\n"
    }
    return data
}