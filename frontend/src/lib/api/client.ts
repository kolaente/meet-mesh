import createClient from 'openapi-fetch'
import type { paths } from './generated/schema'

export const api = createClient<paths>({
  baseUrl: '/api',
})
