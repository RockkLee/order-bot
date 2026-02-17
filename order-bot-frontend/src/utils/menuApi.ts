export type MenuItemPayload = {
  name: string
  price: number
  status: string
}

const API_BASE = import.meta.env.VITE_API_BASE_URL ?? ''

const buildUrl = (path: string) => `${API_BASE}${path}`

export const fetchMenuItems = async () => {
  const response = await fetch(buildUrl('/menu-items'))

  if (!response.ok) {
    throw new Error(`Failed to fetch menu items (${response.status})`)
  }

  const data = (await response.json()) as MenuItemPayload[]
  return data
}

export const submitMenuItems = async (menuItems: MenuItemPayload[]) => {
  const response = await fetch(buildUrl('/menu-items'), {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ menuItems }),
  })

  if (!response.ok) {
    throw new Error(`Failed to submit menu items (${response.status})`)
  }

  return response
}
