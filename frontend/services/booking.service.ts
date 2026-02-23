import { useApiFetch } from './api'

export interface Booking {
    id?: string
    courtId: string
    customerName: string
    customerPhone: string
    date: string
    startTime: string
    endTime: string
    totalPrice?: number
    status?: string
}

export const BookingService = {
    async create(booking: Booking): Promise<Booking> {
        const response = await useApiFetch<{ data: Booking }>('/api/bookings', {
            method: 'POST',
            body: booking
        })
        return response.data
    },

    async getMyBookings(): Promise<Booking[]> {
        const response = await useApiFetch<{ data: Booking[] }>('/api/bookings')
        return response.data || []
    }
}
