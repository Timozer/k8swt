import { useState, useRef, useEffect } from "react";

import searchNextIcon from "../assets/icons/search-next.svg"
import searchPrevIcon from "../assets/icons/search-prev.svg"

import "./SearchBar.css"

export function SearchBar({onNext, onPrev}) {
    const regexRef = useRef(null);
    const wholeWordRef = useRef(null);
    const caseSensitiveRef = useRef(null);
    const textRef = useRef(null);

    const [opts, setOpts] = useState({
        regex: false,
        wholeWord: false,
        caseSensitive: false,
        decorations: {
            matchBackground: "#99c2ff",
            matchBorder: "#99c2ff",
            activeMatchBackground: "#0066ff",
            activeMatchBorder: "#0066ff",
        }
    })

    function onOptsChange(e) {
        console.log(e)
        setOpts({
            ...opts,
            regex: regexRef.current.checked,
            wholeWord: wholeWordRef.current.checked,
            caseSensitive: caseSensitiveRef.current.checked,
        })
    }

    useEffect(() => {
        textRef.current.addEventListener("keyup", function (event) {
            if (event.key === "Enter") {
                onNext(textRef.current.value, opts)
            }
        });
    }, [textRef])

    return (
        <div className="search_bar">
            <div className="search_bar_options">
                <input type="checkbox" id="regex" ref={regexRef} checked={opts.regex} onChange={onOptsChange}/>
                <label for="regex">Regex</label>
                <input type="checkbox" id="wholeWord" ref={wholeWordRef} checked={opts.wholeWord} onChange={onOptsChange}/>
                <label for="wholeWord">WholeWord</label>
                <input type="checkbox" id="caseSensitive" ref={caseSensitiveRef} checked={opts.caseSensitive} onChange={onOptsChange}/>
                <label for="caseSensitive">CaseSensitive</label>
            </div>
            <div className="search_bar_input">
                <input type="text" ref={textRef}/>
            </div>
            <div className="search_bar_buttons">
                <button onClick={() => onNext(textRef.current.value, opts)}>
                    <img src={searchNextIcon} alt="Next" />
                </button>
                <button onClick={() => onPrev(textRef.current.value, opts)}>
                    <img src={searchPrevIcon} alt="Prev" />
                </button>
            </div>
        </div>
    );
};
