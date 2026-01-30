// Auth composable for managing authentication state
export const useAuth = () => {
    const user = useState<any>('user', () => null)
    const token = useState<string>('token', () => '')
    const isLoggedIn = computed(() => !!token.value && !!user.value)

    const config = useRuntimeConfig()

    // Initialize from localStorage on client
    const initAuth = () => {
        if (import.meta.client) {
            const savedToken = localStorage.getItem('auth_token')
            if (savedToken) {
                token.value = savedToken
                fetchCurrentUser()
            }
        }
    }

    // Set token after OAuth callback
    const setToken = (newToken: string) => {
        token.value = newToken
        if (import.meta.client) {
            localStorage.setItem('auth_token', newToken)
        }
        fetchCurrentUser()
    }

    // Fetch current user from API
    const fetchCurrentUser = async () => {
        if (!token.value) return

        try {
            const data = await $fetch<{ data: any }>(`${config.public.apiBase}/auth/me`, {
                headers: {
                    Authorization: `Bearer ${token.value}`
                }
            })
            user.value = data.data
        } catch (error) {
            console.error('Failed to fetch user:', error)
            logout()
        }
    }

    // Login with Google - redirect to backend OAuth
    const loginWithGoogle = () => {
        window.location.href = `${config.public.apiBase}/auth/google`
    }

    // Logout
    const logout = () => {
        user.value = null
        token.value = ''
        if (import.meta.client) {
            localStorage.removeItem('auth_token')
        }
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
        isLoggedIn,
        initAuth,
        setToken,
        loginWithGoogle,
        logout,
        getAuthHeader,
        fetchCurrentUser
    }
}
