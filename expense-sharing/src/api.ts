// Resolve API base URL for both local and production
const API_BASE =
  import.meta.env.VITE_API_BASE ??
  "http://localhost:8080";

if (!API_BASE) {
  throw new Error(
    "VITE_API_BASE is not defined. Set it in .env or Vercel environment variables."
  );
}

async function handle<T>(res: Response): Promise<T> {
  if (!res.ok) {
    throw new Error(await res.text());
  }
  return (await res.json()) as T;
}

export async function get<T>(path: string): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`);
  return handle<T>(res);
}

export async function post<T>(path: string, body: unknown): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  return handle<T>(res);
}
