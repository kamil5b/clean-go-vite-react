/**
 * API Client with automatic token refresh and CSRF protection
 *
 * Features:
 * - Automatic 401 handling with token refresh
 * - CSRF token injection for state-changing requests
 * - Credentials included by default
 * - Prevents refresh loops on auth endpoints
 */

const API_BASE_URL = "/api";

// Track if we're currently refreshing to prevent multiple refresh attempts
let isRefreshing = false;
let refreshPromise: Promise<void> | null = null;

/**
 * Refresh the access token using the refresh token cookie
 */
async function refreshAccessToken(): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
        method: "POST",
        credentials: "include",
    });

    if (!response.ok) {
        throw new Error("Token refresh failed");
    }

    return response.json();
}

/**
 * Get a CSRF token from the backend
 */
async function fetchCSRFToken(): Promise<string> {
    const response = await fetch(`${API_BASE_URL}/csrf`, {
        credentials: "include",
    });

    if (!response.ok) {
        throw new Error("Failed to fetch CSRF token");
    }

    const data = await response.json();
    return data.token;
}

/**
 * Check if the URL is an auth endpoint that shouldn't trigger refresh
 */
function isAuthEndpoint(url: string): boolean {
    const authEndpoints = ["/auth/login", "/auth/register", "/auth/refresh"];
    return authEndpoints.some((endpoint) => url.includes(endpoint));
}

/**
 * Check if the request method requires CSRF protection
 */
function requiresCSRF(method?: string): boolean {
    const statefulMethods = ["POST", "PUT", "DELETE", "PATCH"];
    return statefulMethods.includes(method?.toUpperCase() || "");
}

export interface ApiClientOptions extends RequestInit {
    /**
     * Skip automatic token refresh on 401 for public endpoints
     */
    skipAuthRefresh?: boolean;
}

/**
 * Main fetch wrapper with automatic token refresh and CSRF protection
 *
 * @param url - The URL to fetch (relative to API_BASE_URL or absolute)
 * @param options - Standard fetch options with optional skipAuthRefresh
 * @returns Promise<Response>
 */
export async function apiClient(
    url: string,
    options: ApiClientOptions = {},
): Promise<Response> {
    // Ensure URL is properly formatted
    const fullUrl = url.startsWith("http") ? url : `${API_BASE_URL}${url}`;

    // Default options
    const defaultOptions: RequestInit = {
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
            ...(options.headers || {}),
        },
    };

    // Merge options
    const mergedOptions = {
        ...defaultOptions,
        ...options,
        headers: {
            ...defaultOptions.headers,
            ...(options.headers || {}),
        },
    };

    // Add CSRF token for state-changing requests
    if (requiresCSRF(mergedOptions.method)) {
        try {
            const csrfToken = await fetchCSRFToken();
            (mergedOptions.headers as Record<string, string>)["X-CSRF-Token"] =
                csrfToken;
        } catch (error) {
            console.error("Failed to fetch CSRF token:", error);
            // Continue without CSRF token - let backend reject if needed
        }
    }

    // Make the request
    let response = await fetch(fullUrl, mergedOptions);

    // Handle 401 with token refresh (only for non-auth and non-public endpoints)
    if (
        response.status === 401 &&
        !isAuthEndpoint(fullUrl) &&
        !options.skipAuthRefresh
    ) {
        // If already refreshing, wait for that to complete
        if (isRefreshing && refreshPromise) {
            try {
                await refreshPromise;
                // Retry the original request after refresh completes
                response = await fetch(fullUrl, mergedOptions);
            } catch {
                // Refresh failed, only redirect if not on public pages
                const publicPaths = ["/", "/login", "/register"];
                if (!publicPaths.includes(window.location.pathname)) {
                    window.location.href = "/login";
                }
                throw new Error("Session expired. Please login again.");
            }
        } else {
            // Start a new refresh
            isRefreshing = true;
            refreshPromise = refreshAccessToken()
                .then(() => {
                    isRefreshing = false;
                    refreshPromise = null;
                })
                .catch((error) => {
                    isRefreshing = false;
                    refreshPromise = null;
                    throw error;
                });

            try {
                await refreshPromise;
                // Retry the original request
                response = await fetch(fullUrl, mergedOptions);
            } catch {
                // Refresh failed, only redirect if not on public pages
                const publicPaths = ["/", "/login", "/register"];
                if (!publicPaths.includes(window.location.pathname)) {
                    window.location.href = "/login";
                }
                throw new Error("Session expired. Please login again.");
            }
        }
    }

    return response;
}

/**
 * Convenience wrapper for JSON responses
 */
export async function apiClientJson<T>(
    url: string,
    options: ApiClientOptions = {},
): Promise<T> {
    const response = await apiClient(url, options);

    if (!response.ok) {
        const error = await response.json().catch(() => ({
            error: response.statusText,
        }));
        throw new Error(error.error || `Request failed: ${response.status}`);
    }

    return response.json();
}

/**
 * Export API_BASE_URL for direct use if needed
 */
export { API_BASE_URL };
