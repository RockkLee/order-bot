<template>
  <section class="closed-shell">
    <div class="closed-card">
      <p class="eyebrow">Service hours</p>
      <h1>We’re closed right now</h1>
      <p class="closed-copy">
        This app is only open from 09:00–17:00 (UTC+8), Monday through Friday.
      </p>
      <p class="closed-copy muted">
        Please check back during business hours.
      </p>
    </div>
  </section>
</template>

<style scoped>
.closed-shell {
  min-height: 70vh;
  display: grid;
  place-items: center;
  padding: 32px 16px;
}

.closed-card {
  width: min(640px, 100%);
  background: var(--card);
  padding: 32px;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow);
  display: grid;
  gap: 12px;
  text-align: center;
}

h1 {
  margin: 6px 0 0;
  font-size: clamp(1.8rem, 3.2vw, 2.6rem);
}

.closed-copy {
  margin: 0;
  color: var(--ink);
  font-size: 1.05rem;
}

.muted {
  color: var(--muted);
}

.eyebrow {
  margin: 0;
  font-size: 0.75rem;
  letter-spacing: 0.2em;
  text-transform: uppercase;
  color: var(--muted);
}
</style>
