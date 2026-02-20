import { describe, it, expect, vi } from 'vitest'
import { CourtService } from './court.service'
import { BookingService } from './booking.service'

// Mock useApiFetch directly to avoid Nuxt dependency issues in service tests
vi.mock('./api', () => ({
    useApiFetch: vi.fn()
}))

import { useApiFetch } from './api'

describe('CourtService', () => {
    it('should fetch all courts', async () => {
        const mockData = { data: [{ id: '1', name: 'Court A' }] }
        vi.mocked(useApiFetch).mockResolvedValueOnce(mockData)

        const courts = await CourtService.getAll()
        expect(courts).toHaveLength(1)
        expect(courts[0].name).toBe('Court A')
    })

    it('should fetch court by ID', async () => {
        const mockData = { data: { id: '1', name: 'Court A' } }
        vi.mocked(useApiFetch).mockResolvedValueOnce(mockData)

        const court = await CourtService.getById('1')
        expect(court.name).toBe('Court A')
        expect(useApiFetch).toHaveBeenCalledWith('/api/courts/1')
    })

    it('should handle fetch errors', async () => {
        vi.mocked(useApiFetch).mockRejectedValueOnce(new Error('API Error'))
        await expect(CourtService.getAll()).rejects.toThrow('API Error')
    })
})

describe('BookingService', () => {
    it('should create a booking', async () => {
        const mockResponse = { data: { id: 'b1', totalPrice: 200 } }
        vi.mocked(useApiFetch).mockResolvedValueOnce(mockResponse)

        const booking = await BookingService.create({
            courtId: '1',
            customerName: 'Indra',
            customerPhone: '0812',
            date: '2024-02-21',
            startTime: '10:00',
            endTime: '12:00'
        })

        expect(booking.totalPrice).toBe(200)
        expect(vi.mocked(useApiFetch)).toHaveBeenCalledWith('/api/bookings', expect.objectContaining({ method: 'POST' }))
    })

    it('should fetch user bookings', async () => {
        const mockData = { data: [{ id: 'b1', customerName: 'Indra' }] }
        vi.mocked(useApiFetch).mockResolvedValueOnce(mockData)

        const bookings = await BookingService.getMyBookings()
        expect(bookings).toHaveLength(1)
        expect(bookings[0].customerName).toBe('Indra')
    })

    it('should handle creation errors', async () => {
        vi.mocked(useApiFetch).mockRejectedValueOnce(new Error('Court fully booked'))
        await expect(BookingService.create({} as any)).rejects.toThrow('Court fully booked')
    })
})
