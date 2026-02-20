import { useApiFetch } from './api'

export interface Court {
    id: string
    name: string
    type: string
    location: string
    description: string
    pricePerHour: number
    imageUrl: string
    facilities: string[]
    isAvailable: boolean
}

export const CourtService = {
    async getAll(): Promise<Court[]> {
        const response = await useApiFetch<{ data: Court[] }>('/api/courts')
        return response.data || []
    },

    async getById(id: string): Promise<Court> {
        const response = await useApiFetch<{ data: Court }>(`/api/courts/${id}`)
        return response.data
    },

    async getBookings(id: string): Promise<any[]> {
        const response = await useApiFetch<{ data: any[] }>(`/api/courts/${id}/bookings`)
        return response.data || []
    }
}
