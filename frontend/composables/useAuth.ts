// Auth composable for managing authentication state
import type { FetchOptions, FetchError } from 'ofetch'

export interface User {
    id: string
    googleId?: string
    email: string
    name: string
    picture?: string
    role?: string
    createdAt?: string
    updatedAt?: string
}

interface SessionTokens {
    data: {
        token: string
        refreshToken: string
    }
}

const isUnauthorized = (error: unknown): boolean => {
    const status = (error as FetchError | undefined)?.statusCode ?? (error as FetchError | undefined)?.status
    return status === 401
}

export const useAuth = () => {
    const user = useState<User | null>('user', () => null)
    const token = useState<string>('token', () => '')
    const refreshToken = useState<string>('refreshToken', () => '')
    const isLoggedIn = computed(() => !!token.value && !!user.value)

    const config = useRuntimeConfig()

    // Initialize from localStorage on client
    const initAuth = () => {
        if (import.meta.client) {
            const savedToken = localStorage.getItem('auth_token')
            const savedRefreshToken = localStorage.getItem('auth_refresh_token')
            if (savedToken) {
                token.value = savedToken
                refreshToken.value = savedRefreshToken || ''
                fetchCurrentUser()
            }
        }
    }

    // Set tokens after OAuth callback or session refresh
    const setToken = (newToken: string, newRefreshToken?: string) => {
        token.value = newToken
        if (import.meta.client) {
            localStorage.setItem('auth_token', newToken)
        }
        if (newRefreshToken !== undefined) {
            refreshToken.value = newRefreshToken
            if (import.meta.client) {
                localStorage.setItem('auth_refresh_token', newRefreshToken)
            }
        }
        fetchCurrentUser()
    }

    // Clear all local auth state
    const clearAuthState = () => {
        user.value = null
        token.value = ''
        refreshToken.value = ''
        if (import.meta.client) {
            localStorage.removeItem('auth_token')
            localStorage.removeItem('auth_refresh_token')
        }
    }

    // Exchange the refresh token for a new token pair (both are rotated)
    const refreshSession = async (): Promise<boolean> => {
        if (!refreshToken.value) return false

        try {
            const response = await $fetch<SessionTokens>(`${config.public.apiBase}/auth/refresh`, {
                method: 'POST',
                body: { refreshToken: refreshToken.value }
            })
            token.value = response.data.token
            refreshToken.value = response.data.refreshToken
            if (import.meta.client) {
                localStorage.setItem('auth_token', response.data.token)
                localStorage.setItem('auth_refresh_token', response.data.refreshToken)
            }
            return true
        } catch {
            clearAuthState()
            return false
        }
    }

    // Authenticated fetch: injects the Authorization header and retries once
    // with a refreshed session on 401
    const authFetch = async <T>(url: string, opts: FetchOptions = {}): Promise<T> => {
        const request = (currentToken: string): Promise<T> =>
            $fetch<T>(url, {
                ...opts,
                headers: {
                    ...opts.headers as Record<string, string> | undefined,
                    Authorization: `Bearer ${currentToken}`
                }
            })

        try {
            return await request(token.value)
        } catch (error) {
            if (!isUnauthorized(error)) throw error

            const refreshed = await refreshSession()
            if (!refreshed) {
                clearAuthState()
                throw error
            }

            try {
                return await request(token.value)
            } catch (retryError) {
                clearAuthState()
                throw retryError
            }
        }
    }

    // Fetch current user from API
    const fetchCurrentUser = async () => {
        if (!token.value) return

        try {
            const data = await authFetch<{ data: User }>(`${config.public.apiBase}/auth/me`)
            user.value = data.data
        } catch (error) {
            console.error('Failed to fetch user:', error)
            clearAuthState()
        }
    }

    // Login with Google - redirect to backend OAuth
    const loginWithGoogle = () => {
        window.location.href = `${config.public.apiBase}/auth/google`
    }

    // Logout: revoke the session server-side, then always clear local state
    const logout = async () => {
        if (token.value) {
            try {
                await $fetch(`${config.public.apiBase}/auth/logout`, {
                    method: 'POST',
                    headers: {
                        Authorization: `Bearer ${token.value}`
                    }
                })
            } catch {
                // Server revocation is best-effort; local state is cleared regardless
            }
        }
        clearAuthState()
    }

    // Get auth header for API requests
    const getAuthHeader = () => {
        if (!token.value) return {}
        return {
            Authorization: `Bearer ${token.value}`
        }
    }

    return {
        user,
        token,
        refreshToken,
        isLoggedIn,
        initAuth,
        setToken,
        refreshSession,
        authFetch,
        loginWithGoogle,
        logout,
        getAuthHeader,
        fetchCurrentUser
    }
}
