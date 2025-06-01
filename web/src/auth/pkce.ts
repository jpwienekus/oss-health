export function generateCodeVerifier(): string {
  const array = new Uint32Array(56 / 2)
  window.crypto.getRandomValues(array)

  return Array.from(array, (e) => ("0" + e.toString(16)).slice(-2)).join("")
}

export async function generateCodeChallenge(codeVerifier: string): Promise<string> {
  const data = new TextEncoder().encode(codeVerifier)
  const digest = await window.crypto.subtle.digest("SHA-256", data)

  return btoa(String.fromCharCode(...new Uint8Array(digest)))
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=+$/, "")
}
