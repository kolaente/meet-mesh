type Variant = 'success' | 'error' | 'info'

interface Toast {
  id: string
  variant: Variant
  message: string
}

let toasts = $state<Toast[]>([])

function generateId(): string {
  return Math.random().toString(36).substring(2, 9)
}

function addToast(variant: Variant, message: string) {
  const id = generateId()
  toasts = [...toasts, { id, variant, message }]
  return id
}

export function getToasts() {
  return {
    get all() { return toasts },

    success(message: string) {
      return addToast('success', message)
    },

    error(message: string) {
      return addToast('error', message)
    },

    info(message: string) {
      return addToast('info', message)
    },

    remove(id: string) {
      toasts = toasts.filter(t => t.id !== id)
    }
  }
}
