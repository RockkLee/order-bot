<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const email = ref('')
const password = ref('')

const submit = () => {
  if (!email.value || !password.value) return
  router.push('/b/app')
}
</script>

<template>
  <section class="auth-shell">
    <div class="auth-card">
      <div>
        <p class="eyebrow">B-side</p>
        <h1>Welcome back</h1>
        <p class="auth-subtitle">Sign in to manage menus and orders.</p>
      </div>
      <form class="auth-form" @submit.prevent="submit">
        <label>
          Email
          <input v-model="email" type="email" placeholder="team@orderbot.ai" />
        </label>
        <label>
          Password
          <input v-model="password" type="password" placeholder="••••••••" />
        </label>
        <button type="submit">Login</button>
      </form>
      <p class="auth-footer">
        New here? <RouterLink to="/b/signup">Create an account</RouterLink>
      </p>
    </div>
    <aside class="auth-aside">
      <h2>Live order control</h2>
      <p>Track order events, upload menus, and keep the kitchen in sync.</p>
      <div class="metric-grid">
        <div>
          <h3>2m</h3>
          <p>Avg. response time</p>
        </div>
        <div>
          <h3>98%</h3>
          <p>On-time status</p>
        </div>
      </div>
    </aside>
  </section>
</template>

<style scoped>
.auth-shell {
  display: grid;
  gap: 24px;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
}

.auth-card {
  background: var(--card);
  padding: 28px;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow);
  display: grid;
  gap: 18px;
}

h1 {
  margin: 6px 0 4px;
  font-size: clamp(1.6rem, 2.6vw, 2.4rem);
}

.auth-subtitle {
  margin: 0;
  color: var(--muted);
}

.auth-form {
  display: grid;
  gap: 14px;
}

label {
  display: grid;
  gap: 8px;
  font-weight: 600;
  color: var(--ink);
}

input {
  padding: 12px 14px;
  border-radius: 12px;
  border: 1px solid rgba(19, 32, 28, 0.15);
  font-size: 1rem;
}

button {
  border: none;
  padding: 12px 16px;
  border-radius: 999px;
  background: var(--accent);
  color: #fff7ed;
  font-weight: 600;
  cursor: pointer;
}

.auth-footer {
  margin: 0;
  color: var(--muted);
}

.auth-footer a {
  color: var(--accent-strong);
  font-weight: 600;
}

.auth-aside {
  padding: 28px;
  border-radius: var(--radius-lg);
  background: linear-gradient(140deg, #fff3dc 0%, #d5f7e4 100%);
  display: grid;
  gap: 16px;
}

.metric-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.metric-grid h3 {
  margin: 0;
  font-size: 1.5rem;
}

.metric-grid p {
  margin: 0;
  color: var(--muted);
}

.eyebrow {
  margin: 0;
  font-size: 0.75rem;
  letter-spacing: 0.2em;
  text-transform: uppercase;
  color: var(--muted);
}

@media (max-width: 900px) {
  .auth-shell {
    grid-template-columns: 1fr;
  }
}
</style>
