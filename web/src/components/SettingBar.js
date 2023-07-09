import { useState, useRef, useEffect } from "react";

import "./SettingBar.css"

export function SettingBar({termOpts, onOptsChange}) {

    const fontSizeRef = useRef(null);

    useEffect(() => {
        fontSizeRef.current.value = termOpts.fontSize
    }, [termOpts])

    return (
        <div className="setting_bar">
            <label for="fontSize">FontSize: </label>
            <input ref={fontSizeRef} id="fontSize" type="number" min="6" max="40" name="font size" onChange={()=> { onOptsChange({fontSize: fontSizeRef.current.value}) }}/>
        </div>
    );
}
