const TTL_STR = import.meta.env.VITE_LOCAL_STORAGE_TTL ?? ''
const TTL_MS = Number(TTL_STR)
const DEFAULT_TTL_MS = 0
const ttlMs = Number.isFinite(TTL_MS) && TTL_MS > 0 ? TTL_MS : DEFAULT_TTL_MS

type StoredItem<T> = {
  value: T
  expiry: number | null
}

export function setLocalStorage<T> (key: string, value: T) {
  try {
    const now = Date.now()
    const item: StoredItem<T> = {
      value,
      expiry: ttlMs > 0 ? now + ttlMs : null,
    }
    const s = JSON.stringify(item)

    localStorage.setItem(key, s)

    const readback = localStorage.getItem(key)
  } catch (e) {
    console.error("localStorage.setItem failed:", e)
  }
}

export function getLocalStorage<T> (key: string): T | null {
  console.log("getLocalStorage")
  const itemStr = localStorage.getItem(key)

  if (!itemStr) return null

  let item: StoredItem<T> | null = null
  try {
    item = JSON.parse(itemStr) as StoredItem<T>
  } catch {
    console.log(`failed to parse the localStorage item ${item}`)
    localStorage.removeItem(key)
    return null
  }

  if (!item) return null

  if (typeof item.expiry === 'number' && Date.now() > item.expiry) {
    localStorage.removeItem(key)
    return null
  }

  return item.value ?? null
}
