import {Dispatch, SetStateAction, useEffect, useState} from "react";

// to store the state to local storage
export default function usePersistedState<T>(key: string, defaultValue: T): [T, Dispatch<SetStateAction<T>>] {
  const [state, setState] = useState(() => {
    const parsed: T = JSON.parse(localStorage.getItem(key) || JSON.stringify(defaultValue));
    return parsed;
  });

  useEffect(() => {
    localStorage.setItem(key, JSON.stringify(state));
  }, [key, state]);

  return [state, setState];
}
