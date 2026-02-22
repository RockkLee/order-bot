<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'

import { fetchApi } from '@/utils/api'

const API_BASE = import.meta.env.VITE_ORDER_BOT_BASE_PATH ?? ''
const API_PATH_CHAT = '/chat'
const API_PATH_PUBLISHED_MENU = '/chat/menu'

type Panel = 'dialogue' | 'menu'

type DialogHis = {
  incomingMsg: string
  outgoingMsg: string
}

type ChatRequest = {
  menu_id: string
  bot_id: string
  message: string
}

type ChatResponse = {
  session_id: string
  reply: string
  intent: unknown
  cart: unknown
  order_id: string | null
  menu_results: unknown[]
}

type MenuItem = {
  menu_item_id: string
  name: string
  price: number
}

const route = useRoute()
const botId = ref<string | null>(null)
const menuId = ref<string | null>(null)

const activePanel = ref<Panel>('dialogue')
const userMessage = ref('')
const isSending = ref(false)
const dialogHis = ref<DialogHis[]>([])
const menuItems = ref<MenuItem[]>([])

let sessionId = null

const setPanel = (panel: Panel) => {
  activePanel.value = panel
}

const fetchChkBotIdAndMenuId = async (): Promise<boolean> => {
  const response = await fetchApi<undefined>(
    API_BASE,
    `${API_PATH_PUBLISHED_MENU}/${botId.value}/${menuId.value}`,
    {
      method: 'GET',
      errMsg: 'Failed to check the botId and menuId.',
    },
  )
  const resJson = (await response.json()) as { exists: boolean }
  return resJson.exists
}

const fetchPublishedMenuItems = async (): Promise<void> => {
  const response = await fetchApi<undefined>(
    API_BASE,
    `${API_PATH_PUBLISHED_MENU}/${botId.value}/${menuId.value}/items`,
    {
      method: 'GET',
      errMsg: 'Failed to fetch published menu items.',
    },
  )
  const resJson = (await response.json()) as { published_menu_items: MenuItem[] }
  menuItems.value = resJson.published_menu_items
}

onMounted(async () => {
  const rawBotId = route.params.botId
  const rawMenuId = route.params.menuId
  botId.value = typeof rawBotId === 'string' ? rawBotId : null
  menuId.value = typeof rawMenuId === 'string' ? rawMenuId : null

  if (!botId.value || !menuId.value) {
    alert('Missing botId or menuId in the URL. Expected /c/{botId}/{menuId}.')
    return
  }

  try {
    let exists = await fetchChkBotIdAndMenuId()
    if (!exists) {
      alert('Invalid botId or menuId. Please check your link.')
      return
    }

    await fetchPublishedMenuItems()
  } catch (error) {
    console.error(error)
    alert(error)
  }
})

const sendMessage = async () => {
  const trimmedMessage = userMessage.value.trim()

  if (!trimmedMessage || isSending.value) {
    return
  }

  if (!botId.value || !menuId.value) {
    alert('Missing botId or menuId in the URL. Expected /c/{botId}/{menuId}.')
    return
  }

  try {
    let exists = await fetchChkBotIdAndMenuId()
    if (!exists) {
      alert('Invalid botId or menuId. Please check your link.')
      return
    }
  } catch (error) {
    console.error(error)
    alert(error)
    return
  }

  const reqJson: ChatRequest = {
    menu_id: menuId.value,
    bot_id: botId.value,
    message: trimmedMessage,
  }

  isSending.value = true

  try {
    const response = await fetchApi(API_BASE, API_PATH_CHAT, {
      method: 'POST',
      req: reqJson,
      wrapReq: false,
      headers: {
        ...(sessionId ? { 'Session-Id': `${sessionId}` } : {}),
      },
      errMsg: 'Failed to send message to order bot.',
    })

    const resJson = (await response.json()) as ChatResponse
    sessionId = resJson.session_id

    dialogHis.value.push({
      incomingMsg: reqJson.message,
      outgoingMsg: resJson.reply,
    })

    userMessage.value = ''
  } catch (error) {
    console.error(error)
    alert(error)
  } finally {
    isSending.value = false
  }
}
</script>

