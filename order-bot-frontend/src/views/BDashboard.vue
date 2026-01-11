<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

type Panel = 'menu' | 'events'

const router = useRouter()
const activePanel = ref<Panel>('menu')
const fileInput = ref<HTMLInputElement | null>(null)
const uploadedFile = ref('')

const setPanel = (panel: Panel) => {
  activePanel.value = panel
}

const triggerImport = () => {
  fileInput.value?.click()
}

const onFileChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  uploadedFile.value = target.files?.[0]?.name ?? ''
}

const logout = () => {
  router.push('/b/login')
}
</script>

<template>
  <section class="panel-shell">
    <header class="panel-top">
      <div>
        <p class="eyebrow">B-side</p>
        <h1>Operations control</h1>
        <p class="panel-subtitle">Manage menus, imports, and live orders.</p>
      </div>
      <div class="panel-actions">
        <input ref="fileInput" type="file" accept=".csv" hidden @change="onFileChange" />
        <button type="button" class="ghost-btn" @click="triggerImport">
          Import CSV
        </button>
        <button
          type="button"
          :class="['toggle-btn', activePanel === 'menu' ? 'is-active' : '']"
          @click="setPanel('menu')"
        >
          Show the menu
        </button>
        <button
          type="button"
          :class="['toggle-btn', activePanel === 'events' ? 'is-active' : '']"
          @click="setPanel('events')"
        >
          Show order events
        </button>
        <button type="button" class="ghost-btn" @click="logout">Log out</button>
      </div>
    </header>

    <div class="panel-content">
      <div v-if="activePanel === 'menu'" class="menu-view">
        <div class="menu-card">
          <h2>Current menu</h2>
          <p class="menu-note">Last import: {{ uploadedFile || 'No CSV uploaded yet.' }}</p>
          <ul>
            <li>
              <span>Citrus Chili Bowl</span>
              <span>Available</span>
            </li>
            <li>
              <span>Miso Maple Noodles</span>
              <span>Low stock</span>
            </li>
            <li>
              <span>Charred Corn Salad</span>
              <span>Paused</span>
            </li>
          </ul>
          <div class="menu-actions">
            <button type="button" class="primary-btn">Edit items</button>
            <button type="button" class="ghost-btn">Publish updates</button>
          </div>
        </div>
        <div class="menu-highlight">
          <p class="eyebrow">Inventory sync</p>
          <h3>14 items tracked</h3>
          <p>Auto-pause low-stock items and notify the front-of-house.</p>
        </div>
      </div>

      <div v-else class="events-view">
        <div class="event-card" v-for="event in 4" :key="event">
          <p class="event-time">{{ event }} min ago</p>
          <h3>Order #{{ 1280 + event }}</h3>
          <p>Pickup ready. Runner assigned and notified.</p>
        </div>
        <div class="event-stream">
          <h2>Live status</h2>
          <ul>
            <li>2 orders waiting on prep</li>
            <li>3 orders ready for pickup</li>
            <li>1 delayed order flagged</li>
          </ul>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.panel-shell {
  display: flex;
  flex-direction: column;
  gap: 24px;
  padding: 28px;
  border-radius: var(--radius-lg);
  background: var(--card);
  box-shadow: var(--shadow);
}

.panel-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  flex-wrap: wrap;
}

h1 {
  margin: 6px 0 4px;
  font-size: clamp(1.6rem, 2.6vw, 2.4rem);
}

.panel-subtitle {
  margin: 0;
  color: var(--muted);
}

.panel-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.toggle-btn {
  border: none;
  padding: 10px 18px;
  border-radius: 999px;
  background: #f0e6d7;
  color: var(--ink);
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.toggle-btn.is-active {
  background: var(--accent);
  color: #fff7ed;
  transform: translateY(-1px);
  box-shadow: 0 10px 18px rgba(210, 79, 35, 0.35);
}

.ghost-btn {
  border: 1px solid rgba(19, 32, 28, 0.2);
  padding: 10px 16px;
  border-radius: 999px;
  background: transparent;
  color: var(--ink);
  font-weight: 600;
  cursor: pointer;
}

.panel-content {
  display: grid;
  gap: 20px;
}

.menu-view {
  display: grid;
  gap: 20px;
  grid-template-columns: minmax(0, 2fr) minmax(0, 1fr);
}

.menu-card {
  background: #fffaf2;
  border-radius: 20px;
  padding: 20px 22px;
  display: grid;
  gap: 12px;
}

.menu-card ul {
  list-style: none;
  padding: 0;
  margin: 0;
  display: grid;
  gap: 10px;
}

.menu-card li {
  display: flex;
  justify-content: space-between;
  font-weight: 600;
}

.menu-note {
  margin: 0;
  color: var(--muted);
}

.menu-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.primary-btn {
  border: none;
  padding: 12px 18px;
  border-radius: 999px;
  background: #13201c;
  color: #fff7ed;
  font-weight: 600;
  cursor: pointer;
}

.menu-highlight {
  border-radius: 20px;
  padding: 20px;
  background: linear-gradient(135deg, #d5f7e4 0%, #fffaf2 100%);
  display: grid;
  gap: 8px;
}

.events-view {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
}

.event-card {
  border-radius: 18px;
  padding: 18px;
  background: #fff;
  box-shadow: 0 12px 20px rgba(19, 32, 28, 0.12);
}

.event-time {
  margin: 0;
  color: var(--muted);
  font-size: 0.8rem;
}

.event-stream {
  border-radius: 18px;
  padding: 18px;
  background: #f8efe2;
}

.event-stream ul {
  margin: 10px 0 0;
  padding-left: 18px;
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
  .menu-view {
    grid-template-columns: 1fr;
  }
}
</style>
