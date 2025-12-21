// Resolve API base URL for both local and production
const API_BASE =
  import.meta.env.VITE_API_BASE || "http://localhost:8080";

// Fail fast if API_BASE is somehow still invalid
if (!API_BASE) {
  throw new Error(
    "VITE_API_BASE is not defined. Set it in .env (local) or Vercel environment variables."
  );
}

// Centralized response handler
async function handle<T>(res: Response): Promise<T> {
  if (!res.ok) {
    const message = await res.text();
    throw new Error(message || "Request failed");
  }

  // Handle empty responses (e.g. 204 No Content)
  if (res.status === 204) {
    return undefined as T;
  }

  return (await res.json()) as T;
}

// GET request helper
export async function get<T>(path: string): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    method: "GET",
  });

  return handle<T>(res);
}

// POST request helper
export async function post<T>(path: string, body: unknown): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(body),
  });

  return handle<T>(res);
}
