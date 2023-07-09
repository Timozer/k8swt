import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import { AttachAddon } from "xterm-addon-attach";
import { SerializeAddon } from "xterm-addon-serialize";
import { Unicode11Addon } from "xterm-addon-unicode11";
import { WebLinksAddon } from "xterm-addon-web-links";
import { CanvasAddon } from "xterm-addon-canvas";
import { SearchAddon } from "xterm-addon-search";
import 'xterm/css/xterm.css'
import './Term.css'

import { Collapse } from "./Collapse";
import { SearchBar } from "./SearchBar";
import { SettingBar } from "./SettingBar";

import { GetEnvsAndCmdsFromQueryStr, GetQueryParams } from "../utils/utils";

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
    const searchAddon = new SearchAddon();

    const wsConn = useRef(null);
    const term = useRef(null);
    const termRef = useRef(null);

    const [termOpts, setTermOpts] = useState({});

    function onTermOptsChange(opts) {
        for (var key in opts) {
            term.current.options[key] = opts[key]
        }
        fitAddon.fit()
        setTermOpts(term.current.options)
    }

    function mounted() {
        let url = getWsUrl();
        console.log("socket url:", url)
        let ws = new WebSocket(url);
        ws.timeoutInterval = 5400;
        ws.onopen = openTerm;
        ws.onclose = (e) => { term.current && term.current.dispose() };
        ws.onerror = (e) => { term.current && term.current.dispose(); console.log("socket err", e) };
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
        t.loadAddon(searchAddon);
        t.loadAddon(new Unicode11Addon())
        t.unicode.activeVersion = '11'

        t.onResize(size => {
            sendData('resize', size);
        });

        t.open(termRef.current);
        t.loadAddon(new CanvasAddon())
        t.loadAddon(fitAddon)
        fitAddon.fit()

        setTimeout(() => { wsConn.current.send(GetEnvsAndCmdsFromQueryStr()) }, 100);

        t.focus();
        term.current = t;
        setTermOpts(term.current.options)
    }
    
    function sendData(type, data) {
        wsConn.current && wsConn.current.send(JSON.stringify({ Type: type, Data: data }));
    }

    function onResize() {
        termRef.current && fitAddon.fit();
    }

    useEffect(() => {
        mounted()
        return unmounted;
    }, []);

    return (
        <>
            <div ref={termRef} className="terminal-container code">
            </div>
            <Collapse 
                collapsed={true} 
                onCollapsed={() => {
                    searchAddon.clearDecorations()
                    searchAddon.clearActiveDecoration()
                }}
            >
                <SearchBar onNext={(text, opts) => searchAddon.findNext(text, opts)} onPrev={(text, opts) => searchAddon.findPrevious(text, opts)} />
                {/* <SettingBar termOpts={termOpts} onOptsChange={onTermOptsChange} /> */}
            </Collapse>
        </>
    );
}