<template>
  <section class="panel-shell">
    <header class="panel-top">
      <div>
        <p class="eyebrow">C-side</p>
        <h1>Guest experience</h1>
        <p class="panel-subtitle">No login needed. Tap a view to switch.</p>
      </div>
      <div class="panel-actions">
        <button
          type="button"
          :class="['toggle-btn', activePanel === 'dialogue' ? 'is-active' : '']"
          @click="setPanel('dialogue')"
        >
          Dialogue
        </button>
        <button
          type="button"
          :class="['toggle-btn', activePanel === 'menu' ? 'is-active' : '']"
          @click="setPanel('menu')"
        >
          Show the menu
        </button>
      </div>
    </header>

    <div class="panel-content">
      <div v-if="activePanel === 'dialogue'" class="dialogue-view">
        <div v-for="(dialog, index) in dialogHis" :key="`${index}-${dialog.incomingMsg}`" class="dialog-row">
          <div class="chat-bubble incoming">
            <p>{{ dialog.incomingMsg }}</p>
            <span>You</span>
          </div>
          <div class="chat-bubble outgoing">
            <p>{{ dialog.outgoingMsg }}</p>
            <span>Order Bot</span>
          </div>
        </div>
        <div class="chat-input">
          <input
            v-model="userMessage"
            placeholder="Ask for a recommendation..."
            @keyup.enter="sendMessage"
          />
          <button type="button" :disabled="isSending" @click="sendMessage">
            {{ isSending ? 'Sending...' : 'Send' }}
          </button>
        </div>
      </div>

      <div v-else class="menu-view">
        <div class="menu-card">
          <h2>Today&apos;s menu</h2>
          <ul>
            <li v-for="menuItem in menuItems" :key="menuItem.menu_item_id">
              <span>{{ menuItem.name }}</span>
              <span>${{ menuItem.price }}</span>
            </li>
          </ul>
          <button type="button" class="primary-btn">Start an order</button>
        </div>
        <div class="menu-highlight">
          <p class="eyebrow">Popular right now</p>
          <h3>+120 orders today</h3>
          <p>Personalize spice, swap grains, and save favorites.</p>
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

.panel-content {
  display: grid;
  gap: 20px;
}

.dialogue-view {
  display: grid;
  gap: 16px;
}

.dialog-row {
  display: grid;
  gap: 16px;
}

.chat-bubble {
  max-width: 420px;
  padding: 14px 16px;
  border-radius: 18px;
  display: grid;
  gap: 6px;
  position: relative;
  box-shadow: 0 12px 20px rgba(19, 32, 28, 0.15);
}

.chat-bubble p {
  margin: 0;
  white-space: pre-wrap;
}

.chat-bubble span {
  font-size: 0.75rem;
  color: var(--muted);
}

.outgoing {
  background: #ffffff;
}

.incoming {
  margin-left: auto;
  background: #f7d8c8;
}

.chat-input {
  display: flex;
  gap: 12px;
  padding: 12px;
  background: #fff5e8;
  border-radius: 16px;
}

.chat-input input {
  flex: 1;
  border: none;
  background: transparent;
  font-size: 1rem;
  outline: none;
}

.chat-input button {
  border: none;
  padding: 10px 16px;
  border-radius: 999px;
  background: var(--accent);
  color: #fff7ed;
  font-weight: 600;
  cursor: pointer;
}

.chat-input button:disabled {
  opacity: 0.65;
  cursor: not-allowed;
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
  gap: 16px;
}

.menu-card ul {
  list-style: none;
  padding: 0;
  margin: 0;
  display: grid;
  gap: 12px;
}

.menu-card li {
  display: flex;
  justify-content: space-between;
  font-weight: 600;
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
