<script setup lang="ts">
import { ref } from 'vue'

import { fetchApi } from '@/utils/menuApi'

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

const CHAT_PATH = '/orderbot/chat'
const MENU_ID = '26c89938-82d7-4b5f-9cbf-069bab7c66c5'
const BOT_ID = 'd9407c68-e084-4fb4-b791-55caaf7613fb'
const SESSION_ID = '0250bafa-c421-4326-b14c-d5c6837c3309'

const activePanel = ref<Panel>('dialogue')
const userMessage = ref('')
const isSending = ref(false)
const dialogHis = ref<DialogHis[]>([])

const setPanel = (panel: Panel) => {
  activePanel.value = panel
}

const sendMessage = async () => {
  const trimmedMessage = userMessage.value.trim()

  if (!trimmedMessage || isSending.value) {
    return
  }

  const reqJson: ChatRequest = {
    menu_id: MENU_ID,
    bot_id: BOT_ID,
    message: trimmedMessage,
  }

  isSending.value = true

  try {
    const response = await fetchApi(CHAT_PATH, {
      method: 'POST',
      req: reqJson,
      wrapReq: false,
      headers: {
        'Session-Id': SESSION_ID,
      },
      errMsg: 'Failed to send message to order bot.',
    })

    const resJson = (await response.json()) as ChatResponse

    dialogHis.value.push({
      incomingMsg: reqJson.message,
      outgoingMsg: resJson.reply,
    })

    userMessage.value = ''
  } catch (error) {
    console.error(error)
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
            <p>{{ dialog.outgoingMsg }}</p>
            <span>Order Bot</span>
          </div>
          <div class="chat-bubble outgoing">
            <p>{{ dialog.incomingMsg }}</p>
            <span>You</span>
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
            <li>
              <span>Citrus Chili Bowl</span>
              <span>$12</span>
            </li>
            <li>
              <span>Miso Maple Noodles</span>
              <span>$11</span>
            </li>
            <li>
              <span>Charred Corn Salad</span>
              <span>$8</span>
            </li>
            <li>
              <span>Cold Brew Spritz</span>
              <span>$5</span>
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

.chat-bubble span {
  font-size: 0.75rem;
  color: var(--muted);
}

.incoming {
  background: #ffffff;
}

.outgoing {
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
