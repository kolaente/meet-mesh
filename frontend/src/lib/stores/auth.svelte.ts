import { api } from '$lib/api/client'

type User = { id: number; email: string; name?: string }

let user = $state<User | null>(null)
let loading = $state(true)
let checked = $state(false)

export function getAuth() {
  return {
    get user() { return user },
    get loading() { return loading },
    get isAuthenticated() { return user !== null },

    async check() {
      if (checked) return
      loading = true
      const { data } = await api.GET('/auth/me')
      user = data ?? null
      loading = false
      checked = true
    },

    async logout() {
      await api.POST('/auth/logout')
      user = null
      window.location.href = '/auth/login'
    }
  }
}
