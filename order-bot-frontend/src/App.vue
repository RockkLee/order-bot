<script setup lang="ts"></script>

<template>
  <div class="app-shell">
    <header class="app-header">
      <div class="brand">
        <span class="brand-dot" />
        <div>
          <p class="brand-title">Order Bot</p>
          <p class="brand-subtitle">C-side and B-side interfaces</p>
        </div>
      </div>
      <nav class="app-nav">
        <RouterLink to="/c">C-side</RouterLink>
        <RouterLink to="/b/login">B-side Login</RouterLink>
        <RouterLink to="/b/app">B-side App</RouterLink>
      </nav>
    </header>
    <main class="app-main">
      <RouterView />
    </main>
  </div>
</template>

<style>
@import url('https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@400;500;600;700&family=IBM+Plex+Mono:wght@400;500&display=swap');

:root {
  color-scheme: light;
  font-family: 'Space Grotesk', 'IBM Plex Mono', system-ui, sans-serif;
  line-height: 1.5;
  font-weight: 400;
  color: #10201a;
  background-color: #f7f2ea;
  text-rendering: optimizeLegibility;
  -webkit-font-smoothing: antialiased;
  --ink: #13201c;
  --muted: #5b6d65;
  --accent: #f1663b;
  --accent-strong: #d24f23;
  --mint: #d5f7e4;
  --oat: #f4ede2;
  --card: rgba(255, 255, 255, 0.78);
  --shadow: 0 22px 60px rgba(19, 32, 28, 0.18);
  --radius-lg: 24px;
  --radius-md: 16px;
}

* {
  box-sizing: border-box;
}

body {
  margin: 0;
  min-height: 100vh;
  background: radial-gradient(circle at top, #fbe6b7 0%, #f7f2ea 45%, #eef7f1 100%);
}

a {
  color: inherit;
  text-decoration: none;
}

.app-shell {
  min-height: 100vh;
  padding: 32px clamp(20px, 6vw, 80px) 64px;
  display: flex;
  flex-direction: column;
  gap: 32px;
  position: relative;
  overflow: hidden;
}

.app-shell::before,
.app-shell::after {
  content: '';
  position: absolute;
  border-radius: 999px;
  filter: blur(0);
  z-index: 0;
}

.app-shell::before {
  width: 380px;
  height: 380px;
  right: -120px;
  top: -120px;
  background: radial-gradient(circle, rgba(241, 102, 59, 0.35) 0%, transparent 70%);
}

.app-shell::after {
  width: 340px;
  height: 340px;
  left: -120px;
  bottom: -120px;
  background: radial-gradient(circle, rgba(65, 140, 112, 0.35) 0%, transparent 70%);
}

.app-header,
.app-main {
  position: relative;
  z-index: 1;
}

.app-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  padding: 18px 22px;
  border-radius: var(--radius-lg);
  background: var(--card);
  box-shadow: var(--shadow);
}

.brand {
  display: flex;
  align-items: center;
  gap: 16px;
}

.brand-dot {
  width: 14px;
  height: 14px;
  background: var(--accent);
  border-radius: 50%;
  box-shadow: 0 0 0 6px rgba(241, 102, 59, 0.15);
}

.brand-title {
  margin: 0;
  font-size: 1.2rem;
  font-weight: 700;
}

.brand-subtitle {
  margin: 0;
  color: var(--muted);
  font-size: 0.9rem;
}

.app-nav {
  display: flex;
  gap: 14px;
  flex-wrap: wrap;
}

.app-nav a {
  padding: 8px 14px;
  border-radius: 999px;
  background: #f6efe4;
  border: 1px solid rgba(19, 32, 28, 0.1);
  font-weight: 600;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.app-nav a.router-link-active {
  background: var(--accent);
  color: #fffaf3;
}

.app-nav a:hover {
  transform: translateY(-1px);
  box-shadow: 0 8px 16px rgba(19, 32, 28, 0.18);
}

.app-main {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

@media (max-width: 900px) {
  .app-header {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
