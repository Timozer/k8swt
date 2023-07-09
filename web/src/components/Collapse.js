import { useState } from "react";
import arrowLeftIcon from "../assets/icons/arrow-left.svg"
import arrowRightIcon from "../assets/icons/arrow-right.svg"

import "./Collapse.css"

export function Collapse({ collapsed, children, onCollapsed }) {

    const [isCollapsed, setIsCollapsed] = useState(collapsed)
    const style = {
        collapsed: {
            display: "none"
        },
        expanded: {
            display: "block"
        },
        buttonStyle: {
            display: "block",
            width: "100%"
        }
    };

    if (isCollapsed) {
        onCollapsed();
    }

    return (
        <div className="collapse">
            <button className="collapse-button" style={style.buttonStyle} onClick={() => setIsCollapsed(!isCollapsed)} >
                <img src={isCollapsed ? arrowLeftIcon : arrowRightIcon} alt={isCollapsed ? "Show" : "Hide"} />
            </button>
            <div className="collapse-content" style={isCollapsed ? style.collapsed : style.expanded}>
                {children}
            </div>
        </ div>
    );
}