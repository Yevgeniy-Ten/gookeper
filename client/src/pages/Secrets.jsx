import { useState } from "react";
import { useAuth } from "../layout/AuthContext";
import {
    useSecrets,
    useCreateSecret,
    useDeleteSecret,
} from "../hooks/useSecrets";

function decodeIfNeeded(type, data) {
    if (type === "binary") return data;
    try {
        return decodeURIComponent(escape(window.atob(data)));
    } catch {
        return "[ошибка декодирования]";
    }
}

function toBase64(str) {
    return btoa(unescape(encodeURIComponent(str)));
}

function renderMeta(meta) {
    if (!meta) return null;
    if (typeof meta === "string") {
        try {
            // 1. base64 decode
            const jsonStr = decodeURIComponent(escape(window.atob(meta)));
            // 2. parse JSON
            const obj = JSON.parse(jsonStr);
            // 3. render key-value
            return (
                <ul>
                    {Object.entries(obj).map(([k, v]) => (
                        <li key={k}>
                            <span className="font-semibold">{k}:</span>{" "}
                            {String(v)}
                        </li>
                    ))}
                </ul>
            );
        } catch {
            return <span>{meta}</span>;
        }
    }
    // Если вдруг уже объект
    if (typeof meta === "object") {
        return (
            <ul>
                {Object.entries(meta).map(([k, v]) => (
                    <li key={k}>
                        <span className="font-semibold">{k}:</span> {String(v)}
                    </li>
                ))}
            </ul>
        );
    }
    return <span>{String(meta)}</span>;
}

export default function Secrets() {
    const { logout } = useAuth();
    const { data: secrets, isLoading, error } = useSecrets();
    const createSecret = useCreateSecret();
    const deleteSecret = useDeleteSecret();

    const [type, setType] = useState("password");
    const [data, setData] = useState("");
    const [metaPairs, setMetaPairs] = useState([{ key: "", value: "" }]);

    const handleMetaChange = (idx, field, value) => {
        setMetaPairs((prev) =>
            prev.map((pair, i) =>
                i === idx ? { ...pair, [field]: value } : pair
            )
        );
    };

    const addMetaPair = () =>
        setMetaPairs([...metaPairs, { key: "", value: "" }]);
    const removeMetaPair = (idx) =>
        setMetaPairs(metaPairs.filter((_, i) => i !== idx));

    const handleCreate = (e) => {
        e.preventDefault();
        const metaObj = {};
        metaPairs.forEach(({ key, value }) => {
            if (key) metaObj[key] = value;
        });

        let payloadData = data;
        if (type !== "binary") {
            payloadData = toBase64(data);
        }

        createSecret.mutate(
            {
                type,
                data: payloadData,
                meta: metaObj,
            },
            {
                onSuccess: () => {
                    setData("");
                    setMetaPairs([{ key: "", value: "" }]);
                },
            }
        );
    };

    return (
        <div className="max-w-xl mx-auto py-8">
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-2xl font-bold">Ваши секреты</h1>
                <button
                    onClick={logout}
                    className="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600"
                >
                    Выйти
                </button>
            </div>

            <form
                onSubmit={handleCreate}
                className="bg-white p-4 rounded shadow mb-6 flex flex-col gap-2"
            >
                <select
                    className="border rounded px-2 py-1"
                    value={type}
                    onChange={(e) => setType(e.target.value)}
                >
                    <option value="password">Пароль</option>
                    <option value="text">Текст</option>
                    <option value="binary">Бинарь (base64)</option>
                    <option value="card">Карта</option>
                </select>
                <input
                    className="border rounded px-2 py-1"
                    placeholder="Данные (строка или base64)"
                    value={data}
                    onChange={(e) => setData(e.target.value)}
                    required
                />
                <div>
                    <div className="mb-1 font-semibold">Мета (key-value):</div>
                    {metaPairs.map((pair, idx) => (
                        <div key={idx} className="flex gap-2 mb-1">
                            <input
                                className="border rounded px-2 py-1 flex-1"
                                placeholder="Ключ"
                                value={pair.key}
                                onChange={(e) =>
                                    handleMetaChange(idx, "key", e.target.value)
                                }
                            />
                            <input
                                className="border rounded px-2 py-1 flex-1"
                                placeholder="Значение"
                                value={pair.value}
                                onChange={(e) =>
                                    handleMetaChange(
                                        idx,
                                        "value",
                                        e.target.value
                                    )
                                }
                            />
                            {metaPairs.length > 1 && (
                                <button
                                    type="button"
                                    className="text-red-500 px-2"
                                    onClick={() => removeMetaPair(idx)}
                                >
                                    ×
                                </button>
                            )}
                        </div>
                    ))}
                    <button
                        type="button"
                        className="text-blue-600 hover:underline text-sm"
                        onClick={addMetaPair}
                    >
                        + Добавить поле
                    </button>
                </div>
                <button
                    type="submit"
                    className="bg-blue-600 text-white py-2 rounded hover:bg-blue-700"
                    disabled={createSecret.isLoading}
                >
                    Добавить секрет
                </button>
            </form>

            {isLoading && <div>Загрузка...</div>}
            {error && (
                <div className="text-red-500">
                    Ошибка загрузки: {error.message}
                </div>
            )}
            {createSecret.isError && (
                <div className="text-red-500">
                    Ошибка при создании:{" "}
                    {createSecret.error?.response?.data?.error ||
                        createSecret.error?.message}
                </div>
            )}
            {deleteSecret.isError && (
                <div className="text-red-500">
                    Ошибка при удалении:{" "}
                    {deleteSecret.error?.response?.data?.error ||
                        deleteSecret.error?.message}
                </div>
            )}

            <ul className="space-y-2">
                {secrets &&
                    secrets.map((s) => (
                        <li
                            key={s.id}
                            className="bg-white p-4 rounded shadow flex justify-between items-center"
                        >
                            <div>
                                <div className="font-bold">{s.type}</div>
                                <div className="break-all text-sm text-gray-700">
                                    {decodeIfNeeded(s.type, s.data)}
                                </div>
                                {s.type === "binary" && (
                                    <div className="text-xs text-gray-400">
                                        base64
                                    </div>
                                )}
                                <h2 className="text-lg font-bold">Meta data</h2>
                                <div className="text-xs text-gray-500">
                                    {renderMeta(s.meta)}
                                </div>
                            </div>
                            <button
                                onClick={() => deleteSecret.mutate(s.id)}
                                className="ml-4 bg-red-500 text-white px-3 py-1 rounded hover:bg-red-600"
                            >
                                Удалить
                            </button>
                        </li>
                    ))}
            </ul>
        </div>
    );
}
