<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import { fetchApi, isControlledErr } from '@/utils/api'
import { errMenuNotFound } from '@/models/errs'

const API_BASE = import.meta.env.VITE_ORDER_BOT_MGMT_BASE_PATH ?? ''
const API_PATH_BOT = '/bot/'
const API_PATH_MENUS = '/menus/'

type Panel = 'menu' | 'events'
type MenuAct = 'create' | 'update'

type MenuItemPayload = {
  name: string
  price: number
  status: string
}

type MenuItemRes = {
  id: string
  name: string
  price: number
}

type MenuRes = {
  bot_id: string
  items: MenuItemRes[]
}

type EditableMenuItem = {
  id: number
  name: string
  price: string
  status: string
}

const router = useRouter()
const activePanel = ref<Panel>('menu')
const menuItems = ref<EditableMenuItem[]>([])
const submitState = ref<'idle' | 'submitting' | 'success' | 'error'>('idle')
const submitMessage = ref('')
const publishState = ref<'idle' | 'publishing' | 'success' | 'error'>('idle')
const publishMessage = ref('')
const botId = ref<string | null>(null)

let nextId = 1
let menuAct: MenuAct = 'update'

const toEditableItem = (menuItem: MenuItemPayload): EditableMenuItem => ({
  id: nextId++,
  name: menuItem.name,
  price: String(menuItem.price),
  status: menuItem.status,
})

const getJwt = () => {
  const jwt = localStorage.getItem('access_token')
  if (!jwt) {
    router.push('/b/login')
    return null
  }
  return jwt
}

const fetchBotId = async (jwt: string) => {
  const response = await fetchApi<undefined>(API_BASE, API_PATH_BOT, {
    method: 'GET',
    jwt,
    errMsg: 'Failed to fetch bot id',
  })
  return (await response.json()) as string
}

const loadMenu = async () => {
  submitMessage.value = ''
  publishMessage.value = ''

  try {
    const jwt = getJwt()
    if (!jwt) return

    botId.value = await fetchBotId(jwt)
    const response = await fetchApi<undefined>(API_BASE, `${API_PATH_MENUS}${botId.value}`, {
      method: 'GET',
      jwt,
      errMsg: 'Failed to fetch menu items',
    })
    const body: unknown = await response.json()
    if (isControlledErr(body)) {
      console.log(body)
      if (body.code === errMenuNotFound) {
        menuAct = 'create'
        console.log('Menu not found')
        return
      }
    }

    const fetchedMenu = body as MenuRes
    menuItems.value = fetchedMenu.items.map((item) =>
      toEditableItem({
        name: item.name,
        price: item.price,
        status: 'Available',
      }),
    )
  } catch {
    submitMessage.value = 'Failed to load menu from API.'
  }
}

const setPanel = (panel: Panel) => {
  activePanel.value = panel
}

const addRow = () => {
  menuItems.value.push({
    id: nextId++,
    name: '',
    price: '',
    status: 'Available',
  })
}

const removeRow = (id: number) => {
  menuItems.value = menuItems.value.filter((item) => item.id !== id)
}

const normalizedMenuItems = computed<MenuItemPayload[]>(() =>
  menuItems.value
    .map((item) => ({
      name: item.name.trim(),
      price: Number(item.price),
      status: item.status.trim(),
    }))
    .filter((item) => item.name && Number.isFinite(item.price) && item.status),
)

const canSubmit = computed(
  () =>
    menuItems.value.length > 0 &&
    normalizedMenuItems.value.length === menuItems.value.length &&
    submitState.value !== 'submitting' &&
    publishState.value !== 'publishing',
)

const canPublish = computed(
  () => menuItems.value.length > 0 && submitState.value !== 'submitting' && publishState.value !== 'publishing',
)

const submitMenu = async () => {
  if (!canSubmit.value) {
    submitState.value = 'error'
    submitMessage.value = 'Please complete every row before submitting the full menu.'
    return
  }

  submitState.value = 'submitting'
  submitMessage.value = 'Submitting full menu...'

  try {
    const jwt = getJwt()
    if (!jwt) return

    if (!botId.value) {
      botId.value = await fetchBotId(jwt)
    }

    const reqPayload = {
      bot_id: botId.value,
      items: normalizedMenuItems.value.map((item) => ({
        name: item.name,
        price: item.price,
      })),
    }

    console.log(`submit action: ${menuAct}`)
    await fetchApi<typeof reqPayload>(API_BASE, API_PATH_MENUS, {
      method: menuAct === 'create' ? 'POST' : 'PUT',
      jwt,
      req: reqPayload,
      wrapReq: false,
      errMsg: 'Failed to submit the full menu',
    })
    submitState.value = 'success'
    submitMessage.value = `Submitted ${normalizedMenuItems.value.length} menu items.`
  } catch {
    submitState.value = 'error'
    submitMessage.value = 'Failed to submit the full menu. Please try again.'
  }
}

