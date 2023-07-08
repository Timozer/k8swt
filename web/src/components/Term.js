import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import { AttachAddon } from "xterm-addon-attach";
import { SerializeAddon } from "xterm-addon-serialize";
import { Unicode11Addon } from "xterm-addon-unicode11";
import { WebLinksAddon } from "xterm-addon-web-links";
import { CanvasAddon } from "xterm-addon-canvas";
import 'xterm/css/xterm.css'
import './Term.css'

import { GetQueryParams } from "../utils/getQueryParams";

import { useState, useRef, useEffect } from "react";

const getWsUrl = () => {
    let location = window.location;
    let protocol = (location.protocol === 'https:') ? 'wss://' : 'ws://';
    let url = protocol + location.hostname + ((location.port) ? (':' + location.port) : '') + '/ws';
    let params = GetQueryParams();

    url += '?';
    for (var key in params) {
        url += '&' + key + '=' + encodeURIComponent(params[key])
    }
    return url;
};

export function Term() {
    const fitAddon = new FitAddon();

    const wsConn = useRef(null);
    const term = useRef(null);
    const termRef = useRef(null);

    function mounted() {
        let url = getWsUrl();
        console.log("socket url:", url)
        let ws = new WebSocket(url);
        ws.timeoutInterval = 5400;
        ws.onopen = openTerm;
        ws.onclose = (e) => { term.current && term.current.dispose() };
        ws.onerror = (e) => { term.current && term.current.dispose(); alert(JSON.stringify(e)); console.log("socket err", e) };
        ws.onmessage = (e) => console.log("web socket msg event", e);
        wsConn.current = ws;

        window.addEventListener("resize", onResize)
    }

    function unmounted() {
        wsConn.current && wsConn.current.close();
        window.removeEventListener("resize", onResize);
    }

    function openTerm() {
        let t = new Terminal({
            convertEol: true, scrollback: 1000,
            cursorBlink: true, allowProposedApi: true,
        })
        t.loadAddon(new AttachAddon(wsConn.current))
        t.loadAddon(new SerializeAddon())
        t.loadAddon(new WebLinksAddon())
        t.loadAddon(new Unicode11Addon())
        t.unicode.activeVersion = '11'

        t.onResize(size => {
            sendData('resize', size);
        });

        t.open(termRef.current);
        t.loadAddon(new CanvasAddon())
        t.loadAddon(fitAddon)
        fitAddon.fit()

        setTimeout(() => { sendEnvsAmdCmds() }, 100);

        t.focus();
        term.current = t;
    }

    function sendData(type, data) {
        wsConn.current && wsConn.current.send(JSON.stringify({ Type: type, Data: data }));
    }

    function onResize() {
        termRef.current && fitAddon.fit();
        // let css = getComputedStyle(terminal);
        // console.log(terminal.clientWidth - parseInt(css.paddingLeft) - parseInt(css.paddingRight));
    }

    function sendEnvsAmdCmds() {
        let data = "";
        let params = GetQueryParams();
        for (var key in params) {
            if (key === "envs" || key === "cmds") {
                let values = decodeURI(params[key]).split(";");
                values.forEach(val => {
                    if  (val.length > 0) {
                        switch (key) {
                            case "envs":
                                data += "export " + val + ";"
                                break;
                            case "cmds":
                                data += val + ";"
                        }
                    }
                });
            }
        }
        if (data.length > 0) {
            data += "\n"
        }
        console.log("send envs or cmds", data)
        wsConn.current.send(data)
    }

    useEffect(() => {
        mounted()
        return unmounted;
    });

    return (
        <div ref={termRef} className="terminal-container code">
        </div>
    );
}
