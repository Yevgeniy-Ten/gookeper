import { useState } from "react";
import { useAuth } from "../layout/AuthContext";

export default function RegisterLogin() {
    const { loginUser, registerUser } = useAuth();
    const [isLogin, setIsLogin] = useState(true);
    const [loginStr, setLoginStr] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState("");

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError("");
        try {
            if (isLogin) {
                await loginUser(loginStr, password);
            } else {
                await registerUser(loginStr, password);
                setIsLogin(true);
            }
        } catch (err) {
            setError(err?.response?.data?.error || "Ошибка");
        }
    };

    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100">
            <form
                onSubmit={handleSubmit}
                className="bg-white p-8 rounded shadow-md w-full max-w-xs"
            >
                <h2 className="text-2xl font-bold mb-4 text-center">
                    {isLogin ? "Вход" : "Регистрация"}
                </h2>
                <input
                    className="w-full mb-3 px-3 py-2 border rounded"
                    placeholder="Логин"
                    value={loginStr}
                    onChange={(e) => setLoginStr(e.target.value)}
                    required
                />
                <input
                    className="w-full mb-3 px-3 py-2 border rounded"
                    placeholder="Пароль"
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                />
                {error && (
                    <div className="text-red-500 text-sm mb-2 text-center">
                        {error}
                    </div>
                )}
                <button
                    type="submit"
                    className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700 transition"
                >
                    {isLogin ? "Войти" : "Зарегистрироваться"}
                </button>
                <button
                    type="button"
                    className="w-full mt-2 text-blue-600 hover:underline"
                    onClick={() => setIsLogin((v) => !v)}
                >
                    {isLogin
                        ? "Нет аккаунта? Зарегистрироваться"
                        : "Уже есть аккаунт? Войти"}
                </button>
            </form>
        </div>
    );
}
