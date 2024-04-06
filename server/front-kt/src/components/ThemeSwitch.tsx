import {useEffect} from "react";
import "./theme-switch.scss";
import usePersistedState from "../hooks/usePersistedState";

// https://codesandbox.io/p/sandbox/bootstrap-dark-theme-in-react-0p1qk?file=%2Fsrc%2Fscss%2Ftheme-switch.scss%3A1%2C1-84%2C1

export default function ThemeSwitch() {
    const [darkMode, setDarkMode] = usePersistedState("darkmode", false);

    const switchTheme = () => {
        setDarkMode((prev) => {
            console.log(`upd, new value: ${prev}`);
            return !prev
        });
    }

    useEffect(() => {
        if (darkMode) {
            document.documentElement.removeAttribute("data-bs-theme");
        } else {
            document.documentElement.setAttribute("data-bs-theme", "dark");
        }
    }, [darkMode]);

    return (
        <div id="theme-switch" className="me-5">
            <div className="switch-track">
                <div className="switch-check">
                    <span className="switch-icon">🌙</span>
                </div>
                <div className="switch-x">
                    <span className="switch-icon">🌞</span>
                </div>
                <div className="switch-thumb"></div>
            </div>

            <input
                type="checkbox"
                checked={darkMode}
                onChange={switchTheme}
                aria-label="Switch between dark and light mode"
            />
        </div>
    );
}
