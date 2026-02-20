const API_BASE = import.meta.env.VITE_API_BASE_URL ?? ''

type ApiMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'

type FetchApiOptions<T> = {
  method?: ApiMethod
  req?: T
  jwt?: string
  headers?: Record<string, string>
  wrapReq?: boolean
  errMsg: string
}

export const fetchApi = async <T>(basePath: string, path: string, options: FetchApiOptions<T>) => {
  const { method = 'PUT', req, jwt, headers, wrapReq = true, errMsg } = options

  console.log(`${basePath}${path}`)
  const response = await fetch(`${basePath}${path}`, {
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

export const fetchAuthApi = async <T>(basePath: string, path: string, options: FetchApiOptions<T>) => {
  const { method = 'PUT', req, jwt, headers, wrapReq = true, errMsg } = options

  console.log(`${basePath}${path}`)
  const response = await fetch(`${basePath}${path}`, {
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

  if (response.status === 401 || response.status === 409) {
    return response
  }
  if (!response.ok) {
    throw new Error(errMsg)
  }

  return response
}