const publishMenu = async () => {
  publishState.value = 'publishing'
  publishMessage.value = 'Publishing menu to c-Side...'

  try {
    const jwt = getJwt()
    if (!jwt) return

    if (!botId.value) {
      botId.value = await fetchBotId(jwt)
    }

    await fetchApi<undefined>(API_BASE, `${API_PATH_MENUS}${botId.value}/publish`, {
      method: 'POST',
      jwt,
      errMsg: 'Failed to publish menu to c-Side',
    })

    publishState.value = 'success'
    publishMessage.value = 'Published menu to c-Side.'
  } catch {
    publishState.value = 'error'
    publishMessage.value = 'Failed to publish menu to c-Side. Please try again.'
  }
}

const logout = () => {
  localStorage.removeItem('access_token')
  router.push('/b/login')
}

onMounted(loadMenu)
</script>

<template>
  <section class="panel-shell">
    <header class="panel-top">
      <div>
        <p class="eyebrow">B-side</p>
        <h1>Operations control</h1>
        <p class="panel-subtitle">Manage menu updates and live orders.</p>
      </div>
      <div class="panel-actions">
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
          <div class="menu-heading">
            <h2>Current menu</h2>
            <button type="button" class="ghost-btn" @click="addRow">Insert row</button>
          </div>
          <p class="menu-note">Update rows directly, then submit once to send the full menu.</p>
          <div class="menu-table" role="table" aria-label="Editable menu items">
            <div class="table-head" role="row">
              <span>Name</span>
              <span>Price ($)</span>
              <span>Status</span>
              <span>Action</span>
            </div>
            <div class="table-row" role="row" v-for="item in menuItems" :key="item.id">
              <input v-model="item.name" type="text" placeholder="Item name" />
              <input v-model="item.price" type="number" min="0" step="0.01" placeholder="0.00" />
              <select v-model="item.status">
                <option>Available</option>
                <!--
                <option>Low stock</option>
                <option>Paused</option>
                -->
              </select>
              <button type="button" class="danger-btn" @click="removeRow(item.id)">Delete</button>
            </div>
          </div>
          <div class="menu-actions">
            <button type="button" class="primary-btn" :disabled="!canSubmit" @click="submitMenu">
              {{ submitState === 'submitting' ? 'Submitting...' : 'Submit full menu' }}
            </button>
            <button
              type="button"
              class="icon-btn"
              :disabled="!canPublish"
              aria-label="Publish menu to c-Side"
              title="Publish menu to c-Side"
              @click="publishMenu"
            >
              <svg viewBox="0 0 24 24" aria-hidden="true">
                <path
                  d="M12 3a1 1 0 0 1 1 1v8.59l2.3-2.3a1 1 0 1 1 1.4 1.42l-4 3.99a1 1 0 0 1-1.4 0l-4-4a1 1 0 1 1 1.4-1.42l2.3 2.3V4a1 1 0 0 1 1-1Zm-7 14a1 1 0 0 1 1 1v1h12v-1a1 1 0 1 1 2 0v2a1 1 0 0 1-1 1H5a1 1 0 0 1-1-1v-2a1 1 0 0 1 1-1Z"
                />
              </svg>
            </button>
            <p :class="['submit-note', submitState]">{{ submitMessage }}</p>
            <p :class="['submit-note', publishState]">{{ publishMessage }}</p>
          </div>
        </div>
        <div class="menu-highlight">
          <p class="eyebrow">Bulk update mode</p>
          <h3>{{ menuItems.length }} rows staged</h3>
          <p>All menu rows are sent together only when you press submit.</p>
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

.ghost-btn,
.danger-btn {
  border: 1px solid rgba(19, 32, 28, 0.2);
  padding: 10px 16px;
  border-radius: 999px;
  background: transparent;
  color: var(--ink);
  font-weight: 600;
  cursor: pointer;
}

.danger-btn {
  border-color: rgba(179, 34, 34, 0.25);
  color: #b32222;
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
  gap: 14px;
}

.menu-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
}

.menu-note {
  margin: 0;
  color: var(--muted);
}

.menu-table {
  display: grid;
  gap: 10px;
}

.table-head,
.table-row {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr auto;
  gap: 10px;
  align-items: center;
}

.table-head {
  font-weight: 700;
  font-size: 0.88rem;
}

.table-row input,
.table-row select {
  width: 100%;
  padding: 10px 12px;
  border-radius: 10px;
  border: 1px solid rgba(19, 32, 28, 0.15);
  font-size: 0.95rem;
}

.menu-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  align-items: center;
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

.icon-btn {
  border: 1px solid rgba(19, 32, 28, 0.2);
  width: 42px;
  height: 42px;
  border-radius: 999px;
  background: #fff;
  color: #13201c;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.icon-btn svg {
  width: 18px;
  height: 18px;
  fill: currentColor;
}

.icon-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.primary-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.submit-note {
  margin: 0;
  color: var(--muted);
}

.submit-note.success {
  color: #1f7a4d;
}

.submit-note.error {
  color: #b32222;
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

@media (max-width: 1100px) {
  .table-head,
  .table-row {
    grid-template-columns: 1.5fr 1fr 1fr auto;
  }
}

@media (max-width: 900px) {
  .menu-view {
    grid-template-columns: 1fr;
  }

  .table-head {
    display: none;
  }

  .table-row {
    grid-template-columns: 1fr;
    background: #fff;
    padding: 12px;
    border-radius: 12px;
  }
}
</style>
