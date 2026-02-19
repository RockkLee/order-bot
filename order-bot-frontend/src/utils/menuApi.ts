const API_BASE = import.meta.env.VITE_API_BASE_URL ?? ''

const buildUrl = (path: string) => `${API_BASE}${path}`

type ApiMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'

type FetchApiOptions<T> = {
  method?: ApiMethod
  req?: T
  jwt?: string
  headers?: Record<string, string>
  wrapReq?: boolean
  errMsg: string
}

export const fetchApi = async <T>(path: string, options: FetchApiOptions<T>) => {
  const { method = 'PUT', req, jwt, headers, wrapReq = true, errMsg } = options

  const response = await fetch(buildUrl(path), {
    method,
    headers: {
      'Content-Type': 'application/json',
      ...(jwt ? { Authorization: `Bearer ${jwt}` } : {}),
      ...(headers ?? {}),
    },
    ...(req === undefined
      ? {}
      : { body: JSON.stringify(wrapReq ? { req } : req) }),
  })

  if (!response.ok) {
    throw new Error(errMsg)
  }

  return response
}
