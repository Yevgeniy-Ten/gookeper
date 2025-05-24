import { useAuth } from "./layout/AuthContext";
import RegisterLogin from "./pages/RegisterLogin";
import Secrets from "./pages/Secrets";

function App() {
    const { isAuth } = useAuth();

    return (
        <div className="min-h-screen bg-gray-100">
            {isAuth ? <Secrets /> : <RegisterLogin />}
        </div>
    );
}

export default App;
