const HAS_CLOSE_SCHEDULE = import.meta.env.VITE_HAS_CLOSE_SCHEDULE === "true"
const UTC8_OFFSET_MINUTES = 8 * 60
const OPEN_HOUR_UTC8 = 9
const CLOSE_HOUR_UTC8 = 17

export const isBusinessOpenUtc8 = (): boolean => {
  const now = new Date()
  const utcMs = now.getTime() + now.getTimezoneOffset() * 60_000
  const utc8 = new Date(utcMs + UTC8_OFFSET_MINUTES * 60_000)
  const day = utc8.getDay() // 0=Sun, 1=Mon ... 6=Sat
  const hour = utc8.getHours()
  console.log(utc8)
  console.log(day)
  console.log(hour)
  const isWeekday = day >= 1 && day <= 5
  const isWithinHours = hour >= OPEN_HOUR_UTC8 && hour < CLOSE_HOUR_UTC8

  return (isWeekday && isWithinHours) || !HAS_CLOSE_SCHEDULE
}
