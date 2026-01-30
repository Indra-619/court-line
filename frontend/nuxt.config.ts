// https://nuxt.com/docs/api/configuration/nuxt-config
import { defineNuxtConfig } from 'nuxt/config'

export default defineNuxtConfig({
    devtools: { enabled: true },

    css: ['~/assets/css/main.css'],

    runtimeConfig: {
        public: {
            apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080'
        }
    },

    app: {
        head: {
            title: 'CourtLine - Book Sports Venues Online',
            meta: [
                { charset: 'utf-8' },
                { name: 'viewport', content: 'width=device-width, initial-scale=1' },
                { name: 'description', content: 'Find and book your favorite sports courts online. Fast, easy, and hassle-free booking experience.' }
            ],
            link: [
                { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }
            ]
        }
    },

    compatibilityDate: '2024-01-01'
})
