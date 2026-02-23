export const useApiFetch = async <T>(url: string, options: any = {}) => {
    const config = useRuntimeConfig()
    const { token } = useAuth()

    const headers = {
        ...options.headers,
    }

    if (token.value) {
        headers.Authorization = `Bearer ${token.value}`
    }

    return $fetch<T>(url, {
        baseURL: config.public.apiBase,
        ...options,
        headers,
    })
}
