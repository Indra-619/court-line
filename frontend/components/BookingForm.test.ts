import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import BookingForm from './BookingForm.vue'

// Mock auto-imports
vi.stubGlobal('ref', (val: any) => ({ value: val }))
vi.stubGlobal('computed', (fn: any) => ({ get value() { return fn() } }))
vi.stubGlobal('watch', vi.fn())

// Mock useAuth
vi.mock('../composables/useAuth', () => ({
    useAuth: () => ({
        isLoggedIn: { value: true },
        user: { value: { name: 'Indra' } },
        loginWithGoogle: vi.fn()
    })
}))

// Mock BookingService
vi.mock('../services/booking.service', () => ({
    BookingService: {
        create: vi.fn()
    }
}))

describe('BookingForm.vue', () => {
    it('renders correctly with props', () => {
        const wrapper = mount(BookingForm, {
            props: {
                courtId: '1',
                pricePerHour: 100
            }
        })

        expect(wrapper.find('.form-title').text()).toBe('Book This Court')
        expect(wrapper.find('button[type="submit"]').exists()).toBe(true)
    })

    it('calculates estimated price correctly', async () => {
        const wrapper = mount(BookingForm, {
            props: {
                courtId: '1',
                pricePerHour: 120000
            }
        })

        const startTimeSelect = wrapper.find('#startTime')
        const endTimeSelect = wrapper.find('#endTime')

        await startTimeSelect.setValue('10:00')
        await endTimeSelect.setValue('12:00')

        // Price for 2 hours should be 240,000
        expect(wrapper.find('.price').text()).toContain('240.000')
    })
})
