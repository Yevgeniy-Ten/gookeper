import { createContext, useContext, useState } from "react";
import { login, register } from "../api/auth";

const AuthContext = createContext();

export function AuthProvider({ children }) {
    const [token, setToken] = useState(() => localStorage.getItem("token"));
    const isAuth = !!token;

    const loginUser = async (loginStr, password) => {
        const res = await login(loginStr, password);
        localStorage.setItem("token", res.token);
        setToken(res.token);
    };

    const registerUser = async (loginStr, password) => {
        const res = await register(loginStr, password);
        return res;
    };

    const logout = () => {
        localStorage.removeItem("token");
        setToken(null);
    };

    return (
        <AuthContext.Provider
            value={{ token, isAuth, loginUser, registerUser, logout }}
        >
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth() {
    return useContext(AuthContext);
}
