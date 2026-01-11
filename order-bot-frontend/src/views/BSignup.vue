<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const name = ref('')
const email = ref('')
const password = ref('')

const submit = () => {
  if (!name.value || !email.value || !password.value) return
  router.push('/b/app')
}
</script>

<template>
  <section class="auth-shell">
    <div class="auth-card">
      <div>
        <p class="eyebrow">B-side</p>
        <h1>Create your workspace</h1>
        <p class="auth-subtitle">Set up a team login in under a minute.</p>
      </div>
      <form class="auth-form" @submit.prevent="submit">
        <label>
          Team name
          <input v-model="name" type="text" placeholder="Order Bot Labs" />
        </label>
        <label>
          Email
          <input v-model="email" type="email" placeholder="owner@orderbot.ai" />
        </label>
        <label>
          Password
          <input v-model="password" type="password" placeholder="Create a password" />
        </label>
        <button type="submit">Sign up</button>
      </form>
      <p class="auth-footer">
        Already have access? <RouterLink to="/b/login">Sign in</RouterLink>
      </p>
    </div>
    <aside class="auth-aside">
      <h2>What you get</h2>
      <ul>
        <li>CSV menu imports with instant validation.</li>
        <li>Real-time order event timeline.</li>
        <li>Dedicated control view for the kitchen.</li>
      </ul>
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
  background: linear-gradient(140deg, #f8e0d1 0%, #fff3dc 100%);
  display: grid;
  gap: 16px;
}

.auth-aside ul {
  margin: 0;
  padding-left: 18px;
  color: var(--muted);
  display: grid;
  gap: 10px;
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
