<template>
    <div ref="terminal" class="terminal-container">
    </div>
</template>

<style>
html::-webkit-scrollbar,
body::-webkit-scrollbar,
div::-webkit-scrollbar {
  /* display: none; */
  /* width: 0; */
}

.terminal-container {
    padding: 0px;
    position: fixed;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
}

.xterm-screen {
}

.bg {
  /* background-image: url("@/assets/code.svg"); */
  background-repeat: no-repeat;
  background-size: 33%;
  background-position: center;
}
</style>


<script>

import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import { AttachAddon } from "xterm-addon-attach";
import { SerializeAddon } from "xterm-addon-serialize";
import { Unicode11Addon } from "xterm-addon-unicode11";
import { WebLinksAddon } from "xterm-addon-web-links";
import { AdventureTime } from "xterm-theme";
import 'xterm/css/xterm.css'

const fitAddon = new FitAddon();

export default {
    // eslint-disable-next-line
    name: 'Terminal',
    data() {
        return {
            term: "",
            ws: null,
        }
    },
    mounted() {
        this.init()
    },
    unmounted() {
        if (this.ws) {
            this.ws.close();
        }
        if (this.term) {
            this.term.dispose();
        }
    },
    methods: {
        init() {
            let protocol = (location.protocol === 'https:') ? 'wss://' : 'ws://';
            let url = protocol + location.hostname + ((location.port) ? (':' + location.port) : '') + '/ws';
            let queryParams = this.$route.query;
            if (queryParams.namespace === undefined) {
                queryParams.namespace = "";
            }
            if (queryParams.podname === undefined) {
                queryParams.podname = "";
            }
            if (queryParams.ip === undefined) {
                queryParams.ip = "";
            }
            url += '?podname=' + queryParams.podname + '&ip=' + queryParams.ip + '&namespace=' + queryParams.namespace;

            console.log("socket url:", url)
            this.ws = new WebSocket(url);
            this.ws.timeoutInterval = 5400;
            this.ws.onopen = this.wsOnOpen;
            this.ws.onclose = this.wsOnClose;
            this.ws.onerror = this.wsOnError;
            this.ws.onmessage = this.wsOnMessage;
        },
        wsOnOpen() {
            let term = new Terminal({
                convertEol: true, scrollback: 1000, 
                cursorBlink: true, allowProposedApi: true,
                theme: AdventureTime,
            })
            term.loadAddon(new AttachAddon(this.ws))
            term.loadAddon(new SerializeAddon())
            term.loadAddon(new WebLinksAddon())
            term.loadAddon(new Unicode11Addon())
            term.unicode.activeVersion = '11'

            term.onResize(size => {
                this.onSend('resize', size);
            });
            // term.onData(data => {
                // console.log("term onData: ", data);
                // this.onSend('input', data);
            // });

            term.open(this.$refs.terminal)
            term.loadAddon(fitAddon)
            fitAddon.fit()

            window.addEventListener("resize", this.onResize)

            this.term = term
        },
        wsOnClose() {
            if (this.term) {
                this.term.dispose();
            }
            console.log("web socket closed")
        },
        wsOnError(err) {
            if (this.term) {
                this.term.dispose();
            }
            console.log("socket err:", err)
        },
        wsOnMessage(e) {
            console.log("web socket msg event", e)
        },
        onSend(type, data) {
            let d = {
                Type: type,
                Data: data,
            }
            this.ws.send(JSON.stringify(d));
        },
        onResize() {
            let terminal = this.$refs.terminal
            if (!terminal) {
                return
            }
            try {
                // let css = getComputedStyle(terminal);
                // console.log(terminal.clientWidth - parseInt(css.paddingLeft) - parseInt(css.paddingRight));
                fitAddon.fit()
            } catch (e) {
                console.log("e", e.message)
            }
        }
    }
}


</script>