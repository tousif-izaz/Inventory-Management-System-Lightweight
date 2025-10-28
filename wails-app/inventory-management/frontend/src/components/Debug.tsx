import { useState } from "react";

declare global {
    interface Window {
        go: {
            main: {
                App: {
                    TestDatabaseConnection: () => Promise<string>;
                    TestServerConnection: () => Promise<string>;
                    GetDatabasePath: () => Promise<string>;
                    GetServerPort: () => Promise<string>;
                    CreateSampleData: () => Promise<string>;
                };
            };
        };
    }
}

export const Debug = () => {
    const [dbTest, setDbTest] = useState<string>("");
    const [serverTest, setServerTest] = useState<string>("");
    const [sampleDataResult, setSampleDataResult] = useState<string>("");
    const [dbPath, setDbPath] = useState<string>("");
    const [serverPort, setServerPort] = useState<string>("");
    const [loading, setLoading] = useState(false);

    const runDatabaseTest = async () => {
        setLoading(true);
        try {
            const result = await window.go.main.App.TestDatabaseConnection();
            setDbTest(result);
        } catch (error) {
            setDbTest(`Error: ${error}`);
        }
        setLoading(false);
    };

    const runServerTest = async () => {
        setLoading(true);
        try {
            const result = await window.go.main.App.TestServerConnection();
            setServerTest(result);
        } catch (error) {
            setServerTest(`Error: ${error}`);
        }
        setLoading(false);
    };

    const createSampleData = async () => {
        setLoading(true);
        try {
            const result = await window.go.main.App.CreateSampleData();
            setSampleDataResult(result);
        } catch (error) {
            setSampleDataResult(`Error: ${error}`);
        }
        setLoading(false);
    };

    const getInfo = async () => {
        try {
            const path = await window.go.main.App.GetDatabasePath();
            const port = await window.go.main.App.GetServerPort();
            setDbPath(path);
            setServerPort(port);
        } catch (error) {
            console.error("Error getting info:", error);
        }
    };

    const testFetch = async () => {
        setLoading(true);
        let result = "";

        // Get auth token
        const token = localStorage.getItem('auth_token');
        result += `Auth Token: ${token ? 'Found ✅' : 'Not found ❌'}\n\n`;

        const headers: HeadersInit = {
            'Content-Type': 'application/json',
        };
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        // Test 1: Health endpoint (no auth required)
        try {
            result += "=== Test 1: Health Check ===\n";
            const healthResponse = await fetch("http://localhost:37285/health");
            const healthData = await healthResponse.json();
            result += `✅ Status: ${healthResponse.status}\n`;
            result += `Response: ${JSON.stringify(healthData, null, 2)}\n\n`;
        } catch (error) {
            result += `❌ Health check failed: ${error}\n\n`;
        }

        // Test 2: Products endpoint (auth required)
        try {
            result += "=== Test 2: Products Endpoint ===\n";
            const response = await fetch("http://localhost:37285/api/products", {
                headers,
                credentials: 'include'
            });
            result += `Status: ${response.status}\n`;

            const contentType = response.headers.get('content-type');
            result += `Content-Type: ${contentType}\n`;

            const text = await response.text();
            result += `Raw Response (first 500 chars): ${text.substring(0, 500)}\n`;

            try {
                const data = JSON.parse(text);
                if (Array.isArray(data)) {
                    result += `Parsed: Array with ${data.length} items\n`;
                    if (data.length > 0) {
                        result += `First item: ${JSON.stringify(data[0], null, 2)}\n`;
                    }
                } else {
                    result += `Parsed: ${JSON.stringify(data, null, 2)}\n`;
                }
            } catch (e) {
                result += `Failed to parse JSON: ${e}\n`;
            }
        } catch (error) {
            result += `❌ Products fetch failed: ${error}\n`;
        }

        // Test 3: Categories endpoint
        try {
            result += "\n=== Test 3: Categories Endpoint ===\n";
            const response = await fetch("http://localhost:37285/api/categories", {
                headers,
                credentials: 'include'
            });
            result += `Status: ${response.status}\n`;

            const text = await response.text();
            result += `Raw Response: ${text}\n`;

            try {
                const data = JSON.parse(text);
                if (Array.isArray(data)) {
                    result += `Parsed: Array with ${data.length} items\n`;
                } else {
                    result += `Parsed: ${JSON.stringify(data, null, 2)}\n`;
                }
            } catch (e) {
                result += `Failed to parse JSON: ${e}\n`;
            }
        } catch (error) {
            result += `❌ Categories fetch failed: ${error}\n`;
        }

        setServerTest(result);
        setLoading(false);
    };

    return (
        <div className="container mx-auto px-4 sm:px-8 py-8">
            <div className="mb-4 flex gap-4 text-sm">
                <a href="/" className="text-blue-600 hover:underline">← Home</a>
                <a href="/login" className="text-blue-600 hover:underline">Login</a>
                <a href="/signup" className="text-blue-600 hover:underline">Sign Up</a>
            </div>
            <h2 className="text-2xl font-semibold mb-6">Debug & Diagnostics</h2>

            <div className="mb-6">
                <button
                    onClick={getInfo}
                    className="rounded-lg px-4 py-2 bg-blue-600 text-white hover:bg-blue-700"
                >
                    Get System Info
                </button>
                {dbPath && (
                    <div className="mt-4 p-4 bg-gray-100 rounded">
                        <p><strong>Database Path:</strong> {dbPath}</p>
                        <p><strong>Server Port:</strong> {serverPort}</p>
                    </div>
                )}
            </div>

            <div className="mb-6">
                <button
                    onClick={runDatabaseTest}
                    disabled={loading}
                    className="rounded-lg px-4 py-2 bg-green-600 text-white hover:bg-green-700 disabled:bg-gray-400"
                >
                    Test Database Connection
                </button>
                {dbTest && (
                    <pre className="mt-4 p-4 bg-gray-100 rounded whitespace-pre-wrap font-mono text-sm">
                        {dbTest}
                    </pre>
                )}
            </div>

            <div className="mb-6">
                <button
                    onClick={runServerTest}
                    disabled={loading}
                    className="rounded-lg px-4 py-2 bg-purple-600 text-white hover:bg-purple-700 disabled:bg-gray-400"
                >
                    Test Server Connection
                </button>
                {serverTest && (
                    <pre className="mt-4 p-4 bg-gray-100 rounded whitespace-pre-wrap font-mono text-sm">
                        {serverTest}
                    </pre>
                )}
            </div>

            <div className="mb-6">
                <button
                    onClick={testFetch}
                    disabled={loading}
                    className="rounded-lg px-4 py-2 bg-orange-600 text-white hover:bg-orange-700 disabled:bg-gray-400"
                >
                    Test Frontend Fetch
                </button>
            </div>

            <div className="mb-6">
                <button
                    onClick={createSampleData}
                    disabled={loading}
                    className="rounded-lg px-4 py-2 bg-indigo-600 text-white hover:bg-indigo-700 disabled:bg-gray-400"
                >
                    Create Sample Data
                </button>
                {sampleDataResult && (
                    <pre className="mt-4 p-4 bg-gray-100 rounded whitespace-pre-wrap font-mono text-sm">
                        {sampleDataResult}
                    </pre>
                )}
            </div>

            <div className="mt-8 p-4 bg-yellow-50 border border-yellow-200 rounded">
                <h3 className="font-semibold mb-2">Expected Behavior:</h3>
                <ul className="list-disc ml-5 space-y-1">
                    <li>Database should show all tables with row counts</li>
                    <li>Server should respond with status 401 (unauthorized) for /api/products</li>
                    <li>Frontend fetch should also get 401 or proper response</li>
                </ul>
            </div>
        </div>
    );
};
